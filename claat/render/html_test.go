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
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googlecodelabs/tools/claat/nodes"
)

func TestHTMLEnv(t *testing.T) {
	one := nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "one "})
	one.MutateEnv([]string{"one"})
	two := nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "two "})
	two.MutateEnv([]string{"two"})
	three := nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "three "})
	three.MutateEnv([]string{"one", "three"})

	tests := []struct {
		env    string
		output string
	}{
		{"", "one two three "},
		{"one", "one three "},
		{"two", "two "},
		{"three", "three "},
		{"four", ""},
	}
	for i, test := range tests {
		var ctx Context
		ctx.Env = test.env
		h, err := HTML(ctx, one, two, three)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}
		if v := string(h); v != test.output {
			t.Errorf("%d: v = %q; want %q", i, v, test.output)
		}
	}
}

// TODO: test HTML
// TODO: test writeHTML

func TestReplaceDoubleCurlyBracketsWithEntity(t *testing.T) {
	tests := []struct {
		name  string
		inStr string
		out   string
	}{
		{
			name: "Simple",
		},
		{
			name:  "Zero",
			inStr: "foobar",
			out:   "foobar",
		},
		{
			name:  "Single",
			inStr: "foo{{bar",
			out:   "foo&#123;&#123;bar",
		},
		{
			name:  "Multi",
			inStr: "foo{{bar{{baz",
			out:   "foo&#123;&#123;bar&#123;&#123;baz",
		},
		{
			name:  "OverlapEven",
			inStr: "{{{{{{",
			out:   "&#123;&#123;&#123;&#123;&#123;&#123;",
		},
		{
			name:  "OverlapOdd",
			inStr: "{{{{{{{",
			out:   "&#123;&#123;&#123;&#123;&#123;&#123;{",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := ReplaceDoubleCurlyBracketsWithEntity(tc.inStr)
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("ReplaceDoubleCurlyBracketsWithEntity(%q) got diff (-want +got):\n%s", tc.inStr, diff)
			}
		})
	}
}

func TestMatchEnv(t *testing.T) {
	tests := []struct {
		name  string
		inEnv string
		inV   []string
		out   bool
	}{
		{
			name: "Empty",
			out:  true,
		},
		{
			name:  "NoChecks",
			inEnv: "foo",
			out:   true,
		},
		{
			name: "NoEnv",
			inV:  []string{"foo", "bar", "baz"},
			out:  true,
		},
		{
			name:  "SimpleMatch",
			inEnv: "foo",
			inV:   []string{"foo"},
			out:   true,
		},
		{
			// TODO: should this be false?
			name:  "MultiMatch",
			inEnv: "foo",
			inV:   []string{"foo", "bar", "baz"},
		},
		{
			name:  "NoMatch",
			inEnv: "foo",
			inV:   []string{"bar", "baz", "qux"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer, env: tc.inEnv}
			if out := hw.matchEnv(tc.inV); out != tc.out {
				t.Errorf("hw.matchEnv(%+v) = %t, want %t", tc.inV, out, tc.out)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name    string
		inNodes []nodes.Node
		out     string
	}{
		{
			name: "Empty",
		},
		{
			name: "Text",
			inNodes: []nodes.Node{
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foobar"}),
			},
			out: "foobar",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.write(tc.inNodes...)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.write(%+v) got diff (-want +got):\n%s", tc.inNodes, diff)
			}
		})
	}
}

func TestWriteString(t *testing.T) {
	b := &bytes.Buffer{}
	hw := htmlWriter{
		w: b,
	}
	hw.writeString("foobar")
	out := b.String()
	want := "foobar"
	if out != want {
		t.Errorf("hw.String() = %q, want %q", out, want)
	}
}

func TestWriteStringError(t *testing.T) {
	b := &bytes.Buffer{}
	hw := htmlWriter{
		w:   b,
		err: errors.New("foobar"),
	}
	hw.writeString("foobar")
	out := b.String()
	want := ""
	if out != want {
		t.Errorf("hw.String() = %q, want %q", out, want)
	}
}

func TestWriteFmt(t *testing.T) {
	tests := []struct {
		name   string
		inStr  string
		inArgs []interface{}
		out    string
	}{
		{
			name: "Empty",
		},
		{
			name:  "Simple",
			inStr: "foobar",
			out:   "foobar",
		},
		{
			name:   "Format",
			inStr:  "foo%s",
			inArgs: []interface{}{"bar"},
			out:    "foobar",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.writeFmt(tc.inStr, tc.inArgs...)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.writeFmt(%q) got diff (-want +got):\n%s", tc.inStr, diff)
			}
		})
	}
}

func TestEscape(t *testing.T) {
	tests := []struct {
		name  string
		inStr string
		out   string
	}{
		{
			name: "Empty",
		},
		{
			name:  "NothingToEscape",
			inStr: "foobar",
			out:   "foobar",
		},
		{
			name:  "<",
			inStr: "foo<bar",
			out:   "foo&lt;bar",
		},
		{
			name:  "{{",
			inStr: "foo{{bar",
			out:   "foo&#123;&#123;bar",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := escape(tc.inStr)
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("escape(%q) got diff (-want +got):\n%s", tc.inStr, diff)
			}
		})
	}
}

func TestWriteEscape(t *testing.T) {
	tests := []struct {
		name  string
		inStr string
		out   string
	}{
		{
			name: "Empty",
		},
		{
			name:  "NothingToEscape",
			inStr: "foobar",
			out:   "foobar",
		},
		{
			name:  "<",
			inStr: "foo<bar",
			out:   "foo&lt;bar",
		},
		{
			name:  "{{",
			inStr: "foo{{bar",
			out:   "foo&#123;&#123;bar",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.writeEscape(tc.inStr)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.writeEscape(%q) got diff (-want +got):\n%s", tc.inStr, diff)
			}
		})
	}
}

func TestText(t *testing.T) {
	tests := []struct {
		name   string
		inNode *nodes.TextNode
		out    string
	}{
		{
			name:   "Empty",
			inNode: nodes.NewTextNode(nodes.NewTextNodeOptions{Value: ""}),
		},
		{
			name:   "Simple",
			inNode: nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foobar"}),
			out:    "foobar",
		},
		{
			name: "Bold",
			inNode: nodes.NewTextNode(nodes.NewTextNodeOptions{
				Value: "foobar",
				Bold:  true,
			}),
			out: "<strong>foobar</strong>",
		},
		{
			name: "Italic",
			inNode: nodes.NewTextNode(nodes.NewTextNodeOptions{
				Value:  "foobar",
				Italic: true,
			}),
			out: "<em>foobar</em>",
		},
		{
			name: "Code",
			inNode: nodes.NewTextNode(nodes.NewTextNodeOptions{
				Value: "foobar",
				Code:  true,
			}),
			out: "<code>foobar</code>",
		},
		{
			name: "BoldItalic",
			inNode: nodes.NewTextNode(nodes.NewTextNodeOptions{
				Value:  "foobar",
				Bold:   true,
				Italic: true,
			}),
			out: "<strong><em>foobar</em></strong>",
		},
		{
			name: "ItalicCode",
			inNode: nodes.NewTextNode(nodes.NewTextNodeOptions{
				Value:  "foobar",
				Italic: true,
				Code:   true,
			}),
			out: "<em><code>foobar</code></em>",
		},
		{
			name: "BoldCode",
			inNode: nodes.NewTextNode(nodes.NewTextNodeOptions{
				Value: "foobar",
				Bold:  true,
				Code:  true,
			}),
			out: "<strong><code>foobar</code></strong>",
		},
		{
			name: "BoldItalicCode",
			inNode: nodes.NewTextNode(nodes.NewTextNodeOptions{
				Value:  "foobar",
				Bold:   true,
				Italic: true,
				Code:   true,
			}),
			out: "<strong><em><code>foobar</code></em></strong>",
		},
		{
			name:   "HTMLEscape",
			inNode: nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foo<bar"}),
			out:    "foo&lt;bar",
		},
		{
			name:   "{{Escape",
			inNode: nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foo{{bar"}),
			out:    "foo&#123;&#123;bar",
		},
		{
			name:   "Newline",
			inNode: nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foo\nbar"}),
			out:    "foo<br>bar",
		},
		{
			name:   "NonBreakingSpaceTrim",
			inNode: nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "\uFEFFfoobar\uFEFF"}),
			out:    "foobar",
		},
		{
			name: "CodeEscape",
			inNode: nodes.NewTextNode(nodes.NewTextNodeOptions{
				Value: "foo<bar",
				Code:  true,
			}),
			out: "<code>foo<bar</code>",
		},
		{
			name: "Code{{",
			inNode: nodes.NewTextNode(nodes.NewTextNodeOptions{
				Value: "foo{{bar",
				Code:  true,
			}),
			out: "<code>foo&#123;&#123;bar</code>",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.text(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.text(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}

func TestImage(t *testing.T) {
	tests := []struct {
		name   string
		inNode *nodes.ImageNode
		out    string
	}{
		{
			name: "Simple",
			inNode: nodes.NewImageNode(nodes.NewImageNodeOptions{
				Src: "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
			}),
			out: `<img src="https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png">`,
		},
		{
			name: "Alt",
			inNode: nodes.NewImageNode(nodes.NewImageNodeOptions{
				Src: "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
				Alt: "foo",
			}),
			out: `<img alt="foo" src="https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png">`,
		},
		{
			name: "Title",
			inNode: nodes.NewImageNode(nodes.NewImageNodeOptions{
				Src:   "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
				Title: "bar",
			}),
			out: `<img title="bar" src="https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png">`,
		},
		{
			name: "Width",
			inNode: nodes.NewImageNode(nodes.NewImageNodeOptions{
				Src:   "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
				Width: 5.0,
			}),
			out: `<img style="width: 5.00px" src="https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png">`,
		},
		{
			name: "All",
			inNode: nodes.NewImageNode(nodes.NewImageNodeOptions{
				Alt:   "foo",
				Title: "bar",
				Src:   "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
				Width: 5.0,
			}),
			out: `<img alt="foo" title="bar" style="width: 5.00px" src="https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png">`,
		},
		{
			name: "WidthPrecision",
			inNode: nodes.NewImageNode(nodes.NewImageNodeOptions{
				Src:   "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
				Width: 9.87654321,
			}),
			out: `<img style="width: 9.88px" src="https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png">`,
		},
		{
			name: "WidthNegative",
			inNode: nodes.NewImageNode(nodes.NewImageNodeOptions{
				Src:   "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
				Width: -1.2345,
			}),
			out: `<img src="https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png">`,
		},
		{
			name:   "Empty",
			inNode: nodes.NewImageNode(nodes.NewImageNodeOptions{}),
			out:    `<img src="">`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.image(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.image(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}

func TestURL(t *testing.T) {
	a := nodes.NewURLNode("google.com")
	a.Name = "foobar"

	b := nodes.NewURLNode("google.com")
	b.Name = "foo{{bar"

	c := nodes.NewURLNode("google.com")
	c.Target = "_self"

	d := nodes.NewURLNode("google.com")
	d.Target = "_self{{"

	e := nodes.NewURLNode("google.com")
	e.Name = "foobar"
	e.Target = "_self"

	tests := []struct {
		name   string
		inNode *nodes.URLNode
		out    string
	}{
		{
			name:   "Empty",
			inNode: nodes.NewURLNode("google.com"),
			out:    `<a href="google.com" target="_blank"></a>`,
		},
		{
			name:   "Name",
			inNode: a,
			out:    `<a href="google.com" name="foobar" target="_blank"></a>`,
		},
		{
			name:   "NameEscape",
			inNode: b,
			out:    `<a href="google.com" name="foo&#123;&#123;bar" target="_blank"></a>`,
		},
		{
			name:   "Target",
			inNode: c,
			out:    `<a href="google.com" target="_self"></a>`,
		},
		{
			name:   "TargetEscape",
			inNode: d,
			out:    `<a href="google.com" target="_self&#123;&#123;"></a>`,
		},
		{
			name:   "NameTarget",
			inNode: e,
			out:    `<a href="google.com" name="foobar" target="_self"></a>`,
		},
		{
			name:   "Simple",
			inNode: nodes.NewURLNode("google.com", nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foobar"})),
			out:    `<a href="google.com" target="_blank">foobar</a>`,
		},
		{
			name:   "MultipleContent",
			inNode: nodes.NewURLNode("google.com", nodes.NewHeaderNode(1, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foo"})), nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "bar"})),
			out: `<a href="google.com" target="_blank"><h1 is-upgraded>foo</h1>
bar</a>`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.url(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.url(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}

func TestButton(t *testing.T) {
	tests := []struct {
		name   string
		inNode *nodes.ButtonNode
		out    string
	}{
		{
			name:   "Empty",
			inNode: nodes.NewButtonNode(false, false, false),
			out:    `<paper-button></paper-button>`,
		},
		{
			name:   "NoProperties",
			inNode: nodes.NewButtonNode(false, false, false, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foobar"})),
			out:    `<paper-button>foobar</paper-button>`,
		},
		{
			name:   "Raise",
			inNode: nodes.NewButtonNode(true, false, false, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foobar"})),
			out:    `<paper-button raised>foobar</paper-button>`,
		},
		{
			name:   "Color",
			inNode: nodes.NewButtonNode(false, true, false, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foobar"})),
			out:    `<paper-button class="colored">foobar</paper-button>`,
		},
		{
			name:   "Download",
			inNode: nodes.NewButtonNode(false, false, true, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foobar"})),
			out:    `<paper-button><iron-icon icon="file-download"></iron-icon>foobar</paper-button>`,
		},
		{
			name:   "RaiseColor",
			inNode: nodes.NewButtonNode(true, true, false, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foobar"})),
			out:    `<paper-button class="colored" raised>foobar</paper-button>`,
		},
		{
			name:   "ColorDownload",
			inNode: nodes.NewButtonNode(false, true, true, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foobar"})),
			out:    `<paper-button class="colored"><iron-icon icon="file-download"></iron-icon>foobar</paper-button>`,
		},
		{
			name:   "RaiseDownload",
			inNode: nodes.NewButtonNode(true, false, true, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foobar"})),
			out:    `<paper-button raised><iron-icon icon="file-download"></iron-icon>foobar</paper-button>`,
		},
		{
			name:   "RaiseColorDownload",
			inNode: nodes.NewButtonNode(true, true, true, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foobar"})),
			out:    `<paper-button class="colored" raised><iron-icon icon="file-download"></iron-icon>foobar</paper-button>`,
		},
		{
			name:   "MultipleContent",
			inNode: nodes.NewButtonNode(false, false, false, nodes.NewHeaderNode(2, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foo"})), nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "bar"})),
			out: `<paper-button><h2 is-upgraded>foo</h2>
bar</paper-button>`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.button(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.button(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}

func TestCode(t *testing.T) {
	tests := []struct {
		name   string
		inNode *nodes.CodeNode
		out    string
	}{
		{
			name:   "NoLang",
			inNode: nodes.NewCodeNode("foobar", false, ""),
			out:    `<pre><code>foobar</code></pre>`,
		},
		{
			name:   "Lang",
			inNode: nodes.NewCodeNode("foobar", false, "c"),
			out:    `<pre><code language="c" class="c">foobar</code></pre>`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.code(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.code(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name   string
		inNode *nodes.ListNode
		out    string
	}{
		{
			name:   "Zero",
			inNode: nodes.NewListNode(),
		},
		{
			name: "One",
			inNode: nodes.NewListNode(
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foobar"}),
			),
			out: "foobar",
		},
		{
			name: "Multi",
			inNode: nodes.NewListNode(
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foo"}),
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "bar"}),
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "baz"}),
			),
			out: "foobarbaz",
		},
		{
			name: "ListOfLists",
			inNode: nodes.NewListNode(
				nodes.NewListNode(
					nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foo"}),
					nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "bar"}),
				),
				nodes.NewListNode(
					nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "baz"}),
					nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "qux"}),
				),
			),
			out: "foobar\nbazqux\n",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.list(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.list(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}

func TestOnlyImages(t *testing.T) {
	tests := []struct {
		name    string
		inNodes []nodes.Node
		out     bool
	}{
		{
			name:    "None",
			inNodes: []nodes.Node{},
			out:     true,
		},
		{
			name: "OneImage",
			inNodes: []nodes.Node{
				nodes.NewImageNode(nodes.NewImageNodeOptions{Src: "foobar"}),
			},
			out: true,
		},
		{
			name: "MultiImages",
			inNodes: []nodes.Node{
				nodes.NewImageNode(nodes.NewImageNodeOptions{Src: "foo"}),
				nodes.NewImageNode(nodes.NewImageNodeOptions{Src: "bar"}),
				nodes.NewImageNode(nodes.NewImageNodeOptions{Src: "baz"}),
			},
			out: true,
		},
		{
			name: "OneWhitespace",
			inNodes: []nodes.Node{
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: " "}),
			},
			out: true,
		},
		{
			name: "MultiWhitespace",
			inNodes: []nodes.Node{
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: " "}),
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "\n"}),
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "\t"}),
			},
			out: true,
		},
		{
			name: "ImagesAndWhitespace",
			inNodes: []nodes.Node{
				nodes.NewImageNode(nodes.NewImageNodeOptions{Src: "foo"}),
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: " "}),
				nodes.NewImageNode(nodes.NewImageNodeOptions{Src: "bar"}),
				nodes.NewImageNode(nodes.NewImageNodeOptions{Src: "baz"}),
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "\n"}),
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "\t"}),
			},
			out: true,
		},
		{
			name: "Text",
			inNodes: []nodes.Node{
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "qux"}),
			},
		},
		{
			name: "TextAndImages",
			inNodes: []nodes.Node{
				nodes.NewImageNode(nodes.NewImageNodeOptions{Src: "foo"}),
				nodes.NewImageNode(nodes.NewImageNodeOptions{Src: "bar"}),
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "qux"}),
				nodes.NewImageNode(nodes.NewImageNodeOptions{Src: "baz"}),
			},
		},
		{
			name: "TextAndWhitespace",
			inNodes: []nodes.Node{
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: " "}),
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "\n"}),
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foo"}),
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "\t"}),
			},
		},
		{
			name: "TextImagesAndWhitespace",
			inNodes: []nodes.Node{
				nodes.NewImageNode(nodes.NewImageNodeOptions{Src: "foo"}),
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: " "}),
				nodes.NewImageNode(nodes.NewImageNodeOptions{Src: "bar"}),
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "qux"}),
				nodes.NewImageNode(nodes.NewImageNodeOptions{Src: "baz"}),
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "\n"}),
				nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "\t"}),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := onlyImages(tc.inNodes...); out != tc.out {
				t.Errorf("onlyImages(%v) = %t, want %t", tc.inNodes, out, tc.out)
			}
		})
	}
}

// TODO: test itemsList

func TestGrid(t *testing.T) {
	tests := []struct {
		name   string
		inNode *nodes.GridNode
		out    string
	}{
		{
			name:   "No rows or cells ",
			inNode: nodes.NewGridNode(),
			out:    "<table>\n</table>",
		},
		{
			name: "One row and no cells",
			inNode: nodes.NewGridNode([][]*nodes.GridCell{
				[]*nodes.GridCell{},
			}...),
			out: "<table>\n<tr></tr>\n</table>",
		},
		{
			name: "One row and one cell",
			inNode: nodes.NewGridNode([][]*nodes.GridCell{
				[]*nodes.GridCell{
					&nodes.GridCell{
						Content: nodes.NewListNode([]nodes.Node{
							nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foo"}),
						}...),
					},
				},
			}...),
			out: "<table>\n<tr><td colspan=\"0\" rowspan=\"0\">foo</td></tr>\n</table>",
		},
		{
			name: "One row and multiple cells",
			inNode: nodes.NewGridNode([][]*nodes.GridCell{
				[]*nodes.GridCell{
					&nodes.GridCell{
						Content: nodes.NewListNode([]nodes.Node{
							nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foo"}),
						}...),
					},
					&nodes.GridCell{
						Content: nodes.NewListNode([]nodes.Node{
							nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "bar"}),
						}...),
					},
				},
			}...),
			out: "<table>\n<tr><td colspan=\"0\" rowspan=\"0\">foo</td><td colspan=\"0\" rowspan=\"0\">bar</td></tr>\n</table>",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.grid(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.grid(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}

func TestInfobox(t *testing.T) {
	tests := []struct {
		name   string
		inNode *nodes.InfoboxNode
		out    string
	}{
		{
			name:   "PositiveEmpty",
			inNode: nodes.NewInfoboxNode(nodes.InfoboxPositive),
			out:    `<aside class="special"></aside>`,
		},
		{
			name:   "NegativeEmpty",
			inNode: nodes.NewInfoboxNode(nodes.InfoboxNegative),
			out:    `<aside class="warning"></aside>`,
		},
		{
			name:   "PositiveNonEmpty",
			inNode: nodes.NewInfoboxNode(nodes.InfoboxPositive, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foo"}), nodes.NewHeaderNode(3, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "bar"}))),
			out: `<aside class="special">foo<h3 is-upgraded>bar</h3>
</aside>`,
		},
		{
			name:   "NegativeNonEmpty",
			inNode: nodes.NewInfoboxNode(nodes.InfoboxNegative, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foo"}), nodes.NewHeaderNode(3, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "bar"}))),
			out: `<aside class="warning">foo<h3 is-upgraded>bar</h3>
</aside>`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.infobox(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.infobox(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}

func TestSurvey(t *testing.T) {
	tests := []struct {
		name   string
		inNode *nodes.SurveyNode
		out    string
	}{
		{
			name: "NoID",
			inNode: nodes.NewSurveyNode("",
				&nodes.SurveyGroup{
					Name:    "pick a number",
					Options: []string{"1", "2", "3"},
				},
				&nodes.SurveyGroup{
					Name:    "choose an answer",
					Options: []string{"yes", "no", "probably"},
				},
			),
			out: `<google-codelab-survey survey-id="">
<h4>pick a number</h4>
<paper-radio-group>
<paper-radio-button>1</paper-radio-button>
<paper-radio-button>2</paper-radio-button>
<paper-radio-button>3</paper-radio-button>
</paper-radio-group>
<h4>choose an answer</h4>
<paper-radio-group>
<paper-radio-button>yes</paper-radio-button>
<paper-radio-button>no</paper-radio-button>
<paper-radio-button>probably</paper-radio-button>
</paper-radio-group>
</google-codelab-survey>`,
		},
		{
			name:   "NoGroups",
			inNode: nodes.NewSurveyNode("foobar"),
			out: `<google-codelab-survey survey-id="foobar">
</google-codelab-survey>`,
		},
		{
			name: "Simple",
			inNode: nodes.NewSurveyNode("a simple example",
				&nodes.SurveyGroup{
					Name:    "pick a color",
					Options: []string{"red", "blue", "yellow"},
				}),
			out: `<google-codelab-survey survey-id="a simple example">
<h4>pick a color</h4>
<paper-radio-group>
<paper-radio-button>red</paper-radio-button>
<paper-radio-button>blue</paper-radio-button>
<paper-radio-button>yellow</paper-radio-button>
</paper-radio-group>
</google-codelab-survey>`,
		},
		{
			name: "MultipleGroups",
			inNode: nodes.NewSurveyNode("an example with multiple survey groups",
				&nodes.SurveyGroup{
					Name:    "a",
					Options: []string{"a", "aa", "aaa"},
				},
				&nodes.SurveyGroup{
					Name:    "b",
					Options: []string{"b", "bb", "bbb"},
				},
				&nodes.SurveyGroup{
					Name:    "c",
					Options: []string{"c", "cc", "ccc"},
				}),
			out: `<google-codelab-survey survey-id="an example with multiple survey groups">
<h4>a</h4>
<paper-radio-group>
<paper-radio-button>a</paper-radio-button>
<paper-radio-button>aa</paper-radio-button>
<paper-radio-button>aaa</paper-radio-button>
</paper-radio-group>
<h4>b</h4>
<paper-radio-group>
<paper-radio-button>b</paper-radio-button>
<paper-radio-button>bb</paper-radio-button>
<paper-radio-button>bbb</paper-radio-button>
</paper-radio-group>
<h4>c</h4>
<paper-radio-group>
<paper-radio-button>c</paper-radio-button>
<paper-radio-button>cc</paper-radio-button>
<paper-radio-button>ccc</paper-radio-button>
</paper-radio-group>
</google-codelab-survey>`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.survey(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.survey(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}

func TestHeader(t *testing.T) {
	a1 := nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foo"})
	a1.Italic = true
	a2 := nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "bar"})
	a3 := nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "baz"})
	a3.Code = true

	tests := []struct {
		name   string
		inNode *nodes.HeaderNode
		out    string
	}{
		{
			name:   "SimpleH1",
			inNode: nodes.NewHeaderNode(1, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foobar"})),
			out:    "<h1 is-upgraded>foobar</h1>",
		},
		{
			name:   "LevelOutOfRange",
			inNode: nodes.NewHeaderNode(100, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "foobar"})),
			out:    "<h100 is-upgraded>foobar</h100>",
		},
		{
			name:   "EmptyContent",
			inNode: nodes.NewHeaderNode(2),
			out:    "<h2 is-upgraded></h2>",
		},
		{
			name:   "StyledText",
			inNode: nodes.NewHeaderNode(3, a1, a2, a3),
			out:    "<h3 is-upgraded><em>foo</em>bar<code>baz</code></h3>",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.header(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.header(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}

func TestYouTube(t *testing.T) {
	tests := []struct {
		name   string
		inNode *nodes.YouTubeNode
		out    string
	}{
		{
			name:   "NonEmpty",
			inNode: nodes.NewYouTubeNode("Mlk888FiI8A"),
			out:    `<iframe class="youtube-video" src="https://www.youtube.com/embed/Mlk888FiI8A?rel=0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`,
		},
		{
			name:   "Empty",
			inNode: nodes.NewYouTubeNode(""),
			out:    `<iframe class="youtube-video" src="https://www.youtube.com/embed/?rel=0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.youtube(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.youtube(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}

func TestIframe(t *testing.T) {
	tests := []struct {
		name   string
		inNode *nodes.IframeNode
		out    string
	}{
		{
			name:   "SomeText",
			inNode: nodes.NewIframeNode("maps.google.com"),
			out:    `<iframe class="embedded-iframe" src="maps.google.com"></iframe>`,
		},
		{
			name:   "Escape",
			inNode: nodes.NewIframeNode("ma ps.google.com"),
			out:    `<iframe class="embedded-iframe" src="ma ps.google.com"></iframe>`,
		},
		{
			name:   "Empty",
			inNode: nodes.NewIframeNode(""),
			out:    `<iframe class="embedded-iframe" src=""></iframe>`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.iframe(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.iframe(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}
