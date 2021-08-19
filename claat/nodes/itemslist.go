package nodes

// NewItemsListNode creates a new ItemsListNode of type NodeItemsList,
// which defaults to an unordered list.
// Provide a positive start to make this a numbered list.
// NodeItemsCheck and NodeItemsFAQ are always unnumbered.
func NewItemsListNode(typ string, start int) *ItemsListNode {
	iln := ItemsListNode{
		node: node{typ: NodeItemsList},
		// TODO document this
		ListType: typ,
		Start:    start,
	}
	iln.MutateBlock(true)
	return &iln
}

// ItemsListNode containts sets of ListNode.
// Non-zero ListType indicates an ordered list.
type ItemsListNode struct {
	node
	ListType string
	Start    int
	Items    []*ListNode
}

// Empty returns true if every item has empty content.
func (il *ItemsListNode) Empty() bool {
	for _, i := range il.Items {
		if !i.Empty() {
			return false
		}
	}
	return true
}

// NewItem creates a new ListNode and adds it to il.Items.
func (il *ItemsListNode) NewItem(nodes ...Node) *ListNode {
	n := NewListNode(nodes...)
	il.Items = append(il.Items, n)
	return n
}

// IsItemsList returns true if t is one of ItemsListNode types.
func IsItemsList(t NodeType) bool {
	return t&(NodeItemsList|NodeItemsCheck|NodeItemsFAQ) != 0
}

// MutateType sets the items list's node type if the given type is an items list type.
func (il *ItemsListNode) MutateType(t NodeType) {
	if IsItemsList(t) {
		il.typ = t
	}
}
