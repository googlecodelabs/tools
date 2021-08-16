package nodes

// NewURLNode creates a new Node of type NodeURL with optional content n.
func NewURLNode(url string, n ...Node) *URLNode {
	return &URLNode{
		node:    node{typ: NodeURL},
		URL:     url,
		Target:  "_blank",
		Content: NewListNode(n...),
	}
}

// URLNode represents elements such as <a href="...">
type URLNode struct {
	node
	URL     string
	Name    string
	Target  string
	Content *ListNode
}

// Empty returns true if un content is empty.
func (un *URLNode) Empty() bool {
	return un.Content.Empty()
}
