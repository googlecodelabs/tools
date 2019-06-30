package genrenderer

import (
	"fmt"
	"text/template"

	"github.com/googlecodelabs/tools/third_party"
)

// RenderRepeated returns iterator-friend string outputs of the passed
// repeated proto by recursively rendering their contents.
func RenderRepeated(elSlice interface{}, t *template.Template) []string {
	contents := AssertRepeated(elSlice)
	sz := len(contents)
	renderedEls := make([]string, sz)

	if sz < 1 {
		// debug-friendly panic
		panic(fmt.Sprintf("RenderRepeated empty repeated field %#v", contents))
	}

	for i := 0; i < sz; i++ {
		// Recursive rendering happens here
		renderedEls[i] = ExecuteTemplate(contents[i], t)
	}

	return renderedEls
}

// AssertRepeated turns a generic proto slice into typed-slice that can be
// interated over without reflection. Panics if the passed type is not
// explicitly defined
func AssertRepeated(el interface{}) (guaranteedProtoSlice []interface{}) {
	// Below we convert turn all protos used as repeated fields
	// from interface{} into []interface{}.
	// Generalizable convertion approach that doesn't rely on reflection not found
	switch el.(type) {
	case []*tutorial.InlineContent:
		tempSlice := el.([]*tutorial.InlineContent)
		sz := len(tempSlice)

		guaranteedProtoSlice = make([]interface{}, sz)
		for i := 0; i < sz; i++ {
			guaranteedProtoSlice[i] = tempSlice[i]
		}
	}

	// debug-friendly panic
	if guaranteedProtoSlice == nil {
		panic(TypeNotSupported("AssertRepeated", el))
	}

	return guaranteedProtoSlice
}
