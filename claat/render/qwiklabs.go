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
	"bytes"
	"io"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/googlecodelabs/tools/claat/types"
	htmlTemplate "html/template"
)

// MD renders nodes as markdown for the target env.
func Qwiklabs(env string, nodes ...types.Node) (string, error) {
	var buf bytes.Buffer
	if err := WriteQwiklabs(&buf, env, nodes...); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// WriteMD does the same as MD but outputs rendered markup to w.
func WriteQwiklabs(w io.Writer, env string, nodes ...types.Node) error {
	qw := qwiklabsWriter{w: w, env: env}
	return qw.write(nodes...)
}

type qwiklabsWriter struct {
	w         io.Writer // output writer
	env       string    // target environment
	err       error     // error during any writeXxx methods
	lineStart bool
}


func (qw *qwiklabsWriter) writeBytes(b []byte) {
	if qw.err != nil {
		return
	}
	qw.lineStart = len(b) > 0 && b[len(b)-1] == '\n'
	_, qw.err = qw.w.Write(b)
}

func (qw *qwiklabsWriter) writeString(s string) {
	qw.writeBytes([]byte(s))
}

func (qw *qwiklabsWriter) writeEscape(s string) {
	htmlTemplate.HTMLEscape(qw.w, []byte(s))
}

func (qw *qwiklabsWriter) space() {
	if !qw.lineStart {
		qw.writeString(" ")
	}
}

func (qw *qwiklabsWriter) newBlock() {
	if !qw.lineStart {
		qw.writeBytes(newLine)
	}
	qw.writeBytes(newLine)
}

func (qw *qwiklabsWriter) matchEnv(v []string) bool {
	if len(v) == 0 || qw.env == "" {
		return true
	}
	i := sort.SearchStrings(v, qw.env)
	return i < len(v) && v[i] == qw.env
}

func (qw *qwiklabsWriter) write(nodes ...types.Node) error {
	for _, n := range nodes {
		if !qw.matchEnv(n.Env()) {
			continue
		}
		switch n := n.(type) {
		case *types.TextNode:
			qw.text(n)
		case *types.ImageNode:
			qw.image(n)
		case *types.URLNode:
			qw.url(n)
		case *types.ButtonNode:
			qw.write(n.Content.Nodes...)
		case *types.CodeNode:
			qw.code(n)
		case *types.ListNode:
			qw.list(n)
		case *types.ImportNode:
			if len(n.Content.Nodes) == 0 {
				break
			}
			qw.write(n.Content.Nodes...)
		case *types.ItemsListNode:
			qw.itemsList(n)
		//case *types.GridNode:
		//	qw.grid(n)
		case *types.InfoboxNode:
			qw.infobox(n)
		//case *types.SurveyNode:
		//	qw.survey(n)
		case *types.HeaderNode:
			qw.header(n)
		//case *types.YouTubeNode:
		//	qw.youtube(n)
		}
		if qw.err != nil {
			return qw.err
		}
	}
	return nil
}

func (qw *qwiklabsWriter) text(n *types.TextNode) {
	if n.Bold {
		qw.writeString("__")
	}
	if n.Italic {
		qw.writeString(" *")
	}
	if n.Code {
		qw.writeString("`")
	}
	qw.writeString(n.Value)
	if n.Code {
		qw.writeString("`")
	}
	if n.Italic {
		qw.writeString("* ")
	}
	if n.Bold {
		qw.writeString("__")
	}
}

func (qw *qwiklabsWriter) image(n *types.ImageNode) {
	qw.space()
	qw.writeString("![")
	qw.writeString(path.Base(n.Src))
	qw.writeString("](")
	qw.writeString(n.Src)
	qw.writeString(")")
}

func (qw *qwiklabsWriter) url(n *types.URLNode) {
	qw.space()
	if n.URL != "" {
		qw.writeString("[")
	}
	for _, cn := range n.Content.Nodes {
		if t, ok := cn.(*types.TextNode); ok {
			qw.writeString(t.Value)
		}
	}
	if n.URL != "" {
		qw.writeString("](")
		qw.writeString(n.URL)
		qw.writeString(")")
	}
}

func (qw *qwiklabsWriter) code(n *types.CodeNode) {
	qw.newBlock()
	defer qw.writeBytes(newLine)
	if n.Term {
		var buf bytes.Buffer
		const prefix = "    "
		lineStart := true
		for _, r := range n.Value {
			if lineStart {
				buf.WriteString(prefix)
			}
			buf.WriteRune(r)
			lineStart = r == '\n'
		}
		qw.writeBytes(buf.Bytes())
		return
	}
	qw.writeString("```")
	qw.writeString(n.Lang)
	qw.writeBytes(newLine)
	qw.writeString(n.Value)
	if !qw.lineStart {
		qw.writeBytes(newLine)
	}
	qw.writeString("```")
}

func (qw *qwiklabsWriter) list(n *types.ListNode) {
	if n.Block() == true {
		qw.newBlock()
	}
	qw.write(n.Nodes...)
	if !qw.lineStart {
		qw.writeBytes(newLine)
	}
}

func (qw *qwiklabsWriter) itemsList(n *types.ItemsListNode) {
	qw.newBlock()
	for i, item := range n.Items {
		s := "* "
		if n.Type() == types.NodeItemsList && n.Start > 0 {
			s = strconv.Itoa(i+n.Start) + ". "
		}
		qw.writeString(s)
		qw.write(item.Nodes...)
		if !qw.lineStart {
			qw.writeBytes(newLine)
		}
	}
}

func (qw *qwiklabsWriter) infobox(n *types.InfoboxNode) {
	// Note: There is no defined mapping of a Codelabs info box to any default
	//   Markdown syntax. We have decided to mix raw HTML into our Qwiklabs
	//   Markdown documents.
	// Future work: We may choose to extend the Markdown syntax more rigorously.
	qw.newBlock()
	qw.writeString(`<div class="codelabs-infobox codelabs-infobox-`)
	qw.writeEscape(string(n.Kind))
	qw.writeString(`">`)

	// Take advantage of the existing HTML writer to transform
	// the body of thi infobox.
	WriteHTML(qw.w, qw.env, n.Content.Nodes...)

	qw.writeString("</div>")

	// Alternatively we could use the HTML renderers output for the whole node
	// not just its contents.
	//WriteHTML(qw.w, qw.env, n)
}

func (qw *qwiklabsWriter) header(n *types.HeaderNode) {
	qw.newBlock()
	qw.writeString(strings.Repeat("#", n.Level+1))
	qw.writeString(" ")
	qw.write(n.Content.Nodes...)
	if !qw.lineStart {
		qw.writeBytes(newLine)
	}
}
