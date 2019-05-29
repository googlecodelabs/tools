package devrel_tutorial

import (
	"bytes"
	"text/template"
	"github.com/googlecodelabs/tools/claat/util"
)

// Reviewing comment: Templates can live anywhere within the repo
const (
	mdTmplsDir = "github.com/googlecodelabs/tools/third_party/templates/md/*"
	htmlTmplsDir = "github.com/googlecodelabs/tools/third_party/templates/html/*"
)

var (
	md   *template.Template
	html *template.Template
)

type (
	ProtoRenderer interface {
		Md() string
		Html() string
	}
)

// Base Templating Logic
func executeTemplate(d interface{}, tName string, t *template.Template) string {
	var w bytes.Buffer
	util.LogIfError(t.ExecuteTemplate(&w, tName, d))
	return w.String()
}