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
	"io"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/googlecodelabs/tools/claat/types"
)

// MD renders nodes as markdown for the target env.
func MD(env string, nodes ...types.Node) (string, error) {
	var buf bytes.Buffer
	if err := WriteMD(&buf, env, nodes...); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// WriteMD does the same as MD but outputs rendered markup to w.
func WriteMD(w io.Writer, env string, nodes ...types.Node) error {
	mw := mdWriter{w: w, env: env}
	return mw.write(nodes...)
}

type mdWriter struct {
	w         io.Writer // output writer
	env       string    // target environment
	err       error     // error during any writeXxx methods
	lineStart bool
}

func (mw *mdWriter) writeBytes(b []byte) {
	if mw.err != nil {
		return
	}
	mw.lineStart = len(b) > 0 && b[len(b)-1] == '\n'
	_, mw.err = mw.w.Write(b)
}

func (mw *mdWriter) writeString(s string) {
	mw.writeBytes([]byte(s))
}

func (mw *mdWriter) space() {
	if !mw.lineStart {
		mw.writeString(" ")
	}
}

func (mw *mdWriter) newBlock() {
	if !mw.lineStart {
		mw.writeBytes(newLine)
	}
	mw.writeBytes(newLine)
}

func (mw *mdWriter) matchEnv(v []string) bool {
	if len(v) == 0 || mw.env == "" {
		return true
	}
	i := sort.SearchStrings(v, mw.env)
	return i < len(v) && v[i] == mw.env
}

func (mw *mdWriter) write(nodes ...types.Node) error {
	for _, n := range nodes {
		if !mw.matchEnv(n.Env()) {
			continue
		}
		switch n := n.(type) {
		case *types.TextNode:
			mw.text(n)
		case *types.ImageNode:
			mw.image(n)
		case *types.URLNode:
			mw.url(n)
		case *types.ButtonNode:
			mw.write(n.Content.Nodes...)
		case *types.CodeNode:
			mw.code(n)
		case *types.ListNode:
			mw.list(n)
		case *types.ImportNode:
			if len(n.Content.Nodes) == 0 {
				break
			}
			mw.write(n.Content.Nodes...)
		case *types.ItemsListNode:
			mw.itemsList(n)
		//case *types.GridNode:
		//	mw.grid(n)
		case *types.InfoboxNode:
			mw.infobox(n)
		//case *types.SurveyNode:
		//	mw.survey(n)
		case *types.HeaderNode:
			mw.header(n)
		case *types.YouTubeNode:
			mw.youtube(n)
		}
		if mw.err != nil {
			return mw.err
		}
	}
	return nil
}

func (mw *mdWriter) text(n *types.TextNode) {
	t := strings.TrimSpace(n.Value)
	tl := len([]rune(t))
	nl := len([]rune(n.Value))
	ls := nl - len([]rune(strings.TrimLeft(n.Value, " ")))
	// Don't just copy above and TrimRight instead of TrimLeft to avoid " " counting as 1
	// left space and 1 right space. Instead, number of right spaces is
	// length of whole string - length of string with spaces trimmed - number of left spaces.
	rs := nl - tl - ls

	mw.writeString(strings.Repeat(" ", ls))
	if tl > 0 {
		if n.Bold {
			mw.writeString("**")
		}
		if n.Italic {
			mw.writeString("*")
		}
		if n.Code {
			mw.writeString("`")
		}
	}

	mw.writeString(t)

	if tl > 0 {
		if n.Code {
			mw.writeString("`")
		}
		if n.Italic {
			mw.writeString("*")
		}
		if n.Bold {
			mw.writeString("**")
		}
	}
	mw.writeString(strings.Repeat(" ", rs))
}

func (mw *mdWriter) image(n *types.ImageNode) {
	mw.space()
	mw.writeString("<img ")
	mw.writeString(fmt.Sprintf("src=\"%s\" ", n.Src))

	if n.Alt != "" {
		mw.writeString(fmt.Sprintf("alt=\"%s\" ", n.Alt))
	} else {
		mw.writeString(fmt.Sprintf("alt=\"%s\" ", path.Base(n.Src)))
	}

	if n.Title != "" {
		mw.writeString(fmt.Sprintf("title=\"%q\" ", n.Title))
	}

	// If available append width to the src string of the image.
	if n.Width > 0 {
		mw.writeString(fmt.Sprintf(" width=\"%.2f\" ", n.Width))
	}

	mw.writeString("/>")
}

func (mw *mdWriter) url(n *types.URLNode) {
	mw.space()
	if n.URL != "" {
		// Look-ahead for button syntax.
		if _, ok := n.Content.Nodes[0].(*types.ButtonNode); ok {
			mw.writeString("<button>")
		}
		mw.writeString("[")
	}
	mw.write(n.Content.Nodes...)
	if n.URL != "" {
		mw.writeString("](")
		mw.writeString(n.URL)
		mw.writeString(")")
		if _, ok := n.Content.Nodes[0].(*types.ButtonNode); ok {
			// Look-ahead for button syntax.
			mw.writeString("</button>")
		}
	}
}

func (mw *mdWriter) code(n *types.CodeNode) {
	mw.newBlock()
	defer mw.writeBytes(newLine)
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
		mw.writeBytes(buf.Bytes())
		return
	}
	mw.writeString("```")
	mw.writeString(n.Lang)
	mw.writeBytes(newLine)
	mw.writeString(n.Value)
	if !mw.lineStart {
		mw.writeBytes(newLine)
	}
	mw.writeString("```")
}

func (mw *mdWriter) list(n *types.ListNode) {
	if n.Block() == true {
		mw.newBlock()
	}
	mw.write(n.Nodes...)
	if !mw.lineStart {
		mw.writeBytes(newLine)
	}
}

func (mw *mdWriter) itemsList(n *types.ItemsListNode) {
	mw.newBlock()
	for i, item := range n.Items {
		s := "* "
		if n.Type() == types.NodeItemsList && n.Start > 0 {
			s = strconv.Itoa(i+n.Start) + ". "
		}
		mw.writeString(s)
		mw.write(item.Nodes...)
		if !mw.lineStart {
			mw.writeBytes(newLine)
		}
	}
}

func (mw *mdWriter) infobox(n *types.InfoboxNode) {
	// InfoBoxes are comprised of a ListNode with the contents of the InfoBox. 
	// Writing the ListNode directly results in extra newlines in the md output
	// which breaks the formatting. So instead, write the ListNode's children 
	// directly and don't write the ListNode itself.
	ln, ok := n.Content.Nodes[0].(*types.ListNode)
	if !ok {
		return
	}
	mw.newBlock()
	k := "Positive"
	if n.Kind == types.InfoboxNegative {
		k = "Negative"
	}
	mw.writeString(k)
	mw.writeString("\n: ")
	mw.write(ln.Nodes...)
}

func (mw *mdWriter) header(n *types.HeaderNode) {
	mw.newBlock()
	mw.writeString(strings.Repeat("#", n.Level+1))
	mw.writeString(" ")
	mw.write(n.Content.Nodes...)
	if !mw.lineStart {
		mw.writeBytes(newLine)
	}
}

func (mw *mdWriter) youtube(n *types.YouTubeNode) {
	mw.newBlock()
	mw.writeString(fmt.Sprintf(`<video id="%s"></video>`, n.VideoID))
}
