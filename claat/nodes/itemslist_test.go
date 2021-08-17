package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var cmpOptItemsList = cmp.AllowUnexported(ItemsListNode{}, node{})

func TestNewItemsListNode(t *testing.T) {
	// Only one code path, so this is not a tabular test.
	inTyp := "foobar"
	inStart := 5
	got := NewItemsListNode(inTyp, inStart)
	want := &ItemsListNode{
		node: node{
			typ:   NodeItemsList,
			block: true,
		},
		ListType: "foobar",
		Start:    5,
	}

	if diff := cmp.Diff(want, got, cmpOptItemsList); diff != "" {
		t.Errorf("NewItemsListNode(%q, %d) got diff (-want +got): %s", inTyp, inStart, diff)
	}
}

func TestItemsListNodeEmpty(t *testing.T) {
	a := NewItemsListNode("foobar", 0)
	a.Items = append(a.Items, NewListNode(NewTextNode("")))

	b := NewItemsListNode("foobar", 0)
	b.Items = append(b.Items, NewListNode(NewTextNode("a")))

	c := NewItemsListNode("foobar", 0)
	c.Items = append(c.Items, NewListNode(NewTextNode("")), NewListNode(NewTextNode("")), NewListNode(NewTextNode("")))

	d := NewItemsListNode("foobar", 0)
	d.Items = append(d.Items, NewListNode(NewTextNode("a")), NewListNode(NewTextNode("")), NewListNode(NewTextNode("b")))

	e := NewItemsListNode("foobar", 0)
	e.Items = append(e.Items, NewListNode(NewTextNode("a")), NewListNode(NewTextNode("b")), NewListNode(NewTextNode("c")))

	tests := []struct {
		name   string
		inNode *ItemsListNode
		out    bool
	}{
		{
			name:   "Zero",
			inNode: NewItemsListNode("foobar", 0),
			out:    true,
		},
		{
			name:   "OneEmpty",
			inNode: a,
			out:    true,
		},
		{
			name:   "OneNonEmpty",
			inNode: b,
		},
		{
			name:   "MultipleEmpty",
			inNode: c,
			out:    true,
		},
		{
			name:   "MultipleSomeNonEmpty",
			inNode: d,
		},
		{
			name:   "MultipleAllNonEmpty",
			inNode: e,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := tc.inNode.Empty()
			if out != tc.out {
				t.Errorf("ItemsListNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
