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
	"text/template"

	"github.com/googlecodelabs/tools/third_party"
)

// RenderOneof returns the underyling, rendered, passed oneof element
func RenderOneof(el interface{}, t *template.Template) string {
	// Recursive redering happens here
	return ExecuteTemplate(typeAssertUnderlingOneofType(el), t)
}

// typeAssertUnderlingOneofType turns a generic oneof proto into its underlying typed-proto
func typeAssertUnderlingOneofType(el interface{}) (underlyingType interface{}) {
	// Pure Oneof protos
	switch el.(type) {
	case *tutorial.InlineContent:
		underlyingType = underlyingInlineContentType(el.(*tutorial.InlineContent))
	}

	// debug-friendly panic
	if underlyingType == nil {
		panic(TypeNotSupported("typeAssertUnderlingOneofType", el))
	}

	return underlyingType
}

// underlyingInlineContentType asserts the underlying type of tutorial.InlineContent
func underlyingInlineContentType(el *tutorial.InlineContent) (underlyingType interface{}) {
	switch x := el.Content.(type) {
	case *tutorial.InlineContent_Text:
		// StylizedText
		underlyingType = x.Text
	case *tutorial.InlineContent_Code:
		// InlineCode
		underlyingType = x.Code
	case *tutorial.InlineContent_Link:
		// Link
		underlyingType = x.Link
	}

	// debug-friendly panic
	if underlyingType == nil {
		panic(TypeNotSupported("InnerInline", el))
	}

	return underlyingType
}
