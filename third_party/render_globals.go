package devrel_tutorial

import (
	"bytes"
	"log"
	"text/template"
)

// Reviewing comment: Templates can live anywhere within the repo
const (
  mdTmplsDir = 	"github.com/googlecodelabs/tools/third_party/templates/md/*"
  htmlTmplsDir = 	"github.com/googlecodelabs/tools/third_party/templates/html/*"
)

var (
	md   *template.Template
	html *template.Template
)