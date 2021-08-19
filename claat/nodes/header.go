package nodes

// NewHeaderNode creates a new HeaderNode with optional content nodes n.
func NewHeaderNode(level int, n ...Node) *HeaderNode {
	return &HeaderNode{
		node:    node{typ: NodeHeader},
		Level:   level,
		Content: NewListNode(n...),
	}
}

// HeaderNode is any regular header, a checklist header, or an FAQ header.
type HeaderNode struct {
	node
	Level   int
	Content *ListNode
}

// Empty returns true if header content is empty.
func (hn *HeaderNode) Empty() bool {
	return hn.Content.Empty()
}

// IsHeader returns true if t is one of header types.
func IsHeader(t NodeType) bool {
	return t&(NodeHeader|NodeHeaderCheck|NodeHeaderFAQ) != 0
}

// MutateType sets the header's node type if the given type is a header type.
func (hn *HeaderNode) MutateType(t NodeType) {
	if IsHeader(t) {
		hn.typ = t
	}
}
