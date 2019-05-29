package codelab_renderer

// import (
// 	"html/template"
// 	// "fmt"
// 	"strings"
// )


// func NewCodelabElement(data interface{}) CodelabElement {
// 	// reflective search and assignment...
// }

// func Render(data interface{}) template.HTML {
// 	t := reflect.TypeOf(data).Elem().Name()
// 	return executeTemplate(t, data)
// }

// // func (element *Heading) Html() template.HTML {
// // 	return executeTemplate("HeadingHtml", element)
// // }

// // func (element *StylizedText) Html() template.HTML {
// // 	return executeTemplate("StylizedTextHtml", element)
// // }
// // Missing: Image, YT
// // Missing: InlineCode, CodeBlock, 

// // Missing: TextBlock, InfoBlock
// // Missing: List, Survey, FAQ, CheckList
// // Missing: BlockContent, Step, Tutorial
// // func (element *Link) Html() template.HTML {
// // 	return executeTemplate("LinkHtml", element)
// // }

// // func (element *InlineCode) Html() template.HTML {
// // 	return executeTemplate("InlineCodeHtml", element)
// // }

// // func (element *TextBlock) Html() template.HTML {
// // 	// q := []ProtoRenderer{element.Content}
// // 	d := &CompositeData{
// // 		Children: ParseRepeated(element.Content, "Html"),
// // 		Data:     element,
// // 	}
// // 	return executeTemplate("TextBlockHtml", d)
// // }

// // func (element *InfoBox) Html() template.HTML {
// // 	return executeTemplate("InfoBoxHtml", element)
// // }

// // use same approach here as protos use for oneof stuff ~
// // ... or just have 3 of these! :(
// // just do it on the template...
// func ParseRepeated(children []*InlineContent, ext string) template.HTML {
// 	ext = strings.ToLower(ext)
// 	rendered_content := make([]string, len(children))
// 	for i, c := range children {
// 		if ext == "html" {
// 			rendered_content[i] = string(c.Html())
// 		// } else if ext == "md" {
// 		// 	rendered_content[i] = string(c.Md())
// 		} else {
// 			// log fatal ~
// 			rendered_content[i] = "Nope!"
// 		}
// 	}
// 	return template.HTML(strings.Join(rendered_content, "&#10;"))
// }