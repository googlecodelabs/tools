package devrel_tutorial

import (
	"strings"
	"text/template"
)

func repeatedHeading(level int32) string {
	return strings.Repeat("#", int(level))
}

func init() {
	funcMap := template.FuncMap{
		"repeatedHeading": repeatedHeading,
	}
	md = template.Must(template.New("master").ParseGlob(mdTmplsDir))
}

func (el *Heading) Md() string {
	return executeTemplate(&el, "Heading", md)
}
