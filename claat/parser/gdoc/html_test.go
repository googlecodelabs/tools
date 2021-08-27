package gdoc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func makeBlinkNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Blink,
		Data:     "blink",
	}
}

func makePNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.P,
		Data:     "p",
	}
}

func makeEmNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Em,
		Data:     "em",
	}
}

func makeMarqueeNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Marquee,
		Data:     "marquee",
	}
}

func makeStrongNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Strong,
		Data:     "strong",
	}
}

func makeCodeNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Code,
		Data:     "code",
	}
}

func makeBNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.B,
		Data:     "b",
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

func TestIsHeader(t *testing.T) {
	tests := []struct {
		name string
		in   *html.Node
		out  bool
	}{
		{
			name: "StepTitle",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H1,
				Data:     "h1",
			},
		},
		{
			name: "FirstLevel",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H2,
				Data:     "h2",
			},
			out: true,
		},
		{
			name: "SecondLevel",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H3,
				Data:     "h3",
			},
			out: true,
		},
		{
			name: "ThirdLevel",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H4,
				Data:     "h4",
			},
			out: true,
		},
		{
			name: "FourthLevel",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H5,
				Data:     "h5",
			},
			out: true,
		},
		{
			name: "FifthLevel",
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

// TODO: test isMeta

// TODO: test isBold

// TODO: test isItalic

// TODO: test isConsole

// TODO: test isCode

// TODO: test isButton

// TODO: test isInfobox

// TODO: test isInfoboxNegative

// TODO: test isSurvey

// TODO: test isComment

// TODO: test isTable

// TODO: test isList

// TODO: test countTwo

func TestCountDirect(t *testing.T) {
	a1 := makePNode()
	a2 := makeTextNode("foobar")
	a1.AppendChild(a2)

	b1 := makePNode()
	b2 := makeTextNode("foobar")
	b3 := makeTextNode("foobar2")
	b4 := makeTextNode("foobar3")
	// The nodes should be siblings.
	b1.AppendChild(b2)
	b1.AppendChild(b3)
	b1.AppendChild(b4)

	c1 := makePNode()
	c2 := makeBlinkNode()
	c3 := makeTextNode("foobar")
	c1.AppendChild(c2)
	c2.AppendChild(c3)

	tests := []struct {
		name string
		in   *html.Node
		out  int
	}{
		{
			name: "Zero",
			in:   makePNode(),
			out:  0,
		},
		{
			name: "One",
			in:   a1,
			out:  1,
		},
		{
			name: "MoreThanOne",
			in:   b1,
			out:  3,
		},
		{
			name: "NonRecursive",
			in:   c1,
			out:  1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := countDirect(tc.in); out != tc.out {
				t.Errorf("countDirect(%+v) = %d, want %d", tc.in, out, tc.out)
			}
		})
	}
}

func TestFindAtom(t *testing.T) {
	a1 := makePNode()
	a2 := makeEmNode()
	a3 := makeTextNode("foobar")
	a1.AppendChild(a2)
	a2.AppendChild(a3)

	b1 := makePNode()
	b2 := makeMarqueeNode()
	b3 := makeMarqueeNode()
	b4 := makeBlinkNode()
	// The nodes should be siblings.
	b1.AppendChild(b2)
	b1.AppendChild(b3)
	b1.AppendChild(b4)

	c1 := makePNode()
	c2 := makeEmNode()
	c3 := makeStrongNode()
	c4 := makeTextNode("foobar")
	c1.AppendChild(c2)
	c2.AppendChild(c3)
	c3.AppendChild(c4)

	d1 := makeBlinkNode()

	e1 := makeEmNode()
	e2 := makeStrongNode()
	e3 := makeTextNode("foobar")
	e1.AppendChild(e2)
	e2.AppendChild(e3)

	tests := []struct {
		name   string
		inNode *html.Node
		inAtom atom.Atom
		out    *html.Node
	}{
		{
			name:   "OneMatch",
			inNode: a1,
			inAtom: atom.Em,
			out:    a2,
		},
		{
			name:   "MultipleMatches",
			inNode: b1,
			inAtom: atom.Marquee,
			out:    b2,
		},
		{
			name:   "Recursive",
			inNode: c1,
			inAtom: atom.Strong,
			out:    c3,
		},
		{
			name:   "Self",
			inNode: d1,
			inAtom: atom.Blink,
			out:    d1,
		},
		{
			name:   "NoMatches",
			inNode: e1,
			inAtom: atom.Div,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := findAtom(tc.inNode, tc.inAtom); out != tc.out {
				t.Errorf("findAtom(%+v, %+v) = %+v, want %v", tc.inNode, tc.inAtom, out, tc.out)
			}
		})
	}
}

// TODO: test findChildAtoms

func TestFindParent(t *testing.T) {
	a1 := makePNode()
	a2 := makeStrongNode()
	a3 := makeEmNode()
	a4 := makeCodeNode()
	a5 := makeTextNode("foobar")
	a1.AppendChild(a2)
	a2.AppendChild(a3)
	a3.AppendChild(a4)
	a4.AppendChild(a5)

	tests := []struct {
		name   string
		inNode *html.Node
		inAtom atom.Atom
		out    *html.Node
	}{
		{
			name:   "Parent",
			inNode: a4,
			inAtom: atom.Em,
			out:    a3,
		},
		{
			name:   "DistantAncestor",
			inNode: a4,
			inAtom: atom.P,
			out:    a1,
		},
		{
			name:   "Self",
			inNode: a4,
			inAtom: atom.Code,
			out:    a4,
		},
		{
			name:   "NotFound",
			inNode: a4,
			inAtom: atom.Blink,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if diff := cmp.Diff(tc.out, findParent(tc.inNode, tc.inAtom)); diff != "" {
				t.Errorf("findParent(%+v, %+v) got diff (-want +got):\n%s", tc.inNode, tc.inAtom, diff)
			}
		})
	}
}

func TestFindBlockParent(t *testing.T) {
	// Choice of <p> from blockParents is arbitrary.
	a1 := makePNode()
	a2 := makeBNode()
	a3 := makeINode()
	a4 := makeCodeNode()
	a5 := makeTextNode("foobar")
	a1.AppendChild(a2)
	a2.AppendChild(a3)
	a3.AppendChild(a4)
	a4.AppendChild(a5)

	tests := []struct {
		name string
		in   *html.Node
		out  *html.Node
	}{
		{
			name: "Parent",
			in:   a2,
			out:  a1,
		},
		{
			name: "DistantAncestor",
			in:   a5,
			out:  a1,
		},
		{
			name: "Self",
			in:   a1,
		},
		{
			name: "None",
			in:   makeBlinkNode(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if diff := cmp.Diff(tc.out, findBlockParent(tc.in)); diff != "" {
				t.Errorf("findBlockParent(%+v) got diff (-want +got):\n%s", tc.in, diff)
			}
		})
	}
}

func TestNodeAttr(t *testing.T) {
	a1 := makeBlinkNode()
	a1.Attr = append(a1.Attr, html.Attribute{Key: "keyone", Val: "valone"})
	a1.Attr = append(a1.Attr, html.Attribute{Key: "keytwo", Val: "valtwo"})
	a1.Attr = append(a1.Attr, html.Attribute{Key: "keythree", Val: "valthree"})

	tests := []struct {
		name   string
		inNode *html.Node
		inKey  string
		out    string
	}{
		{
			name:   "Simple",
			inNode: a1,
			inKey:  "keyone",
			out:    "valone",
		},
		{
			name:   "MixedCase",
			inNode: a1,
			inKey:  "KEytWO",
			out:    "valtwo",
		},
		{
			name:   "NotFound",
			inNode: a1,
			inKey:  "nokey",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if diff := cmp.Diff(tc.out, nodeAttr(tc.inNode, tc.inKey)); diff != "" {
				t.Errorf("nodeAttr(%+v, %s) got diff (-want +got):\n%s", tc.inNode, tc.inKey, diff)
			}
		})
	}
}

// TODO: test stringifyNode
