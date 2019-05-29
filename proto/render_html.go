package devrel_tutorial

import (
	"strings"
	htmlTemplate "html/template"
	"text/template"
)

func init() {
	// TODO:
	//   _html_template_htmlescaper/safeHTML template f(x) for StylizedText
	//   https://github.com/golang/go/issues/21844#issuecomment-329128520
	funcMap := template.FuncMap{
		"safeHTML": func(s string) string {
			return htmlTemplate.HTMLEscapeString(s)
		},
	}

	html = template.Must(template.New("master").ParseGlob(htmlTmplsDir))
}

// TODO if possible:
//     Template names are named after their struct type,
//     use `reflect.TypeOf(el).Elem().Name()`
//     as a generic caller once a catch-all type, or
//     field-persistent inferace-passing process is figured out

// Leaf Types
func (el *StylizedText) Html() string {
	return executeTemplate(&el, "StylizedText", html)
}

func (el *InlineCode) Html() string {
	return executeTemplate(&el, "InlineCode", html)
}

// Oneof Types
func (el *InlineContent) Html() string {
	// Delegator to subtype render
	return el.GetInnerContent().Html()
}

func (el *BlockContent) Html() string {
	underlyingEl := el.GetInnerContent()
	data := newCompositeData(nil, underlyingEl.Html())
	return executeTemplate(&data, "BlockContent", html)
}

// Repeated Types
func (el *TextBlock) Html() string {
	desambiguatedContents := make([]string, len(el.Children))
	for i, c := range el.Children {
		// TODO: Concurrent SyncGroup with Call.(MethodName)
		// 			 Perfomance change: O(h) vs O(n)
		underlyingEl := c.GetInnerContent()
		desambiguatedContents[i] = underlyingEl.Html()
	}
	mergedContent := strings.Join(desambiguatedContents, "")

	data := newCompositeData(nil, mergedContent)
	return executeTemplate(&data, "TextBlock", html)
}
