package nodes

import "strings"

type NewTextNodeOptions struct {
	Bold   bool
	Italic bool
	Code   bool
	Value  string
}

// NewTextNode creates a new Node of type NodeText.
func NewTextNode(opts NewTextNodeOptions) *TextNode {
	return &TextNode{
		node:   node{typ: NodeText},
		Bold:   opts.Bold,
		Italic: opts.Italic,
		Code:   opts.Code,
		Value:  opts.Value,
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
