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

package parser

import (
	"strings"
	"unicode"

	"github.com/googlecodelabs/tools/claat/nodes"
)

// blockSquashable returns true if a node of type t can be squash in a block.
func blockSquashable(n nodes.Node) bool {
	if n.Block() == nil {
		return false
	}
	return nodes.IsInline(n.Type())
}

func squashHeadBlock(nodesToSquash []nodes.Node) (squash, remainder []nodes.Node) {
	first := nodesToSquash[0]
	if !blockSquashable(first) {
		return nodesToSquash[:1], nodesToSquash[1:]
	}
	hnodes := []nodes.Node{first}
	for _, n := range nodesToSquash[1:] {
		if !blockSquashable(n) || n.Block() != first.Block() {
			break
		}
		hnodes = append(hnodes, n)
	}
	next := nodesToSquash[len(hnodes):]
	hnodes = trimNodes(hnodes)
	if len(hnodes) == 0 {
		return nil, next
	}
	head := nodes.NewListNode(hnodes...)
	head.MutateBlock(true)
	head.MutateEnv(first.Env())
	return []nodes.Node{head}, next
}

func trimNodes(nodesToTrim []nodes.Node) []nodes.Node {
	trim := make([]nodes.Node, 0, len(nodesToTrim))
	for i, n := range nodesToTrim {
		if n.Type() == nodes.NodeCode && i == 0 {
			cn := n.(*nodes.CodeNode)
			cn.Value = strings.TrimLeft(cn.Value, "\n")
		}
		if !n.Empty() || len(trim) > 0 {
			trim = append(trim, n)
			continue
		}
	}
	return trim
}

func concatNodes(a, b nodes.Node) bool {
	switch {
	case a.Type() == nodes.NodeText && b.Type() == nodes.NodeText:
		return concatText(a, b)
	case a.Type() == nodes.NodeCode && b.Type() == nodes.NodeCode:
		return concatCode(a, b)
	case a.Type() == nodes.NodeCode && b.Type() == nodes.NodeText:
		t := b.(*nodes.TextNode)
		if strings.TrimSpace(t.Value) == "" {
			return true
		}
	case a.Type() == nodes.NodeURL && b.Type() == nodes.NodeURL:
		return concatURL(a, b)
	case nodes.IsItemsList(a.Type()) && nodes.IsItemsList(b.Type()):
		return concatItemsList(a, b)
	}
	return false
}

func concatItemsList(a, b nodes.Node) bool {
	l1 := a.(*nodes.ItemsListNode)
	l2 := b.(*nodes.ItemsListNode)
	if l1.ListType != l2.ListType {
		return false
	}
	if l1.ListType != "" && l1.Start > 0 && l2.Start > 0 && l2.Start-len(l1.Items) != 1 {
		return false
	}
	l1.Items = append(l1.Items, l2.Items...)
	return true
}

func concatText(a, b nodes.Node) bool {
	t1 := a.(*nodes.TextNode)
	t2 := b.(*nodes.TextNode)

	if t1.Block() != t2.Block() {
		return false
	}
	v1, sp1 := splitSpaceRight(t1.Value)
	v2, sp2 := splitSpaceLeft(t2.Value)
	// <code+spaces><non-code>
	if t1.Code && !t2.Code && sp1 != "" {
		t1.Value = v1
		t2.Value = sp1 + t2.Value
		return false
	}
	// <non-code><spaces+code>
	if !t1.Code && t2.Code && sp2 != "" {
		t2.Value = v2
		t1.Value += sp2
		return false
	}
	// <non-code><spaces>
	if !t1.Code && strings.TrimSpace(t2.Value) == "" {
		t1.Value += t2.Value
		return true
	}
	// different text styles: bold, italic or code
	if t1.Code != t2.Code || t1.Bold != t2.Bold || t1.Italic != t2.Italic {
		return false
	}
	// everything else can be concatenated
	t1.Value += t2.Value
	return true
}

func concatCode(a, b nodes.Node) bool {
	c1 := a.(*nodes.CodeNode)
	c2 := b.(*nodes.CodeNode)
	if c1.Block() != c2.Block() || c1.Term != c2.Term || c1.Lang != c2.Lang {
		return false
	}
	c1.Value += c2.Value
	return true
}

func concatURL(a, b nodes.Node) bool {
	u1 := a.(*nodes.URLNode)
	u2 := b.(*nodes.URLNode)
	if u1.Block() != u2.Block() || u1.URL != u2.URL || u1.Name != u2.Name {
		return false
	}
	u1.Content.Append(u2.Content.Nodes...)
	u1.Content.Nodes = CompactNodes(u1.Content.Nodes)
	return true
}

func splitSpaceLeft(s string) (v string, sp string) {
	for i, r := range s {
		if !unicode.IsSpace(r) {
			return s[i:], s[:i]
		}
	}
	return "", s
}

func splitSpaceRight(s string) (v string, sp string) {
	rs := []rune(s)
	for i := len(rs) - 1; i >= 0; i-- {
		if !unicode.IsSpace(rs[i]) {
			return string(rs[:i+1]), string(rs[i+1:])
		}
	}
	return "", string(rs)
}

func requiresSpacer(a, b nodes.Node) bool {
	t1, ok1 := a.(*nodes.TextNode)
	t2, ok2 := b.(*nodes.TextNode)

	if !(ok1 && ok2) {
		return false
	}

	return (t1.Bold && t2.Bold) || (t1.Italic && t2.Italic)
}

// nodeBlocks encapsulates all nodes of the same block into a new ListNode
// with its B field set to true.
// Nodes which are not blockSquashable remain as is.
func BlockNodes(nodesToEncapsulate []nodes.Node) []nodes.Node {
	var blocks []nodes.Node
	for {
		if len(nodesToEncapsulate) == 0 {
			break
		}
		var head []nodes.Node
		head, nodesToEncapsulate = squashHeadBlock(nodesToEncapsulate)
		blocks = append(blocks, head...)
	}
	return blocks
}

// Although the input slice is not modified, its elements are.
func CompactNodes(nodesToCompact []nodes.Node) []nodes.Node {
	res := make([]nodes.Node, 0, len(nodesToCompact))
	var last nodes.Node
	for _, n := range nodesToCompact {
		switch {
		case n.Type() == nodes.NodeList:
			l := n.(*nodes.ListNode)
			l.Nodes = CompactNodes(l.Nodes)
		case nodes.IsItemsList(n.Type()):
			l := n.(*nodes.ItemsListNode)
			for _, it := range l.Items {
				it.Nodes = CompactNodes(it.Nodes)
			}
		}
		if last == nil || !concatNodes(last, n) {
			if requiresSpacer(last, n) {
				// Append non-breaking zero-width space.
				res = append(res, nodes.NewTextNode(nodes.NewTextNodeOptions{Value: string('\uFEFF')}))
			}
			res = append(res, n)

			if n.Type() == nodes.NodeCode {
				c := n.(*nodes.CodeNode)
				c.Value = strings.TrimLeft(c.Value, "\n")
			}

			last = n
		}
	}
	return res
}
