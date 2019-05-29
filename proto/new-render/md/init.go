package codelab_renderer

import (
	"text/template"
	"strings"
)

var (
	t  *template.Template
)

func repeatedHeading(level int32) string {
	return strings.Repeat("#", int(level))
}

func init() {
	funcMap := template.FuncMap{
		"repeatedHeading": repeatedHeading,
	}
	tmplsDir := "../claat/new-render/md/templates/*"
	t = template.Must(template.New("master").Funcs(funcMap).ParseGlob(tmplsDir))
}