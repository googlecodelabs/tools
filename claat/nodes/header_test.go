package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// AllowUnexported option for cmp to make sure we can diff properly.
var cmpOptHeader = cmp.AllowUnexported(HeaderNode{}, node{}, ListNode{}, TextNode{})

func TestNewHeaderNode(t *testing.T) {
	tests := []struct {
		name      string
		inLevel   int
		inContent []Node
		out       *HeaderNode
	}{
		{
			name:    "Empty",
			inLevel: 1,
			out: &HeaderNode{
				node:    node{typ: NodeHeader},
				Level:   1,
				Content: NewListNode(),
			},
		},
		{
			name:    "NonEmpty",
			inLevel: 1,
			inContent: []Node{
				NewTextNode("foo"),
				NewTextNode("bar"),
			},
			out: &HeaderNode{
				node:    node{typ: NodeHeader},
				Level:   1,
				Content: NewListNode(NewTextNode("foo"), NewTextNode("bar")),
			},
		},
		{
			name:    "ValidLevel",
			inLevel: 2,
			out: &HeaderNode{
				node:    node{typ: NodeHeader},
				Level:   2,
				Content: NewListNode(),
			},
		},
		// TODO should the function accept levels that do not correspond to <h_> elements?
		{
			name: "ZeroLevel",
			out: &HeaderNode{
				node:    node{typ: NodeHeader},
				Content: NewListNode(),
			},
		},
		{
			name:    "NegativeLevel",
			inLevel: -1337,
			out: &HeaderNode{
				node:    node{typ: NodeHeader},
				Level:   -1337,
				Content: NewListNode(),
			},
		},
		{
			name:    "VeryHighLevel",
			inLevel: 1337,
			out: &HeaderNode{
				node:    node{typ: NodeHeader},
				Level:   1337,
				Content: NewListNode(),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewHeaderNode(tc.inLevel, tc.inContent...)
			if diff := cmp.Diff(tc.out, out, cmpOptHeader); diff != "" {
				t.Errorf("NewHeaderNode(%d, %+v) got diff (-want +got): %s", tc.inLevel, tc.inContent, diff)
				return
			}
		})
	}
}

func TestHeaderNodeEmpty(t *testing.T) {
	tests := []struct {
		name      string
		inLevel   int
		inContent []Node
		out       bool
	}{
		{
			name:    "Empty",
			inLevel: 1,
			out:     true,
		},
		{
			name:    "NonEmpty",
			inLevel: 1,
			inContent: []Node{
				NewTextNode("foo"),
				NewTextNode("bar"),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewHeaderNode(tc.inLevel, tc.inContent...)
			out := n.Empty()
			if out != tc.out {
				t.Errorf("HeaderNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
