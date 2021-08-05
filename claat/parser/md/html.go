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
	"bytes"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	// headerLevel maps HTML tags to their level in parser.HeaderNode.
	// we -1 as H2 is a new step
	headerLevel = map[atom.Atom]int{
		atom.H3: 2,
		atom.H4: 3,
		atom.H5: 4,
		atom.H6: 5,
	}
)

// isHeader returns true if hn is one of secondary headers.
// Step header is not one of them.
func isHeader(hn *html.Node) bool {
	_, ok := headerLevel[hn.DataAtom]
	return ok
}

// TODO rename, it only captures some meta. Maybe redo the meta system?
func isMeta(hn *html.Node) bool {
	elem := strings.ToLower(hn.Data)
	return strings.HasPrefix(elem, metaDuration+metaSep) || strings.HasPrefix(elem, metaEnvironment+metaSep)
}

func isBold(hn *html.Node) bool {
	if hn.Type == html.TextNode {
		hn = hn.Parent
	} else if hn.DataAtom == atom.Code {
		// Look up as many as 2 levels, to handle the case of e.g. <bold><em><code>
		for i := 0; i < 2; i++ {
			hn = hn.Parent
			if hn.DataAtom == atom.Strong || hn.DataAtom == atom.B {
				return true
			}
		}
		return false
	}
	return hn.DataAtom == atom.Strong ||
		hn.DataAtom == atom.B
}

func isItalic(hn *html.Node) bool {
	if hn.Type == html.TextNode {
		hn = hn.Parent
	} else if hn.DataAtom == atom.Code {
		// Look up as many as 2 levels, to handle the case of e.g. <em><bold><code>
		for i := 0; i < 2; i++ {
			hn = hn.Parent
			if hn.DataAtom == atom.Em || hn.DataAtom == atom.I {
				return true
			}
		}
		return false
	}
	return hn.DataAtom == atom.Em ||
		hn.DataAtom == atom.I
}

// This is different to calling isBold and isItalic separately as we must look
// up an extra level in the tree
func isBoldAndItalic(hn *html.Node) bool {
	if hn.Parent == nil || hn.Parent.Parent == nil {
		return false
	}
	if hn.Type == html.TextNode {
		hn = hn.Parent
	}
	return (isItalic(hn) && isBold(hn.Parent)) || (isItalic(hn.Parent) && isBold(hn))

}

func isConsole(hn *html.Node) bool {
	if hn.Type == html.TextNode {
		hn = hn.Parent
	}
	if hn.DataAtom == atom.Code {
		for _, a := range hn.Attr {
			if a.Key == "class" && a.Val == "language-console" {
				return true
			}
		}
	}
	return false
}

func isCode(hn *html.Node) bool {
	if hn.Type == html.TextNode {
		hn = hn.Parent
	}
	return hn.DataAtom == atom.Code && !isConsole(hn)
}

func isButton(hn *html.Node) bool {
	return hn.DataAtom == atom.Button
}

func isAside(hn *html.Node) bool {
	return hn.DataAtom == atom.Aside
}

func isNewAside(hn *html.Node) bool {
	if hn.DataAtom != atom.Blockquote ||
		hn.FirstChild == nil ||
		hn.FirstChild.NextSibling == nil ||
		hn.FirstChild.NextSibling.FirstChild == nil {
		return false
	}

	asideText := strings.ToLower(hn.FirstChild.NextSibling.FirstChild.Data)
	return strings.HasPrefix(asideText, "aside positive") || strings.HasPrefix(asideText, "aside negative")
}

func isInfobox(hn *html.Node) bool {
	if hn.DataAtom != atom.Dt {
		return false
	}
	return strings.ToLower(hn.FirstChild.Data) == "positive" || isInfoboxNegative(hn)
}

func isInfoboxNegative(hn *html.Node) bool {
	if hn.DataAtom != atom.Dt {
		return false
	}
	return strings.ToLower(hn.FirstChild.Data) == "negative"
}

func isSurvey(hn *html.Node) bool {
	if hn.DataAtom != atom.Form {
		return false
	}
	if findAtom(hn, atom.Name) == nil {
		return false
	}
	if len(findChildAtoms(hn, atom.Input)) == 0 {
		return false
	}
	return true
}

// TODO Write an explanation for why the countTwo checks are necessary.
func isTable(hn *html.Node) bool {
	if hn.DataAtom != atom.Table {
		return false
	}
	// TODO if =1 is fine, can we sub findAtom?
	return countTwo(hn, atom.Tr) >= 1 || countTwo(hn, atom.Td) >= 1
}

func isList(hn *html.Node) bool {
	return hn.DataAtom == atom.Ul || hn.DataAtom == atom.Ol
}

func isYoutube(hn *html.Node) bool {
	return hn.DataAtom == atom.Video
}

func isFragmentImport(hn *html.Node) bool {
	return hn.DataAtom == 0 && strings.HasPrefix(hn.Data, convertedImportsDataPrefix)
}

// countTwo starts counting the number of a Atom children in hn.
// It returns as soon as the count exceeds 1, so the returned value is inexact.
//
// The callers can test for > 1 to verify whether a node contains two
// or more children of the Atom a.
func countTwo(hn *html.Node, a atom.Atom) int {
	var count int
	for c := hn.FirstChild; c != nil; c = c.NextSibling {
		if c.DataAtom == a {
			count++
		} else {
			count += countTwo(c, a)
		}
		if count > 1 {
			break
		}
	}
	return count
}

// countDirect returns the number of immediate children of hn.
func countDirect(hn *html.Node) int {
	var count int
	for c := hn.FirstChild; c != nil; c = c.NextSibling {
		count++
	}
	return count
}

// findAtom returns first child of root which matches a, nil otherwise.
// It returns root if it is the same Atom as a.
func findAtom(root *html.Node, a atom.Atom) *html.Node {
	if root.DataAtom == a {
		return root
	}
	for c := root.FirstChild; c != nil; c = c.NextSibling {
		if v := findAtom(c, a); v != nil {
			return v
		}
	}
	return nil
}

// TODO reuse code with findAtom?
func findChildAtoms(root *html.Node, a atom.Atom) []*html.Node {
	var nodes []*html.Node
	for hn := root.FirstChild; hn != nil; hn = hn.NextSibling {
		if hn.DataAtom == a {
			nodes = append(nodes, hn)
		}
		nodes = append(nodes, findChildAtoms(hn, a)...)
	}
	return nodes
}

type considerSelf int

const (
	doNotConsiderSelf considerSelf = iota
	doConsiderSelf
)

// findNearestAncestor finds the nearest ancestor of the given node of any of the given atoms.
// A pointer to the ancestor is returned, or nil if none are found.
// If doConsiderSelf is passed, the given node itself counts as an ancestor for our purposes.
func findNearestAncestor(n *html.Node, atoms map[atom.Atom]struct{}, cs considerSelf) *html.Node {
	if _, ok := atoms[n.DataAtom]; cs == doConsiderSelf && ok {
		return n
	}
	for p := n.Parent; p != nil; p = p.Parent {
		if _, ok := atoms[p.DataAtom]; ok {
			return p
		}
	}
	return nil
}

var blockParents = map[atom.Atom]struct{}{
	atom.H1:  {},
	atom.H2:  {},
	atom.H3:  {},
	atom.H4:  {},
	atom.H5:  {},
	atom.H6:  {},
	atom.Li:  {},
	atom.P:   {},
	atom.Div: {},
}

// findNearestBlockAncestor finds the nearest ancestor node of a block atom.
// For instance, block parent of "text" in <ul><li>text</li></ul> is <li>,
// while block parent of "text" in <p><span>text</span></p> is <p>.
// The node passed in itself is never considered.
// A pointer to the ancestor is returned, or nil if none are found.
func findNearestBlockAncestor(n *html.Node) *html.Node {
	return findNearestAncestor(n, blockParents, doNotConsiderSelf)
}

// nodeAttr checks the given node's HTML attributes for the given key.
// The corresponding value is returned, or the empty string if the key is not found.
// Keys are case insensitive.
func nodeAttr(n *html.Node, key string) string {
	key = strings.ToLower(key)
	for _, attr := range n.Attr {
		if strings.ToLower(attr.Key) == key {
			return attr.Val
		}
	}
	return ""
}

// TODO divide into smaller functions
// TODO redo comment, more than just text nodes are handled and atom.A
// TODO should we really have trim?
// TODO part of why this is weird is because the processing is split across root and child nodes. could restructure
// stringifyNode extracts and concatenates all text nodes starting with root.
// Line breaks are inserted at <br> and any non-<span> elements.
func stringifyNode(root *html.Node, trim bool) string {
	if root.Type == html.TextNode {
		s := textCleaner.Replace(root.Data)
		s = strings.Replace(s, "\n", " ", -1)
		if !trim {
			return s
		}
		return strings.TrimSpace(s)
	}
	if root.DataAtom == atom.Br && !trim {
		return "\n"
	}
	var buf bytes.Buffer
	for c := root.FirstChild; c != nil; c = c.NextSibling {
		if c.DataAtom == atom.Br {
			buf.WriteRune('\n')
			continue
		}
		if c.Type == html.TextNode {
			buf.WriteString(c.Data)
			continue
		}
		if c.DataAtom != atom.Span && c.DataAtom != atom.A {
			buf.WriteRune('\n')
		}
		buf.WriteString(stringifyNode(c, false))
	}
	s := textCleaner.Replace(buf.String())
	if !trim {
		return s
	}
	return strings.TrimSpace(s)
}
