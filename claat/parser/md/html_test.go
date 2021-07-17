package md

import (
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// The utility functions for these tests are purposefully kept very simple to make it easy to understand what the tests are doing.

func makeStrongNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Strong,
		Data:     "strong",
	}
}

func makeBNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.B,
		Data:     "b",
	}
}

func makeEmNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Em,
		Data:     "em",
	}
}

// <i>, not the filesystem abstraction.
func makeINode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.I,
		Data:     "i",
	}
}

// data is the text in the node.
func makeTextNode(data string) *html.Node {
	return &html.Node{
		Type: html.TextNode,
		Data: data,
	}
}

func makeCodeNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Code,
		Data:     "code",
	}
}

func makePNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.P,
		Data:     "p",
	}
}

func makeButtonNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Button,
		Data:     "button",
	}
}

func makeBlinkNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Blink,
		Data:     "blink",
	}
}

func makeAsideNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Aside,
		Data:     "aside",
	}
}

func makeDtNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Dt,
		Data:     "dt",
	}
}

func makeFormNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Form,
		Data:     "form",
	}
}

func makeNameNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Name,
		Data:     "name",
	}
}

func makeInputNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Input,
		Data:     "input",
	}
}

func TestIsHeader(t *testing.T) {
	tests := []struct {
		name string
		in   *html.Node
		out  bool
	}{
		{
			name: "LabTitle",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H1,
				Data:     "h1",
			},
		},
		{
			name: "StepTitle",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H2,
				Data:     "h2",
			},
		},
		{
			name: "FirstLevel",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H3,
				Data:     "h3",
			},
			out: true,
		},
		{
			name: "SecondLevel",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H4,
				Data:     "h4",
			},
			out: true,
		},
		{
			name: "ThirdLevel",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H5,
				Data:     "h5",
			},
			out: true,
		},
		{
			name: "FourthLevel",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H6,
				Data:     "h6",
			},
			out: true,
		},
		{
			name: "NotAHeader",
			in:   makeBlinkNode(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isHeader(tc.in); out != tc.out {
				t.Errorf("isHeader(%v) = %t, want %t", tc.in, out, tc.out)
			}
		})
	}
}

// TODO TestIsMeta

func TestIsBold(t *testing.T) {
	// <strong>foobar</strong>
	a1 := makeStrongNode()
	a2 := makeTextNode("foobar")
	a1.AppendChild(a2)

	// <b>foobar</b>
	b1 := makeBNode()
	b2 := makeTextNode("foobar")
	b1.AppendChild(b2)

	// <strong><code>foobar</code></strong>
	c1 := makeStrongNode()
	c2 := makeCodeNode()
	c3 := makeTextNode("foobar")
	c1.AppendChild(c2)
	c2.AppendChild(c3)

	// <b><code>foobar</code></b>
	d1 := makeBNode()
	d2 := makeCodeNode()
	d3 := makeTextNode("foobar")
	d1.AppendChild(d2)
	d2.AppendChild(d3)

	// <p>foobar</p>
	e1 := makePNode()
	e2 := makeTextNode("foobar")
	e1.AppendChild(e2)

	// <i>foobar</i>
	f1 := makeINode()
	f2 := makeTextNode("foobar")
	f1.AppendChild(f2)

	tests := []struct {
		name string
		in   *html.Node
		out  bool
	}{
		{
			name: "StrongText_Strong",
			in:   a1,
			out:  true,
		},
		{
			name: "StrongText_Strong",
			in:   a2,
			out:  true,
		},
		{
			name: "BText_B",
			in:   b1,
			out:  true,
		},
		{
			name: "BText_Text",
			in:   b2,
			out:  true,
		},
		{
			name: "StrongCodeText_Strong",
			in:   c1,
			out:  true,
		},
		{
			name: "StrongCodeText_Code",
			in:   c2,
			out:  true,
		},
		/*
			// TODO: I think this should work but it doesn't.
			{
				name: "StrongCodeText_Text",
				in:   c3,
				out:  true,
			},
		*/
		{
			name: "BCodeText_B",
			in:   d1,
			out:  true,
		},
		{
			name: "BCodeText_Code",
			in:   d2,
			out:  true,
		},
		/*
			// TODO: I think this should work but it doesn't
			{
				name: "BCodeText_Text",
				in:   d3,
				out:  true,
			},
		*/
		{
			name: "PText_P",
			in:   e1,
		},
		{
			name: "PText_Text",
			in:   e2,
		},
		{
			name: "IText_I",
			in:   f1,
		},
		{
			name: "IText_Text",
			in:   f2,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isBold(tc.in); out != tc.out {
				t.Errorf("isBold(%v) = %t, want %t", tc.in, out, tc.out)
			}
		})
	}
}

func TestIsItalic(t *testing.T) {
	// <em>foobar</em>
	a1 := makeEmNode()
	a2 := makeTextNode("foobar")
	a1.AppendChild(a2)

	// <i>foobar</i>
	b1 := makeINode()
	b2 := makeTextNode("foobar")
	b1.AppendChild(b2)

	// <em><code>foobar</code></em>
	c1 := makeEmNode()
	c2 := makeCodeNode()
	c3 := makeTextNode("foobar")
	c1.AppendChild(c2)
	c2.AppendChild(c3)

	// <i><code>foobar</code></i>
	d1 := makeINode()
	d2 := makeCodeNode()
	d3 := makeTextNode("foobar")
	d1.AppendChild(d2)
	d2.AppendChild(d3)

	// <p>foobar</p>
	e1 := makePNode()
	e2 := makeTextNode("foobar")
	e1.AppendChild(e2)

	// <b>foobar</b>
	f1 := makeBNode()
	f2 := makeTextNode("foobar")
	f1.AppendChild(f2)

	tests := []struct {
		name string
		in   *html.Node
		out  bool
	}{
		{
			name: "EmText_Em",
			in:   a1,
			out:  true,
		},
		{
			name: "EmText_Em",
			in:   a2,
			out:  true,
		},
		{
			name: "IText_I",
			in:   b1,
			out:  true,
		},
		{
			name: "IText_Text",
			in:   b2,
			out:  true,
		},
		{
			name: "EmCodeText_Em",
			in:   c1,
			out:  true,
		},
		{
			name: "EmCodeText_Code",
			in:   c2,
			out:  true,
		},
		/*
			// TODO: I think this should work but it doesn't.
			{
				name: "EmCodeText_Text",
				in:   c3,
				out:  true,
			},
		*/
		{
			name: "ICodeText_I",
			in:   d1,
			out:  true,
		},
		{
			name: "ICodeText_Code",
			in:   d2,
			out:  true,
		},
		/*
			// TODO: I think this should work but it doesn't
			{
				name: "ICodeText_Text",
				in:   d3,
				out:  true,
			},
		*/
		{
			name: "PText_P",
			in:   e1,
		},
		{
			name: "PText_Text",
			in:   e2,
		},
		{
			name: "BText_B",
			in:   f1,
		},
		{
			name: "BText_Text",
			in:   f2,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isItalic(tc.in); out != tc.out {
				t.Errorf("isItalic(%v) = %t, want %t", tc.in, out, tc.out)
			}
		})
	}
}

func TestIsBoldAndItalic(t *testing.T) {
	// <em><strong>foobar</strong></em>
	a1 := makeEmNode()
	a2 := makeStrongNode()
	a3 := makeTextNode("foobar")
	a1.AppendChild(a2)
	a2.AppendChild(a3)

	// <i><strong>foobar</strong></i>
	b1 := makeINode()
	b2 := makeStrongNode()
	b3 := makeTextNode("foobar")
	b1.AppendChild(b2)
	b2.AppendChild(b3)

	// <em><b>foobar</b></em>
	c1 := makeEmNode()
	c2 := makeBNode()
	c3 := makeTextNode("foobar")
	c1.AppendChild(c2)
	c2.AppendChild(c3)

	// <i><b>foobar</b></i>
	d1 := makeINode()
	d2 := makeBNode()
	d3 := makeTextNode("foobar")
	d1.AppendChild(d2)
	d2.AppendChild(d3)

	// <em><strong><code>foobar</code></strong></em>
	e1 := makeEmNode()
	e2 := makeStrongNode()
	e3 := makeCodeNode()
	e4 := makeTextNode("foobar")
	e1.AppendChild(e2)
	e2.AppendChild(e3)
	e3.AppendChild(e4)

	// <em><code><strong>foobar</strong></code></em>
	f1 := makeEmNode()
	f2 := makeCodeNode()
	f3 := makeStrongNode()
	f4 := makeTextNode("foobar")
	f1.AppendChild(f2)
	f2.AppendChild(f3)
	f3.AppendChild(f4)

	// <strong><em>foobar</em></strong>
	g1 := makeStrongNode()
	g2 := makeEmNode()
	g3 := makeTextNode("foobar")
	g1.AppendChild(g2)
	g2.AppendChild(g3)

	// <strong><i>foobar</i></strong>
	h1 := makeStrongNode()
	h2 := makeINode()
	h3 := makeTextNode("foobar")
	h1.AppendChild(h2)
	h2.AppendChild(h3)

	// <b><em>foobar</em></b>
	// Skipped i and j due to widespread use of <i>
	k1 := makeBNode()
	k2 := makeEmNode()
	k3 := makeTextNode("foobar")
	k1.AppendChild(k2)
	k2.AppendChild(k3)

	// <b><i>foobar</i></b>
	l1 := makeBNode()
	l2 := makeINode()
	l3 := makeTextNode("foobar")
	l1.AppendChild(l2)
	l2.AppendChild(l3)

	// <strong><em><code>foobar</code></em></strong>
	m1 := makeStrongNode()
	m2 := makeEmNode()
	m3 := makeCodeNode()
	m4 := makeTextNode("foobar")
	m1.AppendChild(m2)
	m2.AppendChild(m3)
	m3.AppendChild(m4)

	// <strong><code><em>foobar</em></code></strong>
	n1 := makeStrongNode()
	n2 := makeCodeNode()
	n3 := makeEmNode()
	n4 := makeTextNode("foobar")
	n1.AppendChild(n2)
	n2.AppendChild(n3)
	n3.AppendChild(n4)

	// <p>foobar</p>
	o1 := makePNode()
	o2 := makeTextNode("foobar")
	o1.AppendChild(o2)

	// <em>foobar</em>
	p1 := makeEmNode()
	p2 := makeTextNode("foobar")
	p1.AppendChild(p2)

	// <strong>foobar</strong>
	q1 := makeStrongNode()
	q2 := makeTextNode("foobar")
	q1.AppendChild(q2)

	tests := []struct {
		name string
		in   *html.Node
		out  bool
	}{
		/*
			// TODO: I think this should work but it doesn't
			// Specifically, without loss of generality, passing <em> in <em>foobar</em> returns true, but this behaves differently
			{
				name: "EmStrongText_Strong",
				in:   a2,
				out:  true,
			},
		*/
		{
			name: "EmStrongText_Text",
			in:   a3,
			out:  true,
		},
		/*
			// TODO: I think this should work but it doesn't
			{
				name: "IStrongText_Strong",
				in:   b2,
				out:  true,
			},
		*/
		{
			name: "IStrongText_Text",
			in:   b3,
			out:  true,
		},
		/*
			// TODO: I think this should work but it doesn't
			{
				name: "EmBText_B",
				in:   c2,
				out:  true,
			},
		*/
		{
			name: "EmBText_Text",
			in:   c3,
			out:  true,
		},
		/*
			// TODO: I think this should work but it doesn't
			{
				name: "IBText_B",
				in:   d2,
				out:  true,
			},
		*/
		{
			name: "IBText_Text",
			in:   d3,
			out:  true,
		},
		/*
			// TODO: I (maybe) think this should work but it doesn't
			{
				name: "EmStrongCodeText_Strong",
				in:   e2,
				out:  true,
			},
		*/
		{
			name: "EmStrongCodeText_Code",
			in:   e3,
			out:  true,
		},
		{
			name: "EmStrongCodeText_Text",
			in:   e4,
			out:  true,
		},
		/*
			// TODO: I (maybe) think this should work but it doesn't
			{
				name: "EmCodeStrongText_Code",
				in:   f2,
				out:  true,
			},
		*/
		{
			name: "EmCodeStrongText_Strong",
			in:   f3,
			out:  true,
		},
		{
			name: "EmCodeStrongText_Text",
			in:   f4,
			out:  true,
		},
		/*
			// TODO: I (maybe) think this should work but it doesn't
			{
				name: "StrongEmText_Em",
				in:   g2,
				out:  true,
			},
		*/
		{
			name: "StrongEmText_Text",
			in:   g3,
			out:  true,
		},
		/*
			// TODO: I (maybe) think this should work but it doesn't
			{
				name: "StrongIText_I",
				in:   h2,
				out:  true,
			},
		*/
		{
			name: "strongIText_Text",
			in:   h3,
			out:  true,
		},
		/*
			// TODO: I (maybe) think this should work but it doesn't
			{
				name: "BEmText_Em",
				in:   k2,
				out:  true,
			},
		*/
		{
			name: "BEmText_Text",
			in:   k3,
			out:  true,
		},
		/*
			// TODO: I (maybe) think this should work but it doesn't
			{
				name: "BIText_I",
				in:   l2,
				out:  true,
			},
		*/
		{
			name: "BIText_Text",
			in:   l3,
			out:  true,
		},
		/*
			// TODO: I (maybe) think this should work but it doesn't
			{
				name: "StrongEmCodeText_Em",
				in:   m2,
				out:  true,
			},
		*/
		{
			name: "StrongEmCodeText_Code",
			in:   m3,
			out:  true,
		},
		{
			name: "StrongEmCodeText_Text",
			in:   m4,
			out:  true,
		},
		/*
			// TODO: I (maybe) think this should work but it doesn't
			{
				name: "StrongCodeEmText_Code",
				in:   m2,
				out:  true,
			},
		*/
		{
			name: "StrongCodeEmText_Em",
			in:   n3,
			out:  true,
		},
		{
			name: "StrongCodeEmText_Text",
			in:   n4,
			out:  true,
		},
		{
			name: "PText_P",
			in:   o1,
		},
		{
			name: "PText_Text",
			in:   o2,
		},
		{
			name: "EmText_Em",
			in:   p1,
		},
		{
			name: "EmText_Text",
			in:   p2,
		},
		{
			name: "StrongText_Strong",
			in:   q1,
		},
		{
			name: "StrongText_Text",
			in:   q2,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isBoldAndItalic(tc.in); out != tc.out {
				t.Errorf("isBoldAndItalic(%v) = %t, want %t", tc.in, out, tc.out)
			}
		})
	}
}

func TestIsConsole(t *testing.T) {
	// <code class="language-console">foobar</code>
	a1 := makeCodeNode()
	a1.Attr = append(a1.Attr, html.Attribute{Key: "class", Val: "language-console"})
	a2 := makeTextNode("foobar")
	a1.AppendChild(a2)

	// <code class="language-js">foobar</code>
	b1 := makeCodeNode()
	b1.Attr = append(b1.Attr, html.Attribute{Key: "class", Val: "language-js"})
	b2 := makeTextNode("foobar")
	b1.AppendChild(b2)

	// <code>foobar</code>
	c1 := makeCodeNode()
	c2 := makeTextNode("foobar")
	c1.AppendChild(c2)

	// <p>foobar</p>
	d1 := makePNode()
	d2 := makeTextNode("foobar")
	d1.AppendChild(d2)

	tests := []struct {
		name string
		in   *html.Node
		out  bool
	}{
		{
			name: "ConsoleText_Console",
			in:   a1,
			out:  true,
		},
		{
			name: "ConsoleText_Text",
			in:   a2,
			out:  true,
		},
		{
			name: "JavascriptText_Javascript",
			in:   b1,
		},
		{
			name: "JavascriptText_Text",
			in:   b2,
		},
		{
			name: "CodeText_Code",
			in:   c1,
		},
		{
			name: "CodeText_Text",
			in:   c2,
		},
		{
			name: "PText_P",
			in:   d1,
		},
		{
			name: "PText_Text",
			in:   d2,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isConsole(tc.in); out != tc.out {
				t.Errorf("isConsole(%v) = %t, want %t", tc.in, out, tc.out)
			}
		})
	}
}

func TestIsCode(t *testing.T) {
	// <code class="language-console">foobar</code>
	a1 := makeCodeNode()
	a1.Attr = append(a1.Attr, html.Attribute{Key: "class", Val: "language-console"})
	a2 := makeTextNode("foobar")
	a1.AppendChild(a2)

	// <code class="language-js">foobar</code>
	b1 := makeCodeNode()
	b1.Attr = append(b1.Attr, html.Attribute{Key: "class", Val: "language-js"})
	b2 := makeTextNode("foobar")
	b1.AppendChild(b2)

	// <code>foobar</code>
	c1 := makeCodeNode()
	c2 := makeTextNode("foobar")
	c1.AppendChild(c2)

	// <p>foobar</p>
	d1 := makePNode()
	d2 := makeTextNode("foobar")
	d1.AppendChild(d2)

	tests := []struct {
		name string
		in   *html.Node
		out  bool
	}{
		{
			name: "ConsoleText_Console",
			in:   a1,
		},
		{
			name: "ConsoleText_Text",
			in:   a2,
		},
		{
			name: "JavascriptText_Javascript",
			in:   b1,
			out:  true,
		},
		{
			name: "JavascriptText_Text",
			in:   b2,
			out:  true,
		},
		{
			name: "CodeText_Code",
			in:   c1,
			out:  true,
		},
		{
			name: "CodeText_Text",
			in:   c2,
			out:  true,
		},
		{
			name: "PText_P",
			in:   d1,
		},
		{
			name: "PText_Text",
			in:   d2,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isCode(tc.in); out != tc.out {
				t.Errorf("isCode(%v) = %t, want %t", tc.in, out, tc.out)
			}
		})
	}
}

func TestIsButton(t *testing.T) {
	a1 := makeButtonNode()
	a2 := makeTextNode("foobar")
	a1.AppendChild(a2)

	tests := []struct {
		name string
		in   *html.Node
		out  bool
	}{
		{
			name: "Button",
			in:   makeButtonNode(),
			out:  true,
		},
		{
			name: "ButtonWithText",
			in:   a1,
			out:  true,
		},
		{
			name: "TextInButton",
			in:   a2,
		},
		{
			name: "NotAButton",
			in:   makeBlinkNode(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isButton(tc.in); out != tc.out {
				t.Errorf("isButton(%v) = %t, want %t", tc.in, out, tc.out)
			}
		})
	}
}

func TestIsAside(t *testing.T) {
	a1 := makeAsideNode()
	a2 := makeTextNode("foobar")
	a1.AppendChild(a2)

	tests := []struct {
		name string
		in   *html.Node
		out  bool
	}{
		{
			name: "Aside",
			in:   makeAsideNode(),
			out:  true,
		},
		{
			name: "AsideWithText",
			in:   a1,
			out:  true,
		},
		{
			name: "TextInAside",
			in:   a2,
		},
		{
			name: "NotAnAside",
			in:   makeBlinkNode(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isAside(tc.in); out != tc.out {
				t.Errorf("isAside(%v) = %t, want %t", tc.in, out, tc.out)
			}
		})
	}
}

// TODO: test isNewAside

func TestIsInfobox(t *testing.T) {
	a1 := makeDtNode()
	a2 := makeTextNode("positive")
	a3 := makeTextNode("foobar")
	// The text nodes should be siblings.
	a1.AppendChild(a2)
	a1.AppendChild(a3)

	b1 := makeDtNode()
	b2 := makeTextNode("negative")
	b3 := makeTextNode("foobar")
	// The text nodes should be siblings.
	b1.AppendChild(b2)
	b1.AppendChild(b3)

	c1 := makeDtNode()
	c2 := makeTextNode("foobar")
	c1.AppendChild(c2)

	tests := []struct {
		name string
		in   *html.Node
		out  bool
	}{
		{
			name: "InfoboxPositive",
			in:   a1,
			out:  true,
		},
		{
			name: "TextInInfoboxPositive",
			in:   a3,
		},
		{
			name: "InfoboxNegative",
			in:   b1,
			out:  true,
		},
		{
			name: "TextInInfoboxNegative",
			in:   b3,
		},
		{
			name: "NotAnInfobox",
			in:   makeBlinkNode(),
		},
		// TODO: Is this how this function should work?
		{
			name: "InfoboxNoKind",
			in:   c1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isInfobox(tc.in); out != tc.out {
				t.Errorf("isInfobox(%v) = %t, want %t", tc.in, out, tc.out)
			}
		})
	}
}

func TestIsInfoboxNegative(t *testing.T) {
	a1 := makeDtNode()
	a2 := makeTextNode("positive")
	a3 := makeTextNode("foobar")
	// The text nodes should be siblings.
	a1.AppendChild(a2)
	a1.AppendChild(a3)

	b1 := makeDtNode()
	b2 := makeTextNode("negative")
	b3 := makeTextNode("foobar")
	// The text nodes should be siblings.
	b1.AppendChild(b2)
	b1.AppendChild(b3)

	tests := []struct {
		name string
		in   *html.Node
		out  bool
	}{
		{
			name: "InfoboxPositive",
			in:   a1,
		},
		{
			name: "TextInInfoboxPositive",
			in:   a3,
		},
		{
			name: "InfoboxNegative",
			in:   b1,
			out:  true,
		},
		{
			name: "TextInInfoboxNegative",
			in:   b3,
		},
		{
			name: "NotAnInfobox",
			in:   makeBlinkNode(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isInfoboxNegative(tc.in); out != tc.out {
				t.Errorf("isInfoboxNegative(%v) = %t, want %t", tc.in, out, tc.out)
			}
		})
	}
}

func TestIsSurvey(t *testing.T) {
	a1 := makeFormNode()
	a2 := makeNameNode()
	a3 := makeInputNode()
	// The name and input nodes should be siblings.
	a1.AppendChild(a2)
	a1.AppendChild(a3)

	b1 := makeFormNode()
	b2 := makeInputNode()
	b1.AppendChild(b2)

	c1 := makeFormNode()
	c2 := makeNameNode()
	c1.AppendChild(c2)

	tests := []struct {
		name string
		in   *html.Node
		out  bool
	}{
		{
			name: "ValidSurvey",
			in:   a1,
			out:  true,
		},
		{
			name: "NoName",
			in:   b1,
		},
		{
			name: "NoInputs",
			in:   c1,
		},
		{
			name: "NotAForm",
			in:   makeBlinkNode(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isSurvey(tc.in); out != tc.out {
				t.Errorf("isSurvey(%v) = %t, want %t", tc.in, out, tc.out)
			}
		})
	}
}
