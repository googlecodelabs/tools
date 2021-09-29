package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// AllowUnexported option for cmp to make sure we can diff properly.
var cmpOptURL = cmp.AllowUnexported(URLNode{}, node{}, ListNode{}, TextNode{})

func TestNewURLNode(t *testing.T) {
	tests := []struct {
		name      string
		inURL     string
		inContent []Node
		out       *URLNode
	}{
		{
			name: "Empty",
			out: &URLNode{
				node:    node{typ: NodeURL},
				Target:  "_blank",
				Content: NewListNode(),
			},
		},
		{
			name:  "NonEmpty",
			inURL: "google.com",
			inContent: []Node{
				NewTextNode(NewTextNodeOptions{Value: "foo"}),
				NewTextNode(NewTextNodeOptions{Value: "bar"}),
			},
			out: &URLNode{
				node:    node{typ: NodeURL},
				URL:     "google.com",
				Target:  "_blank",
				Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "foo"}), NewTextNode(NewTextNodeOptions{Value: "bar"})),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewURLNode(tc.inURL, tc.inContent...)
			if diff := cmp.Diff(tc.out, out, cmpOptURL); diff != "" {
				t.Errorf("NewURLNode(%s, %+v) got diff (-want +got): %s", tc.inURL, tc.inContent, diff)
				return
			}
		})
	}
}

func TestURLNodeEmpty(t *testing.T) {
	tests := []struct {
		name      string
		inURL     string
		inContent []Node
		out       bool
	}{
		{
			name: "Empty",
			out:  true,
		},
		{
			name:  "NonEmpty",
			inURL: "google.com",
			inContent: []Node{
				NewTextNode(NewTextNodeOptions{Value: "foo"}),
				NewTextNode(NewTextNodeOptions{Value: "bar"}),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewURLNode(tc.inURL, tc.inContent...)
			out := n.Empty()
			if out != tc.out {
				t.Errorf("URLNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
