package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewTextNode(t *testing.T) {
	tests := []struct {
		name    string
		inValue string
		out     *TextNode
	}{
		{
			name: "Empty",
			out: &TextNode{
				node: node{typ: NodeText},
			},
		},
		{
			name:    "NonEmpty",
			inValue: "foobar",
			out: &TextNode{
				node:  node{typ: NodeText},
				Value: "foobar",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewTextNode(tc.inValue)
			if diff := cmp.Diff(tc.out, out, cmp.AllowUnexported(TextNode{}, node{})); diff != "" {
				t.Errorf("NewTextNode(%q) got diff (-want +got): %s", tc.inValue, diff)
				return
			}
		})
	}
}

func TestTextNodeEmpty(t *testing.T) {
	tests := []struct {
		name    string
		inValue string
		out     bool
	}{
		{
			name: "Empty",
			out:  true,
		},
		{
			name:    "NonEmpty",
			inValue: "foobar",
		},
		{
			name:    "EmptyWithSpaces",
			inValue: "\n \t",
			out:     true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewTextNode(tc.inValue)
			out := n.Empty()
			if out != tc.out {
				t.Errorf("TextNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
