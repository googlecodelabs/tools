// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nodes

// TODO be consistent between Node/*Node

import (
	"sort"
)

// NodeType is type for parsed codelab nodes tree.
type NodeType uint32

// Codelab node kinds.
const (
	NodeInvalid     NodeType = 1 << iota
	NodeList                 // A node which contains a list of other nodes
	NodeGrid                 // Table
	NodeText                 // Simple node with a string as the value
	NodeCode                 // Source code or console (terminal) output
	NodeInfobox              // An aside box for notes or warnings
	NodeSurvey               // Sets of grouped questions
	NodeURL                  // Represents elements such as <a href="...">
	NodeImage                // Image
	NodeButton               // Button
	NodeItemsList            // Set of NodeList items
	NodeItemsCheck           // Special kind of NodeItemsList, checklist
	NodeItemsFAQ             // Special kind of NodeItemsList, FAQ
	NodeHeader               // A header text node
	NodeHeaderCheck          // Special kind of header, checklist
	NodeHeaderFAQ            // Special kind of header, FAQ
	NodeYouTube              // YouTube video
	NodeIframe               // Embedded iframe
	NodeImport               // A node which holds content imported from another resource
)

// Node is an interface common to all node types.
type Node interface {
	// Type returns node type.
	Type() NodeType
	// MutateType changes node type where possible.
	// Only changes within this same category are allowed.
	// For instance, items list or header nodes can change their types
	// to another kind of items list or header.
	MutateType(NodeType)
	// Block returns a source reference of the node.
	Block() interface{}
	// MutateBlock updates source reference of the node.
	MutateBlock(interface{})
	// Empty returns true if the node has no content.
	Empty() bool
	// Env returns node environment
	Env() []string
	// MutateEnv replaces current node environment tags with env.
	MutateEnv(env []string)
}

// IsItemsList returns true if t is one of ItemsListNode types.
func IsItemsList(t NodeType) bool {
	return t&(NodeItemsList|NodeItemsCheck|NodeItemsFAQ) != 0
}

// IsHeader returns true if t is one of header types.
func IsHeader(t NodeType) bool {
	return t&(NodeHeader|NodeHeaderCheck|NodeHeaderFAQ) != 0
}

// IsInline returns true if t is an inline node type.
func IsInline(t NodeType) bool {
	return t&(NodeText|NodeURL|NodeImage|NodeButton) != 0
}

// EmptyNodes returns true if all of nodes are empty.
func EmptyNodes(nodes []Node) bool {
	for _, n := range nodes {
		if !n.Empty() {
			return false
		}
	}
	return true
}

type node struct {
	typ   NodeType
	block interface{}
	env   []string
}

func (b *node) Type() NodeType {
	return b.typ
}

func (b *node) MutateType(t NodeType) {
	if IsItemsList(b.typ) && IsItemsList(t) || IsHeader(b.typ) && IsHeader(t) {
		b.typ = t
	}
}

func (b *node) Block() interface{} {
	return b.block
}

func (b *node) MutateBlock(v interface{}) {
	b.block = v
}

func (b *node) Env() []string {
	return b.env
}

func (b *node) MutateEnv(e []string) {
	b.env = make([]string, len(e))
	copy(b.env, e)
	sort.Strings(b.env)
}

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

// Append appends nodes n to the end of l.Nodes slice.
func (l *ListNode) Append(n ...Node) {
	l.Nodes = append(l.Nodes, n...)
}

// Prepend prepends nodes n at the beginning of l.Nodes slice.
func (l *ListNode) Prepend(n ...Node) {
	l.Nodes = append(n, l.Nodes...)
}

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

// NewGridNode creates a new grid with optional content.
func NewGridNode(rows ...[]*GridCell) *GridNode {
	return &GridNode{
		node: node{typ: NodeGrid},
		Rows: rows,
	}
}

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

// NewItemsListNode creates a new ItemsListNode of type NodeItemsList,
// which defaults to an unordered list.
// Provide a positive start to make this a numbered list.
// NodeItemsCheck and NodeItemsFAQ are always unnumbered.
func NewItemsListNode(typ string, start int) *ItemsListNode {
	iln := ItemsListNode{
		node: node{typ: NodeItemsList},
		// TODO document this
		ListType: typ,
		Start:    start,
	}
	iln.MutateBlock(true)
	return &iln
}

// ItemsListNode containts sets of ListNode.
// Non-zero ListType indicates an ordered list.
type ItemsListNode struct {
	node
	ListType string
	Start    int
	Items    []*ListNode
}

// Empty returns true if every item has empty content.
func (il *ItemsListNode) Empty() bool {
	for _, i := range il.Items {
		if !i.Empty() {
			return false
		}
	}
	return true
}

// NewItem creates a new ListNode and adds it to il.Items.
func (il *ItemsListNode) NewItem(nodes ...Node) *ListNode {
	n := NewListNode(nodes...)
	il.Items = append(il.Items, n)
	return n
}
