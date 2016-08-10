package bundle

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/googlecodelabs/tools/claat/render"
	"github.com/googlecodelabs/tools/claat/types"

	"golang.org/x/net/context"
)

const (
	// TODO: make it configurable?
	assetsDir = "img"

	// driveAPIBase is a base URL for Drive API v2
	driveAPIBase = "https://www.googleapis.com/drive/v2"
	// TODO: this v2, replace with
	// https://www.googleapis.com/drive/v3/files/id/export?mimeType=text/html
	driveAPIExport = "https://docs.google.com/feeds/download/documents/export/Export"

	driveMimeDocument = "application/vnd.google-apps.document"
	driveMimeFolder   = "application/vnd.google-apps.folder"
	driveExportMime   = "text/html"
)

type HTTPClientFunc func(ctx context.Context, url string) *http.Client

type Exporter struct {
	Env             string
	Prefix          string
	GoogleAnalytics string
	Extra           map[string]string
	NoAssets        bool
	NoMeta          bool
	NoLocal         bool
	HTTPClient      HTTPClientFunc
}

func (e *Exporter) Export(ctx context.Context, dst ContentWriter, src string, of types.MarkupFormat) error {
	clab, err := e.Fetch(ctx, src)
	if err != nil {
		return err
	}

	// process codelab's assets, if requested
	// need to do this before exportMarkup, because exportAssets rewrites image URLs
	// in the markup
	if !e.NoAssets {
		if err := e.exportAssets(ctx, dst, clab); err != nil {
			return err
		}
	}
	// create metadata, if requested
	if !e.NoMeta {
		if err := e.exportMeta(ctx, dst, clab, of); err != nil {
			return err
		}
	}
	// finally, instructions markup
	return e.exportMarkup(ctx, dst, clab, of)
}

type Codelab struct {
	*types.Codelab
	Src string             // source location
	Fmt types.MarkupFormat // source format
	Mod time.Time          // last modified timestamp
}

// Fetch retrieves and parses codelab source from url.
// It returns parsed codelab and its source type.
//
// The function will also fetch and parse fragments included
// with types.ImportNode.
func (e *Exporter) Fetch(ctx context.Context, url string) (*Codelab, error) {
	res, err := e.fetch(ctx, url, false)
	if err != nil {
		return nil, err
	}
	defer res.body.Close()
	clab, err := parser.Parse(res.fmt, res.body)
	if err != nil {
		return nil, err
	}

	// fetch imports and parse them as fragments, in parallel
	var imports []*types.ImportNode
	for _, st := range clab.Steps {
		imports = append(imports, importNodes(st.Content.Nodes)...)
	}
	ch := make(chan error, len(imports)) // goroutines error reporting
	for _, imp := range imports {
		go func(n *types.ImportNode) {
			frag, err := e.fetchFragment(ctx, n.URL)
			if err != nil {
				err = fmt.Errorf("%s: %v", n.URL, err)
			}
			n.Content.Nodes = frag
			ch <- err
		}(imp)
	}

	// wait for all goroutines
	done := make(chan struct{})
	var me multiError
	go func() {
		for range imports {
			if err := <-ch; err != nil {
				me = append(me, err)
			}
		}
		close(ch)
		close(done)
	}()
	select {
	case <-ctx.Done():
		// don't worry about other errors and channels
		// the waiting goroutine will eventually close them
		return nil, ctx.Err()
	case <-done:
		// finished all imports
	}
	if len(me) > 0 {
		return nil, me
	}
	return &Codelab{
		Codelab: clab,
		Src:     url,
		Fmt:     res.fmt,
		Mod:     res.mod,
	}, nil
}

func (e *Exporter) exportMarkup(ctx context.Context, dst ContentWriter, clab *Codelab, of types.MarkupFormat) error {
	data := &render.Context{
		Env:      e.Env,
		Prefix:   e.Prefix,
		GlobalGA: e.GoogleAnalytics,
		Meta:     &clab.Meta,
		Steps:    clab.Steps,
		Extra:    e.Extra,
	}
	var buf bytes.Buffer
	if err := render.Execute(&buf, of, data); err != nil {
		return err
	}
	return dst.WriteMarkup(ctx, clab.ID, buf.Bytes())
}

func (e *Exporter) exportMeta(ctx context.Context, dst ContentWriter, clab *Codelab, of types.MarkupFormat) error {
	lastmod := types.ContextTime(clab.Mod)
	cm := &types.ContextMeta{
		Meta: clab.Meta,
		Context: types.Context{
			Env:     e.Env,
			Prefix:  e.Prefix,
			MainGA:  e.GoogleAnalytics,
			Format:  of,
			Source:  clab.Src,
			Updated: &lastmod,
		},
	}
	return dst.WriteMeta(ctx, cm)
}

func (e *Exporter) exportAssets(ctx context.Context, dst ContentWriter, clab *Codelab) error {
	// rewrite image URLs
	// imap keys are file names mapped to image URLs
	imap := make(map[string]string)
	for _, st := range clab.Steps {
		nodes := imageNodes(st.Content.Nodes)
		for _, n := range nodes {
			file := fmt.Sprintf("%x.png", md5.Sum([]byte(n.Src)))
			imap[file] = n.Src
			n.Src = filepath.Join(assetsDir, file)
		}
	}

	// fetch all images in parallel
	ch := make(chan error, len(imap)) // goroutines error reporting
	for name, url := range imap {
		go func(name, url string) {
			res, err := retryGet(ctx, e.httpClient(ctx, url), url)
			if err != nil {
				ch <- err
				return
			}
			defer res.Body.Close()
			ch <- dst.WriteAsset(ctx, clab.ID, name, res.Body)
		}(name, url)
	}

	// wait for all goroutines to finish
	done := make(chan struct{})
	var me multiError
	go func() {
		for range imap {
			if err := <-ch; err != nil {
				me = append(me, err)
			}
		}
		close(ch)
		close(done)
	}()
	select {
	case <-ctx.Done():
		// don't worry about other errors and channels
		// the waiting goroutine will eventually close them
		return ctx.Err()
	case <-done:
		// finished all fetches
	}
	if len(me) > 0 {
		return me
	}
	return nil
}

// rawCodelab represents a codelab markup, loaded from local file
// or fetched from remote location.
// It is used internally to pass around codelab's markup bytes.
type rawCodelab struct {
	fmt  types.MarkupFormat // source type
	body io.ReadCloser      // resource body
	mod  time.Time          // last content update
}

func (e *Exporter) fetch(ctx context.Context, url string, fragment bool) (*rawCodelab, error) {
	const local = "file://"
	if strings.HasPrefix(url, local) {
		if e.NoLocal {
			return nil, errors.New("local resources disabled")
		}
		return fetchLocal(strings.TrimPrefix(url, local))
	}
	return fetchRemote(ctx, e.httpClient(ctx, url), url, fragment)
}

func (e *Exporter) fetchFragment(ctx context.Context, url string) ([]types.Node, error) {
	res, err := e.fetch(ctx, url, true)
	if err != nil {
		return nil, err
	}
	defer res.body.Close()
	return parser.ParseFragment(res.fmt, res.body)
}

func (e *Exporter) httpClient(ctx context.Context, url string) *http.Client {
	var client *http.Client
	if e.HTTPClient != nil {
		client = e.HTTPClient(ctx, url)
	}
	if client == nil {
		client = http.DefaultClient
	}
	return client
}

func fetchLocal(filename string) (*rawCodelab, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	raw := &rawCodelab{
		body: f,
		mod:  time.Now(),
		fmt:  types.FmtMarkdown,
	}
	if fi, err := f.Stat(); err == nil {
		raw.mod = fi.ModTime()
	}
	if filepath.Ext(filename) == ".html" {
		raw.fmt = types.FmtHTML
	}
	return raw, nil
}

// fetchRemote retrieves codelab's markup from urlStr url.
//
// If urlStr is not a URL, i.e. does not have the host part, it is considered to be
// a Google Doc ID and fetched accordingly. Otherwise, a simple GET request
// is used to retrieve the contents.
//
// The caller is responsible for closing returned stream.
// If fragment is true, resource.mod may have zero value.
func fetchRemote(ctx context.Context, client *http.Client, urlStr string, fragment bool) (*rawCodelab, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if u.Host == "" || u.Host == "docs.google.com" {
		return fetchDriveFile(ctx, client, urlStr, fragment)
	}
	return fetchRemoteFile(ctx, client, urlStr)
}

// fetchRemoteFile retrieves codelab resource from an arbitrary url.
func fetchRemoteFile(ctx context.Context, client *http.Client, url string) (*rawCodelab, error) {
	res, err := retryGet(ctx, client, url)
	if err != nil {
		return nil, err
	}
	t, err := http.ParseTime(res.Header.Get("last-modified"))
	if err != nil {
		t = time.Now()
	}
	raw := &rawCodelab{
		body: res.Body,
		mod:  t,
		fmt:  types.FmtMarkdown,
	}
	if filepath.Ext(url) == ".html" {
		raw.fmt = types.FmtHTML
	}
	return raw, nil
}

// fetchDriveFile uses Drive API to retrieve HTML representation of a Google Doc.
// See https://developers.google.com/drive/web/manage-downloads#downloading_google_documents
// for more details.
//
// If nometa is true, resource.mod will have zero value.
func fetchDriveFile(ctx context.Context, client *http.Client, id string, fragment bool) (*rawCodelab, error) {
	id = gdocID(id)

	if fragment {
		q := url.Values{"id": {id}, "exportFormat": {"html"}}
		u := fmt.Sprintf("%s?%s", driveAPIExport, q.Encode())
		res, err := retryGet(ctx, client, u)
		if err != nil {
			return nil, err
		}
		return &rawCodelab{body: res.Body, fmt: types.FmtGoogleDoc}, nil
	}

	u := fmt.Sprintf("%s/files/%s?fields=id,mimeType,exportLinks,modifiedDate", driveAPIBase, id)
	res, err := retryGet(ctx, client, u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	meta := &struct {
		ID          string            `json:"id"`
		MimeType    string            `json:"mimeType"`
		ExportLinks map[string]string `json:"exportLinks"`
		Modified    time.Time         `json:"modifiedDate"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(meta); err != nil {
		return nil, err
	}
	if meta.MimeType != driveMimeDocument {
		return nil, fmt.Errorf("%s: invalid mime type: %s", id, meta.MimeType)
	}
	link := meta.ExportLinks[driveExportMime]
	if link == "" {
		return nil, fmt.Errorf("%s: no %q export link", id, driveExportMime)
	}

	if res, err = retryGet(ctx, client, link); err != nil {
		return nil, err
	}
	return &rawCodelab{
		body: res.Body,
		mod:  meta.Modified,
		fmt:  types.FmtGoogleDoc,
	}, nil
}

// retryGet tries to GET specified url until a successful response is received,
// an unrecoverable error occurs or ctx is cancelled.
func retryGet(ctx context.Context, client *http.Client, url string) (*http.Response, error) {
	for i := 0; ; i++ {
		if i > 0 {
			t := time.Duration((math.Pow(2, float64(i)) + rand.Float64()) * float64(time.Second))
			select {
			case <-ctx.Done():
				// cancelled by the caller
				return nil, ctx.Err()
			case <-time.After(t):
				// do retry
			}
		}
		res, err := client.Get(url)
		if err == nil && res.StatusCode == http.StatusOK {
			return res, nil
		}
		if err != nil {
			// can't continue without res.Body
			continue
		}
		defer res.Body.Close()

		// Google APIs error response format
		var erres struct {
			Error struct {
				Errors []struct{ Reason string }
			}
		}
		// don't care about ReadAll and Unmarshal errors:
		// there's nothing we can do more than a generic error handling below
		b, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(b, &erres)
		var rateLimit bool
		for _, e := range erres.Error.Errors {
			if e.Reason == "rateLimitExceeded" || e.Reason == "userRateLimitExceeded" {
				rateLimit = true
				break
			}
		}
		if !rateLimit && res.StatusCode < 500 {
			// this is neither a rate limit error, nor a server error:
			// retrying is useless
			return nil, fmt.Errorf("retry(%d) %s: %s; %s", i, url, res.Status, b)
		}
	}
}

// imageNodes filters out everything except types.NodeImage nodes, recursively.
func imageNodes(nodes []types.Node) []*types.ImageNode {
	var imgs []*types.ImageNode
	for _, n := range nodes {
		switch n := n.(type) {
		case *types.ImageNode:
			imgs = append(imgs, n)
		case *types.ListNode:
			imgs = append(imgs, imageNodes(n.Nodes)...)
		case *types.ItemsListNode:
			for _, i := range n.Items {
				imgs = append(imgs, imageNodes(i.Nodes)...)
			}
		case *types.HeaderNode:
			imgs = append(imgs, imageNodes(n.Content.Nodes)...)
		case *types.URLNode:
			imgs = append(imgs, imageNodes(n.Content.Nodes)...)
		case *types.ButtonNode:
			imgs = append(imgs, imageNodes(n.Content.Nodes)...)
		case *types.InfoboxNode:
			imgs = append(imgs, imageNodes(n.Content.Nodes)...)
		case *types.GridNode:
			for _, r := range n.Rows {
				for _, c := range r {
					imgs = append(imgs, imageNodes(c.Content.Nodes)...)
				}
			}
		}
	}
	return imgs
}

// importNodes filters out everything except types.NodeImport nodes, recursively.
func importNodes(nodes []types.Node) []*types.ImportNode {
	var imps []*types.ImportNode
	for _, n := range nodes {
		switch n := n.(type) {
		case *types.ImportNode:
			imps = append(imps, n)
		case *types.ListNode:
			imps = append(imps, importNodes(n.Nodes)...)
		case *types.InfoboxNode:
			imps = append(imps, importNodes(n.Content.Nodes)...)
		case *types.GridNode:
			for _, r := range n.Rows {
				for _, c := range r {
					imps = append(imps, importNodes(c.Content.Nodes)...)
				}
			}
		}
	}
	return imps
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
