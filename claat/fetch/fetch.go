// Copyright 2016-2019 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/googlecodelabs/tools/claat/fetch/drive/auth"
	"github.com/googlecodelabs/tools/claat/nodes"
	"github.com/googlecodelabs/tools/claat/parser"
	"github.com/googlecodelabs/tools/claat/types"
	"github.com/googlecodelabs/tools/claat/util"
)

const (
	// supported codelab source types must be registered parsers
	// TODO: define these in claat/parser/..., e.g. in parser/gdoc
	// alternate TODO: make this an iota-based enum?
	SrcInvalid   srcType = ""
	SrcGoogleDoc srcType = "gdoc" // Google Docs doc
	SrcMarkdown  srcType = "md"   // Markdown text

	// driveAPI is a base URL for Drive API
	driveAPI = "https://www.googleapis.com/drive/v3"

	// Minimum image size in bytes for extension detection.
	minImageSize = 11
)

// TODO: create an enum for use with "nometa" for readability's sake

// srcType is codelab source type
type srcType string

// resource is a codelab resource, loaded from local file
// or fetched from remote location.
type resource struct {
	typ  srcType       // source type
	body io.ReadCloser // resource body
	mod  time.Time     // last update of content
}

// codelab wraps types.Codelab, while adding source type
// and modified timestamp fields.
type codelab struct {
	*types.Codelab
	Typ  srcType           //  source type
	Mod  time.Time         // last modified timestamp
	Imgs map[string]string // Slurped local image paths
}

type MemoryFetcher struct {
	passMetadata map[string]bool
}

func NewMemoryFetcher(pm map[string]bool) *MemoryFetcher {
	return &MemoryFetcher{
		passMetadata: pm,
	}
}

func (m *MemoryFetcher) SlurpCodelab(rc io.ReadCloser) (*codelab, error) {
	r := &resource{
		body: rc,
		typ:  SrcMarkdown,
		mod:  time.Now(),
	}
	defer r.body.Close()

	opts := *parser.NewOptions()
	opts.PassMetadata = m.passMetadata

	clab, err := parser.Parse(string(r.typ), r.body, opts)
	if err != nil {
		return nil, err
	}

	return &codelab{
		Codelab: clab,
		Typ:     r.typ,
		Mod:     r.mod,
	}, nil
}

type Fetcher struct {
	authHelper   *auth.Helper
	authToken    string
	crcTable     *crc64.Table
	passMetadata map[string]bool
	roundTripper http.RoundTripper
}

// NewFetcher creates an instance of Fetcher.
func NewFetcher(at string, pm map[string]bool, rt http.RoundTripper) (*Fetcher, error) {
	return &Fetcher{
		authHelper:   nil,
		authToken:    at,
		crcTable:     crc64.MakeTable(crc64.ECMA),
		passMetadata: pm,
		roundTripper: rt,
	}, nil
}

// SlurpCodelab retrieves and parses codelab source.
// It takes the source, plus an auth token and a set of extra metadata to pass along.
// It returns parsed codelab and its source type.
//
// The function will also fetch and parse fragments included
// with nodes.ImportNode.
func (f *Fetcher) SlurpCodelab(src string, output string) (*codelab, error) {
	_, err := os.Stat(src)
	// Only setup oauth if this source is not a local file.
	if os.IsNotExist(err) {
		if f.authHelper == nil {
			f.authHelper, err = auth.NewHelper(f.authToken, auth.ProviderGoogle, f.roundTripper)
			if err != nil {
				return nil, err
			}
		}
	}
	res, err := f.fetch(src)
	if err != nil {
		return nil, err
	}
	defer res.body.Close()

	opts := *parser.NewOptions()
	opts.PassMetadata = f.passMetadata

	clab, err := parser.Parse(string(res.typ), res.body, opts)
	if err != nil {
		return nil, err
	}
	images := make(map[string]string)
	dir := codelabDir(output, &clab.Meta)
	imgDir := filepath.Join(dir, util.ImgDirname)
	if !isStdout(output) {
		// download or copy codelab assets to disk, and rewrite image URLs
		var nodes []nodes.Node
		for _, step := range clab.Steps {
			nodes = append(nodes, step.Content.Nodes...)
		}
		err := f.SlurpImages(src, imgDir, nodes, images)
		if err != nil {
			return nil, err
		}
	}

	// fetch imports and parse them as fragments
	var imports []*nodes.ImportNode
	for _, st := range clab.Steps {
		imports = append(imports, nodes.ImportNodes(st.Content.Nodes)...)
	}
	ch := make(chan error, len(imports))
	defer close(ch)
	for _, imp := range imports {
		go func(n *nodes.ImportNode) {
			frag, err := f.slurpFragment(n.URL)
			if err != nil {
				ch <- fmt.Errorf("%s: %v", n.URL, err)
				return
			}
			if !isStdout(output) {
				// download or copy codelab assets to disk, and rewrite image URLs
				err = f.SlurpImages(gdocID(n.URL), imgDir, frag, images)
				if err != nil {
					return
				}
			}
			n.Content.Nodes = frag
			ch <- nil
		}(imp)
	}
	for range imports {
		if err := <-ch; err != nil {
			return nil, err
		}
	}

	v := &codelab{
		Codelab: clab,
		Typ:     res.typ,
		Mod:     res.mod,
		Imgs:    images,
	}
	return v, nil
}

func (f *Fetcher) SlurpImages(src, dir string, n []nodes.Node, images map[string]string) error {
	// make sure img dir exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	type res struct {
		url, file string
		err       error
	}

	ch := make(chan *res, 100)
	defer close(ch)
	var count int
	imageNodes := nodes.ImageNodes(n)
	count += len(imageNodes)
	for _, imageNode := range imageNodes {
		go func(imageNode *nodes.ImageNode) {
			url := imageNode.Src
			file, err := f.slurpBytes(src, dir, url)
			if err == nil {
				imageNode.Src = filepath.Join(util.ImgDirname, file)
			}
			ch <- &res{url, file, err}
		}(imageNode)
	}
	var errStr string
	for i := 0; i < count; i++ {
		r := <-ch
		images[r.file] = r.url
		if r.err != nil {
			errStr += fmt.Sprintf("%s => %s: %v\n", r.url, r.file, r.err)
		}
	}
	if len(errStr) > 0 {
		return errors.New(errStr)
	}

	return nil
}

func (f *Fetcher) slurpBytes(codelabSrc, dir, imgURL string) (string, error) {
	// images can be local in Markdown cases or remote.
	// Only proceed a simple copy on local reference.
	var b []byte
	var ext string
	u, err := url.Parse(imgURL)
	if err != nil {
		return "", err
	}

	// If the codelab source is being downloaded from the network, then we should interpret
	// the image URL in the same way.
	srcUrl, err := url.Parse(codelabSrc)
	if err == nil && srcUrl.Host != "" {
		u = srcUrl.ResolveReference(u)
	}

	if u.Host == "" {
		if imgURL, err = restrictPathToParent(imgURL, filepath.Dir(codelabSrc)); err != nil {
			return "", err
		}
		if b, err = ioutil.ReadFile(imgURL); err != nil {
			return "", err
		}
		ext = filepath.Ext(imgURL)
	} else {
		if b, err = f.slurpRemoteBytes(u.String(), 5); err != nil {
			return "", fmt.Errorf("Error downloading image at %s: %v", u.String(), err)
		}
		if ext, err = imgExtFromBytes(b); err != nil {
			return "", fmt.Errorf("Error reading image type at %s: %v", u.String(), err)
		}
	}

	crc := crc64.Checksum(b, f.crcTable)
	file := fmt.Sprintf("%x%s", crc, ext)
	dst := filepath.Join(dir, file)
	return file, ioutil.WriteFile(dst, b, 0644)
}

func (f *Fetcher) slurpFragment(url string) ([]nodes.Node, error) {
	res, err := f.fetch(url)
	if err != nil {
		return nil, err
	}
	defer res.body.Close()

	opts := *parser.NewOptions()
	opts.PassMetadata = f.passMetadata

	return parser.ParseFragment(string(res.typ), res.body, opts)
}

// fetch retrieves codelab doc either from local disk
// or a remote location.
// The caller is responsible for closing returned stream.
func (f *Fetcher) fetch(name string) (*resource, error) {
	fi, err := os.Stat(name)
	if os.IsNotExist(err) {
		return f.fetchRemote(name, false)
	}
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return &resource{
		body: r,
		typ:  SrcMarkdown,
		mod:  fi.ModTime(),
	}, nil
}

// fetchRemote retrieves resource r from the network.
//
// If urlStr is not a URL, i.e. does not have the host part, it is considered to be
// a Google Doc ID and fetched accordingly. Otherwise, a simple GET request
// is used to retrieve the contents.
//
// The caller is responsible for closing returned stream.
// If nometa is true, resource.mod may have zero value.
func (f *Fetcher) fetchRemote(urlStr string, nometa bool) (*resource, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if u.Host == "" || u.Host == "docs.google.com" {
		return f.fetchDriveFile(urlStr, nometa)
	}
	return f.fetchRemoteFile(urlStr)
}

// fetchRemoteFile retrieves codelab resource from url.
// It is a special case of fetchRemote function.
func (f *Fetcher) fetchRemoteFile(url string) (*resource, error) {
	res, err := retryGet(f.authHelper.DriveClient(), url, 3)
	if err != nil {
		return nil, err
	}
	t, err := http.ParseTime(res.Header.Get("last-modified"))
	if err != nil {
		t = time.Now()
	}
	return &resource{
		body: res.Body,
		mod:  t,
		typ:  SrcMarkdown,
	}, nil
}

// fetchDriveFile uses Drive API to retrieve HTML representation of a Google Doc.
// See https://developers.google.com/drive/web/manage-downloads#downloading_google_documents
// for more details.
//
// If nometa is true, resource.mod will have zero value.
func (f *Fetcher) fetchDriveFile(id string, nometa bool) (*resource, error) {
	id = gdocID(id)
	exportURL := gdocExportURL(id)

	if nometa {
		res, err := retryGet(f.authHelper.DriveClient(), exportURL, 7)
		if err != nil {
			return nil, err
		}
		return &resource{body: res.Body, typ: SrcGoogleDoc}, nil
	}

	q := url.Values{
		"fields":             {"id,mimeType,modifiedTime"},
		"supportsTeamDrives": {"true"},
	}
	u := fmt.Sprintf("%s/files/%s?%s", driveAPI, id, q.Encode())
	res, err := retryGet(f.authHelper.DriveClient(), u, 7)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	meta := &struct {
		ID       string    `json:"id"`
		MimeType string    `json:"mimeType"`
		Modified time.Time `json:"modifiedTime"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(meta); err != nil {
		return nil, err
	}
	if meta.MimeType != "application/vnd.google-apps.document" {
		return nil, fmt.Errorf("%s: invalid mime type: %s", id, meta.MimeType)
	}

	if res, err = retryGet(f.authHelper.DriveClient(), exportURL, 7); err != nil {
		return nil, err
	}
	return &resource{
		body: res.Body,
		mod:  meta.Modified,
		typ:  SrcGoogleDoc,
	}, nil
}

func (f *Fetcher) slurpRemoteBytes(url string, n int) ([]byte, error) {
	res, err := retryGet(f.authHelper.DriveClient(), url, n)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// retryGet tries to GET specified url up to n times.
// Attempts are spaced out with exponential backoff.
// Default client will be used if not provided.
func retryGet(client *http.Client, url string, n int) (*http.Response, error) {
	if client == nil {
		client = http.DefaultClient
	}
	for i := 0; i <= n; i++ {
		if i > 0 {
			t := time.Duration((math.Pow(2, float64(i)) + rand.Float64()) * float64(time.Second))
			time.Sleep(t)
		}
		res, err := client.Get(url)
		// return early with a good response
		// the rest is error handling
		if err == nil && res.StatusCode == http.StatusOK {
			return res, nil
		}

		// sometimes Drive API wouldn't even start a response,
		// we get net/http: TLS handshake timeout instead:
		// consider this a temporary failure and retry again
		if err != nil {
			continue
		}
		// otherwise, decode error response and check for "rate limit"
		defer res.Body.Close()
		var erres struct {
			Error struct {
				Errors []struct{ Reason string }
			}
		}
		b, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(b, &erres)
		var rateLimit bool
		for _, e := range erres.Error.Errors {
			if e.Reason == "rateLimitExceeded" || e.Reason == "userRateLimitExceeded" {
				rateLimit = true
				break
			}
		}
		// this is neither a rate limit error, nor a server error:
		// retrying is useless
		if !rateLimit && res.StatusCode < http.StatusInternalServerError {
			return nil, fmt.Errorf("fetch %s: %s; %s", url, res.Status, b)
		}
	}
	return nil, fmt.Errorf("%s: failed after %d retries", url, n)
}

func gdocID(url string) string {
	const s = "/document/d/"
	if i := strings.Index(url, s); i >= 0 {
		url = url[i+len(s):]
	}
	if i := strings.IndexRune(url, '/'); i > 0 {
		url = url[:i]
	}
	return url
}

func gdocExportURL(id string) string {
	q := url.Values{
		"mimeType": {"text/html"},
	}
	return fmt.Sprintf("%s/files/%s/export?%s", driveAPI, id, q.Encode())
}

// restrictPathToParent will ensure that assetPath is in parent.
// It will thus return an absolute path to the asset.
func restrictPathToParent(assetPath, parent string) (string, error) {
	parent, err := filepath.Abs(parent)
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(assetPath) {
		assetPath = filepath.Join(parent, assetPath)
	}
	if !strings.HasPrefix(assetPath, parent) {
		return "", fmt.Errorf("%s isn't a subdirectory of %s", assetPath, parent)
	}
	return assetPath, nil
}

// isStdout reports whether filename is stdout.
func isStdout(filename string) bool {
	const stdout = "-"
	return filename == stdout
}

// codelabDir returns codelab root directory.
// The base argument is codelab parent directory.
func codelabDir(base string, m *types.Meta) string {
	return filepath.Join(base, m.ID)
}

func imgExtFromBytes(b []byte) (string, error) {
	if len(b) < minImageSize {
		return "", fmt.Errorf("error parsing image - response \"%s\" is too small (< %d bytes)", b, minImageSize)
	}
	ext := ".png"
	switch {
	case string(b[6:10]) == "JFIF":
		ext = ".jpeg"
	case string(b[0:3]) == "GIF":
		ext = ".gif"
	}
	return ext, nil
}
