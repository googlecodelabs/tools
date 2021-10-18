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

// Package types provide types for format-independent codelab data model.
package types

import (
	"time"

	"github.com/googlecodelabs/tools/claat/nodes"
)

// Meta contains a single codelab metadata.
type Meta struct {
	ID         string            `json:"id"`                 // ID is also part of codelab URL
	Duration   int               `json:"duration"`           // Codelab duration in minutes
	Title      string            `json:"title"`              // Codelab title
	Authors    string            `json:"authors,omitempty"`  // Arbitrary authorship text
	Summary    string            `json:"summary"`            // Short summary
	Source     string            `json:"source"`             // Codelab source doc
	Theme      string            `json:"theme"`              // Usually first item of Categories
	Status     *LegacyStatus     `json:"status"`             // Draft, Published, Hidden, etc.
	Categories []string          `json:"category"`           // Categories from the meta table
	Tags       []string          `json:"tags"`               // All environments supported by the codelab
	Feedback   string            `json:"feedback,omitempty"` // Issues and bugs are sent here
	GA         string            `json:"ga,omitempty"`       // Codelab-specific GA tracking ID
	Extra      map[string]string `json:"extra,omitempty"`    // Extra metadata specified in pass_metadata

	URL string `json:"url"` // Legacy ID; TODO: remove
}

// Codelab is a top-level structure containing metadata and codelab steps.
type Codelab struct {
	Meta
	Steps []*Step
}

func NewCodelab() *Codelab {
	clab := &Codelab{}
	clab.Extra = map[string]string{}

	return clab
}

// NewStep creates a new codelab step, adding it to c.Steps slice.
func (c *Codelab) NewStep(title string) *Step {
	s := &Step{Title: title, Content: nodes.NewListNode()}
	c.Steps = append(c.Steps, s)
	return s
}

// Step is a single codelab step, containing metadata and actual content.
type Step struct {
	Title    string          // Step title
	Tags     []string        // Step environments
	Duration time.Duration   // Duration
	Content  *nodes.ListNode // Root node of the step nodes tree
}
