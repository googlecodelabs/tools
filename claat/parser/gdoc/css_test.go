package gdoc

import (
	"testing"

	"golang.org/x/net/html"
)

func nodeWithStyle(s string) *html.Node {
	n := makePNode()
	n.Attr = append(n.Attr, html.Attribute{Key: "style", Val: s})
	return n
}

// TODO: test parseStyle

// TODO: test classList

// TODO: test hasClass

// TODO: test hasClassStyle

func TestStyleValue(t *testing.T) {
	tests := []struct {
		name   string
		inNode *html.Node
		inName string
		out    string
	}{
		{
			name:   "NoName",
			inNode: makePNode(),
		},
		{
			name:   "NoStyle",
			inNode: makePNode(),
			inName: "foobar",
		},
		{
			name:   "One",
			inNode: nodeWithStyle("position: absolute"),
			inName: "position",
			out:    "absolute",
		},
		{
			name:   "Capitalization",
			inNode: nodeWithStyle("color: #0000FF"),
			inName: "color",
			out:    "#0000ff",
		},
		{
			name:   "Multiple",
			inNode: nodeWithStyle("position: absolute; color: #ff00ff; font-weight: 300"),
			inName: "color",
			out:    "#ff00ff",
		},
		{
			name:   "NotFound",
			inNode: nodeWithStyle("position: absolute; color: #FF00FF; font-weight: 300"),
			inName: "margin-left",
		},
		{
			name:   "NoKVPair",
			inNode: nodeWithStyle("margin-left"),
			inName: "margin-left",
		},
		{
			// TODO should this be the behavior?
			name:   "BadSyntax",
			inNode: nodeWithStyle("margin-left: font-weight: #00ff00"),
			inName: "margin-left",
			out:    "font-weight: #00ff00",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := styleValue(tc.inNode, tc.inName); out != tc.out {
				t.Errorf("styleValue(%+v, %q) = %q, want %q", tc.inNode, tc.inName, out, tc.out)

			}
		})
	}
}

// TODO: test styleFloatValue
