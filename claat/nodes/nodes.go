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

// Default implementation is a no op.
func (b *node) MutateType(t NodeType) {
	return
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
