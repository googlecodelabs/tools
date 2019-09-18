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
	"fmt"
	htmlTemplate "html/template"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/googlecodelabs/tools/claat/types"
)

// TODO: render HTML using golang/x/net/html or template.

var (
	doubleQuote = []byte{'"'}
	lessThan    = []byte{'<'}
	greaterThan = []byte{'>'}
	newLine     = []byte{'\n'}
)

// HTML renders nodes as the markup for the target env.
func HTML(env string, nodes ...types.Node) (htmlTemplate.HTML, error) {
	var buf bytes.Buffer
	if err := WriteHTML(&buf, env, nodes...); err != nil {
		return "", err
	}
	return htmlTemplate.HTML(buf.String()), nil
}

// WriteHTML does the same as HTML but outputs rendered markup to w.
func WriteHTML(w io.Writer, env string, nodes ...types.Node) error {
	hw := htmlWriter{w: w, env: env}
	return hw.write(nodes...)
}

// ReplaceDoubleCurlyBracketsWithEntity replaces Double Curly Brackets with their charater entity.
func ReplaceDoubleCurlyBracketsWithEntity(s string) string {
	return strings.Replace(s, "{{", "&#123;&#123;", -1)
}

type htmlWriter struct {
	w   io.Writer // output writer
	env string    // target environment
	err error     // error during any writeXxx methods
}

func (hw *htmlWriter) matchEnv(v []string) bool {
	if len(v) == 0 || hw.env == "" {
		return true
	}
	i := sort.SearchStrings(v, hw.env)
	return i < len(v) && v[i] == hw.env
}

func (hw *htmlWriter) write(nodes ...types.Node) error {
	for _, n := range nodes {
		if !hw.matchEnv(n.Env()) {
			continue
		}
		switch n := n.(type) {
		case *types.TextNode:
			hw.text(n)
		case *types.ImageNode:
			hw.image(n)
		case *types.URLNode:
			hw.url(n)
		case *types.ButtonNode:
			hw.button(n)
		case *types.CodeNode:
			hw.code(n)
			hw.writeBytes(newLine)
		case *types.ListNode:
			hw.list(n)
			hw.writeBytes(newLine)
		case *types.ImportNode:
			if len(n.Content.Nodes) == 0 {
				break
			}
			hw.list(n.Content)
			hw.writeBytes(newLine)
		case *types.ItemsListNode:
			hw.itemsList(n)
			hw.writeBytes(newLine)
		case *types.GridNode:
			hw.grid(n)
			hw.writeBytes(newLine)
		case *types.InfoboxNode:
			hw.infobox(n)
			hw.writeBytes(newLine)
		case *types.SurveyNode:
			hw.survey(n)
			hw.writeBytes(newLine)
		case *types.HeaderNode:
			hw.header(n)
			hw.writeBytes(newLine)
		case *types.YouTubeNode:
			hw.youtube(n)
			hw.writeBytes(newLine)
		case *types.IframeNode:
			hw.iframe(n)
			hw.writeBytes(newLine)
		}
		if hw.err != nil {
			return hw.err
		}
	}
	return nil
}

func (hw *htmlWriter) writeBytes(b []byte) {
	if hw.err != nil {
		return
	}
	_, hw.err = hw.w.Write(b)
}

func (hw *htmlWriter) writeString(s string) {
	hw.writeBytes([]byte(s))
}

func (hw *htmlWriter) writeFmt(f string, a ...interface{}) {
	hw.writeString(fmt.Sprintf(f, a...))
}

func (hw *htmlWriter) writeEscape(s string) {
	s = htmlTemplate.HTMLEscapeString(s)
	hw.writeString(ReplaceDoubleCurlyBracketsWithEntity(s))
}

func (hw *htmlWriter) text(n *types.TextNode) {
	if n.Bold {
		hw.writeString("<strong>")
	}
	if n.Italic {
		hw.writeString("<em>")
	}
	if n.Code {
		hw.writeString("<code>")
	}
	s := htmlTemplate.HTMLEscapeString(n.Value)
	s = ReplaceDoubleCurlyBracketsWithEntity(s)
	hw.writeString(strings.Replace(s, "\n", "<br>", -1))
	if n.Code {
		hw.writeString("</code>")
	}
	if n.Italic {
		hw.writeString("</em>")
	}
	if n.Bold {
		hw.writeString("</strong>")
	}
}

func (hw *htmlWriter) image(n *types.ImageNode) {
	hw.writeString("<img")
	if n.Alt != "" {
		hw.writeFmt(" alt=%q", n.Alt)
	}
	if n.Title != "" {
		hw.writeFmt(" title=%q", n.Title)
	}
	if n.Width > 0 {
		hw.writeFmt(` style="width: %.2fpx"`, n.Width)
	}
	hw.writeString(` src="`)
	hw.writeString(n.Src)
	hw.writeBytes(doubleQuote)
	hw.writeBytes(greaterThan)
}

func (hw *htmlWriter) url(n *types.URLNode) {
	hw.writeString("<a")
	if n.URL != "" {
		hw.writeString(` href="`)
		hw.writeString(n.URL)
		hw.writeBytes(doubleQuote)
	}
	if n.Name != "" {
		hw.writeString(` name="`)
		hw.writeEscape(n.Name)
		hw.writeBytes(doubleQuote)
	}
	if n.Target != "" {
		hw.writeString(` target="`)
		hw.writeEscape(n.Target)
		hw.writeBytes(doubleQuote)
	}
	hw.writeBytes(greaterThan)
	hw.write(n.Content.Nodes...)
	hw.writeString("</a>")
}

func (hw *htmlWriter) button(n *types.ButtonNode) {
	hw.writeString("<paper-button")
	if n.Colored {
		hw.writeString(` class="colored"`)
	}
	if n.Raised {
		hw.writeString(" raised")
	}
	hw.writeBytes(greaterThan)
	if n.Download {
		hw.writeString(`<iron-icon icon="file-download"></iron-icon>`)
	}
	hw.write(n.Content.Nodes...)
	hw.writeString("</paper-button>")
}

func (hw *htmlWriter) code(n *types.CodeNode) {
	hw.writeString("<pre>")
	if !n.Term {
		hw.writeString("<code")
		if n.Lang != "" {
			hw.writeFmt(" language=%q class=%q", n.Lang, n.Lang)
		}
		hw.writeBytes(greaterThan)
	}
	hw.writeEscape(n.Value)
	if !n.Term {
		hw.writeString("</code>")
	}
	hw.writeString("</pre>")
}

func (hw *htmlWriter) list(n *types.ListNode) {
	wrap := n.Block() == true
	if wrap {
		if onlyImages(n.Nodes...) {
			hw.writeString(`<p class="image-container">`)
		} else {
			hw.writeString("<p>")
		}
	}
	hw.write(n.Nodes...)
	if wrap {
		hw.writeString("</p>")
	}
}

// Returns true if the list of Nodes contains only images or white spaces.
func onlyImages(nodes ...types.Node) bool {
	for _, n := range nodes {
		switch n := n.(type) {
		case *types.TextNode:
			if len(strings.TrimSpace(n.Value)) == 0 {
				continue
			}
			return false
		case *types.ImageNode:
			continue
		default:
			return false
		}
	}
	return true
}

func (hw *htmlWriter) itemsList(n *types.ItemsListNode) {
	tag := "ul"
	if n.Type() == types.NodeItemsList && (n.Start > 0 || n.ListType != "") {
		tag = "ol"
	}
	hw.writeBytes(lessThan)
	hw.writeString(tag)
	switch n.Type() {
	case types.NodeItemsCheck:
		hw.writeString(` class="checklist"`)
	case types.NodeItemsFAQ:
		hw.writeString(` class="faq"`)
	default:
		if n.ListType != "" {
			hw.writeString(` type="`)
			hw.writeString(n.ListType)
			hw.writeBytes(doubleQuote)
		}
		if n.Start > 0 {
			hw.writeFmt(` start="%d"`, n.Start)
		}
	}
	hw.writeBytes(greaterThan)
	hw.writeBytes(newLine)

	for _, i := range n.Items {
		hw.writeString("<li>")
		hw.write(i.Nodes...)
		hw.writeString("</li>\n")
	}

	hw.writeString("</")
	hw.writeString(tag)
	hw.writeBytes(greaterThan)
}

func (hw *htmlWriter) grid(n *types.GridNode) {
	hw.writeString("<table>\n")
	for _, r := range n.Rows {
		hw.writeString("<tr>")
		for _, c := range r {
			hw.writeFmt(`<td colspan="%d" rowspan="%d">`, c.Colspan, c.Rowspan)
			hw.write(c.Content.Nodes...)
			hw.writeString("</td>")
		}
		hw.writeString("</tr>\n")
	}
	hw.writeString("</table>")
}

func (hw *htmlWriter) infobox(n *types.InfoboxNode) {
	hw.writeString(`<aside class="`)
	hw.writeEscape(string(n.Kind))
	hw.writeString(`">`)
	hw.write(n.Content.Nodes...)
	hw.writeString("</aside>")
}

func (hw *htmlWriter) survey(n *types.SurveyNode) {
	hw.writeString(`<google-codelab-survey survey-id="`)
	hw.writeString(n.ID)
	hw.writeBytes(doubleQuote)
	hw.writeString(">\n")
	for _, g := range n.Groups {
		hw.writeString("<h4>")
		hw.writeEscape(g.Name)
		hw.writeString("</h4>\n<paper-radio-group>\n")
		for _, o := range g.Options {
			hw.writeString("<paper-radio-button>")
			hw.writeEscape(o)
			hw.writeString("</paper-radio-button>\n")
		}
		hw.writeString("</paper-radio-group>\n")
	}
	hw.writeString("</google-codelab-survey>")
}

func (hw *htmlWriter) header(n *types.HeaderNode) {
	tag := "h" + strconv.Itoa(n.Level)
	hw.writeBytes(lessThan)
	hw.writeString(tag)
	switch n.Type() {
	case types.NodeHeaderCheck:
		hw.writeString(` class="checklist"`)
	case types.NodeHeaderFAQ:
		hw.writeString(` class="faq"`)

	}
	hw.writeString(` is-upgraded`)
	hw.writeBytes(greaterThan)
	hw.write(n.Content.Nodes...)
	hw.writeString("</")
	hw.writeString(tag)
	hw.writeBytes(greaterThan)
}

func (hw *htmlWriter) youtube(n *types.YouTubeNode) {
	hw.writeFmt(`<iframe class="youtube-video" `+
		`src="https://www.youtube.com/embed/%s?rel=0" allow="accelerometer; `+
		`autoplay; encrypted-media; gyroscope; picture-in-picture" `+
		`allowfullscreen></iframe>`, n.VideoID)
}

func (hw *htmlWriter) iframe(n *types.IframeNode) {
	hw.writeFmt(`<iframe class="embedded-iframe" src="%s"></iframe>`,
		n.URL)
}
