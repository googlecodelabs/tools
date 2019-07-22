// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package genrenderer

import (
	"fmt"
	"text/template"

	"github.com/googlecodelabs/tools/third_party"
)

// RenderRepeated returns iterator-friend string outputs of the passed
// repeated proto by recursively rendering their contents.
func RenderRepeated(elSlice interface{}, t *template.Template) []string {
	contents := typeAssertInterfaceSlice(elSlice)
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

// typeAssertInterfaceSlice turns a generic proto slice into typed-slice that can be
// interated over without reflection. Panics if the passed type is not
// explicitly defined
func typeAssertInterfaceSlice(el interface{}) (protoSlice []interface{}) {
	// Below we convert turn all protos used as repeated fields
	// from interface{} into []interface{}.
	// Generalizable convertion approach that doesn't rely on reflection not found
	switch el.(type) {
	case []*tutorial.InlineContent:
		protoSlice = interfaceSliceInlineContent(el)
	case []*tutorial.StylizedText:
		protoSlice = interfaceSliceStylizedText(el)
	case []*tutorial.Paragraph:
		protoSlice = interfaceSliceParagraph(el)
	}

	// debug-friendly panic
	if protoSlice == nil {
		panic(TypeNotSupported("typeAssertInterfaceSlice", el))
	}

	return protoSlice
}

// []*tutorial.repeatedFields to []interface{} conversion helpers.
// Lack of reflection use ovehead. Only the first line of each is different.

func interfaceSliceInlineContent(elSliceInterface interface{}) []interface{} {
	elSlice := elSliceInterface.([]*tutorial.InlineContent)

	sz := len(elSlice)
	interfaceSlice := make([]interface{}, sz)

	for i := 0; i < sz; i++ {
		interfaceSlice[i] = elSlice[i]
	}
	return interfaceSlice
}

func interfaceSliceStylizedText(elSliceInterface interface{}) []interface{} {
	elSlice := elSliceInterface.([]*tutorial.StylizedText)

	sz := len(elSlice)
	interfaceSlice := make([]interface{}, sz)

	for i := 0; i < sz; i++ {
		interfaceSlice[i] = elSlice[i]
	}
	return interfaceSlice
}

func interfaceSliceParagraph(elSliceInterface interface{}) []interface{} {
	elSlice := elSliceInterface.([]*tutorial.Paragraph)

	sz := len(elSlice)
	interfaceSlice := make([]interface{}, sz)

	for i := 0; i < sz; i++ {
		interfaceSlice[i] = elSlice[i]
	}
	return interfaceSlice
}
