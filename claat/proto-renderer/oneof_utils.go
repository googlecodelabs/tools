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
  case *tutorial.BlockContent:
    underlyingType = InnerContentsBlock(el.(*tutorial.BlockContent))
  }

  // debug-friendly panic
  if guaranteedProtoSlice == nil {
    panic(TypeNotSupported("AssertAndExtractOneof", el))
  }
}

// InnerContentsInline asserts the underlying type of tutorial.InlineContent
func InnerContentsInline(el *tutorial.InlineContent) (underlyingType interface{}) {
  switch x := el.Content.(type) {
  case *tutorial.InlineContent_Text:
    // StylizedText
    return x.Text
  case *tutorial.InlineContent_Code:
    // InlineCode
    return x.Code
  }

  // debug-friendly panic
  if guaranteedProtoSlice == nil {
    panic(TypeNotSupported("InnerContentsInline", el))
  }
}
