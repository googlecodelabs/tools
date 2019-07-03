package genrenderer

import (
  "text/template"

  "github.com/googlecodelabs/tools/third_party"
)

// RenderOneof returns the underyling, rendered, passed oneof element
func RenderOneof(el interface{}, t *template.Template) string {
  // Recursive redering happens here
  return ExecuteTemplate(AssertOneof(el), t)
}

// AssertOneof turns a generic oneof proto into its underlying typed-proto
func AssertOneof(el interface{}) (underlyingType interface{}) {
  // Pure Oneof protos
  switch el.(type) {
  case *tutorial.InlineContent:
    underlyingType = UnderlyingInlineContentType(el.(*tutorial.InlineContent))
  }

  // debug-friendly panic
  if underlyingType == nil {
    panic(TypeNotSupported("AssertOneof", el))
  }

  return underlyingType
}

// UnderlyingInlineContentType asserts the underlying type of tutorial.InlineContent
func UnderlyingInlineContentType(el *tutorial.InlineContent) (underlyingType interface{}) {
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
