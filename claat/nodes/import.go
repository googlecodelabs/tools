package nodes

// NewImportNode creates a new Node of type NodeImport,
// with initialized ImportNode.Content.
func NewImportNode(url string) *ImportNode {
	return &ImportNode{
		node:    node{typ: NodeImport},
		Content: NewListNode(),
		URL:     url,
	}
}

// ImportNode indicates a remote resource available at ImportNode.URL.
type ImportNode struct {
	node
	URL     string
	Content *ListNode
}

// Empty returns the result of in.Content.Empty method.
func (in *ImportNode) Empty() bool {
	return in.Content.Empty()
}

// MutateBlock mutates both in's block marker and that of in.Content.
func (in *ImportNode) MutateBlock(v interface{}) {
	in.node.MutateBlock(v)
	in.Content.MutateBlock(v)
}

// ImportNodes extracts everything except NodeImport nodes, recursively.
func ImportNodes(nodes []Node) []*ImportNode {
	var imps []*ImportNode
	for _, n := range nodes {
		switch n := n.(type) {
		case *ImportNode:
			imps = append(imps, n)
		case *ListNode:
			imps = append(imps, ImportNodes(n.Nodes)...)
		case *InfoboxNode:
			imps = append(imps, ImportNodes(n.Content.Nodes)...)
		case *GridNode:
			for _, r := range n.Rows {
				for _, c := range r {
					imps = append(imps, ImportNodes(c.Content.Nodes)...)
				}
			}
		}
	}
	return imps
}
