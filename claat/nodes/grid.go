package nodes

// NewGridNode creates a new grid with optional content.
func NewGridNode(rows ...[]*GridCell) *GridNode {
	return &GridNode{
		node: node{typ: NodeGrid},
		Rows: rows,
	}
}

// TODO define a convenience type for row
// GridNode is a 2d matrix.
type GridNode struct {
	node
	Rows [][]*GridCell
}

// GridCell is a cell of GridNode.
type GridCell struct {
	Colspan int
	Rowspan int
	Content *ListNode
}

// Empty returns true when every cell has empty content.
func (gn *GridNode) Empty() bool {
	for _, r := range gn.Rows {
		for _, c := range r {
			if !c.Content.Empty() {
				return false
			}
		}
	}
	return true
}
