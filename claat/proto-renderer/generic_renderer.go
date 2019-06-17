package genrenderer

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"
)

// ExecuteTemplate returns the evaluated template per passed templating
// namespace, based on the passed Codelab proto element type
func ExecuteTemplate(el interface{}, t *template.Template) string {
	var w bytes.Buffer
	e := t.ExecuteTemplate(&w, reflect.TypeOf(el).Name(), el)
	if e != nil {
		panic(fmt.Sprintf("Templating err: %s\n", e))
	}
	return w.String()
}
