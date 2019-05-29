package render_base

import (
	"bytes"
	"reflect"
	"text/template"

	"fmt"
)

// Base Templating Logic
func ExecuteTemplate(
	data interface{}, tmpl_name string, t *template.Template) string {
	var w bytes.Buffer
	fmt.Printf("injecting data into tmpl: %s\n", tmpl_name)
	t.ExecuteTemplate(&w, tmpl_name, data)
	return w.String()
}

// Base Rendering Logic
type CodelabElement struct {
	Children []*interface{}
	Content  Oneof
	Text     string
}

type CompositeData struct {
	Data  *CodelabElement
	Text  string
}

func NewCompositeData(d *CodelabElement, txt string) *CompositeData {
	return &CompositeData{
		Data: d,
		Text: txt,
	}
}

type Oneof interface {
	Render()
}

// func (_ *CodelabElement) FieldByName(s string) ([]*reflect.StructField, bool) {
// 	q := []*reflect.StructField{}
// 	return q, false
// }

func IsOneof(el reflect.Type) bool {
	_, hasOneof := el.FieldByName("Content")
	return hasOneof
}

func IsComposite(el reflect.Type) bool {
	_, hasChildre := el.FieldByName("Children")
	return hasChildre
}

// higher util pkg
type TestingBatch struct {
	i string
	o string
}