package genrenderer

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/googlecodelabs/tools/third_party"
)

// templateName Maps protos to their type string name
func templateName(el interface{}) string {
	switch el.(type) {
	case *tutorial.StylizedText, tutorial.StylizedText:
		return "StylizedText"
	case *tutorial.InlineCode, tutorial.InlineCode:
		return "InlineCode"
	case *tutorial.InlineContent, tutorial.InlineContent:
		return "InlineContent"
	}
	// This will cause a debug-friendly panic
	return TypeNotSupported("genrenderer.templateName", el)
}

// ExecuteTemplate returns the evaluated template per passed templating
// namespace, based on the passed tutorial proto type string name
func ExecuteTemplate(el interface{}, t *template.Template) string {
	var w bytes.Buffer
	e := t.ExecuteTemplate(&w, templateName(el), el)
	if e != nil {
		// This method outputs directly to templates. Panicking to surfance errors
		// since we should not handle multiple returns in templates.
		// Errors will be more gracefully handled in output-format packages
		panic(fmt.Sprintf("Templating panic: %s\n", e))
	}
	return w.String()
}
