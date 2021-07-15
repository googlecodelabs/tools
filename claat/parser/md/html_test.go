package md

import (
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

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
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.Blink,
				Data:     "blink",
			},
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

func makeTextNode() *html.Node {
	return &html.Node{
		Type: html.TextNode,
		Data: "foobar",
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

func TestIsBold(t *testing.T) {
	// <strong>foobar</strong>
	a1 := makeStrongNode()
	a2 := makeTextNode()
	a1.AppendChild(a2)

	// <b>foobar</b>
	b1 := makeBNode()
	b2 := makeTextNode()
	b1.AppendChild(b2)

	// <strong><code>foobar</code></strong>
	c1 := makeStrongNode()
	c2 := makeCodeNode()
	c3 := makeTextNode()
	c1.AppendChild(c2)
	c2.AppendChild(c3)

	// <b><code>foobar</code></b>
	d1 := makeBNode()
	d2 := makeCodeNode()
	d3 := makeTextNode()
	d1.AppendChild(d2)
	d2.AppendChild(d3)

	// <p>foobar</p>
	e1 := makePNode()
	e2 := makeTextNode()
	e1.AppendChild(e2)

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
	a2 := makeTextNode()
	a1.AppendChild(a2)

	// <i>foobar</i>
	b1 := makeINode()
	b2 := makeTextNode()
	b1.AppendChild(b2)

	// <em><code>foobar</code></em>
	c1 := makeEmNode()
	c2 := makeCodeNode()
	c3 := makeTextNode()
	c1.AppendChild(c2)
	c2.AppendChild(c3)

	// <i><code>foobar</code></i>
	d1 := makeINode()
	d2 := makeCodeNode()
	d3 := makeTextNode()
	d1.AppendChild(d2)
	d2.AppendChild(d3)

	// <p>foobar</p>
	e1 := makePNode()
	e2 := makeTextNode()
	e1.AppendChild(e2)

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
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isItalic(tc.in); out != tc.out {
				t.Errorf("isItalic(%v) = %t, want %t", tc.in, out, tc.out)
			}
		})
	}
}
