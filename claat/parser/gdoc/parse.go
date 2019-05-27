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
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/googlecodelabs/tools/claat/parser"
	"github.com/googlecodelabs/tools/claat/types"
	"github.com/googlecodelabs/tools/claat/util"
)

func init() {
	parser.Register("gdoc", &Parser{})
}

// Parser is a Google Doc parser.
type Parser struct {
}

// Parse parses a codelab exported in HTML from Google Docs.
func (p *Parser) Parse(r io.Reader) (*types.Codelab, error) {
	// TODO: use html.Tokenizer instead
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return parseDoc(doc)
}

// ParseFragment parses a codelab fragment exported in HTML from Google Docs.
func (p *Parser) ParseFragment(r io.Reader) ([]types.Node, error) {
	// TODO: use html.Tokenizer instead
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return parseFragment(doc)
}

const (
	metaSep         = ":"           // step instruction format, key:value
	metaDuration    = "duration"    // step duration instruction
	metaEnvironment = "environment" // step environment instruction
	metaTagOpen     = "[["          // start of tag-based meta instruction
	metaTagClose    = "]]"          // end of tag-based meta instruction
	metaTagImport   = "import"      // import remote resource instruction

	// possible content of special header nodes in lower case.
	headerLearn = "what you'll learn"
	headerCover = "what we've covered"
	headerFAQ   = "frequently asked questions"

	// google docs comments are links with commentPrefix.
	commentPrefix = "#cmnt"
)

var (
	// durFactor is a slice of duration parser multipliers,
	// ordered after the usage in codelab docs
	durFactor = []time.Duration{time.Hour, time.Minute, time.Second}

	// textCleaner replaces "smart quotes" and other unicode runes
	// with their respective ascii versions.
	textCleaner = strings.NewReplacer(
		"\u2019", "'", "\u201C", `"`, "\u201D", `"`, "\u2026", "...",
		"\u00A0", " ", "\u0085", " ",
	)
)

type stateFlag uint32

// entities to skip during a tree parsing
const (
	fSkipHeader stateFlag = 1 << iota // skip code blocks but not code text
	fSkipCode
	fSkipList
	fSkipTable
	fSkipInfobox
	fSkipSurvey
	fMakeBold
	fMakeItalic
	fMakeCode

	// skip all block structures
	fSkipBlock = fSkipCode | fSkipTable | fSkipInfobox | fSkipSurvey
)

type docState struct {
	clab     *types.Codelab // codelab and its metadata
	totdur   time.Duration  // total codelab duration
	survey   int            // last used survey ID
	css      cssStyle       // styles of the doc
	step     *types.Step    // current codelab step
	lastNode types.Node     // last appended node
	env      []string       // current enviornment
	cur      *html.Node     // current HTML node
	flags    stateFlag      // current flags
	stack    []*stackItem   // cur and flags stack
}

type stackItem struct {
	cur   *html.Node
	flags stateFlag
}

func (ds *docState) push(cur *html.Node, flags stateFlag) {
	if cur == nil {
		cur = ds.cur
	}
	ds.stack = append(ds.stack, &stackItem{ds.cur, ds.flags})
	ds.cur = cur
	ds.flags = flags
}

func (ds *docState) pop() {
	n := len(ds.stack)
	if n == 0 {
		return
	}
	item := ds.stack[n-1]
	ds.stack = ds.stack[:n-1]
	ds.cur = item.cur
	ds.flags = item.flags
}

func (ds *docState) appendNodes(nn ...types.Node) {
	if ds.step == nil || len(nn) == 0 {
		return
	}
	if len(ds.env) != 0 {
		for _, n := range nn {
			n.MutateEnv(append(n.Env(), ds.env...))
		}
	}
	ds.step.Content.Append(nn...)
	ds.lastNode = nn[len(nn)-1]
}

func parseFragment(doc *html.Node) ([]types.Node, error) {
	body := findAtom(doc, atom.Body)
	if body == nil {
		return nil, fmt.Errorf("document without a body")
	}
	style, err := parseStyle(doc)
	if err != nil {
		return nil, err
	}
	ds := &docState{
		clab: &types.Codelab{},
		css:  style,
	}
	ds.step = ds.clab.NewStep("fragment")
	for ds.cur = body.FirstChild; ds.cur != nil; ds.cur = ds.cur.NextSibling {
		if isComment(ds.css, ds.cur) {
			// docs export comments at the end of the body
			break
		}
		parseTop(ds)
	}
	finalizeStep(ds.step)
	return ds.step.Content.Nodes, nil
}

// parseDoc parses codelab doc exported as text/html.
// The doc must contain CSS styles and <body> as exported from Google Doc.
func parseDoc(doc *html.Node) (*types.Codelab, error) {
	body := findAtom(doc, atom.Body)
	if body == nil {
		return nil, fmt.Errorf("document without a body")
	}
	style, err := parseStyle(doc)
	if err != nil {
		return nil, err
	}

	ds := &docState{
		clab: &types.Codelab{},
		css:  style,
	}
	for ds.cur = body.FirstChild; ds.cur != nil; ds.cur = ds.cur.NextSibling {
		if isComment(ds.css, ds.cur) {
			// docs export comments at the end of the body
			break
		}
		switch {
		case hasClass(ds.cur, "title") && ds.step == nil:
			if v := stringifyNode(ds.cur, true, false); v != "" {
				ds.clab.Title = v
			}
			if ds.clab.ID == "" {
				ds.clab.ID = slug(ds.clab.Title)
			}
			continue
		case ds.cur.DataAtom == atom.Table && ds.step == nil:
			metaTable(ds)
			continue
		case ds.cur.DataAtom == atom.H1:
			newStep(ds)
			continue
		}
		// ignore everything else before the first step
		if ds.step != nil {
			parseTop(ds)
		}
	}

	finalizeStep(ds.step) // TODO: last ds.step is never finalized in newStep
	ds.clab.Tags = util.Unique(ds.clab.Tags)
	sort.Strings(ds.clab.Tags)
	ds.clab.Duration = int(ds.totdur.Minutes())
	return ds.clab, nil
}

func finalizeStep(s *types.Step) {
	if s == nil {
		return
	}
	s.Tags = util.Unique(s.Tags)
	sort.Strings(s.Tags)
	s.Content.Nodes = blockNodes(s.Content.Nodes)
	s.Content.Nodes = compactNodes(s.Content.Nodes)
	// TODO: find a better place for the code below
	// find [[directive]] instructions and act accordingly
	for i, n := range s.Content.Nodes {
		if n.Type() != types.NodeList {
			continue
		}
		l := n.(*types.ListNode)
		// [[ directive ... ]]
		if len(l.Nodes) < 4 {
			continue
		}
		// first element is opening [[
		if t, ok := l.Nodes[0].(*types.TextNode); !ok || t.Value != metaTagOpen {
			continue
		}
		// last element is closing ]]
		if t, ok := l.Nodes[len(l.Nodes)-1].(*types.TextNode); !ok || t.Value != metaTagClose {
			continue
		}
		// second element is a text in bold
		t, ok := l.Nodes[1].(*types.TextNode)
		if !ok || !t.Bold || t.Italic || t.Code {
			continue
		}
		// execute transform and replace t with the result
		v := strings.ToLower(strings.TrimSpace(t.Value))
		r := transformNodes(v, l.Nodes[2:len(l.Nodes)-1])
		if r != nil {
			r.MutateEnv(l.Env())
			s.Content.Nodes[i] = r
		}
	}
}

func transformNodes(name string, nodes []types.Node) types.Node {
	if name == metaTagImport && len(nodes) == 1 {
		u, ok := nodes[0].(*types.URLNode)
		if !ok {
			return nil
		}
		return types.NewImportNode(u.URL)
	}
	return nil
}

// parseTop parses nodes tree starting at, and including, ds.cur.
// Parsed nodes are squashed and added to ds.step content.
func parseTop(ds *docState) {
	if n, ok := parseNode(ds); ok {
		if n != nil {
			ds.appendNodes(n)
		}
		return
	}
	ds.push(nil, ds.flags)
	nn := parseSubtree(ds)
	ds.pop()
	ds.appendNodes(compactNodes(nn)...)
}

// parseSubtree parses children of root recursively.
// It may modify ds.cur, so the caller is responsible for wrapping
// this function in ds.push and ds.pop.
func parseSubtree(ds *docState) []types.Node {
	var nodes []types.Node
	for ds.cur = ds.cur.FirstChild; ds.cur != nil; ds.cur = ds.cur.NextSibling {
		if n, ok := parseNode(ds); ok {
			if n != nil {
				nodes = append(nodes, n)
			}
			continue
		}
		ds.push(nil, ds.flags)
		nodes = append(nodes, parseSubtree(ds)...)
		ds.pop()
	}
	return nodes
}

// parseNode parses html node hn if it is a recognized node construction.
// It returns a bool indicating that hn has been accepted and parsed.
// Some nodes result in metadata parsing, in which case the returned bool is still true,
// but resuling types.Node is nil.
//
// The flag argument modifies default behavour of the func.
func parseNode(ds *docState) (types.Node, bool) {
	switch {
	case isMeta(ds.css, ds.cur):
		metaStep(ds)
		return nil, true
	case ds.cur.Type == html.TextNode || ds.cur.DataAtom == atom.Br:
		return text(ds), true
	case ds.cur.DataAtom == atom.A:
		return link(ds), true
	case ds.cur.DataAtom == atom.Img:
		return image(ds), true
	case isButton(ds.css, ds.cur):
		return button(ds), true
	case ds.flags&fSkipHeader == 0 && isHeader(ds.cur):
		return header(ds), true
	case ds.flags&fSkipList == 0 && isList(ds.cur):
		return list(ds), true
	case ds.flags&fSkipCode == 0 && isConsole(ds.css, ds.cur):
		return code(ds, true), true
	case ds.flags&fSkipCode == 0 && isCode(ds.css, ds.cur):
		return code(ds, false), true
	case ds.flags&fSkipInfobox == 0 && isInfobox(ds.css, ds.cur):
		return infobox(ds), true
	case ds.flags&fSkipSurvey == 0 && isSurvey(ds.css, ds.cur):
		return survey(ds), true
	case ds.flags&fSkipTable == 0 && isTable(ds.cur):
		return table(ds), true
	}
	return nil, false
}

// newStep creates a new codelab step from ds.cur
// and finalizes nodes of the previous step.
func newStep(ds *docState) {
	t := stringifyNode(ds.cur, true, false)
	if t == "" {
		return
	}
	finalizeStep(ds.step)
	ds.step = ds.clab.NewStep(t)
	ds.env = nil
}

// metaTable parses the top <table> of a codelab doc
func metaTable(ds *docState) {
	for tr := findAtom(ds.cur, atom.Tr); tr != nil; tr = tr.NextSibling {
		if tr.FirstChild == nil || tr.FirstChild.NextSibling == nil {
			continue
		}
		s := stringifyNode(tr.FirstChild.NextSibling, true, false)
		switch strings.ToLower(stringifyNode(tr.FirstChild, true, false)) {
		case "id", "url":
			ds.clab.ID = s
		case "author", "authors":
			ds.clab.Authors = s
		case "badge", "badge id":
			ds.clab.BadgeID = s
		case "summary":
			ds.clab.Summary = stringifyNode(tr.FirstChild.NextSibling, true, true)
		case "category", "categories":
			ds.clab.Categories = util.Unique(stringSlice(s))
		case "environment", "environments", "tags":
			ds.clab.Tags = stringSlice(s)
			toLowerSlice(ds.clab.Tags)
		case "status", "state":
			v := stringSlice(s)
			toLowerSlice(v)
			sv := types.LegacyStatus(v)
			ds.clab.Status = &sv
		case "feedback", "feedback link":
			ds.clab.Feedback = s
		case "analytics", "analytics account", "google analytics":
			ds.clab.GA = s
		}
	}
	if len(ds.clab.Categories) > 0 {
		ds.clab.Theme = slug(ds.clab.Categories[0])
	}
}

// metaStep parses a codelab step meta instructions.
func metaStep(ds *docState) {
	var text string
	for {
		text += stringifyNode(ds.cur, false, false)
		if ds.cur.NextSibling == nil || !isMeta(ds.css, ds.cur.NextSibling) {
			break
		}
		ds.cur = ds.cur.NextSibling
	}
	meta := strings.SplitN(strings.TrimSpace(text), metaSep, 2)
	if len(meta) != 2 {
		return
	}
	value := strings.TrimSpace(meta[1])
	switch strings.ToLower(strings.TrimSpace(meta[0])) {
	case metaDuration:
		parts := strings.SplitN(value, ":", len(durFactor))
		if len(parts) == 1 {
			parts = append(parts, "0") // default number is minutes
		}
		var d time.Duration
		for i, v := range parts {
			vi, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			d += time.Duration(vi) * durFactor[len(durFactor)-len(parts)+i]
		}
		ds.step.Duration = roundDuration(d)
		ds.totdur += ds.step.Duration
	case metaEnvironment:
		ds.env = util.Unique(stringSlice(value))
		toLowerSlice(ds.env)
		ds.step.Tags = append(ds.step.Tags, ds.env...)
		ds.clab.Tags = append(ds.clab.Tags, ds.env...)
		if ds.lastNode != nil && types.IsHeader(ds.lastNode.Type()) {
			ds.lastNode.MutateEnv(ds.env)
		}
	}
}

// header creates a HeaderNode out of hn.
// It returns nil if header content is empty.
// A non-empty header will always reset ds.env to nil.
//
// Given that headers do not belong to any block, the returned node's B
// field is always nil.
func header(ds *docState) types.Node {
	ds.push(nil, ds.flags|fSkipBlock|fSkipHeader|fSkipList)
	nodes := parseSubtree(ds)
	ds.pop()
	if len(nodes) == 0 {
		return nil
	}
	n := types.NewHeaderNode(headerLevel[ds.cur.DataAtom], nodes...)
	switch strings.ToLower(stringifyNode(ds.cur, true, false)) {
	case headerLearn, headerCover:
		n.MutateType(types.NodeHeaderCheck)
	case headerFAQ:
		n.MutateType(types.NodeHeaderFAQ)
	}
	ds.env = nil
	return n
}

// infobox doesn't have a block parent.
func infobox(ds *docState) types.Node {
	ds.push(nil, ds.flags|fSkipCode|fSkipInfobox|fSkipSurvey)
	nn := parseSubtree(ds)
	nn = blockNodes(nn)
	nn = compactNodes(nn)
	ds.pop()
	if len(nn) == 0 {
		return nil
	}
	kind := types.InfoboxPositive
	if isInfoboxNegative(ds.css, ds.cur) {
		kind = types.InfoboxNegative
	}
	return types.NewInfoboxNode(kind, nn...)
}

// table parses an arbitrary <table> element and its children.
// It may return other elements if the table is just a wrap.
func table(ds *docState) types.Node {
	var rows [][]*types.GridCell
	for _, tr := range findChildAtoms(ds.cur, atom.Tr) {
		ds.push(tr, ds.flags)
		r := tableRow(ds)
		ds.pop()
		rows = append(rows, r)
	}
	if len(rows) == 0 {
		return nil
	}
	return types.NewGridNode(rows...)
}

func tableRow(ds *docState) []*types.GridCell {
	var row []*types.GridCell
	for td := findAtom(ds.cur, atom.Td); td != nil; td = td.NextSibling {
		if td.DataAtom != atom.Td {
			continue
		}
		ds.push(td, ds.flags|fSkipBlock)
		nn := parseSubtree(ds)
		nn = blockNodes(nn)
		nn = compactNodes(nn)
		ds.pop()
		if len(nn) == 0 {
			continue
		}
		cs, err := strconv.Atoi(nodeAttr(td, "colspan"))
		if err != nil {
			cs = 1
		}
		rs, err := strconv.Atoi(nodeAttr(td, "rowspan"))
		if err != nil {
			rs = 1
		}
		cell := &types.GridCell{
			Colspan: cs,
			Rowspan: rs,
			Content: types.NewListNode(nn...),
		}
		row = append(row, cell)
	}
	return row
}

// survey expects a header followed by 1 or more lists.
func survey(ds *docState) types.Node {
	// find direct parent of the survey elements
	hn := findAtom(ds.cur, atom.Ul)
	if hn == nil {
		return nil
	}
	hn = hn.Parent
	// parse survey elements
	var gg []*types.SurveyGroup
	for c := hn.FirstChild; c != nil; {
		if !isHeader(c) {
			c = c.NextSibling
			continue
		}
		opt, next := surveyOpt(c.NextSibling)
		if len(opt) > 0 {
			gg = append(gg, &types.SurveyGroup{
				Name:    stringifyNode(c, true, false),
				Options: opt,
			})
		}
		c = next
	}
	if len(gg) == 0 {
		return nil
	}
	ds.survey++
	id := fmt.Sprintf("%s-%d", ds.clab.ID, ds.survey)
	return types.NewSurveyNode(id, gg...)
}

func surveyOpt(hn *html.Node) ([]string, *html.Node) {
	var opt []string
	for ; hn != nil; hn = hn.NextSibling {
		if isHeader(hn) {
			return opt, hn
		}
		if hn.DataAtom != atom.Ul {
			continue
		}
		for li := findAtom(hn, atom.Li); li != nil; li = li.NextSibling {
			if li.DataAtom != atom.Li {
				continue
			}
			opt = append(opt, stringifyNode(li, true, true))
		}
	}
	return opt, nil
}

// code parses hn as inline or block codes.
// Inline code node will be of type NodeText.
func code(ds *docState, term bool) types.Node {
	td := findParent(ds.cur, atom.Td)
	// inline <code> text
	if td == nil {
		return text(ds)
	}
	// block code or terminal
	v := stringifyNode(ds.cur, false, true)
	if v == "" {
		if countDirect(ds.cur.Parent) > 1 {
			return nil
		}
		v = "\n"
	} else if ds.cur.Parent.FirstChild == ds.cur && ds.cur.Parent.DataAtom != atom.Span {
		v = "\n" + v
	}
	n := types.NewCodeNode(v, term)
	n.MutateBlock(td)
	return n
}

// list parses <ul> and <ol> lists.
// It returns nil if the list has no items.
func list(ds *docState) types.Node {
	typ := nodeAttr(ds.cur, "type")
	if ds.cur.DataAtom == atom.Ol && typ == "" {
		typ = "1"
	}
	start, _ := strconv.Atoi(nodeAttr(ds.cur, "start"))
	list := types.NewItemsListNode(typ, start)
	for hn := findAtom(ds.cur, atom.Li); hn != nil; hn = hn.NextSibling {
		if hn.DataAtom != atom.Li {
			continue
		}
		ds.push(hn, ds.flags)
		nn := parseSubtree(ds)
		nn = compactNodes(nn)
		ds.pop()
		if len(nn) > 0 {
			list.NewItem(nn...)
		}
	}
	if len(list.Items) == 0 {
		return nil
	}
	if ds.lastNode != nil {
		switch ds.lastNode.Type() {
		case types.NodeHeaderCheck:
			list.MutateType(types.NodeItemsCheck)
		case types.NodeHeaderFAQ:
			list.MutateType(types.NodeItemsFAQ)
		}
	}
	return list
}

// image creates a new ImageNode out of hn, parsing its src attribute.
// It returns nil if src is empty.
// It may also return a YouTubeNode if alt property contains specific substring.
// or an IframeNode if the alt property contains a URL other than youtube.
func image(ds *docState) types.Node {
	alt := nodeAttr(ds.cur, "alt")
	errorAlt := ""
	if strings.Contains(alt, "youtube.com/watch") {
		return youtube(ds)
	} else if strings.Contains(alt, "https://") {
		u, err := url.Parse(nodeAttr(ds.cur, "alt"))
		if err != nil {
			return nil
		}
		// For iframe, make sure URL ends in whitelisted domain.
		ok := false
		for _, domain := range types.IframeWhitelist {
			if strings.HasSuffix(u.Hostname(), domain) {
				ok = true
				break
			}
		}
		if ok {
			return iframe(ds)
		} else {
			errorAlt = "The domain of the requested iframe (" + u.Hostname() + ") has not been whitelisted."
			fmt.Fprint(os.Stderr, errorAlt+"\n")
		}
	}
	s := nodeAttr(ds.cur, "src")
	if s == "" {
		return nil
	}
	n := types.NewImageNode(s)
	n.Width = styleFloatValue(ds.cur, "width")
	n.MutateBlock(findBlockParent(ds.cur))
	if errorAlt != "" {
		n.Alt = errorAlt
	} else {
		n.Alt = nodeAttr(ds.cur, "alt")
	}
	n.Title = nodeAttr(ds.cur, "title")
	return n
}

func youtube(ds *docState) types.Node {
	u, err := url.Parse(nodeAttr(ds.cur, "alt"))
	if err != nil {
		return nil
	}
	v := u.Query().Get("v")
	if v == "" {
		return nil
	}
	n := types.NewYouTubeNode(v)
	n.MutateBlock(true)
	return n
}

func iframe(ds *docState) types.Node {
	u, err := url.Parse(nodeAttr(ds.cur, "alt"))
	if err != nil {
		return nil
	}
	// Allow only https.
	if u.Scheme != "https" {
		return nil
	}
	n := types.NewIframeNode(u.String())
	n.MutateBlock(true)
	return n
}

// button returns either a text node, if no <a> child element is present,
// or link node, containing the button.
// It returns nil if no content nodes are present.
func button(ds *docState) types.Node {
	a := findAtom(ds.cur, atom.A)
	if a == nil {
		return text(ds)
	}
	href := cleanURL(nodeAttr(a, "href"))
	if href == "" {
		return nil
	}

	ds.push(a, fSkipBlock|fSkipList)
	nodes := parseSubtree(ds)
	ds.pop()
	if len(nodes) == 0 {
		return nil
	}

	s := strings.ToLower(stringifyNode(a, true, false))
	dl := strings.HasPrefix(s, "download ")
	btn := types.NewButtonNode(true, true, dl, nodes...)

	ln := types.NewURLNode(href, btn)
	ln.MutateBlock(findBlockParent(ds.cur))
	return ln
}

// Link creates a URLNode out of hn, parsing href and name attributes.
// It returns nil if hn contents is empty.
// The resuling link's content is always a single text node.
func link(ds *docState) types.Node {
	href := cleanURL(nodeAttr(ds.cur, "href"))
	if strings.HasPrefix(href, commentPrefix) {
		// doc comments; ignore
		return nil
	}

	text := stringifyNode(ds.cur, false, true)
	if strings.TrimSpace(text) == "" {
		return nil
	}

	t := types.NewTextNode(text)
	if ds.flags&fMakeBold != 0 || isBold(ds.css, ds.cur.Parent) {
		t.Bold = true
	}
	if ds.flags&fMakeItalic != 0 || isItalic(ds.css, ds.cur.Parent) {
		t.Italic = true
	}
	if ds.flags&fMakeCode != 0 || isCode(ds.css, ds.cur.Parent) {
		t.Code = true
	}
	if href == "" || href[0] == '#' {
		t.MutateBlock(findBlockParent(ds.cur))
		return t
	}

	n := types.NewURLNode(href, t)
	n.Name = nodeAttr(ds.cur, "name")
	if v := nodeAttr(ds.cur, "target"); v != "" {
		n.Target = v
	}
	n.MutateBlock(findBlockParent(ds.cur))
	return n
}

// text creates a TextNode using hn.Data as contents.
// It returns nil if hn.Data is empty or contains only space runes.
func text(ds *docState) types.Node {
	bold := isBold(ds.css, ds.cur)
	italic := isItalic(ds.css, ds.cur)
	code := isCode(ds.css, ds.cur) || isConsole(ds.css, ds.cur)

	// TODO: verify whether this actually does anything
	if a := findAtom(ds.cur, atom.A); a != nil {
		f := fSkipBlock
		if bold {
			f |= fMakeBold
		}
		if italic {
			f |= fMakeItalic
		}
		if code {
			f |= fMakeCode
		}
		ds.push(a, f)
		l := link(ds)
		ds.pop()
		if l != nil {
			l.MutateBlock(findBlockParent(ds.cur))
			return l
		}
	}

	v := stringifyNode(ds.cur, false, true)
	n := types.NewTextNode(v)
	n.Bold = bold
	n.Italic = italic
	n.Code = code
	n.MutateBlock(findBlockParent(ds.cur))
	return n
}

// cleanURL extracts original URL from v, where the value
// may be wrapped in https://google.com/url?q=url.
func cleanURL(v string) string {
	if !strings.Contains(v, "google.com/url?") {
		return v
	}
	u, err := url.Parse(v)
	if err != nil {
		return v
	}
	if q, ok := u.Query()["q"]; ok {
		return q[0]
	}
	return v
}

// slug converts any string s to a slug.
// It replaces [^a-z0-9\-] with non-repeating '-'.
func slug(s string) string {
	var buf bytes.Buffer
	dash := true
	for _, r := range strings.ToLower(s) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' && !dash {
			buf.WriteRune(r)
			dash = r == '-'
			continue
		}
		if !dash {
			buf.WriteRune('-')
			dash = true
		}
	}
	return buf.String()
}

// stringSlice splits v by comma "," while ignoring empty elements.
func stringSlice(v string) []string {
	f := strings.Split(v, ",")
	a := f[0:0]
	for _, s := range f {
		s = strings.TrimSpace(s)
		if s != "" {
			a = append(a, s)
		}
	}
	return a
}

func toLowerSlice(a []string) {
	for i, s := range a {
		a[i] = strings.ToLower(s)
	}
}

func roundDuration(d time.Duration) time.Duration {
	rd := time.Duration(d.Minutes()) * time.Minute
	if rd < d {
		rd += time.Minute
	}
	return rd
}
