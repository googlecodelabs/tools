package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var cmpOptImport = cmp.AllowUnexported(ImportNode{}, node{}, ListNode{}, TextNode{})

func TestNewImportNode(t *testing.T) {
	tests := []struct {
		name  string
		inURL string
		out   *ImportNode
	}{
		{
			name: "Empty",
			out: &ImportNode{
				node:    node{typ: NodeImport},
				Content: NewListNode(),
			},
		},
		{
			name:  "HasURL",
			inURL: "google.com",
			out: &ImportNode{
				node:    node{typ: NodeImport},
				URL:     "google.com",
				Content: NewListNode(),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewImportNode(tc.inURL)
			if diff := cmp.Diff(tc.out, out, cmpOptImport); diff != "" {
				t.Errorf("NewImportNode(%q) got diff (-want +got): %s", tc.inURL, diff)
				return
			}
		})
	}
}

func TestImportNodeEmpty(t *testing.T) {
	a := NewImportNode("")
	a.Content.Nodes = append(a.Content.Nodes, NewTextNode("a"))
	b := NewImportNode("foobar")
	b.Content.Nodes = append(b.Content.Nodes, NewTextNode("b"))
	c := NewImportNode("foobar")
	c.Content.Nodes = append(c.Content.Nodes, NewTextNode(""))

	tests := []struct {
		name   string
		inNode *ImportNode
		out    bool
	}{
		{
			name:   "EmptyNoURL",
			inNode: NewImportNode(""),
			out:    true,
		},
		{
			name:   "EmptyWithURL",
			inNode: NewImportNode("google.com"),
			out:    true,
		},
		{
			name:   "NonEmptyNoURL",
			inNode: a,
		},
		{
			name:   "NonEmptyWithURL",
			inNode: b,
		},
		{
			name:   "EmptyWithContent",
			inNode: c,
			out:    true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := tc.inNode.Empty()
			if out != tc.out {
				t.Errorf("ImageNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
