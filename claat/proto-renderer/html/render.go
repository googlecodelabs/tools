package html

import (
	"go/build"
	"path/filepath"
	"text/template"

	"github.com/googlecodelabs/tools/proto-renderer"
)

const (
	tmplsRltvDir = "src/github.com/googlecodelabs/tools/proto-renderer/templates/*"
)

var (
	t *template.Template
)

func init() {
	tmplsAbsDir := filepath.Join(build.Default.GOPATH, tmplsRltvDir)
	t = template.Must(template.New("html").ParseGlob(tmplsAbsDir))
}

// Render returns the rendered HTML representation of a devrel_tutorial proto
// and the first error encountered rendering templates depth-first, if any
func Render(el interface{}) (string, error) {
	// "Catch" panics if they occur
	defer func() {
		err := recover()
		if err != nil {
			return "", err
		}
	}()

	return genrenderer.ExecuteTemplate(el, t), nil
}
