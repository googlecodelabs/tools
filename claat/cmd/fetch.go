// Copyright 2016 Google Inc. All Rights Reserved.
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

package cmd

import (
	"encoding/json"
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

	"github.com/googlecodelabs/tools/claat/parser"
	"github.com/googlecodelabs/tools/claat/types"
)

const (
	// supported codelab source types must be registered parsers
	// TODO: define these in claat/parser/..., e.g. in parser/gdoc
	srcInvalid   srcType = ""
	srcGoogleDoc srcType = "gdoc" // Google Docs doc
	srcMarkdown  srcType = "md"   // Markdown text

	// driveAPI is a base URL for Drive API
	driveAPI = "https://www.googleapis.com/drive/v3"
)

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
	typ srcType   //  source type
	mod time.Time // last modified timestamp
}

// slurpCodelab retrieves and parses codelab source.
// It takes the source, plus an auth token and a set of extra metadata to pass along.
// It returns parsed codelab and its source type.
//
// The function will also fetch and parse fragments included
// with types.ImportNode.
func slurpCodelab(src, authToken string, passMetadata map[string]bool) (*codelab, error) {
	res, err := fetch(src, authToken)
	if err != nil {
		return nil, err
	}
	defer res.body.Close()

	opts := *parser.NewOptions()
	opts.PassMetadata = passMetadata

	clab, err := parser.Parse(string(res.typ), res.body, opts)
	if err != nil {
		return nil, err
	}

	// fetch imports and parse them as fragments
	var imports []*types.ImportNode
	for _, st := range clab.Steps {
		imports = append(imports, importNodes(st.Content.Nodes)...)
	}
	ch := make(chan error, len(imports))
	defer close(ch)
	for _, imp := range imports {
		go func(n *types.ImportNode) {
			frag, err := slurpFragment(n.URL, authToken)
			if err != nil {
				ch <- fmt.Errorf("%s: %v", n.URL, err)
				return
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
		typ:     res.typ,
		mod:     res.mod,
	}
	return v, nil
}

func slurpFragment(url, authToken string) ([]types.Node, error) {
	res, err := fetchRemote(url, authToken, true)
	if err != nil {
		return nil, err
	}
	defer res.body.Close()
	return parser.ParseFragment(string(res.typ), res.body)
}

// fetch retrieves codelab doc either from local disk
// or a remote location.
// The caller is responsible for closing returned stream.
func fetch(name, authToken string) (*resource, error) {
	fi, err := os.Stat(name)
	if os.IsNotExist(err) {
		return fetchRemote(name, authToken, false)
	}
	r, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return &resource{
		body: r,
		typ:  srcMarkdown,
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
func fetchRemote(urlStr, authToken string, nometa bool) (*resource, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if u.Host == "" || u.Host == "docs.google.com" {
		return fetchDriveFile(urlStr, authToken, nometa)
	}
	return fetchRemoteFile(urlStr)
}

// fetchRemoteFile retrieves codelab resource from url.
// It is a special case of fetchRemote function.
func fetchRemoteFile(url string) (*resource, error) {
	res, err := retryGet(nil, url, 3)
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
		typ:  srcMarkdown,
	}, nil
}

// fetchDriveFile uses Drive API to retrieve HTML representation of a Google Doc.
// See https://developers.google.com/drive/web/manage-downloads#downloading_google_documents
// for more details.
//
// If nometa is true, resource.mod will have zero value.
func fetchDriveFile(id, authToken string, nometa bool) (*resource, error) {
	id = gdocID(id)
	exportURL := gdocExportURL(id)
	client, err := driveClient(authToken)
	if err != nil {
		return nil, err
	}

	if nometa {
		res, err := retryGet(client, exportURL, 7)
		if err != nil {
			return nil, err
		}
		return &resource{body: res.Body, typ: srcGoogleDoc}, nil
	}

	q := url.Values{
		"fields":             {"id,mimeType,modifiedTime"},
		"supportsTeamDrives": {"true"},
	}
	u := fmt.Sprintf("%s/files/%s?%s", driveAPI, id, q.Encode())
	res, err := retryGet(client, u, 7)
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

	if res, err = retryGet(client, exportURL, 7); err != nil {
		return nil, err
	}
	return &resource{
		body: res.Body,
		mod:  meta.Modified,
		typ:  srcGoogleDoc,
	}, nil
}

var crcTable = crc64.MakeTable(crc64.ECMA)

func slurpBytes(client *http.Client, codelabSrc, dir, imgURL string) (string, error) {
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
		b, err = ioutil.ReadFile(imgURL)
		ext = filepath.Ext(imgURL)
	} else {
		b, err = slurpRemoteBytes(client, u.String(), 5)
		if string(b[6:10]) == "JFIF" {
			ext = ".jpeg"
		} else if string(b[0:3]) == "GIF" {
			ext = ".gif"
		} else {
			ext = ".png"
		}
	}
	if err != nil {
		return "", err
	}

	crc := crc64.Checksum(b, crcTable)
	file := fmt.Sprintf("%x%s", crc, ext)
	dst := filepath.Join(dir, file)
	return file, ioutil.WriteFile(dst, b, 0644)
}

func slurpRemoteBytes(client *http.Client, url string, n int) ([]byte, error) {
	res, err := retryGet(client, url, n)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// retryGet tries to GET specified url up to n times.
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
	return fmt.Sprintf("%s/files/%s/export?mimeType=text/html", driveAPI, id)
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
