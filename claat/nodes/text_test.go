package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewTextNode(t *testing.T) {
	tests := []struct {
		name   string
		inOpts NewTextNodeOptions
		out    *TextNode
	}{
		{
			name: "Empty",
			out: &TextNode{
				node: node{typ: NodeText},
			},
		},
		{
			name: "NonEmpty",
			inOpts: NewTextNodeOptions{
				Value: "foobar",
			},
			out: &TextNode{
				node:  node{typ: NodeText},
				Value: "foobar",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewTextNode(tc.inOpts)
			if diff := cmp.Diff(tc.out, out, cmp.AllowUnexported(TextNode{}, node{})); diff != "" {
				t.Errorf("NewTextNode(%+v) got diff (-want +got): %s", tc.inOpts, diff)
				return
			}
		})
	}
}

func TestTextNodeEmpty(t *testing.T) {
	tests := []struct {
		name   string
		inOpts NewTextNodeOptions
		out    bool
	}{
		{
			name: "Empty",
			out:  true,
		},
		{
			name: "NonEmpty",
			inOpts: NewTextNodeOptions{
				Value: "foobar",
			},
		},
		{
			name: "EmptyWithSpaces",
			inOpts: NewTextNodeOptions{
				Value: "\n \t",
			},
			out: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewTextNode(tc.inOpts)
			out := n.Empty()
			if out != tc.out {
				t.Errorf("TextNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
