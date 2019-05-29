package devrel_tutorial

import (
	"bytes"
	"log"
	"text/template"
)

const (
  mdTmplsDir = "../claat/new-render/templates/md/*"
  htmlTmplsDir = "../claat/new-render/templates/html/*"
)

var (
	md   *template.Template
	html *template.Template
)

// Base Templating Logic
func executeTemplate(d interface{}, tName string, t *template.Template) string {
	var w bytes.Buffer
	log.Printf("injecting d into tmpl: %s\n", tName)
	t.ExecuteTemplate(&w, tName, d)
	return w.String()
}

// Base Rendering Wrapper Types
type (
	ProtoRenderer interface {
		Md() string
		Html() string
	}

	oneof interface {
		GetInnerContent() ProtoRenderer
	}

	// Helper type template rendering
	compositeData struct {
		Data interface{} // Any message-specific field
		Text string      // Rendered Md or HTML
	}
)

// Helpers (Renderer, Parses, Testing)
func newCompositeData(d interface{}, txt string) *compositeData {
	return &compositeData{
		Data: d,
		Text: txt,
	}
}
