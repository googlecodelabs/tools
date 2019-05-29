package devrel_tutorial

import (
	"strings"
	"text/template"
)

func init() {
	html = template.Must(template.New("master").ParseGlob(htmlTmplsDir))
}
