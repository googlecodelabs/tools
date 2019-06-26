package genrenderer

import (
  "text/template"

  "github.com/googlecodelabs/tools/third_party"
)

// RenderOneof Returns the underyling, rendered, passed oneof element
func RenderOneof(el interface{}, t *template.Template) string {
  // Recursive redering happens here
  return ExecuteTemplate(AssertAndExtractOneof(el), t)
}

// AssertAndExtractOneof ...
func AssertAndExtractOneof(el interface{}) (underlyingType interface{}) {
  // Pure Oneof protos
  switch el.(type) {
  case *tutorial.InlineContent:
    underlyingType = InnerContentsInline(el.(*tutorial.InlineContent))
  }

  // debug-friendly panic
  if underlyingType == nil {
    panic(TypeNotSupported("AssertAndExtractOneof", el))
  }

  return underlyingType
}

// InnerContentsInline asserts the underlying type of tutorial.InlineContent
func InnerContentsInline(el *tutorial.InlineContent) (underlyingType interface{}) {
  switch x := el.Content.(type) {
  case *tutorial.InlineContent_Text:
    // StylizedText
    underlyingType = x.Text
  case *tutorial.InlineContent_Code:
    // InlineCode
    underlyingType = x.Code
  }

  // debug-friendly panic
  if underlyingType == nil {
    panic(TypeNotSupported("InnerContentsInline", el))
  }

  return underlyingType
}
