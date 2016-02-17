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

package md

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/googlecodelabs/tools/claat/types"

	"golang.org/x/net/html"
)

func buildParserWithStep(markup string) *parserState {
	ps := parserState{
		tzr: html.NewTokenizer(strings.NewReader(markup)),
		c:   &types.Codelab{},
	}
	ps.currentStep = ps.c.NewStep("title")
	return &ps
}

func TestHandleCodelabTitle(t *testing.T) {
	// Set up.
	title := "Egret"
	ps := buildParserWithStep(fmt.Sprintf("<h1>%s</h1>", title))
	ps.advance()
	handleCodelabTitle(ps)

	if ps.c.Title != title {
		t.Errorf("[%q] got %v, want %v", title, ps.c.Title, title)
	}
}

func TestProcessDuration(t *testing.T) {
	tests := []struct {
		in  string
		out time.Duration
		ok  bool
	}{
		// Test an easy case.
		{"1:00", time.Hour, true},
		// Test an easy case with period delimiter.
		{"1.00", time.Hour, true},
		// Test a weird number case.
		{"13:37", 817 * time.Minute, true},
		// Test that minutes can be longer than an hour.
		{"0:90", 90 * time.Minute, true},
		// Test zero.
		{"0", 0, true},
		// Test more than two digits in the hours place.
		{"100:00", 100 * time.Hour, true},
		// Test that hours must be present.
		{":25", 0, false},
		// Test that two minutes digits are required.
		{"2:2", 0, false},
		// Test that two minutes digits are required.
		{"4:444", 0, false},
		// Test that zero is the only simple scalar allowed.
		{"654321", 0, false},
		// Test an empty string.
		{"", 0, false},
		// Test complete nonsense.
		{"This isn't even remotely close to resembling a correct duration string.", 0, false},
	}

	for i, tc := range tests {
		out, err := processDuration(tc.in)
		switch {
		case !tc.ok && err == nil:
			t.Errorf("%d: processDuration(%q) = %v; want error", i, tc.in, out)
		case tc.ok && err != nil:
			t.Errorf("%d: processDuration(%q) = %v; want %v", i, tc.in, err, tc.out)
		case out != tc.out:
			t.Errorf("%d: processDuration(%q) = %v; want %v", i, tc.in, out, tc.out)
		}
	}

}

func TestHandleDurationHint(t *testing.T) {
	tests := []struct {
		in  string
		out time.Duration
	}{
		// Test various forms of the "Duration" stub.
		{"<p>Duration: 2:00</p>", 2 * time.Hour},
		{"<p>Duration 2:00</p>", 2 * time.Hour},
		{"<p>duration: 2:00</p>", 2 * time.Hour},
		{"<p>duration 2:00</p>", 2 * time.Hour},
		// Test the zero case.
		{"<p>Duration: 0</p>", 0},
		// Test false positive cases.
		{"<p>Redwood is my favorite kind of tree.</p>", 0},
		{"<p><h6>Subsubsubsubsub Header</h6></p>", 0},
	}

	for i, tc := range tests {
		ps := buildParserWithStep(tc.in)
		ps.advance()
		handleDurationHint(ps)
		if ps.currentStep.Duration != tc.out {
			t.Errorf("%d: [%q] got %v, want %v", i, tc.in, ps.currentStep.Duration, tc.out)
		}
	}
}

func TestComputeTotalDuration(t *testing.T) {
	tests := []struct {
		in  []time.Duration
		out int
	}{
		{[]time.Duration{45 * time.Minute, 90 * time.Minute, 15 * time.Minute}, 150},
		{[]time.Duration{0, 0, 0, 0}, 0},
		{nil, 0},
	}

	for i, tc := range tests {
		ps := parserState{
			c: &types.Codelab{},
		}
		for _, v := range tc.in {
			ps.currentStep = ps.c.NewStep("")
			ps.currentStep.Duration = v
		}
		computeTotalDuration(ps.c)
		if tc.out != ps.c.Duration {
			t.Errorf("%d: [%q] got %v, want %v", i, tc.in, ps.c.Duration, tc.out)
		}
	}
}

func TestStandardSplit(t *testing.T) {
	tests := []struct {
		in  string
		out []string
	}{
		{"qwe,rty,ui,op", []string{"qwe", "rty", "ui", "op"}},
		{"AsD,fGH,jkL;", []string{"asd", "fgh", "jkl;"}},
		{"zxc, vb,\tnm", []string{"zxc", "vb", "nm"}},
		{"QwE,   rt, YUIOP", []string{"qwe", "rt", "yuiop"}},
		{"asdfghjkl;", []string{"asdfghjkl;"}},
	}
	for i, tc := range tests {
		out := standardSplit(tc.in)
		if len(out) != len(tc.out) {
			t.Errorf("%d: standardSplit(%v) got %v, want %v", i, tc.in, out, tc.out)
		}
		for k, v := range out {
			if v != tc.out[k] {
				t.Errorf("%d: standardSplit(%v) got %v, want %v", i, tc.in, out, tc.out)
			}
		}
	}
}

func TestAddMetadataToCodelab(t *testing.T) {
	tempStatus := types.LegacyStatus([]string{"draft"})
	tests := []struct {
		in  map[string]string
		out types.Codelab
	}{
		{
			map[string]string{
				"summary":           "abcdefghij",
				"id":                "zyxwvut",
				"categories":        "not, really",
				"environments":      "web, kiosk",
				"status":            "draft",
				"feedback link":     "https://www.google.com",
				"analytics account": "12345",
			},
			types.Codelab{
				Meta: types.Meta{
					Summary:    "abcdefghij",
					ID:         "zyxwvut",
					Categories: []string{"not", "really"},
					Tags:       []string{"web", "kiosk"},
					Status:     &tempStatus,
					Feedback:   "https://www.google.com",
					GA:         "12345",
				},
			},
		},
	}
	for i, tc := range tests {
		out := types.Codelab{}
		addMetadataToCodelab(tc.in, &out)
		if !reflect.DeepEqual(tc.out, out) {
			t.Errorf("%d: [%q] got %v, want %v", i, tc.in, out, tc.out)
		}
	}
}

func TestNewBreaklessTextNode(t *testing.T) {
	tests := []struct {
		in  string
		out *types.TextNode
	}{
		{"one\ntwo\nthree", types.NewTextNode("one two three")},
		{"four fivesix", types.NewTextNode("four fivesix")},
	}
	for i, tc := range tests {
		out := newBreaklessTextNode(tc.in)
		if !reflect.DeepEqual(out, tc.out) {
			t.Errorf("%d: newBreaklessTextNode(%v) got %v, want %v", i, tc.in, out, tc.out)
		}
	}
}
