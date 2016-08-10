package bundle

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"

	"golang.org/x/net/context"

	_ "github.com/googlecodelabs/tools/claat/parser/gdoc"
	"github.com/googlecodelabs/tools/claat/types"
)

type testTransport struct {
	roundTripper func(*http.Request) (*http.Response, error)
}

func (tt *testTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return tt.roundTripper(r)
}

type testContentWriter struct {
	sync.Mutex
	assets map[string]map[string][]byte
	markup map[string][]byte
	meta   map[string]*types.ContextMeta
}

func (cw *testContentWriter) WriteAsset(ctx context.Context, clab, name string, body io.Reader) error {
	cw.Lock()
	defer cw.Unlock()
	if cw.assets == nil {
		cw.assets = make(map[string]map[string][]byte)
	}
	if cw.assets[clab] == nil {
		cw.assets[clab] = make(map[string][]byte)
	}
	var err error
	cw.assets[clab][name], err = ioutil.ReadAll(body)
	return err
}

func (cw *testContentWriter) WriteMarkup(ctx context.Context, clab string, body []byte) error {
	cw.Lock()
	defer cw.Unlock()
	if cw.markup == nil {
		cw.markup = make(map[string][]byte)
	}
	cw.markup[clab] = body
	return nil
}

func (cw *testContentWriter) WriteMeta(ctx context.Context, cmeta *types.ContextMeta) error {
	cw.Lock()
	defer cw.Unlock()
	if cw.meta == nil {
		cw.meta = make(map[string]*types.ContextMeta)
	}
	cw.meta[cmeta.ID] = cmeta
	return nil
}

func TestExporterExport(t *testing.T) {
	const (
		clabID    = "test-lab"
		clabTitle = "Test Codelab"
		markup    = `
		<html><head><style>
			.bold { font-weight: bold }
		</style></head>
		<body>
			<p class="title">Test Codelab</p>
			<table><tr><td>id</td><td>test-lab</td></tr></table>
			<h1>Overview</h1>
			<p>[[<span class="bold">import</span>
				<a href="https://docs.google.com/document/d/import1">import1</a>]]</p>
			<p>[[<span class="bold">import</span>
				<a href="https://docs.google.com/document/d/import2">import2</a>]]</p>
			<img src="https://example.com/image1.png">
			<img src="https://example.com/image2.png">
		</body>
		</html>
		`
		docJSON = `{
			"id": "doc-123",
			"mimeType": "application/vnd.google-apps.document",
			"exportLinks": {"text/html": "https://docs.google.com/export"},
			"modifiedDate": "2016-05-09T14:03:10.399Z"
		}
	`
	)
	cl := &http.Client{Transport: &testTransport{func(r *http.Request) (*http.Response, error) {
		var body io.Reader
		switch {
		case strings.HasSuffix(r.URL.Path, "/files/doc-123"):
			body = strings.NewReader(docJSON)
		case r.URL.Path == "/export":
			body = strings.NewReader(markup)
		case r.FormValue("id") == "import1":
			body = strings.NewReader("I am import 1")
		case r.FormValue("id") == "import2":
			body = strings.NewReader("I am import 2")
		default:
			body = strings.NewReader(r.URL.String())
		}
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(body),
		}, nil
	}}}

	// cnt will hold all exported assets and markup
	cnt := &testContentWriter{}
	e := &Exporter{
		HTTPClient: func(context.Context, string) *http.Client { return cl },
	}
	err := e.Export(context.Background(), cnt, "doc-123", types.FmtHTML)
	if err != nil {
		t.Fatalf("e.Export: %v", err)
	}

	// html is the codelab in HTML
	html := string(cnt.markup[clabID])
	want := []string{
		"Test Codelab",
		"I am import 1",
		"I am import 2",
	}
	var missing bool
	// verify title and imported fragments
	for _, v := range want {
		if !strings.Contains(html, v) {
			missing = true
			t.Errorf("%q is missing", v)
		}
	}
	// verify assets (images)
	if len(cnt.assets[clabID]) != 2 {
		t.Errorf("len(cnt.assets) = %d; want 2", len(cnt.assets[clabID]))
	}
	for k, v := range cnt.assets[clabID] {
		want := "https://example.com/image"
		if !strings.HasPrefix(string(v), want) {
			t.Errorf("cnt.assets[%q][%q] = %q; want to start with %q", clabID, k, v, want)
		}
		img := fmt.Sprintf(`<img src="img/%s"`, k)
		if !strings.Contains(html, img) {
			missing = true
			t.Errorf("%s is missing", img)
		}
	}
	if missing {
		t.Errorf("html: %s", html)
	}

	// verify metadata
	meta := cnt.meta[clabID]
	if meta == nil {
		t.Fatalf("cnt.meta[%q] = nil", clabID)
	}
	if meta.Title != clabTitle {
		t.Errorf("meta.Title = %q; want %q", meta.Title, clabTitle)
	}
}

func TestGdocID(t *testing.T) {
	tests := []struct{ in, out string }{
		{"https://docs.google.com/document/d/foo", "foo"},
		{"https://docs.google.com/document/d/foo/edit", "foo"},
		{"https://docs.google.com/document/d/foo/edit#abc", "foo"},
		{"https://docs.google.com/document/d/foo/edit?bar=baz#abc", "foo"},
		{"foo", "foo"},
	}
	for i, test := range tests {
		out := gdocID(test.in)
		if out != test.out {
			t.Errorf("%d: gdocID(%q) = %q; want %q", i, test.in, out, test.out)
		}
	}
}
