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
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/x1ddos/csslex"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	metaColor       = "#b7b7b7"     // step meta instruction
	buttonColor     = "#6aa84f"     // button background color
	fontCode        = "courier new" // source code format in original doc
	fontConsole     = "consolas"    // terminal text format in original doc
	ibPositiveColor = "#d9ead3"     // positive infobox background
	ibNegativeColor = "#fce5cd"     // negative infobox background
	surveyColor     = "#cfe2f3"     // survey background color
)

// cssStyle represents styles of an exported Google Doc.
type cssStyle map[string]map[string]string

// parseStyle parses styles found in doc.
// The argument can be anything which is, or contains as a child,
// a <style> element.
func parseStyle(doc *html.Node) (cssStyle, error) {
	style := make(cssStyle)
	node := findAtom(doc, atom.Style)
	if node == nil {
		return style, nil
	}
	css := stringifyNode(node, true, true)

	var skip bool
	var sel []string
	for item := range csslex.Lex(string(css)) {
		switch item.Typ {
		case csslex.ItemError:
			return nil, fmt.Errorf("error parsing CSS at %d: %v", item.Pos, item.Val)
		case csslex.ItemAtRuleIdent, csslex.ItemAtRule:
			continue
		case csslex.ItemAtRuleBlockStart:
			skip = true
			continue
		case csslex.ItemAtRuleBlockEnd:
			skip = false
			continue
		case csslex.ItemBlockEnd:
			sel = nil
		case csslex.ItemSelector:
			if skip || !strings.HasPrefix(item.Val, ".") || strings.ContainsRune(item.Val, ':') {
				continue
			}
			sel = append(sel, item.Val)
			if _, ok := style[item.Val]; !ok {
				style[item.Val] = map[string]string{}
			}
		case csslex.ItemDecl:
			if skip || len(sel) == 0 {
				continue
			}
			decl := strings.SplitN(item.Val, ":", 2)
			decl[0] = strings.ToLower(strings.TrimSpace(decl[0]))
			decl[1] = strings.ToLower(strings.Trim(strings.TrimSpace(decl[1]), `"`))
			for _, s := range sel {
				style[s][decl[0]] = decl[1]
			}
		}
	}
	return style, nil
}

// classList returns a slice of all CSS classes of node hn.
func classList(hn *html.Node) []string {
	var cls string
	for _, a := range hn.Attr {
		if a.Key == "class" {
			cls = a.Val
			break
		}
	}
	a := strings.Split(cls, " ")
	sort.Strings(a)
	return a
}

// hasClass returns true if the node hn has CSS class name.
func hasClass(hn *html.Node, name string) bool {
	cls := classList(hn)
	i := sort.SearchStrings(cls, name)
	return i < len(cls) && cls[i] == name
}

// hasClassStyle returns true if the node hn has a CSS class style property key
// with the value val.
func hasClassStyle(css cssStyle, hn *html.Node, key, val string) bool {
	for _, c := range classList(hn) {
		s, ok := css["."+c]
		if !ok {
			continue
		}
		if s[key] == val {
			return true
		}
	}
	// no class style, try inline style
	return styleValue(hn, key) == val
}

func styleValue(hn *html.Node, name string) string {
	name = strings.ToLower(name)
	var s string
	for _, a := range hn.Attr {
		if a.Key == "style" {
			s = a.Val
			break
		}
	}
	for _, s = range strings.Split(html.UnescapeString(s), ";") {
		v := strings.SplitN(s, ":", 2)
		if len(v) != 2 {
			continue
		}
		if strings.TrimSpace(strings.ToLower(v[0])) == name {
			return strings.ToLower(strings.Trim(v[1], " \""))
		}
	}
	return ""
}

func styleFloatValue(hn *html.Node, name string) float32 {
	s := styleValue(hn, name)
	if s == "" {
		return 0
	}
	for i, r := range s {
		if r >= '0' && r <= '9' || r == '.' {
			continue
		}
		s = s[:i]
		break
	}
	if f, err := strconv.ParseFloat(s, 32); err == nil {
		return float32(f)
	}
	return -1
}
