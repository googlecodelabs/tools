package nodes

// TODO this is a long arg signature. Maybe use an options type?
// NewButtonNode creates a new button with optional content nodes n.
func NewButtonNode(raise, color, download bool, n ...Node) *ButtonNode {
	return &ButtonNode{
		node:     node{typ: NodeButton},
		Raise:    raise,
		Color:    color,
		Download: download,
		Content:  NewListNode(n...),
	}
}

// ButtonNode represents a button, e.g. "Download Zip".
type ButtonNode struct {
	node
	Raise    bool
	Color    bool
	Download bool
	Content  *ListNode
}

// Empty returns true if its content is empty.
func (bn *ButtonNode) Empty() bool {
	return bn.Content.Empty()
}
