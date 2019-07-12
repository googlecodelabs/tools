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
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"testing/quick"

	"github.com/googlecodelabs/tools/claat/render"
	"github.com/googlecodelabs/tools/claat/types"
)

type testTransport struct {
	roundTripper func(*http.Request) (*http.Response, error)
}

func (tt *testTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return tt.roundTripper(r)
}

func TestFetchRemote(t *testing.T) {
	const f = "/file.txt"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("r.Method = %q; want GET", r.Method)
		}
		if r.URL.Path != f {
			t.Errorf("r.URL.Path = %q; want %q", r.URL.Path, f)
		}
		w.Write([]byte("test"))
	}))
	defer ts.Close()

	res, err := fetchRemote(ts.URL+f, "", false)
	if err != nil {
		t.Fatal(err)
	}
	defer res.body.Close()
	if res.typ != srcMarkdown {
		t.Errorf("typ = %q; want %q", res.typ, srcMarkdown)
	}
	b, _ := ioutil.ReadAll(res.body)
	if s := string(b); s != "test" {
		t.Errorf("res = %q; want 'test'", s)
	}
}

func TestFetchRemoteDrive(t *testing.T) {
	const driveHost = "http://dummy"
	rt := &testTransport{func(r *http.Request) (*http.Response, error) {
		if r.Method != "GET" {
			t.Errorf("r.Method = %q; want GET", r.Method)
		}
		// metadata request
		if strings.HasSuffix(r.URL.Path, "/files/doc-123") {
			if v := r.FormValue("supportsTeamDrives"); v != "true" {
				t.Errorf("supportsTeamDrives = %q; want 'true'", v)
			}
			b := ioutil.NopCloser(strings.NewReader(`{
				"mimeType": "application/vnd.google-apps.document"
			}`))
			return &http.Response{Body: b, StatusCode: http.StatusOK}, nil
		}
		// export request
		if !strings.HasSuffix(r.URL.Path, "/doc-123/export") {
			t.Errorf("r.URL.Path = %q; want /doc-123/export suffix", r.URL.Path)
		}
		b := ioutil.NopCloser(strings.NewReader("test"))
		return &http.Response{Body: b, StatusCode: http.StatusOK}, nil
	}}
	clients[providerGoogle] = &http.Client{Transport: rt}

	res, err := fetchRemote("doc-123", "", false)
	if err != nil {
		t.Fatal(err)
	}
	defer res.body.Close()
	if res.typ != srcGoogleDoc {
		t.Errorf("typ = %q; want %q", res.typ, srcGoogleDoc)
	}
	b, _ := ioutil.ReadAll(res.body)
	if s := string(b); s != "test" {
		t.Errorf("res = %q; want 'test'", s)
	}
}

func TestSlurpWithFragment(t *testing.T) {
	dochtml, err := ioutil.ReadFile("testdata/gdoc.html")
	if err != nil {
		t.Fatal(err)
	}
	rt := &testTransport{func(r *http.Request) (*http.Response, error) {
		// metadata request
		if r.URL.Path == "/drive/v3/files/doc-123" {
			b := ioutil.NopCloser(strings.NewReader(`{
				"mimeType": "application/vnd.google-apps.document"
			}`))
			return &http.Response{Body: b, StatusCode: http.StatusOK}, nil
		}
		// main doc export request
		if r.URL.Path == "/drive/v3/files/doc-123/export" {
			b := ioutil.NopCloser(bytes.NewReader(dochtml))
			return &http.Response{Body: b, StatusCode: http.StatusOK}, nil
		}
		// import doc request, referenced in testdata/gdoc.html
		if r.URL.Path == "/drive/v3/files/import/export" {
			b := ioutil.NopCloser(strings.NewReader(`
				<p>I'm imported from elsewhere.</p>
			`))
			return &http.Response{Body: b, StatusCode: http.StatusOK}, nil
		}
		return &http.Response{
			Body:       ioutil.NopCloser(strings.NewReader(r.URL.String())),
			StatusCode: http.StatusBadRequest,
		}, nil
	}}
	clients[providerGoogle] = &http.Client{Transport: rt}

	clab, err := slurpCodelab("doc-123", "", map[string]bool{})
	if err != nil {
		t.Fatal(err)
	}
	var node *types.ImportNode
	for _, st := range clab.Steps {
		for _, n := range st.Content.Nodes {
			if n.Type() == types.NodeImport {
				node = n.(*types.ImportNode)
				break
			}
		}
	}
	if node == nil {
		t.Fatal("no import node found")
	}
	html, err := render.HTML("", node.Content)
	if err != nil {
		t.Fatal(err)
	}
	want := "imported from elsewhere"
	if !strings.Contains(string(html), want) {
		t.Errorf("%s does not contain %q", html, want)
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

func TestRestrictPathToParent(t *testing.T) {
	tests := []struct {
		asset  string
		parent string

		wantPath string
		wantErr  bool
	}{
		{"imgroot.png", ".", "imgroot.png", false},
		{"imgroot.png", "foo/", "foo/imgroot.png", false},
		{"img/sub.png", "foo/", "foo/img/sub.png", false},
		{"imgroot.png", "/tmp/foo", "/tmp/foo/imgroot.png", false},
		{"/tmp/imgabs.png", "foo/", "", true},
		{"../imgup.png", "foo/", "", true},
		{"../imgup.png", "..", "", true},
		{"imgroot.png", "", "imgroot.png", false},
		{"", ".", ".", false},
		{"", "", ".", false},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("asset: %s, parent: %s", tc.asset, tc.parent), func(t *testing.T) {
			tc.wantPath = safeAbs(t, tc.wantPath)

			p, err := restrictPathToParent(tc.asset, tc.parent)

			if err != nil != tc.wantErr {
				t.Errorf("restrictPathToParent() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if p != tc.wantPath {
				t.Errorf("restrictPathToParent() return: got %s, wanted %s", p, tc.wantPath)
			}
		})
	}
}

func TestFuzzRestrictPathToParent(t *testing.T) {
	checkInParent := func(elem, parent string) bool {
		_, err := restrictPathToParent(elem, parent)

		parent = safeAbs(t, parent)
		if !strings.HasPrefix(elem, "/") {
			elem = filepath.Join(parent, elem)
		}
		shouldOk := strings.HasPrefix(elem, parent)
		return shouldOk == (err == nil)
	}

	if err := quick.Check(checkInParent, nil); err != nil {
		t.Error(err)
	}
}

// safeAbs compute Abs of p and fail the test if not valid.
// Empty string return empty path.
func safeAbs(t *testing.T, p string) string {
	if p == "" {
		return p
	}
	p, err := filepath.Abs(p)
	if err != nil {
		t.Fatalf("Error in converting %s to abs path: %v", p, err)
	}
	return p
}
