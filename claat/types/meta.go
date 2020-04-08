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
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Meta contains a single codelab metadata.
type Meta struct {
	ID         string            `json:"id"`                   // ID is also part of codelab URL
	Duration   int               `json:"duration"`             // Codelab duration in minutes
	Title      string            `json:"title"`                // Codelab title
	Authors    string            `json:"authors,omitempty"`    // Arbitrary authorship text
	BadgePath  string            `json:"badge_path,omitempty"` // Path of the Badge to grant on codelab completion on devsite
	Summary    string            `json:"summary"`              // Short summary
	Source     string            `json:"source"`               // Codelab source doc
	Theme      string            `json:"theme"`                // Usually first item of Categories
	Status     *LegacyStatus     `json:"status"`               // Draft, Published, Hidden, etc.
	Categories []string          `json:"category"`             // Categories from the meta table
	Tags       []string          `json:"tags"`                 // All environments supported by the codelab
	Feedback   string            `json:"feedback,omitempty"`   // Issues and bugs are sent here
	GA         string            `json:"ga,omitempty"`         // Codelab-specific GA tracking ID
	Extra      map[string]string `json:"extra,omitempty"`      // Extra metadata specified in pass_metadata

	URL string `json:"url"` // Legacy ID; TODO: remove
}

// Context is an export context.
// It is defined in this package so that it can be used by both cli and a server.
type Context struct {
	Env     string       `json:"environment"`       // Current export environment
	Format  string       `json:"format"`            // Output format, e.g. "html"
	Prefix  string       `json:"prefix,omitempty"`  // Assets URL prefix for HTML-based formats
	MainGA  string       `json:"mainga,omitempty"`  // Global Google Analytics ID
	Updated *ContextTime `json:"updated,omitempty"` // Last update timestamp
}

// ContextMeta is a composition of export context and meta data.
type ContextMeta struct {
	Context
	Meta
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
	s := &Step{Title: title, Content: NewListNode()}
	c.Steps = append(c.Steps, s)
	return s
}

// Step is a single codelab step, containing metadata and actual content.
type Step struct {
	Title    string        // Step title
	Tags     []string      // Step environments
	Duration time.Duration // Duration
	Content  *ListNode     // Root node of the step nodes tree
}

// ContextTime is codelab metadata timestamp.
// It can be of "YYYY-MM-DD" or RFC3339 formats but marshaling
// always uses RFC3339 format.
type ContextTime time.Time

// MarshalJSON implements Marshaler interface.
func (ct ContextTime) MarshalJSON() ([]byte, error) {
	v := time.Time(ct).Format(time.RFC3339)
	b := make([]byte, len(v)+2)
	b[0] = '"'
	b[len(b)-1] = '"'
	copy(b[1:], v)
	return b, nil
}

// UnmarshalJSON implements Unmarshaler interface.
// Accepted format is "YYYY-MM-DD" or RFC3339.
func (ct *ContextTime) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	t, err := time.Parse(time.RFC3339, string(b))
	if err != nil {
		t, err = time.Parse("2006-01-02", string(b))
	}
	if err != nil {
		return err
	}
	*ct = ContextTime(t)
	return nil
}

// LegacyStatus supports legacy status values which are strings
// as opposed to an array, e.g. "['one', u'two', ...]".
type LegacyStatus []string

// MarshalJSON implements Marshaler interface.
func (s LegacyStatus) MarshalJSON() ([]byte, error) {
	if len(s) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal([]string(s))
}

// UnmarshalJSON implements Unmarshaler interface.
func (s *LegacyStatus) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	if b[0] == '"' {
		// legacy status: "['s1', u's2', ...]"
		// assume no status value contains single quotes
		b = bytes.Trim(b, `"`)
		b = bytes.Replace(b, []byte("u'"), []byte(`"`), -1)
		b = bytes.Replace(b, []byte("'"), []byte(`"`), -1)
	}
	var v []string
	if err := json.Unmarshal(b, &v); err != nil {
		return fmt.Errorf("%v: %s", err, b)
	}
	*s = LegacyStatus(v)
	return nil
}

// String turns a status into a string
func (s LegacyStatus) String() string {
	ss := []string(s)
	if len(ss) == 0 {
		return ""
	}

	return "[" + strings.Join(ss, ",") + "]"
}
