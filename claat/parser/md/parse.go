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

// Package md implements a parser for CLaaT. It expects, as input, the output of running a Markdown file through
// the Markdown processor.
package md

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"github.com/stoewer/go-strcase"

	"github.com/googlecodelabs/tools/claat/nodes"
	"github.com/googlecodelabs/tools/claat/parser"
	"github.com/googlecodelabs/tools/claat/types"
	"github.com/googlecodelabs/tools/claat/util"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	gmhtml "github.com/yuin/goldmark/renderer/html"
)

// Metadata constants for the YAML header
const (
	MetaAuthors          = "authors"
	MetaSummary          = "summary"
	MetaID               = "id"
	MetaCategories       = "categories"
	MetaEnvironments     = "environments"
	MetaStatus           = "status"
	MetaFeedbackLink     = "feedback_link"
	MetaAnalyticsAccount = "analytics_account"
	MetaTags             = "tags"
	MetaSource           = "source"
	MetaDuration         = "duration"
)

const (
	metaSep         = ":"           // step instruction format, key:value
	metaDuration    = "duration"    // step duration instruction
	metaEnvironment = "environment" // step environment instruction
	metaTagImport   = "import"      // import remote resource instruction

	// possible content of special header nodes in lower case.
	headerLearn = "what you'll learn"
	headerCover = "what we've covered"
	headerFAQ   = "frequently asked questions"
)

var (
	importsTagRegexp           = regexp.MustCompile("^<<([^<>()]+.md)>>\\s*$")
	convertedImportsDataPrefix = "__unsupported_import_zmcgv2epyv="
	convertedImportsPrefix     = []byte("<!--" + convertedImportsDataPrefix)
	convertedImportsSuffix     = []byte("-->")
)

var metadataRegexp = regexp.MustCompile(`(.+?):(.+)`)
var languageRegexp = regexp.MustCompile(`language-(.+)`)

var (
	// durFactor is a slice of duration parser multipliers,
	// ordered after the usage in codelab docs
	durFactor = []time.Duration{time.Hour, time.Minute, time.Second}

	// TODO make more readable
	// textCleaner replaces "smart quotes" and other unicode runes
	// with their respective ascii versions.
	textCleaner = strings.NewReplacer(
		"\u2019", "'", "\u201C", `"`, "\u201D", `"`, "\u2026", "...",
		"\u00A0", " ", "\u0085", " ",
	)
)

var (
	// ErrForbiddenFragmentImports means importing another markdown file in a markdown fragment is forbidden.
	ErrForbiddenFragmentImports = errors.New("importing content in a fragment is forbidden")
	// ErrForbiddenFragmentSteps means declaring extra codelabs step in a markdown fragment is forbidden.
	ErrForbiddenFragmentSteps = errors.New("defining steps in a fragment is forbidden")
)

// init registers this parser so it is available to CLaaT.
func init() {
	parser.Register("md", &Parser{})
}

// Parser is a Markdown parser.
type Parser struct {
}

// Parse parses a codelab written in Markdown.
func (p *Parser) Parse(r io.Reader, opts parser.Options) (*types.Codelab, error) {
	// Convert Markdown to HTML for easy parsing.
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	b, err = renderToHTML(b)
	if err != nil {
		return nil, err
	}
	h := bytes.NewBuffer(b)
	doc, err := html.Parse(h)
	if err != nil {
		return nil, err
	}
	// Parse the markup.
	return parseMarkup(doc, opts)
}

// ParseFragment parses a codelab fragment written in Markdown.
func (p *Parser) ParseFragment(r io.Reader, opts parser.Options) ([]nodes.Node, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	b, err = renderToHTML(b)
	if err != nil {
		return nil, err
	}
	h := bytes.NewBuffer(b)
	doc, err := html.Parse(h)
	if err != nil {
		return nil, err
	}

	return parsePartialMarkup(doc)
}

func parsePartialMarkup(root *html.Node) ([]nodes.Node, error) {
	body := findAtom(root, atom.Body)
	if body == nil {
		return nil, fmt.Errorf("document without a body")
	}

	ds := newDocState()
	ds.step = ds.clab.NewStep("fragment")
	for ds.cur = body.FirstChild; ds.cur != nil; ds.cur = ds.cur.NextSibling {
		switch {
		case ds.cur.DataAtom == atom.H1:
			return nil, ErrForbiddenFragmentSteps
		case ds.cur.DataAtom == atom.H2:
			return nil, ErrForbiddenFragmentSteps
		}

		parseTop(ds)
	}

	finalizeStep(ds.step)
	if hasImport(ds) {
		return nil, ErrForbiddenFragmentImports
	}

	return ds.step.Content.Nodes, nil
}

type docState struct {
	clab     *types.Codelab // codelab and its metadata
	totdur   time.Duration  // total codelab duration
	survey   int            // last used survey ID
	step     *types.Step    // current codelab step
	lastNode nodes.Node     // last appended node
	env      []string       // current enviornment
	cur      *html.Node     // current HTML node
	stack    []*stackItem   // cur and flags stack
}

type stackItem struct {
	cur *html.Node
}

func newDocState() *docState {
	return &docState{
		clab: types.NewCodelab(),
	}
}

func (ds *docState) push(cur *html.Node) {
	if cur == nil {
		cur = ds.cur
	}
	ds.stack = append(ds.stack, &stackItem{ds.cur})
	ds.cur = cur
}

func (ds *docState) pop() {
	n := len(ds.stack)
	if n == 0 {
		return
	}
	item := ds.stack[n-1]
	ds.stack = ds.stack[:n-1]
	ds.cur = item.cur
}

func (ds *docState) appendNodes(nn ...nodes.Node) {
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

// renderToHTML preprocesses Markdown bytes and then calls a Markdown parser on the Markdown.
// It takes a raw markdown bytes and outputs parsed xhtml in bytes.
func renderToHTML(b []byte) ([]byte, error) {
	b = convertImports(b)
	gmParser := goldmark.New(goldmark.WithRendererOptions(gmhtml.WithUnsafe()), goldmark.WithExtensions(extension.Typographer, extension.Table))
	var out bytes.Buffer
	if err := gmParser.Convert(b, &out); err != nil {
		panic(err)
	}
	return out.Bytes(), nil
}

// parseMarkup accepts html nodes to markup created by the Markdown parser. It returns a pointer to a codelab object, or an error if one occurs.
func parseMarkup(markup *html.Node, opts parser.Options) (*types.Codelab, error) {
	body := findAtom(markup, atom.Body)
	if body == nil {
		return nil, fmt.Errorf("document without a body")
	}

	ds := newDocState()

	for ds.cur = body.FirstChild; ds.cur != nil; ds.cur = ds.cur.NextSibling {
		switch {
		// metadata first
		case ds.cur.DataAtom == atom.H1 && ds.clab.Title == "":
			if v := stringifyNode(ds.cur, true); v != "" {
				ds.clab.Title = v
			}
			continue
		case ds.cur.DataAtom == atom.P && ds.clab.ID == "":
			if err := parseMetadata(ds, opts); err != nil {
				return nil, err
			}
			continue
		case ds.cur.DataAtom == atom.H2:
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
	s.Content.Nodes = parser.BlockNodes(s.Content.Nodes)
	s.Content.Nodes = parser.CompactNodes(s.Content.Nodes)
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
	ds.push(nil)
	nn := parseSubtree(ds)
	ds.pop()
	ds.appendNodes(parser.CompactNodes(nn)...)
}

// parseSubtree parses children of root recursively.
// It may modify ds.cur, so the caller is responsible for wrapping
// this function in ds.push and ds.pop.
func parseSubtree(ds *docState) []nodes.Node {
	var nodes []nodes.Node
	for ds.cur = ds.cur.FirstChild; ds.cur != nil; ds.cur = ds.cur.NextSibling {
		if n, ok := parseNode(ds); ok {
			if n != nil {
				nodes = append(nodes, n)
			}
			continue
		}
		ds.push(nil)
		nodes = append(nodes, parseSubtree(ds)...)
		ds.pop()
	}
	return nodes
}

// parseNode parses html node hn if it is a recognized node construction.
// It returns a bool indicating that hn has been accepted and parsed.
// Some nodes result in metadata parsing, in which case the returned bool is still true,
// but resuling nodes.Node is nil.
//
// The flag argument modifies default behavour of the func.
func parseNode(ds *docState) (nodes.Node, bool) {
	// we have \n end of line nodes after each tag from the blackfriday parser.
	// We just want to ignore them as it makes previous node detection fuzzy.
	if ds.cur.Type == html.TextNode && ds.cur.Data == "\n" {
		return nil, true
	}
	switch {
	case isMeta(ds.cur):
		metaStep(ds)
		return nil, true
	case ds.cur.Type == html.TextNode || ds.cur.DataAtom == atom.Br:
		return text(ds), true
	case ds.cur.DataAtom == atom.A:
		return link(ds), true
	case ds.cur.DataAtom == atom.Img:
		return image(ds), true
	case isButton(ds.cur):
		return button(ds), true
	case isHeader(ds.cur):
		return header(ds), true
	case isList(ds.cur):
		return list(ds), true
	case isConsole(ds.cur):
		return code(ds, true), true
	case isCode(ds.cur):
		return code(ds, false), true
	case isAside(ds.cur):
		return aside(ds), true
	case isNewAside(ds.cur):
		return newAside(ds), true
	case isInfobox(ds.cur):
		return infobox(ds), true
	case isSurvey(ds.cur):
		return survey(ds), true
	case isTable(ds.cur):
		return table(ds), true
	case isYoutube(ds.cur):
		return youtube(ds), true
	case isFragmentImport(ds.cur):
		return fragmentImport(ds), true
	}
	return nil, false
}

// newStep creates a new codelab step from ds.cur
// and finalizes nodes of the previous step.
func newStep(ds *docState) {
	t := stringifyNode(ds.cur, true)
	if t == "" {
		return
	}
	finalizeStep(ds.step)
	ds.step = ds.clab.NewStep(t)
	ds.env = nil
}

// parseMetadata parses the first <p> of a codelab doc to populate metadata.
func parseMetadata(ds *docState, opts parser.Options) error {
	m := map[string]string{}
	// Split the keys from values.
	d := ds.cur.FirstChild.Data
	scanner := bufio.NewScanner(strings.NewReader(d))
	for scanner.Scan() {
		s := metadataRegexp.FindStringSubmatch(scanner.Text())
		if len(s) != 3 {
			continue
		}

		k := strings.ToLower(strings.TrimSpace(s[1]))
		v := strings.TrimSpace(s[2])
		m[k] = v

	}
	if _, ok := m["id"]; !ok || m["id"] == "" {
		return fmt.Errorf("invalid metadata format, missing at least id: %v", m)
	}
	return addMetadataToCodelab(m, ds.clab, opts)
}

// addMetadataToCodelab takes a map of strings to strings, a pointer to a Codelab, and an options struct. It reads the keys of the map,
// and assigns the values to any keys that match a codelab metadata field as defined by the meta* constants.
func addMetadataToCodelab(m map[string]string, c *types.Codelab, opts parser.Options) error {
	for k, v := range m {
		switch strcase.SnakeCase(k) {
		case MetaAuthors:
			// Directly assign the summary to the codelab field.
			c.Authors = v
		case MetaSummary:
			// Directly assign the summary to the codelab field.
			c.Summary = v
		case MetaID:
			// Directly assign the ID to the codelab field.
			c.ID = v
		case MetaCategories:
			// Standardize the categories and append to codelab field.
			c.Categories = append(c.Categories, util.NormalizedSplit(v)...)
		case MetaEnvironments:
			// Standardize the tags and append to the codelab field.
			c.Tags = append(c.Tags, util.NormalizedSplit(v)...)
		case MetaStatus:
			// Standardize the statuses and append to the codelab field.
			statuses := util.NormalizedSplit(v)
			statusesAsLegacy := types.LegacyStatus(statuses)
			c.Status = &statusesAsLegacy
		case MetaFeedbackLink:
			// Directly assign the feedback link to the codelab field.
			c.Feedback = v
		case MetaAnalyticsAccount:
			// Directly assign the GA id to the codelab field.
			c.GA = v
		case MetaTags:
			// Standardize the tags and append to the codelab field.
			c.Tags = append(c.Tags, util.NormalizedSplit(v)...)
		case MetaSource:
			// Directly assign the source doc ID to the source field.
			c.Source = v
		case MetaDuration:
			// Convert the duration to an integer and assign to the duration field.
			duration, err := strconv.Atoi(v)
			if err == nil {
				c.Duration = duration
			}
		default:
			// If not explicitly parsed, it might be a pass_metadata value.
			if _, ok := opts.PassMetadata[k]; ok {
				c.Extra[k] = v
			}
		}
	}
	return nil
}

// metaStep parses a codelab step meta instructions.
func metaStep(ds *docState) {
	var text string
	for {
		text += stringifyNode(ds.cur, false)
		if ds.cur.NextSibling == nil || !isMeta(ds.cur.NextSibling) {
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
		if ds.lastNode != nil && nodes.IsHeader(ds.lastNode.Type()) {
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
func header(ds *docState) nodes.Node {
	ds.push(nil)
	n := parseSubtree(ds)
	ds.pop()
	if len(n) == 0 {
		return nil
	}
	nn := nodes.NewHeaderNode(headerLevel[ds.cur.DataAtom], n...)
	switch strings.ToLower(stringifyNode(ds.cur, true)) {
	case headerLearn, headerCover:
		nn.MutateType(nodes.NodeHeaderCheck)
	case headerFAQ:
		nn.MutateType(nodes.NodeHeaderFAQ)
	}
	ds.env = nil
	return nn
}

// aside produces an infobox.
func aside(ds *docState) nodes.Node {
	kind := nodes.InfoboxPositive
	for _, v := range ds.cur.Attr {
		// If class "negative" is given, set the infobox type.
		if v.Key == "class" && v.Val == "negative" {
			kind = nodes.InfoboxNegative
		}
	}

	ds.push(nil)
	nn := parseSubtree(ds)
	nn = parser.BlockNodes(nn)
	nn = parser.CompactNodes(nn)
	ds.pop()
	if len(nn) == 0 {
		return nil
	}
	return nodes.NewInfoboxNode(kind, nn...)
}

// new style aside, to produce an infobox
func newAside(ds *docState) nodes.Node {
	kind := nodes.InfoboxPositive
	s := ds.cur.FirstChild.NextSibling.FirstChild.Data
	if strings.HasPrefix(s, "aside negative") {
		ds.cur.FirstChild.NextSibling.FirstChild.Data = strings.TrimPrefix(s, "aside negative")
		kind = nodes.InfoboxNegative
	} else {
		ds.cur.FirstChild.NextSibling.FirstChild.Data = strings.TrimPrefix(s, "aside positive")
	}

	ds.push(nil)
	nn := parseSubtree(ds)
	nn = parser.BlockNodes(nn)
	nn = parser.CompactNodes(nn)
	ds.pop()
	if len(nn) == 0 {
		return nil
	}
	return nodes.NewInfoboxNode(kind, nn...)
}

// infobox doesn't have a block parent.
func infobox(ds *docState) nodes.Node {
	negativeInfoBox := isInfoboxNegative(ds.cur)
	// iterate twice on next sibling as there is a \n node in between
	ds.cur = ds.cur.NextSibling.NextSibling
	ds.push(nil)
	nn := parseSubtree(ds)
	nn = parser.BlockNodes(nn)
	nn = parser.CompactNodes(nn)
	ds.pop()
	if len(nn) == 0 {
		return nil
	}
	kind := nodes.InfoboxPositive
	if negativeInfoBox {
		kind = nodes.InfoboxNegative
	}
	return nodes.NewInfoboxNode(kind, nn...)
}

// table parses an arbitrary <table> element and its children.
// It may return other elements if the table is just a wrap.
func table(ds *docState) nodes.Node {
	var rows [][]*nodes.GridCell
	for _, tr := range findChildAtoms(ds.cur, atom.Tr) {
		ds.push(tr)
		r := tableRow(ds)
		ds.pop()
		rows = append(rows, r)
	}
	if len(rows) == 0 {
		return nil
	}
	return nodes.NewGridNode(rows...)
}

func tableRow(ds *docState) []*nodes.GridCell {
	var row []*nodes.GridCell
	firstChild := findAtom(ds.cur, atom.Td)
	// If there is no Td child found, could be table header so look for Th
	if firstChild == nil {
		firstChild = findAtom(ds.cur, atom.Th)
	}

	for td := firstChild; td != nil; td = td.NextSibling {
		if td.DataAtom != atom.Td && td.DataAtom != atom.Th {
			continue
		}
		ds.push(td)
		nn := parseSubtree(ds)
		nn = parser.BlockNodes(nn)
		nn = parser.CompactNodes(nn)
		ds.pop()
		cs, err := strconv.Atoi(nodeAttr(td, "colspan"))
		if err != nil {
			cs = 1
			for ns := td.NextSibling; ns != nil; ns = ns.NextSibling {
				if ns.DataAtom != atom.Td && ns.DataAtom != atom.Th {
					continue
				}
				if ns.FirstChild != nil {
					break
				}
				cs++
			}
		}
		rs, err := strconv.Atoi(nodeAttr(td, "rowspan"))
		if err != nil {
			rs = 1
		}
		cell := &nodes.GridCell{
			Colspan: cs,
			Rowspan: rs,
			Content: nodes.NewListNode(nn...),
		}
		row = append(row, cell)
	}
	return row
}

// survey expects 1 or more name Nodes followed by 1 or more input Nodes.
// Each input node is expected to have a value attribute.
func survey(ds *docState) nodes.Node {
	var gg []*nodes.SurveyGroup
	ns := findChildAtoms(ds.cur, atom.Name)
	for _, n := range ns {
		var inputs []*html.Node
		for hn := n.NextSibling; hn != nil; hn = hn.NextSibling {
			if hn.DataAtom == atom.Input {
				inputs = append(inputs, hn)
			} else if hn.DataAtom == atom.Name {
				break
			}
		}
		opt := surveyOpt(inputs)
		if len(opt) > 0 {
			gg = append(gg, &nodes.SurveyGroup{
				Name:    strings.TrimSpace(n.FirstChild.Data),
				Options: opt,
			})
		}
	}
	if len(gg) == 0 {
		return nil
	}
	ds.survey++
	id := fmt.Sprintf("%s-%d", ds.clab.ID, ds.survey)
	return nodes.NewSurveyNode(id, gg...)
}

func surveyOpt(inputs []*html.Node) []string {
	var opt []string
	for _, input := range inputs {
		for _, attr := range input.Attr {
			if attr.Key == "value" {
				opt = append(opt, attr.Val)
			}
		}
	}
	return opt
}

// code parses hn as inline or block codes.
// Inline code node will be of type NodeText.
func code(ds *docState, term bool) nodes.Node {
	elem := findNearestAncestor(ds.cur, map[atom.Atom]struct{}{atom.Pre: {}}, doConsiderSelf)
	// inline <code> text
	if elem == nil {
		return text(ds)
	}
	// block code or terminal
	v := stringifyNode(ds.cur, false)
	if v == "" {
		if countDirect(ds.cur.Parent) > 1 {
			return nil
		}
		v = "\n"
	} else if ds.cur.Parent.FirstChild == ds.cur && ds.cur.Parent.DataAtom != atom.Span {
		v = "\n" + v
	}
	// get the language hint
	var lan string
	if !term {
		for _, a := range ds.cur.Attr {
			if a.Key == "class" && strings.HasPrefix(a.Val, "language-") {
				lan = strings.Replace(a.Val, "language-", "", 0)
			}
		}
	}
	n := nodes.NewCodeNode(v, term, lan)
	n.MutateBlock(elem)
	return n
}

// list parses <ul> and <ol> lists.
// It returns nil if the list has no items.
func list(ds *docState) nodes.Node {
	typ := nodeAttr(ds.cur, "type")
	if ds.cur.DataAtom == atom.Ol && typ == "" {
		typ = "1"
	}
	start, _ := strconv.Atoi(nodeAttr(ds.cur, "start"))
	list := nodes.NewItemsListNode(typ, start)
	for hn := findAtom(ds.cur, atom.Li); hn != nil; hn = hn.NextSibling {
		if hn.DataAtom != atom.Li {
			continue
		}
		ds.push(hn)
		nn := parseSubtree(ds)
		nn = parser.CompactNodes(nn)
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
		case nodes.NodeHeaderCheck:
			list.MutateType(nodes.NodeItemsCheck)
		case nodes.NodeHeaderFAQ:
			list.MutateType(nodes.NodeItemsFAQ)
		}
	}
	return list
}

// image creates a new ImageNode out of hn, parsing its src attribute.
// It returns nil if src is empty.
// It may also return a YouTubeNode if alt property contains specific substring.
func image(ds *docState) nodes.Node {
	alt := nodeAttr(ds.cur, "alt")
	// Author-added double quotes in attributes break html syntax
	alt = html.EscapeString(alt)
	if strings.Contains(alt, "youtube.com/watch") {
		return youtube(ds)
	} else if strings.Contains(alt, "https://") {
		u, err := url.Parse(alt)
		if err != nil {
			return nil
		}
		// For iframe, make sure URL ends in allowlisted domain.
		ok := false
		for _, domain := range nodes.IframeAllowlist {
			if u.Hostname() == domain {
				ok = true
				break
			}
		}
		if ok {
			return iframe(ds)
		}
	}
	s := nodeAttr(ds.cur, "src")
	if s == "" {
		return nil
	}

	n := nodes.NewImageNode(nodes.NewImageNodeOptions{Src: s})

	if alt != "" {
		n.Alt = alt
	}

	if title := nodeAttr(ds.cur, "title"); title != "" {
		// Author-added double quotes in attributes break html syntax
		n.Title = html.EscapeString(title)
	}

	if ws := nodeAttr(ds.cur, "width"); ws != "" {
		w, err := strconv.ParseFloat(ws, 64)
		if err != nil {
			return nil
		}
		n.Width = float32(w)
	}

	n.MutateBlock(findNearestBlockAncestor(ds.cur))
	return n
}

func youtube(ds *docState) nodes.Node {
	for _, attr := range ds.cur.Attr {
		if attr.Key == "id" {
			n := nodes.NewYouTubeNode(attr.Val)
			n.MutateBlock(true)
			return n
		}
	}
	return nil
}

func fragmentImport(ds *docState) nodes.Node {
	if url := strings.TrimPrefix(ds.cur.Data, convertedImportsDataPrefix); url != "" {
		return nodes.NewImportNode(url)
	}

	return nil
}

func iframe(ds *docState) nodes.Node {
	u, err := url.Parse(nodeAttr(ds.cur, "alt"))
	if err != nil {
		return nil
	}
	// Allow only https.
	if u.Scheme != "https" {
		return nil
	}
	n := nodes.NewIframeNode(u.String())
	n.MutateBlock(true)
	return n
}

// button returns either a text node, if no <a> child element is present,
// or link node, containing the button.
// It returns nil if no content nodes are present.
func button(ds *docState) nodes.Node {
	a := findAtom(ds.cur, atom.A)
	if a == nil {
		return text(ds)
	}
	href := nodeAttr(a, "href")
	if href == "" {
		return nil
	}

	ds.push(a)
	n := parseSubtree(ds)
	ds.pop()
	if len(n) == 0 {
		return nil
	}

	s := strings.ToLower(stringifyNode(a, true))
	dl := strings.HasPrefix(s, "download ")
	btn := nodes.NewButtonNode(true, true, dl, n...)

	ln := nodes.NewURLNode(href, btn)
	ln.MutateBlock(findNearestBlockAncestor(ds.cur))
	return ln
}

// Link creates a URLNode out of hn, parsing href and name attributes.
// It returns nil if hn contents is empty.
// The resuling link's content is always a single text node.
func link(ds *docState) nodes.Node {
	href := nodeAttr(ds.cur, "href")

	ds.push(nil)
	parsedChildNodes := parseSubtree(ds)
	ds.pop()

	// Check outside styles
	outsideBold := isBold(ds.cur.Parent)
	outsideItalic := isItalic(ds.cur.Parent)
	if isBoldAndItalic(ds.cur.Parent) {
		outsideBold = true
		outsideItalic = true
	}
	// Apply outside styles to inside parsed (text) nodes
	for _, node := range parsedChildNodes {
		if textNode, ok := node.(*nodes.TextNode); ok {
			textNode.Bold = textNode.Bold || outsideBold
			textNode.Italic = textNode.Italic || outsideItalic
		}
	}

	n := nodes.NewURLNode(href, parsedChildNodes...)
	n.Name = nodeAttr(ds.cur, "name")
	if v := nodeAttr(ds.cur, "target"); v != "" {
		n.Target = v
	}
	n.MutateBlock(findNearestBlockAncestor(ds.cur))
	return n
}

// text creates a TextNode using hn.Data as contents.
// It returns nil if hn.Data is empty or contains only space runes.
func text(ds *docState) nodes.Node {
	bold := isBold(ds.cur)
	italic := isItalic(ds.cur)
	// We must call this to look up an extra level in the node tree to obtain both styles
	if isBoldAndItalic(ds.cur) {
		bold = true
		italic = true
	}
	code := isCode(ds.cur) || isConsole(ds.cur)

	// TODO: verify whether this actually does anything
	if a := findAtom(ds.cur, atom.A); a != nil {
		ds.push(a)
		l := link(ds)
		ds.pop()
		if l != nil {
			l.MutateBlock(findNearestBlockAncestor(ds.cur))
			return l
		}
	}

	n := nodes.NewTextNode(nodes.NewTextNodeOptions{
		Value:  stringifyNode(ds.cur, false),
		Bold:   bold,
		Italic: italic,
		Code:   code,
	})
	n.MutateBlock(findNearestBlockAncestor(ds.cur))
	return n
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

// roundDuration rounds time to the nearest minute, always rounding
// up when there is any fractional portion of a minute.
// Ex:
//  59s --> 1m
//  60s --> 1m
//  61s --> 2m
func roundDuration(d time.Duration) time.Duration {
	rd := time.Duration(d.Minutes()) * time.Minute
	if rd < d {
		rd += time.Minute
	}
	return rd
}

func convertImports(content []byte) []byte {
	slices := bytes.Split(content, []byte("\n"))
	escaped := [][]byte{}
	for _, slice := range slices {
		if matches := importsTagRegexp.FindSubmatch(slice); len(matches) > 0 {
			if len(matches) > 1 {
				url := string(matches[1])
				slice = bytes.Join([][]byte{
					convertedImportsPrefix,
					[]byte(html.EscapeString(url)),
					convertedImportsSuffix,
				}, []byte(""))
			}
		}

		escaped = append(escaped, slice)
	}

	return bytes.Join(escaped, []byte("\n"))
}

func hasImport(ds *docState) bool {
	for _, step := range ds.clab.Steps {
		if len(nodes.ImportNodes(step.Content.Nodes)) > 0 {
			return true
		}
	}

	return false
}
