package devrel_tutorial

import (
	"strings"
	"text/template"
)

func init() {
	md = template.Must(template.New("master").ParseGlob(mdTmplsDir))
}
