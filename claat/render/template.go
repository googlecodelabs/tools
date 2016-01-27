// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package render

import (
	htmlTemplate "html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"sort"
	textTemplate "text/template"

	"github.com/googlecodelabs/tools/claat/types"
)

// Context is a template context during execution.
type Context struct {
	Env       string
	Prefix    string
	GlobalGA  string
	Meta      *types.Meta
	Steps     []*types.Step
	ExtraVars map[string]string // Extra variables passed from the command line.
}

// Execute renders a template of the fmt format into w.
// The fmt argument can also be a path to a local file.
func Execute(w io.Writer, fmt string, ctx *Context) error {
	t, err := parseTemplate(fmt)
	if err != nil {
		return err
	}
	sort.Strings(ctx.Meta.Tags)
	return t.Execute(w, ctx)
}

// executer satisfies both html/template and text/template.
type executer interface {
	Execute(io.Writer, interface{}) error
}

// funcMap are exposted to the templates.
var funcMap = map[string]interface{}{
	"renderHTML": HTML,
	"matchEnv": func(tags []string, t string) bool {
		if len(tags) == 0 || t == "" {
			return true
		}
		i := sort.SearchStrings(tags, t)
		return i < len(tags) && tags[i] == t
	},
}

//go:generate go run gen-tmpldata.go

type template struct {
	bytes []byte
	html  bool
}

// parseTemplate parses template name defined either in tmpldata
// or a local file.
//
// A local file template is parsed as HTML if file extension is ".html",
// text otherwise.
func parseTemplate(name string) (executer, error) {
	tmpl := tmpldata[name] // defined in pre-generated tmpldata.go
	if tmpl == nil {
		// TODO: add templates in-mem caching
		var err error
		if tmpl, err = readTemplate(name); err != nil {
			return nil, err
		}
	}
	if tmpl.html {
		return htmlTemplate.New(name).
			Funcs(funcMap).
			Parse(string(tmpl.bytes))
	}
	return textTemplate.New(name).
		Funcs(funcMap).
		Parse(string(tmpl.bytes))
}

func readTemplate(name string) (*template, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return &template{
		bytes: b,
		html:  filepath.Ext(name) == ".html",
	}, nil
}
