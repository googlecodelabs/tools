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
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/googlecodelabs/tools/claat/nodes"
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

func parseFragment(markup string) ([]nodes.Node, error) {
	r := strings.NewReader(markup)
	p := &Parser{}

	opts := *parser.NewOptions()
	return p.ParseFragment(r, opts)
}

func stringify(nodesToStringify []nodes.Node, level string) string {
	var content []string
	for _, node := range nodesToStringify {
		base := fmt.Sprintf("%+v", node)
		if node.Type() == nodes.NodeItemsList {
			children := []nodes.Node{}
			for _, list := range node.(*nodes.ItemsListNode).Items {
				children = append(children, list)
			}

			base += "\n" + level + " Child Nodes: vvvv \n" + stringify(children, level+">") + "\n" + level + " Child Nodes: ^^^^"
		}

		if node.Type() == nodes.NodeList {
			base += "\n" + level + " Child Nodes: vvvv \n" + stringify(node.(*nodes.ListNode).Nodes, level+">") + "\n" + level + " Child Nodes: ^^^^"
		}

		content = append(content, base)
	}

	return strings.Join(content, "\n")
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
analytics_account: 12345
feedback_link: https://www.google.com

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
			"extra_field_two": "bbbbb",
		},
	}

	content := `---
id: zyxwvut
authors: john smith
summary: abcdefghij
categories: not, really
environments: kiosk, web
analytics_account: 12345
feedback_link: https://www.google.com
extra_field_one: aaaaa
extra_field_two: bbbbb

---
`
	content += ("# " + title)

	opts := *parser.NewOptions()
	opts.PassMetadata = map[string]bool{
		"extra_field_two": true,
	}

	c := mustParseCodelab(content, opts)
	if !reflect.DeepEqual(c.Meta, wantMeta) {
		t.Errorf("\ngot:\n%+v\nwant:\n%+v", c.Meta, wantMeta)
	}
}

func TestParseFragment(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error

		skipReason string
	}{
		{
			name:  "Empty String",
			input: "",
		},
		{
			name:  "Kinda broken markdown",
			input: `**this is kinda broken*`,
		},
		{
			name: "Valid Markdown",
			input: `
### This is cool
![Image Title](image path)
#### Level more
<video id="dQw4w9WgXcQ"></video>
<!--won't show-->
`,
		},
		{
			name: "Forbidden Nested Imports",
			input: `
I want nested imports
<</path/to/something else.md>>`,
			wantErr: ErrForbiddenFragmentImports,
		},
		{
			name: "If something looks like metadata, it is treated as text content",
			input: `
---
id: zyxwvut
authors: john smith
summary: abcdefghij
categories: not, really
environments: kiosk, web
analytics_account: 12345
feedback_link: https://www.google.com
extrafieldone: aaaaa
extrafieldtwo: bbbbb

---
We don't parse the above!`,
		},
		{
			name:    "Forbidden Steps",
			input:   `## This is not allowed`,
			wantErr: ErrForbiddenFragmentSteps,
		},
		{
			name:    "Forbidden Top level",
			input:   `# This is not allowed`,
			wantErr: ErrForbiddenFragmentSteps,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.skipReason != "" {
				t.Skip(test.skipReason)
			}

			got, err := parseFragment(test.input)
			t.Logf("parseFragment(\n%q\n) =\n%s\n, %+v", test.input, stringify(got, " >"), err)
			if err != test.wantErr {
				t.Errorf("want error = %+v", test.wantErr)
			}
		})
	}
}

func TestParseWithImport(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name: "valid imports",
			input: stdHeader + `
## Step 1
before

<<import1.md>>

after
<<example/import2.md>>
with space
<<space/is allowed.md>>
## Step 2
<<import_another_file.md>>`,
			want: []string{"import1.md", "example/import2.md", "space/is allowed.md", "import_another_file.md"},
		},
		{
			name: "import not in steps",
			input: stdHeader + `
<<should_not_work.md>>

## Step 1
<<allowed.md>>
		`,
			want: []string{"allowed.md"},
		},
		{
			name: "import not on its own line",
			input: stdHeader + `
## Step 1 <<This is not allowed>>
<<not like this>><<not like this>>
<<this is ok.md>>
<<but not this>>this line

<<strange case is here and should not be allowed>>## Step 2
<<you cannot do this ## Step 3>> Otherwise it's really broken.
		`,
			want: []string{"this is ok.md"},
		},
		{
			name: "import inside code block should not be considered",
			input: stdHeader + `
## Step 1
		` + "```" + `
<<I guess we should consider it here.md>>
		` + "```" + `
		`,
			want: []string{"I guess we should consider it here.md"},
		},
		{
			name: "HTML injection is not allowed",
			input: stdHeader + `
## Step 1
I'm going to inject some HTML
<<-->alert("yup")>>
<</>})<script>alert("gotcha")</script>>>>
<<"});alert("how aobut this?">>>
<script>
<<--document.write("random stuff")>>
</script>
`,
		},
		{
			name: "nonmarkdown file is currently not supported",
			input: stdHeader + `
## Step 1
<<nonmd file.gdoc>>
<<somemd.md>>`,
			want: []string{"somemd.md"},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			lab := mustParseCodelab(test.input, *parser.NewOptions())
			var got []string
			for _, s := range lab.Steps {
				for _, n := range nodes.ImportNodes(s.Content.Nodes) {
					got = append(got, n.URL)
				}
			}

			// make consistent ordering
			sort.StringSlice(got).Sort()
			sort.StringSlice(test.want).Sort()
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("Parsing\n%s\nGot Imports:\n%s\nWant Imports:\n%s\n", test.input, got, test.want)
			}
		})
	}
}
