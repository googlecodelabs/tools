package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewImageNode(t *testing.T) {
	tests := []struct {
		name   string
		inOpts NewImageNodeOptions
		out    *ImageNode
	}{
		{
			name: "Empty",
			out: &ImageNode{
				node: node{typ: NodeImage},
			},
		},
		{
			name: "NonEmpty",
			inOpts: NewImageNodeOptions{
				Src:   "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
				Width: 1.0,
				Title: "foo",
				Alt:   "bar",
			},
			out: &ImageNode{
				node:  node{typ: NodeImage},
				Src:   "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
				Width: 1.0,
				Title: "foo",
				Alt:   "bar",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewImageNode(tc.inOpts)
			if diff := cmp.Diff(tc.out, out, cmp.AllowUnexported(ImageNode{}, node{})); diff != "" {
				t.Errorf("NewImageNode(%+v) got diff (-want +got): %s", tc.inOpts, diff)
				return
			}
		})
	}
}

func TestImageNodeEmpty(t *testing.T) {
	tests := []struct {
		name   string
		inOpts NewImageNodeOptions
		out    bool
	}{
		{
			name: "Empty",
			out:  true,
		},
		{
			name: "NonEmpty",
			inOpts: NewImageNodeOptions{
				Src: "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
			},
		},
		{
			name: "EmptyWithSpaces",
			inOpts: NewImageNodeOptions{
				Src: "\n \t",
			},
			out: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewImageNode(tc.inOpts)
			out := n.Empty()
			if out != tc.out {
				t.Errorf("ImageNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}

func TestImageNodes(t *testing.T) {
	a1 := NewImageNode(NewImageNodeOptions{Src: "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"})
	a2 := NewImageNode(NewImageNodeOptions{Src: "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"})
	a3 := NewImageNode(NewImageNodeOptions{Src: "https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png"})

	b1 := NewItemsListNode("", 1)
	b1.Items = append(b1.Items, NewListNode(a1, a2, NewTextNode(NewTextNodeOptions{Value: "foobar"}), a3))

	c1 := NewGridNode(
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

	d1 := NewURLNode("google.com", NewTextNode(NewTextNodeOptions{Value: "aaa"}), a1, NewTextNode(NewTextNodeOptions{Value: "bbb"}))
	d2 := NewButtonNode(true, true, true, d1)
	d3 := NewURLNode("google.com", a2)
	d4 := NewHeaderNode(2, d3)
	d5 := NewInfoboxNode(InfoboxNegative, d4)
	d6 := NewListNode(d2, d5, a3)

	tests := []struct {
		name    string
		inNodes []Node
		out     []*ImageNode
	}{
		{
			name:    "JustImage",
			inNodes: []Node{a1},
			out:     []*ImageNode{a1},
		},
		{
			name:    "Multiple",
			inNodes: []Node{a1, a2, a3},
			out:     []*ImageNode{a1, a2, a3},
		},
		{
			name: "List",
			inNodes: []Node{
				NewListNode(a1, NewTextNode(NewTextNodeOptions{Value: "foobar"}), a2, a3),
			},
			out: []*ImageNode{a1, a2, a3},
		},
		{
			name:    "ItemsList",
			inNodes: []Node{b1},
			out:     []*ImageNode{a1, a2, a3},
		},
		{
			name: "Header",
			inNodes: []Node{
				NewHeaderNode(1, a1, NewTextNode(NewTextNodeOptions{Value: "foobar"})),
			},
			out: []*ImageNode{a1},
		},
		{
			name: "URL",
			inNodes: []Node{
				NewURLNode("google.com", a2, NewTextNode(NewTextNodeOptions{Value: "foobar"}), a3),
			},
			out: []*ImageNode{a2, a3},
		},
		{
			name: "Button",
			inNodes: []Node{
				NewButtonNode(true, true, true, NewTextNode(NewTextNodeOptions{Value: "foobar"}), a3, a1),
			},
			out: []*ImageNode{a3, a1},
		},
		{
			name: "Infobox",
			inNodes: []Node{
				NewInfoboxNode(InfoboxPositive, a2, a1, NewTextNode(NewTextNodeOptions{Value: "foobar"})),
			},
			out: []*ImageNode{a2, a1},
		},
		{
			name:    "Grid",
			inNodes: []Node{c1},
			out:     []*ImageNode{a1, a3, a2},
		},
		{
			name: "Text",
			inNodes: []Node{
				NewTextNode(NewTextNodeOptions{Value: "foo"}),
				NewTextNode(NewTextNodeOptions{Value: "bar"}),
			},
		},
		{
			name:    "NontrivialStructure",
			inNodes: []Node{d6},
			out:     []*ImageNode{a1, a2, a3},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := ImageNodes(tc.inNodes)
			if diff := cmp.Diff(tc.out, out, cmp.AllowUnexported(ImageNode{}, node{})); diff != "" {
				t.Errorf("ImageNodes(%+v) got diff (-want +got): %s", tc.inNodes, diff)
				return
			}
		})
	}
}
