package devrel_tutorial

import (
	"strings"
	"text/template"
)

func init() {
	funcMap := template.FuncMap{
		"safeHTML": func(s string) htmlTemplate.HTML {
			return template.HTML(s)
		},
	}
	html = template.Must(template.New("master").Funcs(funcMap).ParseGlob(htmlTmplsDir))
}
