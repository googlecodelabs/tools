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
	a.Content.Nodes = append(a.Content.Nodes, NewTextNode(NewTextNodeOptions{Value: "a"}))
	b := NewImportNode("foobar")
	b.Content.Nodes = append(b.Content.Nodes, NewTextNode(NewTextNodeOptions{Value: "b"}))
	c := NewImportNode("foobar")
	c.Content.Nodes = append(c.Content.Nodes, NewTextNode(NewTextNodeOptions{Value: ""}))

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
				t.Errorf("ImportNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}

func TestImportNodeMutateBlock(t *testing.T) {
	n := NewImportNode("")
	mValue := "foobar"

	n.MutateBlock(mValue)

	if n.node.block != mValue {
		t.Errorf("ImportNode.node.block = %+v, want %q", n.node.block, mValue)
	}
	if n.Content.node.block != mValue {
		t.Errorf("ImportNode.Content.node.block = %+v, want %q", n.Content.node.block, mValue)
	}
}

func TestImportNodes(t *testing.T) {
	a1 := NewImportNode("google.com")
	a2 := NewImportNode("youtube.com")
	a3 := NewImportNode("google.com/calendar")

	b1 := NewGridNode(
		[]*GridCell{
			&GridCell{
				Rowspan: 1,
				Colspan: 1,
				Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "aaa"}), NewTextNode(NewTextNodeOptions{Value: "bbb"})),
			},
			&GridCell{
				Rowspan: 1,
				Colspan: 1,
				Content: NewListNode(a1, NewTextNode(NewTextNodeOptions{Value: "ccc"})),
			},
		},
		[]*GridCell{
			&GridCell{
				Rowspan: 1,
				Colspan: 1,
				Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "ddd"}), a3),
			},
			&GridCell{
				Rowspan: 1,
				Colspan: 1,
				Content: NewListNode(a2, NewTextNode(NewTextNodeOptions{Value: "eee"})),
			},
		},
	)

	c1 := NewInfoboxNode(InfoboxNegative, a1, NewTextNode(NewTextNodeOptions{Value: "foobar"}))
	c2 := NewListNode(a2, NewButtonNode(false, false, false, NewTextNode(NewTextNodeOptions{Value: "foobar"})))
	c3 := NewListNode(c1, c2, a3)

	tests := []struct {
		name    string
		inNodes []Node
		out     []*ImportNode
	}{
		{
			name:    "JustImport",
			inNodes: []Node{a1},
			out:     []*ImportNode{a1},
		},
		{
			name:    "Multiple",
			inNodes: []Node{a1, NewTextNode(NewTextNodeOptions{Value: "foo"}), a2, NewTextNode(NewTextNodeOptions{Value: "bar"}), a3},
			out:     []*ImportNode{a1, a2, a3},
		},
		{
			name:    "List",
			inNodes: []Node{NewListNode(a1, a2, a3)},
			out:     []*ImportNode{a1, a2, a3},
		},
		{
			name:    "Infobox",
			inNodes: []Node{NewInfoboxNode(InfoboxPositive, a1, a2, NewTextNode(NewTextNodeOptions{Value: "foobar"}), a3)},
			out:     []*ImportNode{a1, a2, a3},
		},
		{
			name:    "Grid",
			inNodes: []Node{b1},
			out:     []*ImportNode{a1, a3, a2},
		},
		{
			name:    "Button",
			inNodes: []Node{NewButtonNode(true, true, true, a3, a2, a1)},
		},
		{
			name:    "Text",
			inNodes: []Node{NewTextNode(NewTextNodeOptions{Value: "foobar"})},
		},
		{
			name:    "NontrivialStructure",
			inNodes: []Node{c3},
			out:     []*ImportNode{a1, a2, a3},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := ImportNodes(tc.inNodes)
			if diff := cmp.Diff(tc.out, out, cmpOptImport); diff != "" {
				t.Errorf("ImportNodes(%+v) got diff (-want +got): %s", tc.inNodes, diff)
				return
			}
		})
	}
}
