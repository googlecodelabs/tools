// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package html

import (
	"go/build"
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/googlecodelabs/tools/claat/proto/renderer"
	"github.com/googlecodelabs/tools/third_party"
)

const tmplsRltvDir = "src/github.com/googlecodelabs/tools/claat/proto/renderer/html/templates/*"

var (
	tmplNmspc   *template.Template
	tmplsAbsDir = filepath.Join(build.Default.GOPATH, tmplsRltvDir)
	funcMap     = template.FuncMap{
		"asString":            asString,
		"renderOneof":         renderOneof,
		"renderRepeated":      renderRepeated,
		"listVarietyToTag":    listVarietyToTag,
		"listFormattingClass": listFormattingClass,
	}
)

func init() {
	// Defining namespace after initial compilation to avoid initialization loop
	tmplNmspc = template.Must(template.New("html").Funcs(funcMap).ParseGlob(tmplsAbsDir))
}

// Render returns the rendered HTML representation of a tutorial proto,
// or the first error encountered rendering templates depth-first, if any.
// Possible recursive descent follows proto definitions
func Render(el interface{}) (out io.Reader, err error) {
	// "Catches" first nested panic and delegates handling to caller
	defer func() {
		r := recover()
		if r != nil {
			out = nil
			err = genrenderer.AssertError(r)
		}
	}()

	out = strings.NewReader(genrenderer.ExecuteTemplate(el, tmplNmspc))
	return out, err
}

// renderOneof is a self-referential template function used
// in all templates of protos with oneof fields
func renderOneof(contents interface{}) string {
	return genrenderer.RenderOneof(contents, tmplNmspc)
}

// renderRepeated is a self-referential template function used
// in all templates of protos with repeated fields
func renderRepeated(contents interface{}) []string {
	return genrenderer.RenderRepeated(contents, tmplNmspc)
}

// asString concatenates the output of renderRepeated into a string.
// Used in protos with non-custom repeated field rendering
// Similar to Nested templates/proto usage.
// Avoids the use of an extra, single-use grouping proto, yet shortens
// and simplifies templates.
func asString(contents []string) string {
	return strings.Join(contents, "")
}

// listVarietyToTag maps 'ListVariety' enums to their HTML tags
func listVarietyToTag(v tutorial.List_ListVariety) string {
	switch v.String() {
	case "UNORDERED":
		return "ul"
	case "ORDERED":
		return "ol"
	default:
		return "unknown-list-variety"
	}
}

// listFormattingClass maps 'ListStyle' enums to their CSS classes
func listFormattingClass(s tutorial.List_ListStyle) string {
	v := s.String()

	if strings.HasPrefix(v, "UNKNOWN") {
		return ""
	}
	return strings.ToLower(v)
}
