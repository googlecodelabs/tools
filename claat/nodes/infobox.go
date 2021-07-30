package nodes

// InfoboxKind defines kind type for InfoboxNode.
type InfoboxKind string

// InfoboxNode variants.
const (
	InfoboxPositive InfoboxKind = "special"
	InfoboxNegative InfoboxKind = "warning"
)

// InfoboxNode is any regular header, a checklist header, or an FAQ header.
type InfoboxNode struct {
	node
	Kind    InfoboxKind
	Content *ListNode
}

// NewInfoboxNode creates a new infobox node with specified kind and optional content.
func NewInfoboxNode(k InfoboxKind, n ...Node) *InfoboxNode {
	return &InfoboxNode{
		node:    node{typ: NodeInfobox},
		Kind:    k,
		Content: NewListNode(n...),
	}
}

// Empty returns true if ib content is empty.
func (ib *InfoboxNode) Empty() bool {
	return ib.Content.Empty()
}
