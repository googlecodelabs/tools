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
// the Devsite Markdown processor.
package md

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/googlecodelabs/tools/claat/parser"
	"github.com/googlecodelabs/tools/claat/types"
	"github.com/russross/blackfriday"
)

const (
	metaAuthor           = "author"
	metaSummary          = "summary"
	metaID               = "id"
	metaCategories       = "categories"
	metaEnvironments     = "environments"
	metaStatus           = "status"
	metaFeedbackLink     = "feedback link"
	metaAnalyticsAccount = "analytics account"
	metaTags             = "tags"
)

var metadataRegexp = regexp.MustCompile(`(.+?):(.+)`)
var languageRegexp = regexp.MustCompile(`language-(.+)`)
var durationHintRegexp = regexp.MustCompile(`(?i)Duration:? (.+)`)
var durationRegexp = regexp.MustCompile(`(\d+)[:.](\d{2})$`)
var downloadButtonRegexp = regexp.MustCompile(`^(?i)Download(.+)$`)

// init registers this parser so it is available to CLaaT.
func init() {
	parser.Register("md", &Parser{})
}

// Parser is a Markdown parser.
type Parser struct {
}

// Parse parses a codelab writtet in Markdown.
func (p *Parser) Parse(r io.Reader) (*types.Codelab, error) {
	// Convert Markdown to HTML for easy parsing.
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	b = claatMarkdown(b)
	h := bytes.NewBuffer(b)
	// Parse the markup.
	return parseMarkup(h)
}

// ParseFragment parses a codelab fragment writtet in Markdown.
func (p *Parser) ParseFragment(r io.Reader) ([]types.Node, error) {
	return nil, errors.New("fragment parser not implemented")
}

// parserState encapsulates the state of the parser at any given step.
type parserState struct {
	tzr *html.Tokenizer
	c   *types.Codelab
	t   html.Token

	currentStep *types.Step
}

// emit accepts a node, and either writes the node directly to the current step, or writes the node to the node buffer.
func (ps *parserState) emit(n types.Node) {
	ps.currentStep.Content.Append(n)
}

// advance moves the tokenizer to the next token and updates the token convenience variable.
func (ps *parserState) advance() {
	ps.tzr.Next()
	ps.t = ps.tzr.Token()
}

// multiAdvance is a convenience method for repeatedly advancing.
func (ps *parserState) multiAdvance(n int) {
	for i := 0; i < n; i++ {
		ps.advance()
	}
}

// claatMarkdown calls the Blackfriday Markdown parser with some special addons selected. It takes a byte slice as a parameter,
// and returns its result as a byte slice.
func claatMarkdown(b []byte) []byte {
	htmlFlags := blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	extns := blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_DEFINITION_LISTS |
		blackfriday.EXTENSION_TABLES
	o := blackfriday.Options{
		Extensions: extns,
	}
	r := blackfriday.HtmlRenderer(htmlFlags, "", "")
	return blackfriday.MarkdownOptions(b, r, o)
}

// parseMarkup accepts an io.Reader to markup created by the Devsite Markdown parser. It returns a pointer to a codelab object, or an error if one occurs.
func parseMarkup(markup io.Reader) (*types.Codelab, error) {
	// Avoid global vars by encapsulating state.
	ps := parserState{
		tzr: html.NewTokenizer(markup),
		c:   &types.Codelab{},
	}

	var inStepTitle bool

	// Advance through the tokenized input, one token at a time.
	for ps.advance(); ps.t.Type != html.ErrorToken; ps.advance() {
		if ps.c.Title == "" {
			// If we have a <p> tag, it's a metadata section.
			if ps.t.Type == html.StartTagToken && ps.t.DataAtom == atom.P {
				if err := parseMetadata(&ps); err != nil {
					return nil, err
				}
			}
			// We need to handle the title in this pass, as the next advance call will move us past the <h1>.
			if ps.t.Type == html.StartTagToken && ps.t.DataAtom == atom.H1 {
				handleCodelabTitle(&ps)

				// If we just finished parsing a step or the title, we are left possibly pointing to the opening
				// <h2> of another step. Update the flag accordingly.
				inStepTitle = (ps.t.Type == html.StartTagToken && ps.t.DataAtom == atom.H2)
			}
			continue
		}
		if inStepTitle {
			// This is the beginning of a new codelab step, and we are pointing to the contents of the title.
			stepTitle := ps.t.Data
			ps.advance()
			// Emit a step object.
			ps.currentStep = ps.c.NewStep(stepTitle)
			parseStep(&ps)

			// If we just finished parsing a step or the title, we are left possibly pointing to the opening
			// <h2> of another step. Update the flag accordingly.
			inStepTitle = (ps.t.Type == html.StartTagToken && ps.t.DataAtom == atom.H2)
		} else {
			// If we had some intermediate blank lines between step, and are out of one, check if we don't reenter
			// into a new step.
			inStepTitle = (ps.t.Type == html.StartTagToken && ps.t.DataAtom == atom.H2)
		}

	}
	// If not EOF, an error occurred in tokenization.
	if err := ps.tzr.Err(); err != io.EOF {
		return nil, err
	}

	finalizeCodelab(&ps)
	return ps.c, nil
}

// parseMetadata handles the metadata section preceding a codelab.
// It assumes the tokenizer is pointing to the first <p>/.
// It returns any errors it encounters, and leaves the tokenizer pointing at the <h1>
// starting the codelab title, or at io.EOF.
func parseMetadata(ps *parserState) error {
	m := map[string]string{}
	// Iterate over the metadata elements, constructing a map of the metadata.
	for ; ps.t.Type != html.ErrorToken && ps.t.DataAtom != atom.H1; ps.advance() {
		if ps.t.Type == html.StartTagToken && ps.t.DataAtom == atom.P {
			// Advance to text.
			ps.advance()
			// Split the keys from values.
			s := metadataRegexp.FindStringSubmatch(ps.t.Data)
			if len(s) != 3 {
				return fmt.Errorf("invalid metadata format: %v", s)
			}
			k := strings.ToLower(strings.TrimSpace(s[1]))
			v := strings.TrimSpace(s[2])
			m[k] = v
		}
	}
	addMetadataToCodelab(m, ps.c)
	return nil
}

// parseStep handles an entire step in a codelab. It assumes the tokenizer is pointing to the </h2> that ends the step's title.
// It returns any errors it encounters, and leaves the tokenizer pointing at the <h2> starting the next step, or at io.EOF.
func parseStep(ps *parserState) error {
	// Initially check for a duration tag.
	ps.advance()
	//Sometimes raw newlines are injected into the output before the <p> - we need to bypass them if present.
	if ps.t.Type == html.TextToken {
		ps.advance()
	}
	if ps.t.Type == html.StartTagToken && ps.t.DataAtom == atom.P {
		if err := handleDurationHint(ps); err != nil {
			return err
		}
	}

	// Track text styling settings.
	var bold, italic bool

	// Continue reading tokens in order, stopping on an error or the beginning of another step.
	for ; ps.t.Type != html.ErrorToken && !(ps.t.Type == html.StartTagToken && ps.t.DataAtom == atom.H2); ps.advance() {

		// Handle <h3> through <h6>.
		if ps.t.Type == html.StartTagToken && (ps.t.DataAtom == atom.H3 || ps.t.DataAtom == atom.H4 || ps.t.DataAtom == atom.H5 || ps.t.DataAtom == atom.H6) {
			handleHeader(ps)
		}
		// Handle <pre>.
		if ps.t.Type == html.StartTagToken && ps.t.DataAtom == atom.Pre {
			handleFencedCodeBlock(ps)
		}
		// Handle <code>.
		if ps.t.Type == html.StartTagToken && ps.t.DataAtom == atom.Code {
			handleInlineCodeBlock(ps)
		}
		// Handle <ul> and <ol>.
		if ps.t.Type == html.StartTagToken && (ps.t.DataAtom == atom.Ul || ps.t.DataAtom == atom.Ol) {
			handleList(ps)
		}
		// Handle <dt>.
		if ps.t.Type == html.StartTagToken && ps.t.DataAtom == atom.Dt {
			handleInfobox(ps)
		}
		// Handle <em>.
		if ps.t.DataAtom == atom.Em {
			italic = ps.t.Type == html.StartTagToken
		}
		// Hande <strong>.
		if ps.t.DataAtom == atom.Strong {
			bold = ps.t.Type == html.StartTagToken
		}
		// Handle <img>.
		if ps.t.DataAtom == atom.Img {
			handleImage(ps)
		}
		// Handle <a>.
		if ps.t.DataAtom == atom.A && ps.t.Type == html.StartTagToken {
			handleLink(ps)
		}
		// Handle text.
		if ps.t.Type == html.TextToken {
			n := newBreaklessTextNode(ps.t.Data)
			n.Bold = bold
			n.Italic = italic
			ps.emit(n)
		}
	}
	return nil
}

// handleCodelabTitle takes care of setting the title for the codelab. It assumes the tokenizer is pointing to <h1>.
func handleCodelabTitle(ps *parserState) {
	ps.advance()
	ps.c.Title = ps.t.Data
	ps.advance()
}

// processDuration expects a string expressing a duration in hours and minutes, delimited by a colon or a period, such as 1:30, 2.45, or 0:90.
// It returns an equivalent time.Duration, or an error if one occurs.
func processDuration(d string) (time.Duration, error) {
	s := durationRegexp.FindStringSubmatch(d)
	if len(s) == 3 {
		h, err := strconv.ParseInt(s[1], 10, 64)
		if err != nil {
			return 0, err
		}
		m, err := strconv.ParseInt(s[2], 10, 64)
		if err != nil {
			return 0, err
		}
		return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute, nil
	}
	if d == "0" {
		return 0, nil
	}
	return 0, errors.New("unrecognized duration string")
}

// handleDurationHint parses the optional duration string at the beginning of a codelab step. It assumes the tokenizer is
// pointing at the inital <p> of the step. It returns any errors it encounters in the process.
func handleDurationHint(ps *parserState) error {
	ps.advance()
	// If this isn't text, this also isnt a duration string.
	if ps.t.Type != html.TextToken {
		return nil
	}
	s := durationHintRegexp.FindStringSubmatch(ps.t.Data)
	// This is possibly not a duration string, so bail out if we don't have strong indications that it is.
	var err error
	if len(s) >= 2 {
		ps.currentStep.Duration, err = processDuration(s[1])
		if err != nil {
			return err
		}
	}
	ps.advance()
	// Now we're on the closing p tag of the duration string.
	return nil
}

// finalizeCodelab takes care of all work that should be performed after the entire input is parsed.
// It takes a pointer to a parserState object, and acts on the codelab referenced in it.
func finalizeCodelab(ps *parserState) {
	computeTotalDuration(ps.c)
}

// computeTotalDuration computes the total duration for a codelab by summing the duration of each step.
// It takes a pointer to a codelab, and acts on that codelab.
func computeTotalDuration(c *types.Codelab) {
	for _, v := range c.Steps {
		c.Duration += int(v.Duration.Minutes())
	}
}

// handleHeader handles header tags, h3-h6. It assumes the tokenizer is pointing to <h_>.
// It takes a pointer to a parserState object.
func handleHeader(ps *parserState) {
	var l int
	switch ps.t.DataAtom {
	case atom.H3:
		l = 3
		break
	case atom.H4:
		l = 4
		break
	case atom.H5:
		l = 5
		break
	case atom.H6:
		l = 6
		break
	}
	ps.advance()
	n := types.NewHeaderNode(l, newBreaklessTextNode(ps.t.Data))
	ps.emit(n)
	ps.advance() // Now we are on the closing tag.
}

// handleList handles both ordered and unordered lists. It assumes the tokenizer is pointing to <ul> or <ol>.
// It takes a pointer to a parserState object.
func handleList(ps *parserState) {
	start := 0
	if ps.t.DataAtom == atom.Ol {
		start = 1
	}
	iln := types.NewItemsListNode("", start)

	for ps.advance(); ps.t.Type != html.ErrorToken && !(ps.t.Type == html.EndTagToken && (ps.t.DataAtom == atom.Ul || ps.t.DataAtom == atom.Ol)); ps.advance() {
		if ps.t.Type == html.StartTagToken && ps.t.DataAtom == atom.Li {
			ps.advance()
			iln.NewItem(newBreaklessTextNode(ps.t.Data))
		}
	}
	ps.emit(iln)
}

// handleInfobox handles the colored call-out boxes in codelabs. It assumes the tokenizer is pointing to <dt>.
func handleInfobox(ps *parserState) {
	// Advance to <dt>'s text content.
	ps.advance()
	// Deduce the kind of infobox.
	var kind types.InfoboxKind
	sentiment := strings.ToLower(ps.t.Data)
	if sentiment == "positive" {
		kind = types.InfoboxPositive
	} else if sentiment == "negative" {
		kind = types.InfoboxNegative
	}

	// Advance to <dd>'s text content.
	ps.multiAdvance(4)
	n := types.NewInfoboxNode(kind, newBreaklessTextNode(ps.t.Data))
	ps.emit(n)

	// Advance to </dd>.
	ps.advance()
}

// handleImage handles <img> tags. It assumes the tokenizer is pointing to the <img> tag itself.
func handleImage(ps *parserState) {
	for _, v := range ps.t.Attr {
		if v.Key == "src" {
			ps.emit(types.NewImageNode(v.Val))
			break
		}
	}
}

// handleLink handles links and download buttons, both of which appear as <a> elements.
// It assumes the tokenizer is pointing to the <a> tag itself.
func handleLink(ps *parserState) {
	var href string
	for _, v := range ps.t.Attr {
		if v.Key == "href" {
			href = v.Val
		}
	}
	// Advance to text.
	ps.advance()
	// Check for the download button case.
	s := downloadButtonRegexp.FindStringSubmatch(ps.t.Data)
	if len(s) >= 2 {
		// It's a button, emit a button element with all the pretty styling enabled.
		ps.emit(types.NewButtonNode(true, true, true, newBreaklessTextNode(s[1])))
	} else {
		// It's not a button, emit an ordinary link.
		ps.emit(types.NewURLNode(href, newBreaklessTextNode(ps.t.Data)))
	}
	// Advance to </a>.
	ps.advance()
}

// handleFencedCodeBlock handles all code elements wrapped in ```s.
// It assumes the tokenizer is pointing to the <pre> tag establishing the block.
func handleFencedCodeBlock(ps *parserState) {
	// Advance to <code>.
	ps.advance()
	// Check for the presence of a language hint.
	var lang string
	for _, v := range ps.t.Attr {
		if v.Key == "class" {
			// Try to extract a valid language string from the class.
			s := languageRegexp.FindStringSubmatch(v.Val)
			if len(s) == 2 {
				lang = s[1]
			}
		}
	}
	// Advance to text content.
	ps.advance()
	n := types.NewCodeNode(ps.t.Data, false)
	n.Lang = lang
	ps.emit(n)
	// Advance to </pre>.
	ps.multiAdvance(2)
}

// It assumes the tokenizer is pointing to the <code> tag establishing the block.
func handleInlineCodeBlock(ps *parserState) {
	// Advance to text content.
	ps.advance()
	// Inlined code is actually a text node with special formatting.
	n := types.NewTextNode(ps.t.Data)
	n.Code = true
	ps.emit(n)
	// Advance to </code>.
	ps.advance()
}

// standardSplit takes a string, splits it along a comma delimiter, then on each fragment, trims Unicode spaces
// from both ends and converts them to lowercase. It returns a slice of the processed strings.
func standardSplit(s string) []string {
	strs := strings.Split(s, ",")
	for k, v := range strs {
		strs[k] = strings.ToLower(strings.TrimSpace(v))
	}
	return strs
}

// addMetadataToCodelab takes a map of strings to strings, and a pointer to a Codelab. It reads the keys of the map,
// and assigns the values to any keys that match a codelab metadata field as defined by the meta* constants.
func addMetadataToCodelab(m map[string]string, c *types.Codelab) {
	for k, v := range m {
		switch k {
		case metaAuthor:
			// Directly assign the summary to the codelab field.
			c.Author = v
		case metaSummary:
			// Directly assign the summary to the codelab field.
			c.Summary = v
			break
		case metaID:
			// Directly assign the ID to the codelab field.
			c.ID = v
			break
		case metaCategories:
			// Standardize the categories and append to codelab field.
			c.Categories = append(c.Categories, standardSplit(v)...)
			break
		case metaEnvironments:
			// Standardize the tags and append to the codelab field.
			c.Tags = append(c.Tags, standardSplit(v)...)
			break
		case metaStatus:
			// Standardize the statuses and append to the codelab field.
			statuses := standardSplit(v)
			statusesAsLegacy := types.LegacyStatus(statuses)
			c.Status = &statusesAsLegacy
			break
		case metaFeedbackLink:
			// Directly assign the feedback link to the codelab field.
			c.Feedback = v
			break
		case metaAnalyticsAccount:
			// Directly assign the GA id to the codelab field.
			c.GA = v
			break
		case metaTags:
			// Standardize the tags and append to the codelab field.
			c.Tags = append(c.Tags, standardSplit(v)...)
			break
		default:
			break
		}
	}
}

// newBreaklessTextNode accepts a string, and constructs a new TextNode containing the string,
// but replaces all line breaks in the string with spaces first. It returns a pointer to the created node.
func newBreaklessTextNode(s string) *types.TextNode {
	s = strings.Replace(s, "\n", " ", -1)
	return types.NewTextNode(s)
}
