package codelab_renderer

// func renderRepeatedMd(el *render_base.CodelabElement, ext_del string) {
// 	// escape this somehow later ~
// 	render_base.RenderRepeated(el, "\n")
// }

// import (
// 	"text/template"
// 	// "fmt"
// 	// "strings"
// )

// // Leaf Types
// func (element *Heading) Md() string {
// 	return executeTemplate("HeadingMd", element)
// }

// func (element *StylizedText) Md() string {
// 	return executeTemplate("StylizedTextMd", element)
// }
// // Missing: Image, YT
// // Missing: InlineCode, CodeBlock, 

// // Composed types
// func (element *Link) Md() string {
// 	return executeTemplate("LinkMd", element)
// }

// func (element *InlineCode) Md() string {
// 	return executeTemplate("InlineCodeMd", element)
// }

// func (element *InfoBox) Md() string {
// 	return executeTemplate("InfoBoxMd", element)
// }
// // Missing: TextBlock, InfoBlock
// // Missing: List, Survey, FAQ, CheckList
// // Missing: BlockContent, Step, Tutorial