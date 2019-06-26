package html

import (
	"go/build"
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/googlecodelabs/tools/claat/proto-renderer"
)

const tmplsRltvDir = "src/github.com/googlecodelabs/tools/claat/proto-renderer/html/templates/*"

var (
	tmplsAbsDir = filepath.Join(build.Default.GOPATH, tmplsRltvDir)
	tmplNmspc   = template.Must(template.New("html-pkg").ParseGlob(tmplsAbsDir))
)

// Render returns the rendered HTML representation of a tutorial proto,
// or the first error encountered rendering templates depth-first, if any
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
