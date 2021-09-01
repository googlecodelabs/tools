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

package gdoc

import (
	"bytes"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	// headerLevel maps HTML tags to their level in parser.HeaderNode.
	headerLevel = map[atom.Atom]int{
		atom.H2: 2,
		atom.H3: 3,
		atom.H4: 4,
		atom.H5: 5,
		atom.H6: 6,
	}
)

// isHeader returns true if hn is one of secondary headers.
// Step header is not one of them.
func isHeader(hn *html.Node) bool {
	_, ok := headerLevel[hn.DataAtom]
	return ok
}

func isMeta(css cssStyle, hn *html.Node) bool {
	return hasClassStyle(css, hn, "color", metaColor)
}

func isBold(css cssStyle, hn *html.Node) bool {
	if hn.Type == html.TextNode {
		hn = hn.Parent
	}
	return hn.DataAtom == atom.Strong ||
		hn.DataAtom == atom.B ||
		hasClassStyle(css, hn, "font-weight", "bold") ||
		hasClassStyle(css, hn, "font-weight", "700")
}

func isItalic(css cssStyle, hn *html.Node) bool {
	if hn.Type == html.TextNode {
		hn = hn.Parent
	}
	return hn.DataAtom == atom.Em ||
		hn.DataAtom == atom.I ||
		hasClassStyle(css, hn, "font-style", "italic")
}

func isConsole(css cssStyle, hn *html.Node) bool {
	if hn.Type == html.TextNode {
		hn = hn.Parent
	}
	return hasClassStyle(css, hn, "font-family", fontConsole)
}

func isCode(css cssStyle, hn *html.Node) bool {
	if hn.Type == html.TextNode {
		hn = hn.Parent
	}
	return hasClassStyle(css, hn, "font-family", fontCode)
}

func isButton(css cssStyle, hn *html.Node) bool {
	return hasClassStyle(css, hn, "background-color", buttonColor)
}

func isInfobox(css cssStyle, hn *html.Node) bool {
	if hn.DataAtom != atom.Td {
		return false
	}
	return hasClassStyle(css, hn, "background-color", ibPositiveColor) ||
		isInfoboxNegative(css, hn)
}

func isInfoboxNegative(css cssStyle, hn *html.Node) bool {
	if hn.DataAtom != atom.Td {
		return false
	}
	return hasClassStyle(css, hn, "background-color", ibNegativeColor)
}

func isSurvey(css cssStyle, hn *html.Node) bool {
	if hn.DataAtom != atom.Td {
		return false
	}
	return hasClassStyle(css, hn, "background-color", surveyColor)
}

func isComment(css cssStyle, hn *html.Node) bool {
	if hn.DataAtom != atom.Div {
		return false
	}
	return hasClassStyle(css, hn, "border", "1px solid black")
}

func isTable(hn *html.Node) bool {
	if hn.DataAtom != atom.Table {
		return false
	}
	return countTwo(hn, atom.Tr) > 1 || countTwo(hn, atom.Td) > 1
}

func isList(hn *html.Node) bool {
	return hn.DataAtom == atom.Ul || hn.DataAtom == atom.Ol
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

// findParent is like findAtom but search is in the opposite direction.
// It is faster to look for parent than child lookup in findAtom.
func findParent(root *html.Node, a atom.Atom) *html.Node {
	if root.DataAtom == a {
		return root
	}
	for c := root.Parent; c != nil; c = c.Parent {
		if c.DataAtom == a {
			return c
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

// findBlockParent looks up nearest block parent node of hn.
// For instance, block parent of "text" in <ul><li>text</li></ul> is <li>,
// while block parent of "text" in <p><span>text</span></p> is <p>.
func findBlockParent(hn *html.Node) *html.Node {
	for p := hn.Parent; p != nil; p = p.Parent {
		if _, ok := blockParents[p.DataAtom]; ok {
			return p
		}
	}
	return nil
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

// stringifyNode extracts and concatenates all text nodes starting with root.
// Line breaks are inserted at <br> and any non-<span> elements if requested.
func stringifyNode(root *html.Node, trim bool, lineBreak bool) string {
	if root.Type == html.TextNode {
		s := textCleaner.Replace(root.Data)
		if !trim {
			return s
		}
		return strings.TrimSpace(s)
	}
	if root.DataAtom == atom.Br && !trim {
		if lineBreak {
			return "\n"
		}
		return ""
	}
	var buf bytes.Buffer
	for c := root.FirstChild; c != nil; c = c.NextSibling {
		if c.DataAtom == atom.A {
			href := cleanURL(nodeAttr(c, "href"))
			if strings.HasPrefix(href, commentPrefix) {
				// skip doc comments
				continue
			}
		}
		if c.DataAtom == atom.Br {
			if lineBreak {
				buf.WriteRune('\n')
			}
			continue
		}
		if c.Type == html.TextNode {
			buf.WriteString(c.Data)
			continue
		}
		if c.DataAtom != atom.Span && c.DataAtom != atom.A {
			if lineBreak {
				buf.WriteRune('\n')
			}
		}
		buf.WriteString(stringifyNode(c, false, lineBreak))
	}
	s := textCleaner.Replace(buf.String())
	if !trim {
		return s
	}
	return strings.TrimSpace(s)
}
