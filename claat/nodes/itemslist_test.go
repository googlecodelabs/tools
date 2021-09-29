package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var cmpOptItemsList = cmp.AllowUnexported(ItemsListNode{}, node{}, ListNode{}, TextNode{})

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
	a.Items = append(a.Items, NewListNode(NewTextNode(NewTextNodeOptions{Value: ""})))

	b := NewItemsListNode("foobar", 0)
	b.Items = append(b.Items, NewListNode(NewTextNode(NewTextNodeOptions{Value: "a"})))

	c := NewItemsListNode("foobar", 0)
	c.Items = append(c.Items, NewListNode(NewTextNode(NewTextNodeOptions{Value: ""})), NewListNode(NewTextNode(NewTextNodeOptions{Value: ""})), NewListNode(NewTextNode(NewTextNodeOptions{Value: ""})))

	d := NewItemsListNode("foobar", 0)
	d.Items = append(d.Items, NewListNode(NewTextNode(NewTextNodeOptions{Value: "a"})), NewListNode(NewTextNode(NewTextNodeOptions{Value: ""})), NewListNode(NewTextNode(NewTextNodeOptions{Value: "b"})))

	e := NewItemsListNode("foobar", 0)
	e.Items = append(e.Items, NewListNode(NewTextNode(NewTextNodeOptions{Value: "a"})), NewListNode(NewTextNode(NewTextNodeOptions{Value: "b"})), NewListNode(NewTextNode(NewTextNodeOptions{Value: "c"})))

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

func TestItemsListNewItem(t *testing.T) {
	// Only one code path, so this is not a tabular test.
	a := NewTextNode(NewTextNodeOptions{Value: "a"})
	b := NewTextNode(NewTextNodeOptions{Value: "b"})
	c := NewTextNode(NewTextNodeOptions{Value: "c"})

	iln := NewItemsListNode("foobar", 0)

	got := iln.NewItem(a, b, c)
	want := NewListNode(a, b, c)

	if diff := cmp.Diff(want, got, cmpOptItemsList); diff != "" {
		t.Errorf("ItemsListNode.NewItem() got diff (-want +got): %s", diff)
	}

	wantItemsListNode := &ItemsListNode{
		node: node{
			typ:   NodeItemsList,
			block: true,
		},
		ListType: "foobar",
		Items: []*ListNode{
			&ListNode{
				node:  node{typ: NodeList},
				Nodes: []Node{a, b, c},
			},
		},
	}
	if diff := cmp.Diff(wantItemsListNode, iln, cmpOptItemsList); diff != "" {
		t.Errorf("ItemsListNode after NewItem got diff ((-want +got): %s", diff)
	}
}

func TestItemsListMutateType(t *testing.T) {
	tests := []struct {
		name   string
		inType NodeType
		out    NodeType
	}{
		{
			name:   "ItemsList",
			inType: NodeItemsList,
			out:    NodeItemsList,
		},
		{
			name:   "AlternateItemsListType",
			inType: NodeItemsCheck,
			out:    NodeItemsCheck,
		},
		{
			name:   "NotAItemsList",
			inType: NodeButton,
			out:    NodeItemsList,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewItemsListNode("foobar", 0) // Args chosen arbitrarily.
			n.MutateType(tc.inType)
			if n.typ != tc.out {
				t.Errorf("ItemsListNode.typ after MutateType = %v, want %v", n.typ, tc.out)
				return
			}
		})
	}
}
