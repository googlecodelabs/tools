package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var cmpOptList = cmp.AllowUnexported(ListNode{}, node{}, TextNode{})

func TestNewListNode(t *testing.T) {
	tests := []struct {
		name    string
		inNodes []Node
		out     *ListNode
	}{
		{
			name: "Empty",
			out: &ListNode{
				node: node{typ: NodeList},
			},
		},
		{
			name: "One",
			inNodes: []Node{
				NewTextNode(NewTextNodeOptions{Value: "foo"}),
			},
			out: &ListNode{
				node: node{typ: NodeList},
				Nodes: []Node{
					NewTextNode(NewTextNodeOptions{Value: "foo"}),
				},
			},
		},
		{
			name: "Multiple",
			inNodes: []Node{
				NewTextNode(NewTextNodeOptions{Value: "foo"}),
				NewTextNode(NewTextNodeOptions{Value: "bar"}),
				NewTextNode(NewTextNodeOptions{Value: "baz"}),
			},
			out: &ListNode{
				node: node{typ: NodeList},
				Nodes: []Node{
					NewTextNode(NewTextNodeOptions{Value: "foo"}),
					NewTextNode(NewTextNodeOptions{Value: "bar"}),
					NewTextNode(NewTextNodeOptions{Value: "baz"}),
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewListNode(tc.inNodes...)
			if diff := cmp.Diff(tc.out, out, cmpOptList); diff != "" {
				t.Errorf("NewListNode(%v) got diff (-want +got): %s", tc.inNodes, diff)
				return
			}
		})
	}
}

func TestListNodeEmpty(t *testing.T) {
	tests := []struct {
		name    string
		inNodes []Node
		out     bool
	}{
		{
			name: "Empty",
			out:  true,
		},
		{
			name: "NonEmpty",
			inNodes: []Node{
				NewTextNode(NewTextNodeOptions{Value: "foo"}),
				NewTextNode(NewTextNodeOptions{Value: "bar"}),
				NewTextNode(NewTextNodeOptions{Value: "baz"}),
			},
		},
		{
			name: "EmptyWithNodes",
			inNodes: []Node{
				NewTextNode(NewTextNodeOptions{Value: ""}),
				NewTextNode(NewTextNodeOptions{Value: ""}),
				NewTextNode(NewTextNodeOptions{Value: ""}),
			},
			out: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewListNode(tc.inNodes...)
			out := n.Empty()
			if out != tc.out {
				t.Errorf("ListNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
