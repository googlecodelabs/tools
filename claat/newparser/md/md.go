package md

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/googlecodelabs/tools/claat/types"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"gopkg.in/russross/blackfriday.v2"
)

type preludeParseMode int

const (
	normal preludeParseMode = iota
	title
	metadata
)

type parser struct {
	metadata map[string]string
	codelab  *types.Codelab

	currentPreludeMode preludeParseMode
}

func newParser() *parser {
	return &parser{
		metadata: map[string]string{},
		codelab:  types.NewCodelab(),
	}
}

// Parse parses a codelab written in the new Markdown format.
// It takes in the input bytes as an io.Reader.
// It returns the parsed codelab, the metadata as a separate map, or an error if one occurs.
func Parse(in io.Reader) (*types.Codelab, map[string]string, error) {
	inSlice, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, nil, err
	}

	rawHTML := blackfriday.Run(inSlice)
	rawHTMLbuf := bytes.NewBuffer(blackfriday.Run(rawHTML))

	rootNode, err := html.Parse(rawHTMLbuf)
	if err != nil {
		return nil, nil, err
	}

	// find the <body> tag which is the earliest thing we care about
	bodyNode, err := seek(rootNode, atom.Body)
	if err != nil {
		return nil, nil, err
	}

	p := newParser()

	// Parse the prelude to locate the first step.
	step, err := p.parsePrelude(bodyNode)
	if err != nil {
		return nil, nil, err
	}

	// Iterate over steps until there are none left.
	for step != nil {
		p.codelab.NewStep("")
		step, err = p.parseStep(step)
		if err != nil {
			return nil, nil, err
		}
	}

	return p.codelab, p.metadata, nil
}

// seek takes an *html.Node rooting an HTML tree, and an atom.Atom.
// It preorder-traverses the tree until the first occurrence of a node of the given atom.
// It returns that node, or an error if no such nodes are found.
func seek(in *html.Node, a atom.Atom) (*html.Node, error) {
	if in.DataAtom == a {
		return in, nil
	}
	for next := in.FirstChild; next != nil; next = next.NextSibling {
		result, err := seek(next, a)
		// i.e. if a result is found
		if err == nil {
			return result, nil
		}
	}
	return nil, fmt.Errorf("could not find a node with atom %q", a.String())
}

// parsePrelude expects an *html.Node that is the <body> of the processed Markdown.
// It handles the prelude, which is all codelab content occurring prior to the first step.
// This is done by preorder-traversing the tree, handling nodes as they are encountered, and stopping at the first <h2>.
// If it finds the <h2>, it will return that <h2>'s node.
// If it encounters an error before finding the <h2>, it will return an error.
// If no <h2> or error is found, it will return nil for both values.
func (p *parser) parsePrelude(node *html.Node) (*html.Node, error) {
	// Test for errors.
	if node.Type == html.ErrorNode {
		return nil, errors.New("encountered an error node while parsing")
	}

	// Handle the node based on the current mode.
	switch p.currentPreludeMode {
	case title:
		// We're handling the codelab title. Pass text along and ignore everything else.
		if node.Type == html.TextNode {
			p.codelab.Title += node.Data
		}
	case metadata:
		p.metadata = merge(p.metadata, extractMetadata(node.Data))
		break
	case normal:
		switch node.DataAtom {
		case atom.H1:
			p.currentPreludeMode = title
		case atom.P:
			p.currentPreludeMode = metadata
		case atom.H2:
			// The start of the first step.
			return node, nil
		}
	}

	// The node itself has been handled; now we handle the child nodes recursively.
	for next := node.FirstChild; next != nil; next = next.NextSibling {
		result, err := p.parsePrelude(next)
		if err != nil || result != nil {
			return result, err
		}
	}

	// Exit prelude modes.
	if p.currentPreludeMode == title && node.DataAtom == atom.H1 || p.currentPreludeMode == metadata && node.DataAtom == atom.P {
		p.currentPreludeMode = normal
	}

	// Nothing found.
	return nil, nil
}

// Merges two string:string maps.
// Behavior is undefined when the input maps have two different values for a single key k.
// Neither input map is modified.
func merge(m1 map[string]string, m2 map[string]string) map[string]string {
	m := map[string]string{}
	for k, v := range m1 {
		m[k] = v
	}
	for k, v := range m2 {
		m[k] = v
	}
	return m
}

// extractMetadata extracts metadata key-value pairs from string input.
// If no metadata is found, the empty map is returned.
func extractMetadata(s string) map[string]string {
	meta := map[string]string{}
	// We're handling something that might be metadata.
	// First, divide the possible metadata by lines.
	potentialMetadataLines := strings.Split(s, "\n")
	for _, l := range potentialMetadataLines {
		// Metadata values are formatted like key:value.
		metadatum := strings.SplitN(l, ":", 2)
		// If the line fits this description, it's metadata.
		if len(metadatum) == 2 {
			meta[strings.TrimSpace(metadatum[0])] = strings.TrimSpace(metadatum[1])
		}
	}
	return meta
}

// parseStep expects an *html.Node that is an <h2> in the processed Markdown.
// It handles the step defined by the given node, up to the beginning of the next step, or the end of the content.
// This is done by the preorder-traversing the tree, handling nodes as they are encountered, and stopping at the next <h2>, or at the end of input.
// If another step comes after this one, this function will return the <h2> starting that step.
// If the end of the content is found, it will return nil.
// If it encounters an error before finding an <h2> or the end of the content, it will return an error.
func (p *parser) parseStep(node *html.Node) (*html.Node, error) {

	// TODO: actually handle step content

	var nextStep *html.Node

	// We don't really handle step content yet, so we can skip ahead.
	for next := node.FirstChild; next != nil; next = next.NextSibling {
		nextStep, err := seek(next, atom.H2)
		if err != nil {
			return nil, err
		}
		if nextStep != nil {
			break
		}
	}

	return nextStep, nil
}
