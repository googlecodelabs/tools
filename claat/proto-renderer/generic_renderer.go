package genrenderer

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/googlecodelabs/tools/third_party"
)

// Maps protos to their type string name
func templateName(el interface{}) string {
	switch el.(type) {
	case *tutorial.StylizedText:
		return "StylizedText"
	}
	// This will cause a debug-friendly panic
	// TODO: Modularize into 'error_utils'
	return fmt.Sprintf("type not supported: %T %#v", el, el)
}

// ExecuteTemplate returns the evaluated template per passed templating
// namespace, based on the passed Codelab proto element type
func ExecuteTemplate(el interface{}, t *template.Template) string {
	var w bytes.Buffer
	e := t.ExecuteTemplate(&w, templateName(el), el)
	if e != nil {
		// This method outputs directly to templates. Panicking to surfance errors
		// since we cannot handle multiple returns in templates.
		// Panics will be gracefully handled by a wrapper function
		panic(fmt.Sprintf("Templating panic: %s\n", e))
	}
	return w.String()
}
