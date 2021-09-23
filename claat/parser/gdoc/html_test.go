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

func makeBrNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Br,
		Data:     "br",
	}
}

func makeSpanNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Span,
		Data:     "span",
	}
}

func makeANode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.A,
		Data:     "a",
	}
}

func makeTdNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Td,
		Data:     "td",
	}
}

func makeDivNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Div,
		Data:     "td",
	}
}

func makeOlNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Ol,
		Data:     "ol",
	}
}

func makeUlNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Ul,
		Data:     "ul",
	}
}

func makeTableNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Table,
		Data:     "table",
	}
}

func makeTrNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Tr,
		Data:     "tr",
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

func TestIsMeta(t *testing.T) {
	metaStyleText := `.meta {
	color: #b7b7b7;
}`
	metaStyle, err := parseStyle(makeStyleNode(metaStyleText))
	if err != nil {
		t.Fatalf("parseStyle(makeStyleNode(%q)) = %+v", metaStyleText, err)
		return
	}

	a := nodeWithAttrs(map[string]string{"class": "meta"})
	a.AppendChild(makeTextNode("foobar"))

	b := makePNode()
	b.AppendChild(makeTextNode("foobar"))

	tests := []struct {
		name   string
		inNode *html.Node
		out    bool
	}{
		{
			name:   "Meta",
			inNode: a,
			out:    true,
		},
		{
			name:   "NonMeta",
			inNode: b,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isMeta(metaStyle, tc.inNode); out != tc.out {
				t.Errorf("isMeta(css, %+v) = %t, want %t", tc.inNode, out, tc.out)
			}
		})
	}
}

func TestIsBold(t *testing.T) {
	boldStyleText := `.literalbold {
	font-weight: bold;
}

.weightbold {
	font-weight: 700;
}
`
	boldStyle, err := parseStyle(makeStyleNode(boldStyleText))
	if err != nil {
		t.Fatalf("parseStyle(makeStyleNode(%q)) = %+v", boldStyleText, err)
		return
	}

	a1 := makeStrongNode()
	a2 := makeTextNode("foobar")
	a1.AppendChild(a2)

	b := makeBNode()
	b.AppendChild(makeTextNode("foobar"))

	c := nodeWithAttrs(map[string]string{"class": "literalbold"})
	c.AppendChild(makeTextNode("foobar"))

	d := nodeWithAttrs(map[string]string{"class": "weightbold"})
	d.AppendChild(makeTextNode("foobar"))

	e1 := makeEmNode()
	e2 := makeTextNode("foobar")
	e1.AppendChild(e2)

	tests := []struct {
		name   string
		inNode *html.Node
		out    bool
	}{
		{
			name:   "Strong",
			inNode: a1,
			out:    true,
		},
		{
			name:   "B",
			inNode: b,
			out:    true,
		},
		{
			name:   "FontWeightBold",
			inNode: c,
			out:    true,
		},
		{
			name:   "FontWeight700",
			inNode: d,
			out:    true,
		},
		{
			name:   "TextNodeBold",
			inNode: a2,
			out:    true,
		},
		{
			name:   "TextNodeNonBold",
			inNode: e2,
		},
		{
			name:   "Em",
			inNode: e1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isBold(boldStyle, tc.inNode); out != tc.out {
				t.Errorf("isBold(css, %+v) = %t, want %t", tc.inNode, out, tc.out)
			}
		})
	}
}

func TestIsItalic(t *testing.T) {
	italicStyleText := `.literalitalic {
	font-style: italic;
}
`
	italicStyle, err := parseStyle(makeStyleNode(italicStyleText))
	if err != nil {
		t.Fatalf("parseStyle(makeStyleNode(%q)) = %+v", italicStyleText, err)
		return
	}

	a1 := makeEmNode()
	a2 := makeTextNode("foobar")
	a1.AppendChild(a2)

	b := makeINode()
	b.AppendChild(makeTextNode("foobar"))

	c := nodeWithAttrs(map[string]string{"class": "literalitalic"})
	c.AppendChild(makeTextNode("foobar"))

	d1 := makeStrongNode()
	d2 := makeTextNode("foobar")
	d1.AppendChild(d2)

	tests := []struct {
		name   string
		inNode *html.Node
		out    bool
	}{
		{
			name:   "Em",
			inNode: a1,
			out:    true,
		},
		{
			name:   "I",
			inNode: b,
			out:    true,
		},
		{
			name:   "FontStyleItalic",
			inNode: c,
			out:    true,
		},
		{
			name:   "TextNodeItalic",
			inNode: a2,
			out:    true,
		},
		{
			name:   "TextNodeNonItalic",
			inNode: d2,
		},
		{
			name:   "Strong",
			inNode: d1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isItalic(italicStyle, tc.inNode); out != tc.out {
				t.Errorf("isItalic(css, %+v) = %t, want %t", tc.inNode, out, tc.out)
			}
		})
	}
}

func TestIsConsole(t *testing.T) {
	consoleStyleText := `.console {
	font-family: consolas;
}

.code {
	font-family: courier new;
}
`
	consoleStyle, err := parseStyle(makeStyleNode(consoleStyleText))
	if err != nil {
		t.Fatalf("parseStyle(makeStyleNode(%q)) = %+v", consoleStyleText, err)
		return
	}

	a1 := nodeWithAttrs(map[string]string{"class": "console"})
	a2 := makeTextNode("foobar")
	a1.AppendChild(a2)

	b1 := makePNode()
	b2 := makeTextNode("foobar")
	b1.AppendChild(b2)

	c1 := nodeWithAttrs(map[string]string{"class": "courier new"})
	c2 := makeTextNode("foobar")
	c1.AppendChild(c2)

	tests := []struct {
		name   string
		inNode *html.Node
		out    bool
	}{
		{
			name:   "ConsoleNonText",
			inNode: a1,
			out:    true,
		},
		{
			name:   "ConsoleText",
			inNode: a2,
			out:    true,
		},
		{
			name:   "NonConsoleNonText",
			inNode: b1,
		},
		{
			name:   "NonConsoleText",
			inNode: b2,
		},
		{
			name:   "CodeNonText",
			inNode: c1,
		},
		{
			name:   "CodeText",
			inNode: c2,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isConsole(consoleStyle, tc.inNode); out != tc.out {
				t.Errorf("isConsole(css, %+v) = %t, want %t", tc.inNode, out, tc.out)
			}
		})
	}
}

func TestIsCode(t *testing.T) {
	codeStyleText := `.console {
	font-family: consolas;
}

.code {
	font-family: courier new;
}
`
	codeStyle, err := parseStyle(makeStyleNode(codeStyleText))
	if err != nil {
		t.Fatalf("parseStyle(makeStyleNode(%q)) = %+v", codeStyleText, err)
		return
	}

	a1 := nodeWithAttrs(map[string]string{"class": "console"})
	a2 := makeTextNode("foobar")
	a1.AppendChild(a2)

	b1 := makePNode()
	b2 := makeTextNode("foobar")
	b1.AppendChild(b2)

	c1 := nodeWithAttrs(map[string]string{"class": "code"})
	c2 := makeTextNode("foobar")
	c1.AppendChild(c2)

	tests := []struct {
		name   string
		inNode *html.Node
		out    bool
	}{
		{
			name:   "ConsoleNonText",
			inNode: a1,
		},
		{
			name:   "ConsoleText",
			inNode: a2,
		},
		{
			name:   "NonCodeNonText",
			inNode: b1,
		},
		{
			name:   "NonCodeText",
			inNode: b2,
		},
		{
			name:   "CodeNonText",
			inNode: c1,
			out:    true,
		},
		{
			name:   "CodeText",
			inNode: c2,
			out:    true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isCode(codeStyle, tc.inNode); out != tc.out {
				t.Errorf("isCode(css, %+v) = %t, want %t", tc.inNode, out, tc.out)
			}
		})
	}
}

func TestIsButton(t *testing.T) {
	buttonStyleText := `.button {
	background-color: #6aa84f;
}`
	buttonStyle, err := parseStyle(makeStyleNode(buttonStyleText))
	if err != nil {
		t.Fatalf("parseStyle(makeStyleNode(%q)) = %+v", buttonStyleText, err)
		return
	}

	a := nodeWithAttrs(map[string]string{"class": "button"})
	a.AppendChild(makeTextNode("foobar"))

	b := makePNode()
	b.AppendChild(makeTextNode("foobar"))

	tests := []struct {
		name   string
		inNode *html.Node
		out    bool
	}{
		{
			name:   "Button",
			inNode: a,
			out:    true,
		},
		{
			name:   "NonButton",
			inNode: b,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isButton(buttonStyle, tc.inNode); out != tc.out {
				t.Errorf("isButton(css, %+v) = %t, want %t", tc.inNode, out, tc.out)
			}
		})
	}
}

func TestIsInfobox(t *testing.T) {
	infoboxStyleText := `.infoboxNegative {
	background-color: #fce5cd;
}

.infoboxPositive {
	background-color: #d9ead3;
}
`
	infoboxStyle, err := parseStyle(makeStyleNode(infoboxStyleText))
	if err != nil {
		t.Fatalf("parseStyle(makeStyleNode(%q)) = %+v", infoboxStyleText, err)
		return
	}

	a := makeTdNode()
	a.Attr = append(a.Attr, html.Attribute{Key: "class", Val: "infoboxNegative"})
	a.AppendChild(makeTextNode("foobar"))

	b := makeTdNode()
	b.AppendChild(makeTextNode("foobar"))

	c := nodeWithAttrs(map[string]string{"class": "infoboxNegative"})
	c.AppendChild(makeTextNode("foobar"))

	d := makeTdNode()
	d.Attr = append(d.Attr, html.Attribute{Key: "class", Val: "infoboxPositive"})
	d.AppendChild(makeTextNode("foobar"))

	e := nodeWithAttrs(map[string]string{"class": "infoboxPositive"})
	e.AppendChild(makeTextNode("foobar"))

	tests := []struct {
		name   string
		inNode *html.Node
		out    bool
	}{
		{
			name:   "TdInfoboxNegative",
			inNode: a,
			out:    true,
		},
		{
			name:   "TdNonInfobox",
			inNode: b,
		},
		{
			name:   "NonTdNegative",
			inNode: c,
		},
		{
			name:   "TdInfoboxPositive",
			inNode: d,
			out:    true,
		},
		{
			name:   "NonTdPositive",
			inNode: e,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isInfobox(infoboxStyle, tc.inNode); out != tc.out {
				t.Errorf("isInfobox(css, %+v) = %t, want %t", tc.inNode, out, tc.out)
			}
		})
	}
}

func TestIsInfoboxNegative(t *testing.T) {
	infoboxNegativeStyleText := `.infoboxNegative {
	background-color: #fce5cd;
}`
	infoboxNegativeStyle, err := parseStyle(makeStyleNode(infoboxNegativeStyleText))
	if err != nil {
		t.Fatalf("parseStyle(makeStyleNode(%q)) = %+v", infoboxNegativeStyleText, err)
		return
	}

	a := makeTdNode()
	a.Attr = append(a.Attr, html.Attribute{Key: "class", Val: "infoboxNegative"})
	a.AppendChild(makeTextNode("foobar"))

	b := makeTdNode()
	b.AppendChild(makeTextNode("foobar"))

	c := nodeWithAttrs(map[string]string{"class": "infoboxNegative"})
	c.AppendChild(makeTextNode("foobar"))

	tests := []struct {
		name   string
		inNode *html.Node
		out    bool
	}{
		{
			name:   "TdInfoboxNegative",
			inNode: a,
			out:    true,
		},
		{
			name:   "TdNonInfoboxNegative",
			inNode: b,
		},
		{
			name:   "NonTd",
			inNode: c,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isInfoboxNegative(infoboxNegativeStyle, tc.inNode); out != tc.out {
				t.Errorf("isInfoboxNegative(css, %+v) = %t, want %t", tc.inNode, out, tc.out)
			}
		})
	}
}

func TestIsSurvey(t *testing.T) {
	surveyStyleText := `.survey {
	background-color: #cfe2f3;
}`
	surveyStyle, err := parseStyle(makeStyleNode(surveyStyleText))
	if err != nil {
		t.Fatalf("parseStyle(makeStyleNode(%q)) = %+v", surveyStyleText, err)
		return
	}

	a := makeTdNode()
	a.Attr = append(a.Attr, html.Attribute{Key: "class", Val: "survey"})
	a.AppendChild(makeTextNode("foobar"))

	b := makeTdNode()
	b.AppendChild(makeTextNode("foobar"))

	c := nodeWithAttrs(map[string]string{"class": "survey"})
	c.AppendChild(makeTextNode("foobar"))

	tests := []struct {
		name   string
		inNode *html.Node
		out    bool
	}{
		{
			name:   "TdSurvey",
			inNode: a,
			out:    true,
		},
		{
			name:   "TdNonSurvey",
			inNode: b,
		},
		{
			name:   "NonTd",
			inNode: c,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isSurvey(surveyStyle, tc.inNode); out != tc.out {
				t.Errorf("isSurvey(css, %+v) = %t, want %t", tc.inNode, out, tc.out)
			}
		})
	}
}

func TestIsComment(t *testing.T) {
	commentStyleText := `.comment {
	border: 1px solid black;
}`
	commentStyle, err := parseStyle(makeStyleNode(commentStyleText))
	if err != nil {
		t.Fatalf("parseStyle(makeStyleNode(%q)) = %+v", commentStyleText, err)
		return
	}

	a := makeDivNode()
	a.Attr = append(a.Attr, html.Attribute{Key: "class", Val: "comment"})
	a.AppendChild(makeTextNode("foobar"))

	b := makeDivNode()
	b.AppendChild(makeTextNode("foobar"))

	c := nodeWithAttrs(map[string]string{"class": "comment"})
	c.AppendChild(makeTextNode("foobar"))

	tests := []struct {
		name   string
		inNode *html.Node
		out    bool
	}{
		{
			name:   "DivComment",
			inNode: a,
			out:    true,
		},
		{
			name:   "DivNonComment",
			inNode: b,
		},
		{
			name:   "NonDiv",
			inNode: c,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isComment(commentStyle, tc.inNode); out != tc.out {
				t.Errorf("isComment(css, %+v) = %t, want %t", tc.inNode, out, tc.out)
			}
		})
	}
}

func TestIsTable(t *testing.T) {
	a := makeTableNode()
	a.AppendChild(makeTrNode())
	a.AppendChild(makeTrNode())
	a.AppendChild(makeTdNode())
	a.AppendChild(makeTdNode())

	b := makeTableNode()
	b.AppendChild(makeTrNode())
	b.AppendChild(makeTdNode())
	b.AppendChild(makeTdNode())

	c := makeTableNode()
	c.AppendChild(makeTrNode())
	c.AppendChild(makeTrNode())
	c.AppendChild(makeTdNode())

	d := makeTableNode()
	d.AppendChild(makeTrNode())
	d.AppendChild(makeTdNode())

	e := makeTableNode()
	e.AppendChild(makeTdNode())

	f := makeTableNode()
	f.AppendChild(makeTrNode())

	g := makeMarqueeNode()
	g.AppendChild(makeTrNode())
	g.AppendChild(makeTrNode())
	g.AppendChild(makeTdNode())
	g.AppendChild(makeTdNode())

	tests := []struct {
		name string
		in   *html.Node
		out  bool
	}{
		{
			name: "Table2Rows2Data",
			in:   a,
			out:  true,
		},
		{
			name: "Table1Row2Data",
			in:   b,
			out:  true,
		},
		{
			name: "Table2Rows1Data",
			in:   c,
			out:  true,
		},
		{
			name: "Table1Row1Data",
			in:   d,
		},
		{
			name: "Table0Rows1Data",
			in:   e,
		},
		{
			name: "Table1Row0Data",
			in:   f,
		},
		{
			name: "TableNone",
			in:   makeTableNode(),
		},
		{
			name: "NonTableAtom",
			in:   makeMarqueeNode(),
		},
		{
			name: "NonTableAtomRowsAndData",
			in:   g,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isTable(tc.in); out != tc.out {
				t.Errorf("isTable(%v) = %t, want %t", tc.in, out, tc.out)
			}
		})
	}
}

func TestIsList(t *testing.T) {
	a1 := makeOlNode()
	a2 := makeTextNode("aaa")
	a3 := makeTextNode("bbb")
	a4 := makeTextNode("ccc")
	// The name and input nodes should be siblings.
	a1.AppendChild(a2)
	a1.AppendChild(a3)
	a1.AppendChild(a4)

	b1 := makeUlNode()
	b2 := makeTextNode("aaa")
	b3 := makeTextNode("bbb")
	b4 := makeTextNode("ccc")
	// The name and input nodes should be siblings.
	b1.AppendChild(b2)
	b1.AppendChild(b3)
	b1.AppendChild(b4)

	tests := []struct {
		name string
		in   *html.Node
		out  bool
	}{
		{
			name: "OrderedListWithElements",
			in:   a1,
			out:  true,
		},
		{
			name: "UnorderedListWithElements",
			in:   b1,
			out:  true,
		},
		// TODO: Should a list of no elements be considered an error?
		{
			name: "OrderedListWithoutElements",
			in:   makeOlNode(),
			out:  true,
		},
		{
			name: "UnorderedListWithoutElements",
			in:   makeUlNode(),
			out:  true,
		},
		{
			name: "NotAList",
			in:   makeBlinkNode(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isList(tc.in); out != tc.out {
				t.Errorf("isList(%v) = %t, want %t", tc.in, out, tc.out)
			}
		})
	}
}

func TestCountTwo(t *testing.T) {
	a1 := makePNode()
	a2 := makeBlinkNode()
	a3 := makeTextNode("foobar")
	a1.AppendChild(a2)
	a2.AppendChild(a3)

	b1 := makePNode()
	b2 := makeTextNode("foobar")
	b3 := makeMarqueeNode()
	// The nodes should be siblings.
	b1.AppendChild(b2)
	b1.AppendChild(b3)

	c1 := makePNode()
	c2 := makeTextNode("foobar")
	c3 := makeMarqueeNode()
	c4 := makeTextNode("foobar2")
	c5 := makeMarqueeNode()
	// The nodes should be siblings.
	c1.AppendChild(c2)
	c1.AppendChild(c3)
	c1.AppendChild(c4)
	c1.AppendChild(c5)

	d1 := makePNode()
	d2 := makeTextNode("foobar")
	d3 := makeMarqueeNode()
	d4 := makeTextNode("foobar2")
	d5 := makeMarqueeNode()
	d6 := makeMarqueeNode()
	d7 := makeMarqueeNode()
	// The nodes should be siblings.
	d1.AppendChild(d2)
	d1.AppendChild(d3)
	d1.AppendChild(d4)
	d1.AppendChild(d5)
	d1.AppendChild(d6)
	d1.AppendChild(d7)

	tests := []struct {
		name   string
		inNode *html.Node
		inAtom atom.Atom
		out    int
	}{
		{
			name:   "Zero",
			inNode: a1,
			inAtom: atom.Marquee,
			out:    0,
		},
		{
			name:   "One",
			inNode: b1,
			inAtom: atom.Marquee,
			out:    1,
		},
		{
			name:   "Two",
			inNode: c1,
			inAtom: atom.Marquee,
			out:    2,
		},
		{
			name:   "MoreThanTwo",
			inNode: d1,
			inAtom: atom.Marquee,
			out:    2,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := countTwo(tc.inNode, tc.inAtom); out != tc.out {
				t.Errorf("countTwo(%+v, %+v) = %d, want %d", tc.inNode, tc.inAtom, out, tc.out)
			}
		})
	}
}

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

func TestFindChildAtoms(t *testing.T) {
	a1 := makePNode()
	a2 := makeEmNode()
	a3 := makeTextNode("foobar")
	a1.AppendChild(a2)
	a2.AppendChild(a3)

	b1 := makePNode()
	b2 := makeCodeNode()
	b3 := makeEmNode()
	b4 := makeStrongNode()
	b5 := makeTextNode("foobar")
	b1.AppendChild(b2)
	b2.AppendChild(b3)
	b3.AppendChild(b4)
	b4.AppendChild(b5)

	c1 := makePNode()
	c2 := makeCodeNode()
	c3 := makeTextNode("foobar1")
	c4 := makeEmNode()
	c5 := makeTextNode("foobar2")
	c6 := makeStrongNode()
	c7 := makeCodeNode()
	c8 := makeTextNode("foobar3")
	//<p><code>foobar1</code><em>foobar2</em><strong><code>foobar3</code></strong></p>
	c1.AppendChild(c2)
	c2.AppendChild(c3)
	c1.AppendChild(c4)
	c4.AppendChild(c5)
	c1.AppendChild(c6)
	c6.AppendChild(c7)
	c7.AppendChild(c8)

	tests := []struct {
		name   string
		inNode *html.Node
		inAtom atom.Atom
		out    []*html.Node
	}{
		{
			name:   "One",
			inNode: a1,
			inAtom: atom.Em,
			out:    []*html.Node{a2},
		},
		{
			name:   "DistantDescendant",
			inNode: b1,
			inAtom: atom.Strong,
			out:    []*html.Node{b4},
		},
		{
			name:   "Multi",
			inNode: c1,
			inAtom: atom.Code,
			out:    []*html.Node{c2, c7},
		},
		{
			name:   "None",
			inNode: a1,
			inAtom: atom.Marquee,
		},
		{
			name:   "Self",
			inNode: a1,
			inAtom: atom.P,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if diff := cmp.Diff(tc.out, findChildAtoms(tc.inNode, tc.inAtom)); diff != "" {
				t.Errorf("findChildAtoms(%+v, %+v) got diff (-want +got):\n%s", tc.inNode, tc.inAtom, diff)
			}
		})
	}
}

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

func TestStringifyNode(t *testing.T) {
	a1 := makePNode()
	a2 := makeTextNode("1")
	a3 := makeBNode()
	a4 := makeTextNode("2 ")
	a5 := makeTextNode("3")
	a6 := makeINode()
	a7 := makeTextNode("  4\n")
	a1.AppendChild(a2)
	a1.AppendChild(a3)
	a1.AppendChild(a5)
	a1.AppendChild(a6)
	a3.AppendChild(a4)
	a6.AppendChild(a7)

	b1 := makePNode()
	b2 := makeTextNode("foo")
	b3 := makeBrNode()
	b4 := makeTextNode("bar")
	b1.AppendChild(b2)
	b1.AppendChild(b3)
	b1.AppendChild(b4)

	c1 := makePNode()
	c2 := makeTextNode("foo")
	c3 := makeSpanNode()
	c4 := makeTextNode("bar")
	c1.AppendChild(c2)
	c1.AppendChild(c3)
	c1.AppendChild(c4)

	d1 := makePNode()
	d2 := makeANode()
	d2.Attr = append(d2.Attr, html.Attribute{Key: "href", Val: "google.com"})
	d3 := makeTextNode("foobar")
	d1.AppendChild(d2)
	d2.AppendChild(d3)

	e1 := makePNode()
	e2 := makeANode()
	e2.Attr = append(e2.Attr, html.Attribute{Key: "href", Val: "#cmnt"})
	e3 := makeTextNode("foobar")
	e1.AppendChild(e2)
	e2.AppendChild(e3)

	tests := []struct {
		name        string
		inRoot      *html.Node
		inTrim      bool
		inLineBreak bool
		out         string
	}{
		{
			name:   "TextRoot",
			inRoot: makeTextNode(" foo bar"),
			out:    " foo bar",
		},
		{
			name:   "TextRootTrim",
			inRoot: makeTextNode(" foo bar"),
			inTrim: true,
			out:    "foo bar",
		},
		{
			name:   "StyledText",
			inRoot: a1,
			out:    "12 3  4\n",
		},
		{
			name:   "StyledTextTrim",
			inRoot: a1,
			inTrim: true,
			out:    "12 3  4",
		},
		{
			name:   "BrNonRoot",
			inRoot: b1,
			out:    "foobar",
		},
		{
			name:        "BrNonRootLineBreak",
			inRoot:      b1,
			inLineBreak: true,
			out:         "foo\nbar",
		},
		{
			name:   "SpanNonRoot",
			inRoot: c1,
			out:    "foobar",
		},
		{
			name:        "SpanNonRootLineBreak",
			inRoot:      c1,
			inLineBreak: true,
			out:         "foobar",
		},
		{
			name:   "AComment",
			inRoot: d1,
			out:    "foobar",
		},
		{
			name:   "ANonComment",
			inRoot: e1,
		},
		{
			name:   "BrRoot",
			inRoot: makeBrNode(),
		},
		{
			name:   "BrRootTrim",
			inRoot: makeBrNode(),
			inTrim: true,
		},
		{
			name:        "BrRootLineBreak",
			inRoot:      makeBrNode(),
			inLineBreak: true,
			out:         "\n",
		},
		{
			name:        "BrRootTrimLineBreak",
			inRoot:      makeBrNode(),
			inTrim:      true,
			inLineBreak: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if diff := cmp.Diff(tc.out, stringifyNode(tc.inRoot, tc.inTrim, tc.inLineBreak)); diff != "" {
				t.Errorf("stringifyNode(%+v, %t, %t) got diff (-want +got):\n%s", tc.inRoot, tc.inTrim, tc.inLineBreak, diff)
			}
		})
	}
}
