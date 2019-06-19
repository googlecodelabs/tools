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
		// This method outputs directly to templates. Panicking to surfance errors
		// since we cannot handle multiple returns in templates.
		// Panics will be gracefully handled by a wrapper function
		panic(fmt.Sprintf("Templating panic: %s\n", e))
	}
	return w.String()
}
