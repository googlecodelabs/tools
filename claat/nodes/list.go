package nodes

// NewListNode creates a new Node of type NodeList.
func NewListNode(nodes ...Node) *ListNode {
	n := &ListNode{node: node{typ: NodeList}}
	n.Append(nodes...)
	return n
}

// ListNode contains other nodes.
type ListNode struct {
	node
	Nodes []Node
}

// Empty returns true if all l.Nodes are empty.
func (l *ListNode) Empty() bool {
	return EmptyNodes(l.Nodes)
}

// TODO remove
// Append appends nodes n to the end of l.Nodes slice.
func (l *ListNode) Append(n ...Node) {
	l.Nodes = append(l.Nodes, n...)
}
