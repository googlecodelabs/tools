package nodes

import "strings"

// NewTextNode creates a new Node of type NodeText.
func NewTextNode(v string) *TextNode {
	return &TextNode{
		node:  node{typ: NodeText},
		Value: v,
	}
}

// TextNode is a simple node containing text as a string value.
type TextNode struct {
	node
	Bold   bool
	Italic bool
	Code   bool
	Value  string
}

// Empty returns true if tn.Value is zero, excluding space runes.
func (tn *TextNode) Empty() bool {
	return strings.TrimSpace(tn.Value) == ""
}
