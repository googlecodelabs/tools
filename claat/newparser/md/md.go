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

type parseMode int

const (
	preludeBase parseMode = iota
	codelabTitle
	codelabMetadata

	stepBase
	stepTitle
)

type parser struct {
	currentNode *html.Node

	metadata map[string]string
	codelab  *types.Codelab

	mode parseMode
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
	p.currentNode = bodyNode

	// Start the recursion
	p.processCurrentNode()

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

func (p *parser) inPrelude() bool {
	switch p.mode {
	case preludeBase, codelabTitle, codelabMetadata:
		return true
	default:
		return false
	}
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

func (p *parser) processCurrentNode() error {
	fmt.Printf("%+v\n", p.currentNode.DataAtom)
	// Test for errors.
	if p.currentNode.Type == html.ErrorNode {
		return errors.New("encountered an error node while parsing")
	}

	var popAction func()
	if p.inPrelude() {
		popAction = p.processPreludeNode()
	} else {
		popAction = p.processStepNode()
	}

	for next := p.currentNode.FirstChild; next != nil; next = next.NextSibling {
		p.currentNode = next
		p.processCurrentNode()
	}

	if popAction != nil {
		popAction()
	}

	return nil
}

func (p *parser) processPreludeNode() func() {
	var popAction func()
	// Handle the node based on the current mode.
	switch p.mode {
	case codelabTitle:
		// We're handling the codelab title. Pass text along and ignore everything else.
		if p.currentNode.Type == html.TextNode {
			p.codelab.Title += p.currentNode.Data
		}
	case codelabMetadata:
		p.metadata = merge(p.metadata, extractMetadata(p.currentNode.Data))
		break
	case preludeBase:
		switch p.currentNode.DataAtom {
		case atom.H1:
			p.mode = codelabTitle
			popAction = func() {
				p.mode = preludeBase
			}
		case atom.P:
			p.mode = codelabMetadata
			popAction = func() {
				p.mode = preludeBase
			}
		case atom.H2:
			// The start of the first step.
			popAction = p.processStepNode()
		}
	}

	// Nothing found.
	return popAction
}

func (p *parser) processStepNode() func() {
	// TODO actually handle step content
	var popAction func()
	switch p.mode {
	case stepTitle:
		if p.currentNode.Type == html.TextNode {
			p.codelab.Steps[len(p.codelab.Steps)-1].Title += p.currentNode.Data
		}
	}

	if p.currentNode.DataAtom == atom.H2 {
		p.codelab.NewStep("")
		p.mode = stepTitle
		popAction = func() {
			p.mode = stepBase
		}
	}
	return popAction
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
