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
	md = template.Must(template.New("master").Funcs(funcMap).ParseGlob(mdTmplsDir))
}

// TODO if possible: Template names follow its calling struct type,
//                   use `reflect.TypeOf(el).Elem().Name()`
//                   as a generic caller, once a catch-all type is figured out
func (el *Heading) Md() string {
	return executeTemplate(&el, "Heading", md)
}

func (el *StylizedText) Md() string {
	return executeTemplate(&el, "StylizedText", md)
}