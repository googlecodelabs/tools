// Copyright 2018 Google Inc. All Rights Reserved.
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
