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
			name:   "CapitalizationKeyStyle",
			inNode: nodeWithStyle("Position: relative"),
			inName: "position",
			out:    "relative",
		},
		{
			name:   "CapitalizationValueStyle",
			inNode: nodeWithStyle("color: #0000FF"),
			inName: "color",
			out:    "#0000ff",
		},
		{
			name:   "CapitalizationKeyInput",
			inNode: nodeWithStyle("position: relative"),
			inName: "Position",
			out:    "relative",
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

func TestStyleFloatValue(t *testing.T) {
	tests := []struct {
		name   string
		inNode *html.Node
		inName string
		out    float32
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
			name:   "Simple",
			inNode: nodeWithStyle("margin-top: 3.14em"),
			inName: "margin-top",
			out:    3.14,
		},
		{
			name:   "NoDecimalPlaces",
			inNode: nodeWithStyle("margin-left: 2in"),
			inName: "margin-left",
			out:    2,
		},
		{
			name:   "DecimalZeroes",
			inNode: nodeWithStyle("margin-right: 1.0px"),
			inName: "margin-right",
			out:    1,
		},
		{
			name:   "NoUnit",
			inNode: nodeWithStyle("margin-bottom: 4"),
			inName: "margin-bottom",
			out:    4,
		},
		{
			name:   "Multiple",
			inNode: nodeWithStyle("padding-top: 1.2; padding-left: 3.4; padding-right: 5.6"),
			inName: "padding-left",
			out:    3.4,
		},
		{
			name:   "NotFound",
			inNode: nodeWithStyle("border-top: 7.8; border-left: 0.9"),
			inName: "border-right",
		},
		{
			name:   "NoKVPair",
			inNode: nodeWithStyle("margin-left"),
			inName: "margin-left",
		},
		{
			name:   "BadSyntax",
			inNode: nodeWithStyle("margin-left: margin-top: 1.234em"),
			inName: "margin-left",
			out:    -1,
		},
		{
			// TODO should this be the behavior?
			name:   "BadSyntaxMiddle",
			inNode: nodeWithStyle("margin-left: margin-top: 1.234em"),
			inName: "margin-top",
		},
		{
			// TODO should this be the behavior?
			name:   "BadValue",
			inNode: nodeWithStyle("margin-left: 7jv9ue4if4.21"),
			inName: "margin-left",
			out:    7,
		},
		{
			name:   "CapitalizationKeyStyle",
			inNode: nodeWithStyle("Margin-Left: 2.3px"),
			inName: "margin-left",
			out:    2.3,
		},
		{
			name:   "CapitalizationKeyInput",
			inNode: nodeWithStyle("margin-left: 4.5px"),
			inName: "Margin-Left",
			out:    4.5,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := styleFloatValue(tc.inNode, tc.inName); out != tc.out {
				t.Errorf("styleFloatValue(%+v, %q) = %f, want %f", tc.inNode, tc.inName, out, tc.out)

			}
		})
	}
}
