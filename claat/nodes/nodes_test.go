package nodes

import (
	"testing"
)

func TestEmptyNodes(t *testing.T) {
	tests := []struct {
		name    string
		inNodes []Node
		out     bool
	}{
		{
			name: "Zero",
			out:  true,
		},
		{
			name: "OneEmpty",
			inNodes: []Node{
				NewTextNode(""),
			},
			out: true,
		},
		{
			name: "OneNonEmpty",
			inNodes: []Node{
				NewTextNode("foo"),
			},
		},
		{
			name: "MultipleEmpty",
			inNodes: []Node{
				NewTextNode(""),
				NewTextNode(""),
				NewTextNode(""),
			},
			out: true,
		},
		{
			name: "MultipleSomeNonEmpty",
			inNodes: []Node{
				NewTextNode("foo"),
				NewTextNode(""),
				NewTextNode("bar"),
			},
		},
		{
			name: "MultipleAllNonEmpty",
			inNodes: []Node{
				NewTextNode("foo"),
				NewTextNode("bar"),
				NewTextNode("baz"),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := EmptyNodes(tc.inNodes)
			if out != tc.out {
				t.Errorf("EmptyNodes(%+v) = %t, want %t", tc.inNodes, out, tc.out)
				return
			}
		})
	}
}
