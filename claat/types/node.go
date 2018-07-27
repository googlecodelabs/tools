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

package types

import (
	"sort"
	"strings"
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
	// AsImageNodes filters and returns all NodeImage nodes recursively
	AsImageNodes() []*ImageNode
	// AsImportNodes filters and returns all ImportNode nodes recursively
	AsImportNodes() []*ImportNode
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

func (b *node) AsImageNodes() []*ImageNode {
	return []*ImageNode{}
}

func (b *node) AsImportNodes() []*ImportNode {
	return []*ImportNode{}
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

// AsImageNodes for ListNode recurses over children
func (l *ListNode) AsImageNodes() []*ImageNode {
	return ImageNodesToList(l.Nodes)
}

// AsImportNodes for ListNode imports the list of Nodes contained
func (l *ListNode) AsImportNodes() []*ImportNode {
	return ImportNodesToList(l.Nodes)
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

// AsImportNodes for ImportNode simply returns a list of itself
func (in *ImportNode) AsImportNodes() []*ImportNode {
	return []*ImportNode{in}
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

// AsImageNodes for GridNode recurses within grid
func (gn *GridNode) AsImageNodes() []*ImageNode {
	var imps []*ImageNode
	for _, r := range gn.Rows {
		for _, c := range r {
			imps = append(imps, ImageNodesToList(c.Content.Nodes)...)
		}
	}
	return imps
}

// AsImportNodes imports all imports in the grid
func (gn *GridNode) AsImportNodes() []*ImportNode {
	var imps []*ImportNode
	for _, r := range gn.Rows {
		for _, c := range r {
			imps = append(imps, ImportNodesToList(c.Content.Nodes)...)
		}
	}
	return imps
}

// NewItemsListNode creates a new ItemsListNode of type NodeItemsList,
// which defaults to an unordered list.
// Provide a positive start to make this a numbered list.
// NodeItemsCheck and NodeItemsFAQ are always unnumbered.
func NewItemsListNode(typ string, start int) *ItemsListNode {
	return &ItemsListNode{
		node:     node{typ: NodeItemsList},
		ListType: typ,
		Start:    start,
	}
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

// AsImageNodes for ItemsListNode recurses in list
func (il *ItemsListNode) AsImageNodes() []*ImageNode {
	var imgs []*ImageNode
	for _, i := range il.Items {
		imgs = append(imgs, ImageNodesToList(i.Nodes)...)
	}
	return imgs
}

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

// NewCodeNode creates a new Node of type NodeCode.
// Use term argument to specify a terminal output.
func NewCodeNode(v string, term bool) *CodeNode {
	return &CodeNode{
		node:  node{typ: NodeCode},
		Value: v,
		Term:  term,
	}
}

// CodeNode is either a source code snippet or a terminal output.
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

// NewHeaderNode creates a new HeaderNode with optional content nodes n.
func NewHeaderNode(level int, n ...Node) *HeaderNode {
	return &HeaderNode{
		node:    node{typ: NodeHeader},
		Level:   level,
		Content: NewListNode(n...),
	}
}

// HeaderNode is any regular header, a checklist header, or an FAQ header.
type HeaderNode struct {
	node
	Level   int
	Content *ListNode
}

// Empty returns true if header content is empty.
func (hn *HeaderNode) Empty() bool {
	return hn.Content.Empty()
}

// AsImageNodes for HeaderNode recurses within Nodes
func (hn *HeaderNode) AsImageNodes() []*ImageNode {
	return ImageNodesToList(hn.Content.Nodes)
}

// NewURLNode creates a new Node of type NodeURL with optinal content n.
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

// AsImageNodes for URLNode recurses within Nodes
func (un *URLNode) AsImageNodes() []*ImageNode {
	return ImageNodesToList(un.Content.Nodes)
}

// NewImageNode creates a new ImageNode  with the give src.
func NewImageNode(src string) *ImageNode {
	return &ImageNode{
		node: node{typ: NodeImage},
		Src:  src,
	}
}

// ImageNode represents a single image.
type ImageNode struct {
	node
	Src      string
	MaxWidth float32
	Alt      string
	Title    string
}

// Empty returns true if its Src is zero, excluding space runes.
func (in *ImageNode) Empty() bool {
	return strings.TrimSpace(in.Src) == ""
}

// AsImageNodes for ImageNode simply returns a list of itself
func (in *ImageNode) AsImageNodes() []*ImageNode {
	return []*ImageNode{in}
}

// NewButtonNode creates a new button with optional content nodes n.
func NewButtonNode(raised, colored, download bool, n ...Node) *ButtonNode {
	return &ButtonNode{
		node:     node{typ: NodeButton},
		Raised:   raised,
		Colored:  colored,
		Download: download,
		Content:  NewListNode(n...),
	}
}

// ButtonNode represents a button, e.g. "Download Zip".
type ButtonNode struct {
	node
	Raised   bool
	Colored  bool
	Download bool
	Content  *ListNode
}

// Empty returns true if its content is empty.
func (bn *ButtonNode) Empty() bool {
	return bn.Content.Empty()
}

// AsImageNodes for ButtonNode recurses within Nodes
func (bn *ButtonNode) AsImageNodes() []*ImageNode {
	return ImageNodesToList(bn.Content.Nodes)
}

// NewSurveyNode creates a new survey node with optional questions.
// If survey is nil, a new empty map will be created.
func NewSurveyNode(id string, groups ...*SurveyGroup) *SurveyNode {
	return &SurveyNode{
		node:   node{typ: NodeSurvey},
		ID:     id,
		Groups: groups,
	}
}

// SurveyNode contains groups of questions. Each group name is the Survey key.
type SurveyNode struct {
	node
	ID     string
	Groups []*SurveyGroup
}

// SurveyGroup contains group name/question and possible answers.
type SurveyGroup struct {
	Name    string
	Options []string
}

// Empty returns true if each group has 0 options.
func (sn *SurveyNode) Empty() bool {
	for _, g := range sn.Groups {
		if len(g.Options) > 0 {
			return false
		}
	}
	return true
}

// NewInfoboxNode creates a new infobox node with specified kind and optional content.
func NewInfoboxNode(k InfoboxKind, n ...Node) *InfoboxNode {
	return &InfoboxNode{
		node:    node{typ: NodeInfobox},
		Kind:    k,
		Content: NewListNode(n...),
	}
}

// InfoboxKind defines kind type for InfoboxNode.
type InfoboxKind string

// InfoboxNode variations.
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

// Empty returns true if ib content is empty.
func (ib *InfoboxNode) Empty() bool {
	return ib.Content.Empty()
}

// AsImageNodes for InfoboxNode recurses within Nodes
func (ib *InfoboxNode) AsImageNodes() []*ImageNode {
	return ImageNodesToList(ib.Content.Nodes)
}

// AsImportNodes imports nodes within the infobox
func (ib *InfoboxNode) AsImportNodes() []*ImportNode {
	return ImportNodesToList(ib.Content.Nodes)
}

// NewYouTubeNode creates a new YouTube video node.
func NewYouTubeNode(vid string) *YouTubeNode {
	return &YouTubeNode{
		node:    node{typ: NodeYouTube},
		VideoID: vid,
	}
}

// YouTubeNode is a YouTube video.
type YouTubeNode struct {
	node
	VideoID string
}

// Empty returns true if yt's VideoID field is zero.
func (yt *YouTubeNode) Empty() bool {
	return yt.VideoID != ""
}

// ImageNodesToList filters out everything except types.NodeImport nodes, recursively.
func ImageNodesToList(nodes []Node) []*ImageNode {
	var imps []*ImageNode
	for _, n := range nodes {
		if list := n.AsImageNodes(); len(list) > 0 {
			imps = append(imps, list...)
		}
	}
	return imps
}

// ImportNodesToList filters out everything except types.NodeImport nodes, recursively.
func ImportNodesToList(nodes []Node) []*ImportNode {
	var imps []*ImportNode
	for _, n := range nodes {
		if list := n.AsImportNodes(); len(list) > 0 {
			imps = append(imps, list...)
		}
	}
	return imps
}
