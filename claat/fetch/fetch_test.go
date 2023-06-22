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
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"testing/quick"

	_ "github.com/googlecodelabs/tools/claat/parser/gdoc" // Explicitly register gdoc parser
)

type testTransport struct {
	roundTripper func(*http.Request) (*http.Response, error)
}

func (tt *testTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return tt.roundTripper(r)
}

// TODO: add tests of core functionality.
// Must be able to fake out disk access, which probably involves a refactor focused on dependency injection?

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

func TestImgExtFromBytes(t *testing.T) {
	tests := []struct {
		bytes []byte

		wantExt string
		wantErr bool
	}{
		{[]byte("012345JFIF0"), ".jpeg", false},
		{[]byte("GIF34567890"), ".gif", false},
		{[]byte("SOMETHINGELSE"), ".png", false},
		{[]byte("GIF345JFIF0"), ".jpeg", false},
		{[]byte("toosmall"), "", true},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("bytes: %s", tc.bytes), func(t *testing.T) {
			ext, err := imgExtFromBytes(tc.bytes)

			if err != nil != tc.wantErr {
				t.Errorf("imgExtFromBytes() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if ext != tc.wantExt {
				t.Errorf("imgExtFromBytes() return: got %s, wanted %s", ext, tc.wantExt)
			}
		})
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
