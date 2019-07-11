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
