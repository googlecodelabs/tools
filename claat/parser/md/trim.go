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

package md

import (
	"strings"
	"unicode"

	"github.com/googlecodelabs/tools/claat/types"
)

// blockSquashable returns true if a node of type t can be squash in a block.
func blockSquashable(n types.Node) bool {
	if n.Block() == nil {
		return false
	}
	return types.IsInline(n.Type())
}

func squashHeadBlock(nodes []types.Node) (squash, remainder []types.Node) {
	first := nodes[0]
	if !blockSquashable(first) {
		return nodes[:1], nodes[1:]
	}
	hnodes := []types.Node{first}
	for _, n := range nodes[1:] {
		if !blockSquashable(n) || n.Block() != first.Block() {
			break
		}
		hnodes = append(hnodes, n)
	}
	next := nodes[len(hnodes):]
	hnodes = trimNodes(hnodes)
	if len(hnodes) == 0 {
		return nil, next
	}
	head := types.NewListNode(hnodes...)
	head.MutateBlock(true)
	head.MutateEnv(first.Env())
	return []types.Node{head}, next
}

func trimNodes(nodes []types.Node) []types.Node {
	trim := make([]types.Node, 0, len(nodes))
	for i, n := range nodes {
		if n.Type() == types.NodeCode && i == 0 {
			cn := n.(*types.CodeNode)
			cn.Value = strings.TrimLeft(cn.Value, "\n")
		}
		if !n.Empty() || len(trim) > 0 {
			trim = append(trim, n)
			continue
		}
	}
	return trim
}

func concatNodes(a, b types.Node) bool {
	switch {
	case a.Type() == types.NodeText && b.Type() == types.NodeText:
		return concatText(a, b)
	case a.Type() == types.NodeCode && b.Type() == types.NodeCode:
		return concatCode(a, b)
	case a.Type() == types.NodeCode && b.Type() == types.NodeText:
		t := b.(*types.TextNode)
		if strings.TrimSpace(t.Value) == "" {
			return true
		}
	case a.Type() == types.NodeURL && b.Type() == types.NodeURL:
		return concatURL(a, b)
	case types.IsItemsList(a.Type()) && types.IsItemsList(b.Type()):
		return concatItemsList(a, b)
	}
	return false
}

func concatItemsList(a, b types.Node) bool {
	l1 := a.(*types.ItemsListNode)
	l2 := b.(*types.ItemsListNode)
	if l1.ListType != l2.ListType {
		return false
	}
	if l1.ListType != "" && l1.Start > 0 && l2.Start > 0 && l2.Start-len(l1.Items) != 1 {
		return false
	}
	l1.Items = append(l1.Items, l2.Items...)
	return true
}

func concatText(a, b types.Node) bool {
	t1 := a.(*types.TextNode)
	t2 := b.(*types.TextNode)

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

func concatCode(a, b types.Node) bool {
	c1 := a.(*types.CodeNode)
	c2 := b.(*types.CodeNode)
	if c1.Block() != c2.Block() || c1.Term != c2.Term || c1.Lang != c2.Lang {
		return false
	}
	c1.Value += c2.Value
	return true
}

func concatURL(a, b types.Node) bool {
	u1 := a.(*types.URLNode)
	u2 := b.(*types.URLNode)
	if u1.Block() != u2.Block() || u1.URL != u2.URL || u1.Name != u2.Name {
		return false
	}
	u1.Content.Append(u2.Content.Nodes...)
	return true
}

func splitSpaceLeft(s string) (v string, sp string) {
	for i, r := range s {
		if !unicode.IsSpace(r) {
			return s[i:], s[:i]
		}
	}
	return s, ""
}

func splitSpaceRight(s string) (v string, sp string) {
	rs := []rune(s)
	for i := len(rs) - 1; i >= 0; i-- {
		if !unicode.IsSpace(rs[i]) {
			return s[:i+1], s[i+1:]
		}
	}
	return s, ""
}

// nodeBlocks encapsulates all nodes of the same block into a new ListNode
// with its B field set to true.
// Nodes which are not blockSquashable remain as is.
func blockNodes(nodes []types.Node) []types.Node {
	var blocks []types.Node
	for {
		if len(nodes) == 0 {
			break
		}
		var head []types.Node
		head, nodes = squashHeadBlock(nodes)
		blocks = append(blocks, head...)
	}
	return blocks
}

// Although nodes slice is not modified, its elements are.
func compactNodes(nodes []types.Node) []types.Node {
	res := make([]types.Node, 0, len(nodes))
	var last types.Node
	for _, n := range nodes {
		switch {
		case n.Type() == types.NodeList:
			l := n.(*types.ListNode)
			l.Nodes = compactNodes(l.Nodes)
		case types.IsItemsList(n.Type()):
			l := n.(*types.ItemsListNode)
			for _, it := range l.Items {
				it.Nodes = compactNodes(it.Nodes)
			}
		}
		if last == nil || !concatNodes(last, n) {
			last = n
			res = append(res, n)

			if n.Type() == types.NodeCode {
				c := n.(*types.CodeNode)
				c.Value = strings.TrimLeft(c.Value, "\n")
			}
		}
	}
	return res
}
