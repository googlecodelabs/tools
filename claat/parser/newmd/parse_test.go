// Copyright 2018-2019 Google LLC. All Rights Reserved.
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

package md

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/googlecodelabs/tools/claat/parser"
	"github.com/googlecodelabs/tools/claat/render"
	"github.com/googlecodelabs/tools/claat/types"
)

const stdMeta = `---
id: codelab
summary: summary

---`

const stdHeader = stdMeta + `
# Codelab Title
`

func mustParseCodelab(markup string, opts parser.Options) *types.Codelab {
	c, err := parseCodelab(markup, opts)
	if err != nil {
		log.Fatalf("Error parsing markup %v: %v", markup, err)
	}

	return c
}

func parseCodelab(markup string, opts parser.Options) (*types.Codelab, error) {
	r := strings.NewReader(markup)
	p := &Parser{}

	return p.Parse(r, opts)
}

func TestHandleCodelabTitle(t *testing.T) {
	// Set up.
	title := "Egret"
	c := mustParseCodelab(fmt.Sprintf("# %s", title), *parser.NewOptions())

	if c.Title != title {
		t.Errorf("[%q] got %v, want %v", title, c.Title, title)
	}
}

func TestParseCodelab(t *testing.T) {
	tests := []struct {
		in  string
		out *types.Codelab
		ok  bool
	}{
		{
			in: `#titlegoeshere

id: idvalue
tags: foo, bar, baz

non-metadata content before the first step

## Title Of The First Step
Duration: 3:00:00

first step text content

here is some more content

## Step Two
Duration: 2:00

whatever

## Step 3

Duration: 05:00

more content
`,
			out: &types.Codelab{
				Meta: types.Meta{
					Title:    "titlegoeshere",
					Tags:     []string{"bar", "baz", "foo"},
					ID:       "idvalue",
					Extra:    map[string]string{},
					Duration: 187, // Minutes. This is an int, not a duration.
				},
				Steps: []*types.Step{
					&types.Step{
						Duration: 3 * time.Hour,
						Tags:     []string{},
						Title:    "Title Of The First Step",
					},
					&types.Step{
						Duration: 2 * time.Minute,
						Tags:     []string{},
						Title:    "Step Two",
					},
					&types.Step{
						Duration: 5 * time.Minute,
						Tags:     []string{},
						Title:    "Step 3",
					},
				},
			},
			ok: true,
		},
	}
	for _, tc := range tests {
		out, err := parseCodelab(tc.in, *parser.NewOptions())

		if !tc.ok && err == nil {
			t.Errorf("parseCodelab(string, options) = %+v, want err", out)
			continue
		}
		if tc.ok && err != nil {
			t.Errorf("parseCodelab(string, options) = %+v, want %+v", err, tc.out)
			continue
		}
		// Compare steps separately from meta, as DeepEqual can't handle the fact that they're pointers.
		if tc.ok && !reflect.DeepEqual(out.Meta, tc.out.Meta) {
			t.Errorf("parseCodelab(string, options).Meta = %#v, want %#v", out.Meta, tc.out.Meta)
		}
		if tc.ok && len(out.Steps) != len(tc.out.Steps) {
			t.Errorf("parseCodelab(string, options).Steps has length %d, want %d", len(out.Steps), len(tc.out.Steps))
		}
		for i, v := range out.Steps {
			// Comparing content fields is prohibitively difficult due to the recursion and cross references.
			v.Content = nil
			if !reflect.DeepEqual(v, tc.out.Steps[i]) {
				t.Errorf("parseCodelab(string, options).Steps[%d] = %#v, want %#v", i, v, tc.out.Steps[i])
			}
		}
	}
}

func TestParseCodelabContent(t *testing.T) {
	const input = `# this is the title

id: idvalue
tags: foo, bar, baz

non-metadata content before the first step

## Title Of The First Step
Duration: 3:00:00

first step text content

here is some more content

## Step Two
Duration: 2:00

whatever

## Step 3

Duration: 05:00

more content
`
	expected := []string{
		strings.TrimSpace(`<p>first step text content</p>
<p>here is some more content</p>`),
		strings.TrimSpace(`<p>whatever</p>`),
		strings.TrimSpace(`<p>more content</p>`),
	}
	p := &Parser{}
	cl, err := p.Parse(strings.NewReader(input), *parser.NewOptions())
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(cl.Steps) != len(expected) {
		t.Errorf("got %d steps, want %d", len(cl.Steps), len(expected))
	}

	for i, v := range cl.Steps {
		markup, err := render.HTML("", v.Content)
		if err != nil {
			t.Errorf(err.Error())
		}
		trimmedMarkup := strings.TrimSpace(string(markup))
		if trimmedMarkup != expected[i] {
			t.Errorf("For step at index %d, got %s, want %s", i, trimmedMarkup, expected[i])
		}
	}
}

// The parser assumes
//   * Any single number is minutes
//   * xx:yy is mm:ss
//   * Hours only appear if you have three parts
func TestProcessDuration(t *testing.T) {
	tests := []struct {
		in  string
		out time.Duration
	}{
		// Test an easy case.
		{"1:00", 1 * time.Minute},
		// Test a weird number case.
		{"13:37", 14 * time.Minute},
		// Test that seconds can be longer than a minute.
		{"00:90", 2 * time.Minute},
		// Test that minutes can be longer than an hour.
		{"00:90:00", 90 * time.Minute},
		// Test zero.
		{"0", 0},
		// Test more than two digits in the hours place.
		{"100:00:00", 100 * time.Hour},
		// Test an empty string.
		{"", 0},
		// Test complete nonsense.
		{"Complete nonsense.", 0},
	}

	for i, tc := range tests {
		content := fmt.Sprintf(stdHeader+"\n## Step Title\nDuration: %v\n", tc.in)
		c := mustParseCodelab(content, *parser.NewOptions())
		out := time.Duration(c.Duration) * time.Minute

		if out != tc.out {
			t.Errorf("%d: got duration %v from %q, wanted %v", i, out, tc.in, tc.out)
		}
	}
}

func TestComputeTotalDuration(t *testing.T) {
	tmp := `
## Step Title
Duration: %v
`

	tests :=
		[]struct {
			in  []string
			out int
		}{
			{[]string{"45:00", "90:00", "15:00"}, 150},
			{[]string{"0", "00", "00:00", "00:00:00"}, 0},
		}

	for i, tc := range tests {
		content := stdHeader
		for _, dur := range tc.in {
			content += fmt.Sprintf(tmp, dur)
		}

		c := mustParseCodelab(content, *parser.NewOptions())
		if c.Duration != tc.out {
			t.Errorf("%d: wanted duration %d but got %d", i, c.Duration, tc.out)
		}
	}
}

func TestParseMetadata(t *testing.T) {
	title := "Codelab Title"
	wantMeta := types.Meta{
		Title:      title,
		ID:         "zyxwvut",
		Authors:    "john smith",
		Summary:    "abcdefghij",
		Categories: []string{"not", "really"},
		Tags:       []string{"kiosk", "web"},
		Feedback:   "https://www.google.com",
		GA:         "12345",
		Extra:      map[string]string{},
	}

	content := `---
id: zyxwvut
authors: john smith
summary: abcdefghij
categories: not, really
environments: kiosk, web
analytics account: 12345
feedback link: https://www.google.com

---
`
	content += ("# " + title)

	c := mustParseCodelab(content, *parser.NewOptions())
	if !reflect.DeepEqual(c.Meta, wantMeta) {
		t.Errorf("\ngot:\n%+v\nwant:\n%+v", c.Meta, wantMeta)
	}
}

func TestParseMetadataPassMetadata(t *testing.T) {
	title := "Codelab Title"
	wantMeta := types.Meta{
		Title:      title,
		ID:         "zyxwvut",
		Authors:    "john smith",
		Summary:    "abcdefghij",
		Categories: []string{"not", "really"},
		Tags:       []string{"kiosk", "web"},
		Feedback:   "https://www.google.com",
		GA:         "12345",
		Extra: map[string]string{
			"extrafieldtwo": "bbbbb",
		},
	}

	content := `---
id: zyxwvut
authors: john smith
summary: abcdefghij
categories: not, really
environments: kiosk, web
analytics account: 12345
feedback link: https://www.google.com
extrafieldone: aaaaa
extrafieldtwo: bbbbb

---
`
	content += ("# " + title)

	opts := *parser.NewOptions()
	opts.PassMetadata = map[string]bool{
		"extrafieldtwo": true,
	}

	c := mustParseCodelab(content, opts)
	if !reflect.DeepEqual(c.Meta, wantMeta) {
		t.Errorf("\ngot:\n%+v\nwant:\n%+v", c.Meta, wantMeta)
	}
}
