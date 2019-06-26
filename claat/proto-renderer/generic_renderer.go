package genrenderer

import (
	"bytes"
	"fmt"
	"text/template"
)

// TODO: update to be proto dependent on next PR
func templateName(el interface{}) string {
	switch el.(type) {
	// This PR only comment: returns 'sampleProtoTemplate' to not break
	// 'ExecuteTemplate' tests
	case *SampleProtoTemplate:
		return "sampleProtoTemplate"
	}
	return fmt.Sprintf("type not supported: %#v", el)
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
