package nodes

import "strings"

// NewCodeNode creates a new Node of type NodeCode.
// Use term argument to specify a terminal output.
func NewCodeNode(v string, term bool, lang string) *CodeNode {
	return &CodeNode{
		node:  node{typ: NodeCode},
		Value: v,
		Term:  term,
		Lang:  lang,
	}
}

// CodeNode is either a source code snippet or a terminal output.
// TODO is there any room to consolidate Term and Lang?
type CodeNode struct {
	node
	Term  bool
	Lang  string
	Value string
}

// Empty returns true if cn.Value is zero, exluding space runes.
func (cn *CodeNode) Empty() bool {
	return strings.TrimSpace(cn.Value) == ""
}
