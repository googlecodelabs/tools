package genrenderer

import (
	"text/template"

	"github.com/googlecodelabs/tools/third_party"
)

// RenderRepeated returns iterator-friend string outputs of the passed
// repeated proto by recursively rendering their contents.
func RenderRepeated(els interface{}, t *template.Template) []string {
	contents := AssertRepeated(els)
	sz := len(contents)
	rendered := make([]string, sz)

	for i := 0; i < sz; i++ {
		// Recursive rendering happens here
		rendered[i] = ExecuteTemplate(contents[i], t)
	}

	return rendered
}

// AssertRepeated turns a generic proto slice into typed-slice that can be
// interated over without reflection. Panics if the passed type is not
// explicitly defined
func AssertRepeated(el interface{}) (guaranteedProtoSlice []interface{}) {
	// These are the protos used as repeated fields
	switch el.(type) {
	case []*tutorial.InlineContent:
		guaranteedProtoSlice = interfaceSlice(el.([]*tutorial.InlineContent))
	}

	// debug-friendly panic
	if guaranteedProtoSlice != nil {
		panic(TypeNotSupported("AssertRepeated", el))
	}
}

// interfaceSlice turns an interface{} into []interface{}.
// Static type-assertions of this kind are illegal.
// Also reduces type-assertion boilerplate for 'AssertRepeated'
func interfaceSlice(elSlice ...interface{}) []interface{} {
	sz := len(elSlice)
	newSlice := make([]interface{}, sz)

	for i := 0; i < sz; i++ {
		newSlice[i] = elSlice[i]
	}
	return newSlice
}
