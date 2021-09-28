package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// AllowUnexported option for cmp to make sure we can diff properly.
var cmpOptButton = cmp.AllowUnexported(ButtonNode{}, node{}, ListNode{}, TextNode{})

//func NewButtonNode(raise, color, download bool, n ...Node) *ButtonNode {
func TestNewButtonNode(t *testing.T) {
	tests := []struct {
		name       string
		inRaise    bool
		inColor    bool
		inDownload bool
		inContent  []Node
		out        *ButtonNode
	}{
		{
			name: "Empty",
			out: &ButtonNode{
				node:    node{typ: NodeButton},
				Content: NewListNode(),
			},
		},
		{
			name:    "Raise",
			inRaise: true,
			out: &ButtonNode{
				node:    node{typ: NodeButton},
				Raise:   true,
				Content: NewListNode(),
			},
		},
		{
			name:    "Color",
			inColor: true,
			out: &ButtonNode{
				node:    node{typ: NodeButton},
				Color:   true,
				Content: NewListNode(),
			},
		},
		{
			name:       "Download",
			inDownload: true,
			out: &ButtonNode{
				node:     node{typ: NodeButton},
				Download: true,
				Content:  NewListNode(),
			},
		},
		{
			name: "ContentNoSettings",
			inContent: []Node{
				NewTextNode(NewTextNodeOptions{Value: "foo"}),
				NewTextNode(NewTextNodeOptions{Value: "bar"}),
			},
			out: &ButtonNode{
				node:    node{typ: NodeButton},
				Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "foo"}), NewTextNode(NewTextNodeOptions{Value: "bar"})),
			},
		},
		{
			name:       "ContentAllSettings",
			inRaise:    true,
			inColor:    true,
			inDownload: true,
			inContent: []Node{
				NewTextNode(NewTextNodeOptions{Value: "foo"}),
				NewTextNode(NewTextNodeOptions{Value: "bar"}),
			},
			out: &ButtonNode{
				node:     node{typ: NodeButton},
				Raise:    true,
				Color:    true,
				Download: true,
				Content:  NewListNode(NewTextNode(NewTextNodeOptions{Value: "foo"}), NewTextNode(NewTextNodeOptions{Value: "bar"})),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewButtonNode(tc.inRaise, tc.inColor, tc.inDownload, tc.inContent...)
			if diff := cmp.Diff(tc.out, out, cmpOptButton); diff != "" {
				t.Errorf("NewButtonNode(%t, %t, %t,%+v) got diff (-want +got): %s", tc.inRaise, tc.inColor, tc.inDownload, tc.inContent, diff)
				return
			}
		})
	}
}

func TestButtonNodeEmpty(t *testing.T) {
	tests := []struct {
		name       string
		inRaise    bool
		inColor    bool
		inDownload bool
		inContent  []Node
		out        bool
	}{
		{
			name:       "Empty",
			inRaise:    true,
			inColor:    true,
			inDownload: true,
			out:        true,
		},
		{
			name: "NonEmpty",
			inContent: []Node{
				NewTextNode(NewTextNodeOptions{Value: "foo"}),
				NewTextNode(NewTextNodeOptions{Value: "bar"}),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewButtonNode(tc.inRaise, tc.inColor, tc.inDownload, tc.inContent...)
			out := n.Empty()
			if out != tc.out {
				t.Errorf("ButtonNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
