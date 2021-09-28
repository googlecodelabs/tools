package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// AllowUnexported option for cmp to make sure we can diff properly.
var cmpOptInfobox = cmp.AllowUnexported(InfoboxNode{}, node{}, ListNode{}, TextNode{})

func TestNewInfoboxNode(t *testing.T) {
	tests := []struct {
		name      string
		inKind    InfoboxKind
		inContent []Node
		out       *InfoboxNode
	}{
		{
			name:   "PositiveEmpty",
			inKind: InfoboxPositive,
			out: &InfoboxNode{
				node: node{typ: NodeInfobox},
				Kind: InfoboxPositive,
				// TODO: Do we really want this to not be nil?
				Content: NewListNode(),
			},
		},
		{
			name:      "PositiveOneContent",
			inKind:    InfoboxPositive,
			inContent: []Node{NewTextNode(NewTextNodeOptions{Value: "hello"})},
			out: &InfoboxNode{
				node:    node{typ: NodeInfobox},
				Kind:    InfoboxPositive,
				Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "hello"})),
			},
		},
		{
			name:      "PositiveMultiContent",
			inKind:    InfoboxPositive,
			inContent: []Node{NewTextNode(NewTextNodeOptions{Value: "orange"}), NewTextNode(NewTextNodeOptions{Value: "strawberry"}), NewTextNode(NewTextNodeOptions{Value: "pineapple"})},
			out: &InfoboxNode{
				node:    node{typ: NodeInfobox},
				Kind:    InfoboxPositive,
				Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "orange"}), NewTextNode(NewTextNodeOptions{Value: "strawberry"}), NewTextNode(NewTextNodeOptions{Value: "pineapple"})),
			},
		},
		{
			name:   "NegativeEmpty",
			inKind: InfoboxNegative,
			out: &InfoboxNode{
				node: node{typ: NodeInfobox},
				Kind: InfoboxNegative,
				// TODO: Do we really want this to not be nil?
				Content: NewListNode(),
			},
		},
		{
			name:      "NegativeOneContent",
			inKind:    InfoboxNegative,
			inContent: []Node{NewTextNode(NewTextNodeOptions{Value: "hello"})},
			out: &InfoboxNode{
				node:    node{typ: NodeInfobox},
				Kind:    InfoboxNegative,
				Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "hello"})),
			},
		},
		{
			name:      "NegativeMultiContent",
			inKind:    InfoboxNegative,
			inContent: []Node{NewTextNode(NewTextNodeOptions{Value: "orange"}), NewTextNode(NewTextNodeOptions{Value: "strawberry"}), NewTextNode(NewTextNodeOptions{Value: "pineapple"})},
			out: &InfoboxNode{
				node:    node{typ: NodeInfobox},
				Kind:    InfoboxNegative,
				Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "orange"}), NewTextNode(NewTextNodeOptions{Value: "strawberry"}), NewTextNode(NewTextNodeOptions{Value: "pineapple"})),
			},
		},
		{
			// TODO: Should we set a default value?
			name:      "NoKind",
			inContent: []Node{NewTextNode(NewTextNodeOptions{Value: "orange"})},
			out: &InfoboxNode{
				node:    node{typ: NodeInfobox},
				Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "orange"})),
			},
		},
		{
			// TODO: Should this return an error instead?
			name:      "UnsupportedKind",
			inKind:    "this is not a valid kind of infobox",
			inContent: []Node{NewTextNode(NewTextNodeOptions{Value: "orange"})},
			out: &InfoboxNode{
				node:    node{typ: NodeInfobox},
				Kind:    "this is not a valid kind of infobox",
				Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "orange"})),
			},
		},
		{
			name:      "ListOfOneList",
			inKind:    InfoboxPositive,
			inContent: []Node{NewListNode(NewTextNode(NewTextNodeOptions{Value: "a"}), NewTextNode(NewTextNodeOptions{Value: "b"}))},
			out: &InfoboxNode{
				node:    node{typ: NodeInfobox},
				Kind:    InfoboxPositive,
				Content: NewListNode(NewListNode(NewTextNode(NewTextNodeOptions{Value: "a"}), NewTextNode(NewTextNodeOptions{Value: "b"}))),
			},
		},
		{
			name: "Empty",
			out: &InfoboxNode{
				node:    node{typ: NodeInfobox},
				Content: NewListNode(),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewInfoboxNode(tc.inKind, tc.inContent...)
			if diff := cmp.Diff(tc.out, out, cmpOptInfobox); diff != "" {
				t.Errorf("NewInfoboxNode(%q, %v) got diff (-want +got): %s", tc.inKind, tc.inContent, diff)
				return
			}
		})
	}
}

func TestInfoboxNodeEmpty(t *testing.T) {
	tests := []struct {
		name      string
		inKind    InfoboxKind
		inContent []Node
		out       bool
	}{
		{
			name:   "Empty",
			inKind: InfoboxPositive,
			out:    true,
		},
		{
			name:      "NonEmpty",
			inKind:    InfoboxPositive,
			inContent: []Node{NewTextNode(NewTextNodeOptions{Value: "a"})},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewInfoboxNode(tc.inKind, tc.inContent...)
			out := n.Empty()
			if out != tc.out {
				t.Errorf("InfoboxNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
