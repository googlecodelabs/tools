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
	"sync"

	"github.com/googlecodelabs/tools/claat/types"
)

// ParseFunc parses source r into a Codelab for the specified environment env.
type ParseFunc func(r io.Reader) (*types.Codelab, error)

var (
	parsersMu sync.Mutex // guards parsers
	parsers   = make(map[string]ParseFunc)
)

// Register registers a new parser f under specified name.
// It panics if another parser is already registered under the same name.
func Register(name string, f ParseFunc) {
	parsersMu.Lock()
	defer parsersMu.Unlock()
	if _, exists := parsers[name]; exists {
		panic(fmt.Sprintf("parser %q already registered", name))
	}
	parsers[name] = f
}

// Parsers returns a slice of all registered parser names.
func Parsers() []string {
	parsersMu.Lock()
	defer parsersMu.Unlock()
	p := make([]string, 0, len(parsers))
	for k := range parsers {
		p = append(p, k)
	}
	return p
}

// Parse parses source r into a Codelab using a parser registered with
// the specified name.
func Parse(name string, r io.Reader) (*types.Codelab, error) {
	parsersMu.Lock()
	p, ok := parsers[name]
	parsersMu.Unlock()
	if !ok {
		return nil, fmt.Errorf("no parser named %q", name)
	}
	c, err := p(r)
	if err != nil {
		return nil, err
	}
	c.URL = c.ID
	return c, err
}
