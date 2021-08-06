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
	"fmt"
	"io"

	"github.com/googlecodelabs/tools/claat/nodes"
	"github.com/googlecodelabs/tools/claat/types"
)

// Parser parses a codelab in specific resource format.
// Each parser needs to call Register to become a known parser.
type Parser interface {
	// Parse parses source r into a Codelab for the specified environment env.
	Parse(r io.Reader, opts Options) (*types.Codelab, error)

	// ParseFragment is similar to Parse except it doesn't parse codelab metadata.
	ParseFragment(r io.Reader, opts Options) ([]nodes.Node, error)
}

// Container for parsing options.
type Options struct {
	PassMetadata map[string]bool
}

func NewOptions() *Options {
	return &Options{
		PassMetadata: map[string]bool{},
	}
}

var parsers = map[string]Parser{}

// Register registers a new parser f under specified name.
// It panics if another parser is already registered under the same name.
func Register(name string, p Parser) {
	if _, exists := parsers[name]; exists {
		panic(fmt.Sprintf("parser %q already registered", name))
	}
	parsers[name] = p
}

// Parsers returns a slice of all registered parser names.
func Parsers() []string {
	p := make([]string, 0, len(parsers))
	for k := range parsers {
		p = append(p, k)
	}
	return p
}

// Parse parses source r into a Codelab using a parser registered with
// the specified name.
func Parse(name string, r io.Reader, opts Options) (*types.Codelab, error) {
	p, ok := parsers[name]
	if !ok {
		return nil, fmt.Errorf("no parser named %q", name)
	}
	c, err := p.Parse(r, opts)
	if err != nil {
		return nil, err
	}
	c.URL = c.ID
	return c, err
}

// ParseFragment parses a codelab fragment provided in r, using a parser
// registered with the specified name.
func ParseFragment(name string, r io.Reader, opts Options) ([]nodes.Node, error) {
	p, ok := parsers[name]
	if !ok {
		return nil, fmt.Errorf("no parser named %q", name)
	}
	return p.ParseFragment(r, opts)
}
