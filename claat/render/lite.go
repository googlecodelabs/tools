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

	"github.com/googlecodelabs/tools/claat/nodes"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Lite renders nodes as a standard HTML markup, without Custom Elements.
func Lite(ctx Context, nodes ...nodes.Node) (htmlTemplate.HTML, error) {
	var buf bytes.Buffer
	if err := WriteLite(&buf, ctx.Env, nodes...); err != nil {
		return "", err
	}
	return htmlTemplate.HTML(buf.String()), nil
}

// WriteLite does the same as Lite but outputs rendered markup to w.
func WriteLite(w io.Writer, env string, nodes ...nodes.Node) error {
	lw := liteWriter{w: w, env: env}
	return lw.write(nodes...)
}

type liteWriter struct {
	w   io.Writer // output writer
	env string    // target environment
	err error     // error during any writeXxx methods
}

func (lw *liteWriter) matchEnv(v []string) bool {
	if len(v) == 0 || lw.env == "" {
		return true
	}
	i := sort.SearchStrings(v, lw.env)
	return i < len(v) && v[i] == lw.env
}

func (lw *liteWriter) write(nodes ...nodes.Node) error {
	doc := &html.Node{Type: html.DocumentNode}
	for _, n := range nodes {
		if hn := lw.htmlnode(n); hn != nil {
			doc.AppendChild(hn)
		}
	}
	return html.Render(lw.w, doc)
}

func (lw *liteWriter) htmlnode(n nodes.Node) *html.Node {
	if !lw.matchEnv(n.Env()) {
		return nil
	}
	var hn *html.Node
	switch n := n.(type) {
	case *nodes.TextNode:
		hn = lw.text(n)
	case *nodes.ImageNode:
		hn = lw.image(n)
	case *nodes.URLNode:
		hn = lw.alink(n)
	case *nodes.ButtonNode:
		hn = lw.button(n)
	case *nodes.CodeNode:
		hn = lw.code(n)
	case *nodes.ListNode:
		hn = lw.list(n)
	case *nodes.ImportNode:
		if len(n.Content.Nodes) > 0 {
			hn = lw.list(n.Content)
		}
	case *nodes.ItemsListNode:
		hn = lw.itemsList(n)
	case *nodes.GridNode:
		hn = lw.grid(n)
	case *nodes.InfoboxNode:
		hn = lw.infobox(n)
	case *nodes.SurveyNode:
		hn = lw.survey(n)
	case *nodes.HeaderNode:
		hn = lw.header(n)
	case *nodes.YouTubeNode:
		hn = lw.youtube(n)
	}
	return hn
}

func (lw *liteWriter) text(n *nodes.TextNode) *html.Node {
	top := &html.Node{Type: html.TextNode, Data: n.Value}
	if n.Bold {
		hn := &html.Node{Type: html.ElementNode, Data: atom.Strong.String()}
		hn.AppendChild(top)
		top = hn
	}
	if n.Italic {
		hn := &html.Node{Type: html.ElementNode, Data: atom.Em.String()}
		hn.AppendChild(top)
		top = hn
	}
	if n.Code {
		hn := &html.Node{Type: html.ElementNode, Data: atom.Code.String()}
		hn.AppendChild(top)
		top = hn
	}
	return top
}

func (lw *liteWriter) image(n *nodes.ImageNode) *html.Node {
	hn := &html.Node{
		Type: html.ElementNode,
		Data: atom.Img.String(),
		Attr: []html.Attribute{{Key: "src", Val: n.Src}},
	}
	if n.Width > 0 {
		hn.Attr = append(hn.Attr, html.Attribute{
			Key: "style",
			Val: fmt.Sprintf("width: %.2fpx", n.Width),
		})
	}
	return hn
}

func (lw *liteWriter) alink(n *nodes.URLNode) *html.Node {
	top := &html.Node{Type: html.ElementNode, Data: atom.A.String()}
	if n.URL != "" {
		top.Attr = append(top.Attr, html.Attribute{Key: "href", Val: n.URL})
	}
	if n.Name != "" {
		top.Attr = append(top.Attr, html.Attribute{Key: "name", Val: n.Name})
	}
	if n.Target != "" {
		top.Attr = append(top.Attr, html.Attribute{Key: "target", Val: n.Target})
	}
	for _, cn := range n.Content.Nodes {
		if hn := lw.htmlnode(cn); hn != nil {
			top.AppendChild(hn)
		}
	}
	return top
}

func (lw *liteWriter) button(n *nodes.ButtonNode) *html.Node {
	cls := []string{"step__button"}
	if n.Color {
		cls = append(cls, "button--colored")
	}
	if n.Raise {
		cls = append(cls, "button--raised")
	}
	if n.Download {
		cls = append(cls, "button--download")
	}
	top := &html.Node{
		Type: html.ElementNode,
		Data: atom.A.String(),
		Attr: []html.Attribute{{Key: "class", Val: strings.Join(cls, " ")}},
	}
	for _, cn := range n.Content.Nodes {
		if hn := lw.htmlnode(cn); hn != nil {
			top.AppendChild(hn)
		}
	}
	return top
}

func (lw *liteWriter) code(n *nodes.CodeNode) *html.Node {
	top := &html.Node{Type: html.TextNode, Data: n.Value}

	if !n.Term {
		hn := &html.Node{Type: html.ElementNode, Data: atom.Code.String()}
		if n.Lang != "" {
			hn.Attr = append(hn.Attr, html.Attribute{
				Key: "language",
				Val: n.Lang,
			})
			hn.Attr = append(hn.Attr, html.Attribute{
				Key: "class",
				Val: n.Lang,
			})
		}
		hn.AppendChild(top)
		top = hn
	}

	hn := &html.Node{Type: html.ElementNode, Data: atom.Pre.String()}
	hn.AppendChild(top)
	top = hn

	return top
}

func (lw *liteWriter) list(n *nodes.ListNode) *html.Node {
	a := atom.P
	if n.Block() != true {
		a = atom.Div
	}
	top := &html.Node{Type: html.ElementNode, Data: a.String()}
	for _, cn := range n.Nodes {
		if hn := lw.htmlnode(cn); hn != nil {
			top.AppendChild(hn)
		}
	}
	return top
}

func (lw *liteWriter) itemsList(n *nodes.ItemsListNode) *html.Node {
	a := atom.Ul
	if n.Type() == nodes.NodeItemsList && n.Start > 0 {
		a = atom.Ol
	}
	top := &html.Node{Type: html.ElementNode, Data: a.String()}
	var itemCls string
	switch n.Type() {
	case nodes.NodeItemsCheck:
		itemCls = "checklist__item"
		top.Attr = append(top.Attr, html.Attribute{
			Key: "class",
			Val: "step__checklist",
		})
	case nodes.NodeItemsFAQ:
		itemCls = "faq__item"
		top.Attr = append(top.Attr, html.Attribute{
			Key: "class",
			Val: "step__faq",
		})
	default:
		if n.ListType != "" {
			top.Attr = append(top.Attr, html.Attribute{
				Key: "type",
				Val: n.ListType,
			})
		}
		if n.Start > 0 {
			top.Attr = append(top.Attr, html.Attribute{
				Key: "start",
				Val: strconv.Itoa(n.Start),
			})
		}
	}
	for _, item := range n.Items {
		li := &html.Node{Type: html.ElementNode, Data: atom.Li.String()}
		if itemCls != "" {
			li.Attr = append(li.Attr, html.Attribute{Key: "class", Val: itemCls})
		}
		for _, cn := range item.Nodes {
			if hn := lw.htmlnode(cn); hn != nil {
				li.AppendChild(hn)
			}
		}
		top.AppendChild(li)
	}
	return top
}

func (lw *liteWriter) grid(n *nodes.GridNode) *html.Node {
	top := &html.Node{Type: html.ElementNode, Data: atom.Table.String()}
	for _, r := range n.Rows {
		tr := &html.Node{Type: html.ElementNode, Data: atom.Tr.String()}
		for _, c := range r {
			td := &html.Node{
				Type: html.ElementNode,
				Data: atom.Td.String(),
				Attr: []html.Attribute{
					{Key: "colspan", Val: strconv.Itoa(c.Colspan)},
					{Key: "rowspan", Val: strconv.Itoa(c.Rowspan)},
				},
			}
			for _, cn := range c.Content.Nodes {
				if hn := lw.htmlnode(cn); hn != nil {
					td.AppendChild(hn)
				}
			}
			tr.AppendChild(td)
		}
		top.AppendChild(tr)
	}
	return top
}

func (lw *liteWriter) infobox(n *nodes.InfoboxNode) *html.Node {
	top := &html.Node{
		Type: html.ElementNode,
		Data: atom.Div.String(),
		Attr: []html.Attribute{{
			Key: "class",
			Val: fmt.Sprintf("step__note note--%s", n.Kind),
		}},
	}
	for _, cn := range n.Content.Nodes {
		if hn := lw.htmlnode(cn); hn != nil {
			top.AppendChild(hn)
		}
	}
	return top
}

func (lw *liteWriter) survey(n *nodes.SurveyNode) *html.Node {
	top := &html.Node{
		Type: html.ElementNode,
		Data: atom.Div.String(),
		Attr: []html.Attribute{
			{Key: "class", Val: "step__survey"},
			{Key: "data-survey-id", Val: n.ID},
		},
	}
	for i, g := range n.Groups {
		h4 := &html.Node{
			Type: html.ElementNode,
			Data: atom.H4.String(),
			Attr: []html.Attribute{{Key: "class", Val: "survey__q"}},
		}
		h4.AppendChild(&html.Node{Type: html.TextNode, Data: g.Name})
		top.AppendChild(h4)
		id := fmt.Sprintf("%s-%d", n.ID, i)
		for _, o := range g.Options {
			oh := &html.Node{
				Type: html.ElementNode,
				Data: atom.Input.String(),
				Attr: []html.Attribute{
					{Key: "type", Val: "radio"},
					{Key: "name", Val: id},
					{Key: "value", Val: o},
				},
			}
			lab := &html.Node{
				Type: html.ElementNode,
				Data: atom.Label.String(),
				Attr: []html.Attribute{{Key: "class", Val: "survey__a"}},
			}
			lab.AppendChild(oh)
			lab.AppendChild(&html.Node{Type: html.TextNode, Data: o})
			top.AppendChild(lab)
		}
	}
	return top
}

func (lw *liteWriter) header(n *nodes.HeaderNode) *html.Node {
	var cls string
	switch n.Type() {
	case nodes.NodeHeaderCheck:
		cls = "checklist"
	case nodes.NodeHeaderFAQ:
		cls = "faq"
	}
	top := &html.Node{
		Type: html.ElementNode,
		Data: "h" + strconv.Itoa(n.Level),
	}
	if cls != "" {
		top.Attr = append(top.Attr, html.Attribute{Key: "class", Val: cls})
	}
	for _, cn := range n.Content.Nodes {
		if hn := lw.htmlnode(cn); hn != nil {
			top.AppendChild(hn)
		}
	}
	return top
}

func (lw *liteWriter) youtube(n *nodes.YouTubeNode) *html.Node {
	top := &html.Node{
		Type: html.ElementNode,
		Data: atom.Div.String(),
		Attr: []html.Attribute{{Key: "class", Val: "keep-ar"}},
	}
	pad := &html.Node{
		Type: html.ElementNode,
		Data: atom.Div.String(),
		Attr: []html.Attribute{{Key: "class", Val: "keep-ar__pad"}},
	}
	box := &html.Node{
		Type: html.ElementNode,
		Data: atom.Iframe.String(),
		Attr: []html.Attribute{
			{Key: "src", Val: fmt.Sprintf("https://www.youtube.com/embed/%s?rel=0", n.VideoID)},
			{Key: "allow", Val: "accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture"},
			{Key: "allowfullscreen", Val: "1"},
			{Key: "class", Val: "keep-ar__box"},
		},
	}
	top.AppendChild(pad)
	pad.AppendChild(box)
	return top
}
