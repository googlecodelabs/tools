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

package gdoc

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/html"

	"github.com/googlecodelabs/tools/claat/nodes"
	"github.com/googlecodelabs/tools/claat/parser"
	"github.com/googlecodelabs/tools/claat/render"
	"github.com/googlecodelabs/tools/claat/types"
)

func trimMarkupSpace(s string) string {
	var buf bytes.Buffer
	for _, l := range strings.Split(s, "\n") {
		buf.WriteString(strings.TrimSpace(l))
	}
	return buf.String()
}

func markupReader(s string) io.Reader {
	s = trimMarkupSpace(s)
	return strings.NewReader(s)
}

func TestParseStepDuration(t *testing.T) {
	tests := []struct {
		markup string
		dur    time.Duration
	}{
		{`<p><span class="c9">Duration: 1:30</span></p>`, 2 * time.Minute},
		{`<p><span class="c9">Duration: 1:30 </span></p>`, 2 * time.Minute},
		{`<p><span class="c9">Duration: 1:30</span> </p>`, 2 * time.Minute},
		{`<p><span class="c9">Duration : 5</span></p>`, 5 * time.Minute},
		{`<p><span class="c9">duration: 1</span></p>`, time.Minute},
		{`<p><span class="c9">Duration:</span><span class="c9">&nbsp;66:</span><span class="c9">1</span><span class="c9">0</span></p>`, 67 * time.Minute},
	}
	for i, test := range tests {
		doc, err := html.Parse(strings.NewReader(test.markup))
		if err != nil {
			t.Errorf("%d: Parse(%q): %v", i, test.markup, err)
		}
		ds := &docState{
			step: &types.Step{Content: nodes.NewListNode()},
			css:  cssStyle{".c9": {"color": metaColor}},
			cur:  doc.FirstChild,
		}
		parseTop(ds)
		if ds.step.Duration != test.dur {
			t.Errorf("%d: ds.step.Duration = %v; want %v", i, ds.step.Duration, test.dur)
		}
	}
}

func TestParseTopCodeBlock(t *testing.T) {
	const markup = `
	<table cellpadding="0" cellspacing="0"><tbody><tr>
	<td colspan="1" rowspan="1">
		<p><span class="code">start func() {<br>}</span><span class="code"></span></p>
		<p><span class="code"></span></p>
		<p><span class="code">func2() {<br>}</span><span class="code">&nbsp;// comment</span></p>
	</td>
	</tr></tbody></table>

	<table cellpadding="0" cellspacing="0"><tbody><tr>
	<td colspan="1" rowspan="1">
		<p><span class="term">adb shell am start -a VIEW \</span></p>
		<p><span class="term">-d &quot;http://host&quot; app</span></p>
	</td>
	</tr></tbody></table>
	`

	code := "start func() {\n}\n\nfunc2() {\n} // comment"
	term := "adb shell am start -a VIEW \\\n-d \"http://host\" app"
	content := nodes.NewListNode()
	var lang string
	content.Append(nodes.NewCodeNode(code, false, lang))
	content.Append(nodes.NewCodeNode(term, true, lang))

	doc, err := html.Parse(markupReader(markup))
	if err != nil {
		t.Fatal(err)
	}
	ds := &docState{
		step: &types.Step{Content: nodes.NewListNode()},
		css: cssStyle{
			".code": {"font-family": fontCode},
			".term": {"font-family": fontConsole},
		},
		cur: doc.FirstChild,
	}
	parseTop(ds)
	var ctx render.Context
	html1, _ := render.HTML(ctx, ds.step.Content)
	html2, _ := render.HTML(ctx, content)
	s1 := strings.TrimSpace(string(html1))
	s2 := strings.TrimSpace(string(html2))
	if s1 != s2 {
		t.Errorf("step.Content:\n\n%s\nwant:\n\n%s", s1, s2)
	}
}

func TestMetaTable(t *testing.T) {
	const markup = `
	<html>
	<body>
		<table>
			<tr>
				<td>Summary</td>
				<td>Test summary</td>
			</tr>
			<tr>
				<td>Authors</td>
				<td>John Smith &lt;user@example.com&gt;</td>
			</tr>
			<tr>
				<td>Category</td>
				<td>Foo, Bar</td>
			</tr>
			<tr>
				<td>Environment</td>
				<td>Web, Kiosk</td>
			</tr>
			<tr>
				<td>Status</td>
				<td>Final</td>
			</tr>
			<tr>
				<td>Feedback</td>
				<td>https://example.com/issues</td>
			</tr>
			<tr>
				<td>Analytics</td>
				<td>GA-12345</td>
			</tr>
		</table>
	</body>
	</html>
	`

	p := &Parser{}
	clab, err := p.Parse(markupReader(markup), *parser.NewOptions())
	if err != nil {
		t.Fatal(err)
	}
	meta := types.Meta{
		Summary:    "Test summary",
		Authors:    "John Smith <user@example.com>",
		Categories: []string{"foo", "bar"},
		Theme:      "foo",
		Status:     clab.Meta.Status, // verified separately
		Feedback:   "https://example.com/issues",
		GA:         "GA-12345",
		// Tags are always sorted.
		// TODO: move sorting to Parse of the parser package
		Tags:  []string{"kiosk", "web"},
		Extra: map[string]string{},
	}
	if !reflect.DeepEqual(clab.Meta, meta) {
		t.Errorf("Meta: \n%+v\nwant:\n%+v", clab.Meta, meta)
	}
	status := types.LegacyStatus([]string{"final"})
	if clab.Meta.Status == nil {
		t.Fatalf("Meta.Status is nil; want %q", status)
	}
	if !reflect.DeepEqual(clab.Meta.Status, &status) {
		t.Errorf("Meta.Status: %q; want %q", *clab.Meta.Status, status)
	}
}

func TestMetaTablePassMetadata(t *testing.T) {
	const markup = `
	<html>
	<body>
		<table>
			<tr>
				<td>Summary</td>
				<td>Test summary</td>
			</tr>
			<tr>
				<td>Authors</td>
				<td>John Smith &lt;user@example.com&gt;</td>
			</tr>
			<tr>
				<td>Category</td>
				<td>Foo, Bar</td>
			</tr>
			<tr>
				<td>Environment</td>
				<td>Web, Kiosk</td>
			</tr>
			<tr>
				<td>Status</td>
				<td>Final</td>
			</tr>
			<tr>
				<td>Feedback Link</td>
				<td>https://example.com/issues</td>
			</tr>
			<tr>
				<td>Analytics</td>
				<td>GA-12345</td>
			</tr>
			<tr>
				<td>ExtraFieldOne</td>
				<td>11111</td>
			</tr>
			<tr>
				<td>ExtraFieldTwo</td>
				<td>22222</td>
			</tr>
		</table>
	</body>
	</html>
	`

	p := &Parser{}
	opts := *parser.NewOptions()
	opts.PassMetadata = map[string]bool{
		"extra_field_one": true,
	}

	clab, err := p.Parse(markupReader(markup), opts)
	if err != nil {
		t.Fatal(err)
	}
	meta := types.Meta{
		Summary:    "Test summary",
		Authors:    "John Smith <user@example.com>",
		Categories: []string{"foo", "bar"},
		Theme:      "foo",
		Status:     clab.Meta.Status, // verified separately
		Feedback:   "https://example.com/issues",
		GA:         "GA-12345",
		// Tags are always sorted.
		// TODO: move sorting to Parse of the parser package
		Tags: []string{"kiosk", "web"},
		Extra: map[string]string{
			"extra_field_one": "11111",
		},
	}
	if !reflect.DeepEqual(clab.Meta, meta) {
		t.Errorf("Meta: \n%+v\nwant:\n%+v", clab.Meta, meta)
	}
	status := types.LegacyStatus([]string{"final"})
	if clab.Meta.Status == nil {
		t.Fatalf("Meta.Status is nil; want %q", status)
	}
	if !reflect.DeepEqual(clab.Meta.Status, &status) {
		t.Errorf("Meta.Status: %q; want %q", *clab.Meta.Status, status)
	}
}

func TestParseDoc(t *testing.T) {
	const markup = `
	<html><head><style>
		.meta { color: #b7b7b7 }
		.code { font-family: "Courier New" }
		.term { font-family: "Consolas" }
		.btn { background-color: #6aa84f }
		.bold { font-weight: bold }
		.ita { font-style: italic }
		.nibox { background-color: #fce5cd }
		.survey { background-color: #cfe2f3 }
		.comment { border: 1px solid black }
	</style></head>
	<body>
		<p class="title"><a name="a1"></a><span>Test Codelab</span></p>

		<p>this should be ignored</p>

		<h1><a name="a2"></a><span>Overview</span></h1>
		<p><span class="meta">Duration: 1:00</span></p>

		<p><span>[[</span><span class="bold">import</span><span>&nbsp;</span><span><a href="https://example.com/import">shared</a></span><span>]]</span></p>

		<img src="https://host/image.png" alt="alt text" title="title text">
		<p><img src="https://host/small.png" style="height: 10px; width: 25.5px"> icon.</p>

		<p><img alt="https://www.youtube.com/watch?v=vid" src="https://yt.com/vid.jpg"></p>
		<p><img alt="https://repl.it/?foo=bar" src="https://host/image.png"></p>
		<p><img alt="https://example.com/?foo=bar" src="https://host/image.png"></p>

		<h3><a name="a3"></a><span>What you&rsquo;ll learn</span></h3>
		<ul class="start">
		<li><span>First </span><span>One</span><sup><a href="#cmnt1" name="cmnt_ref1" target="_blank">[a]</a></sup></li>
		<li><span>Two </span><span><a href="https://google.com/url?q=http%3A%2F%2Fexample.com">Link</a></span></li>
		</ul>
		<ul><li><span>Three</span></li></ul>

		<p>This is<span class="code"> code</span>.</p>
		<p>Just <span>a</span> paragraph.</p>
		<p><a href="url">one</a><a href="url"> url</a></p>
		<p><span class="btn"><a href="http://example.com">Download Zip</a></span></p>
		<p>
			<span class="bold">Bo</span><span>&nbsp;</span><span class="bold">ld</span>
			<span class="ita"> italic</span> text <span class="bold ita">or both.</span></p>

		<h3><a href="http://host/file.java">a file</a></h3>
		<table cellpadding="0" cellspacing="0"><tbody><tr>
		<td colspan="1" rowspan="1">
			<p><span class="code">start func() {<br>}</span></p>
			<p><span class="code"></span></p>
			<p><span class="code">func2() {<br>}</span><span class="code">&nbsp;// comment</span></p>
		</td>
		</tr></tbody></table>

		<table cellpadding="0" cellspacing="0"><tbody><tr>
		<td colspan="1" rowspan="1">
			<p><span class="term">adb shell am start -a VIEW \</span></p>
			<p><span style="font-family:&quot;Consolas&quot;;">-d &quot;http://host&quot; app</span></p>
		</td>
		</tr></tbody></table>

		<table cellpadding="0" cellspacing="0"><tbody><tr>
		<td class="nibox" colspan="1" rowspan="1">
			<p><span class="bold">warning</span></p>
			<p><span>negative box.</span></p>
		</td>
		</tr></tbody></table>

		<table cellpadding="0" cellspacing="0"><tbody><tr>
		<td class="survey" colspan="1" rowspan="1">
		<h4><a name="x"></a><span class="code">How</span><span class="ita">&nbsp;will you use it?</span></h4>
		<ul><li class="bold"><span class="c5">Read it</span></li></ul>
		<ul><li class="c23 c47"><span class="c5">Read and complete</span></li></ul>
		<p class="c23 c44"><span class="c5"></span></p>
		<h4><a name="asd"></a><span>How</span><span>&nbsp;would you rate?</span></h4>
		<ul>
			<li class="c19 c47"><span class="c5">Novice</span></li>
			<li class="c19 c47"><span class="c5">Intermediate</span></li>
			<li class="c19 c47"><span class="c5">Proficient</span></li>
		</ul>
		<p class="c23 c44"><span class="c5"></span></p>
		</td>
		</tr></tbody></table>
		<div class="comment">
		<p><a href="#cmnt_ref1" name="cmnt1">[a]</a><span class="c16 c8">Test comment.</span></p>
		</div>
	</body>
	</html>
	`

	p := &Parser{}
	c, err := p.Parse(markupReader(markup), *parser.NewOptions())
	if err != nil {
		t.Fatal(err)
	}
	if c.Meta.Title != "Test Codelab" {
		t.Errorf("c.Meta.Title = %q; want Test Codelab", c.Meta.Title)
	}
	if c.Meta.ID != "test-codelab" {
		t.Errorf("c.ID = %q; want test-codelab", c.Meta.ID)
	}
	if len(c.Steps) == 0 {
		t.Fatalf("len(c.Steps) = 0")
	}
	step := c.Steps[0]
	if step.Title != "Overview" {
		t.Errorf("step.Title = %q; want Overview", step.Title)
	}
	if len(step.Content.Nodes) == 0 {
		t.Fatal("step.Content.Nodes is empty")
	}
	want := "https://example.com/import"
	in, ok := step.Content.Nodes[0].(*nodes.ImportNode)
	if !ok {
		t.Errorf("step.Content.Nodes[0] = %+v; want nodes.ImportNode", step.Content.Nodes[0])
	}
	if ok && in.URL != want {
		t.Errorf("in.URL = %q; want %q", in.URL, want)
	}
	if in.Block() != nil {
		t.Errorf("in.Block = %+v (%T); want nil", in.Block(), in.Block())
	}

	content := nodes.NewListNode()

	img := nodes.NewImageNode(nodes.NewImageNodeOptions{
		Src:   "https://host/image.png",
		Alt:   "alt text",
		Title: "title text",
	})
	para := nodes.NewListNode(img)
	para.MutateBlock(true)
	content.Append(para)

	img = nodes.NewImageNode(nodes.NewImageNodeOptions{
		Src:   "https://host/small.png",
		Width: 25.5,
	})
	para = nodes.NewListNode(img, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: " icon."}))
	para.MutateBlock(true)
	content.Append(para)

	yt := nodes.NewYouTubeNode("vid")
	yt.MutateBlock(true)
	content.Append(yt)

	iframe := nodes.NewIframeNode("https://repl.it/?foo=bar")
	iframe.MutateBlock(true)
	content.Append(iframe)

	img = nodes.NewImageNode(nodes.NewImageNodeOptions{
		Src: "https://host/image.png",
		Alt: "The domain of the requested iframe (example.com) has not been whitelisted.",
	})
	para = nodes.NewListNode(img)
	para.MutateBlock(true)
	content.Append(para)

	h := nodes.NewHeaderNode(3, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "What you'll learn"}))
	h.MutateType(nodes.NodeHeaderCheck)
	content.Append(h)
	list := nodes.NewItemsListNode("", 0)
	list.MutateType(nodes.NodeItemsCheck)
	list.NewItem().Append(nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "First One"}))
	item := list.NewItem()
	item.Append(nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "Two "}))
	item.Append(nodes.NewURLNode("http://example.com", nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "Link"})))
	list.NewItem().Append(nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "Three"}))
	content.Append(list)

	para = nodes.NewListNode()
	para.MutateBlock(true)
	para.Append(nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "This is "}))
	txt := nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "code"})
	txt.Code = true
	para.Append(txt)
	para.Append(nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "."}))
	content.Append(para)

	para = nodes.NewListNode()
	para.MutateBlock(true)
	para.Append(nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "Just a paragraph."}))
	content.Append(para)

	u := nodes.NewURLNode("url", nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "one url"}))
	para = nodes.NewListNode(u)
	para.MutateBlock(true)
	content.Append(para)

	btn := nodes.NewButtonNode(true, true, true, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "Download Zip"}))
	dl := nodes.NewURLNode("http://example.com", btn)
	para = nodes.NewListNode(dl)
	para.MutateBlock(true)
	content.Append(para)

	b := nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "Bo ld"})
	b.Bold = true
	i := nodes.NewTextNode(nodes.NewTextNodeOptions{Value: " italic"})
	i.Italic = true
	bi := nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "or both."})
	bi.Bold = true
	bi.Italic = true
	para = nodes.NewListNode(b, i, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: " text "}), bi)
	para.MutateBlock(true)
	content.Append(para)

	h = nodes.NewHeaderNode(3, nodes.NewURLNode(
		"http://host/file.java", nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "a file"})))
	content.Append(h)

	var lang string
	code := "start func() {\n}\n\nfunc2() {\n} // comment"
	cn := nodes.NewCodeNode(code, false, lang)
	cn.MutateBlock(1)
	content.Append(cn)

	term := "adb shell am start -a VIEW \\\n-d \"http://host\" app"
	tn := nodes.NewCodeNode(term, true, lang)
	tn.MutateBlock(2)
	content.Append(tn)

	b = nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "warning"})
	b.Bold = true
	n1 := nodes.NewListNode(b)
	n1.MutateBlock(true)
	n2 := nodes.NewListNode(nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "negative box."}))
	n2.MutateBlock(true)
	box := nodes.NewInfoboxNode(nodes.InfoboxNegative, n1, n2)
	content.Append(box)

	sv := nodes.NewSurveyNode("test-codelab-1")
	sv.Groups = append(sv.Groups, &nodes.SurveyGroup{
		Name:    "How will you use it?",
		Options: []string{"Read it", "Read and complete"},
	})
	sv.Groups = append(sv.Groups, &nodes.SurveyGroup{
		Name:    "How would you rate?",
		Options: []string{"Novice", "Intermediate", "Proficient"},
	})
	content.Append(sv)

	var ctx render.Context
	html1, _ := render.HTML(ctx, step.Content)
	html2, _ := render.HTML(ctx, content)
	if html1 != html2 {
		t.Errorf("step.Content:\n\n%s\nwant:\n\n%s", html1, html2)
	}
}

func TestParseFragment(t *testing.T) {
	const markup = `
	<html><head><style>
		.meta { color: #b7b7b7 }
		.code { font-family: "Courier New" }
		.term { font-family: "Consolas" }
		.btn { background-color: #6aa84f }
		.bold { font-weight: bold }
		.ita { font-style: italic }
		.nibox { background-color: #fce5cd }
		.survey { background-color: #cfe2f3 }
		.comment { border: 1px solid black }
	</style></head>
	<body>
		<p class="title"><a name="a1"></a><span>Test Codelab</span></p>
		<p>this should not be ignored</p>
		<p><img src="https://host/image.png"></p>
		<span class="c17 c7"><a class="c11" href="https://www.google.com/url?q=https://www.example.com/%2B/test;l%3D68&amp;sa=D">Test redirector.</a></span>
		<div class="comment">
		<p><a href="#cmnt_ref1" name="cmnt1">[a]</a><span class="c16 c8">Test comment.</span></p>
		</div>
	</body>
	</html>
	`

	p := &Parser{}
	opts := *parser.NewOptions()
	fragmentNodes, err := p.ParseFragment(markupReader(markup), opts)
	if err != nil {
		t.Fatal(err)
	}

	var want []nodes.Node

	para := nodes.NewListNode()
	para.MutateBlock(true)
	para.Append(nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "Test Codelab"}))
	want = append(want, para)

	para = nodes.NewListNode()
	para.MutateBlock(true)
	para.Append(nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "this should not be ignored"}))
	want = append(want, para)

	img := nodes.NewImageNode(nodes.NewImageNodeOptions{Src: "https://host/image.png"})
	para = nodes.NewListNode(img)
	para.MutateBlock(true)
	want = append(want, para)

	tn := nodes.NewTextNode(nodes.NewTextNodeOptions{
		Value: "Test redirector.",
	})
	rlink := nodes.NewURLNode("https://www.example.com/+/test;l=68&sa=D", tn)
	para = nodes.NewListNode(rlink)
	para.MutateBlock(true)
	want = append(want, para)

	var ctx render.Context
	html1, _ := render.HTML(ctx, fragmentNodes...)
	html2, _ := render.HTML(ctx, want...)
	if html1 != html2 {
		t.Errorf("nodes:\n\n%s\nwant:\n\n%s", html1, html2)
	}
}
