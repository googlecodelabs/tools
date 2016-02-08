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
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/googlecodelabs/tools/claat/types"

	"golang.org/x/net/html"
)

const empty = ""

var errNonNil = errors.New(empty)

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
		t.Errorf("[%q] expected %v, got %v", title, title, ps.c.Title)
	}
}

func TestProcessDuration(t *testing.T) {
	type TestCase struct {
		in  string
		out time.Duration
		err error
	}
	tests := []TestCase{
		// Test an easy case.
		TestCase{
			in:  "1:00",
			out: time.Duration(1) * time.Hour,
			err: nil,
		},
		// Test an easy case with period delimiter.
		TestCase{
			in:  "1.00",
			out: time.Duration(1) * time.Hour,
			err: nil,
		},
		// Test a weird number case.
		TestCase{
			in:  "13:37",
			out: time.Duration(817) * time.Minute,
			err: nil,
		},
		// Test that minutes can be longer than an hour.
		TestCase{
			in:  "0:90",
			out: time.Duration(90) * time.Minute,
			err: nil,
		},
		// Test zero.
		TestCase{
			in:  "0",
			out: time.Duration(0),
			err: nil,
		},
		// Test more than two digits in the hours place.
		TestCase{
			in:  "100:00",
			out: time.Duration(100) * time.Hour,
			err: nil,
		},
		// Test that hours must be present.
		TestCase{
			in:  ":25",
			out: 0,
			err: errNonNil,
		},
		// Test that two minutes digits are required.
		TestCase{
			in:  "2:2",
			out: 0,
			err: errNonNil,
		},
		// Test that two minutes digits are required.
		TestCase{
			in:  "4:444",
			out: 0,
			err: errNonNil,
		},
		// Test that zero is the only simple scalar allowed.
		TestCase{
			in:  "600613",
			out: 0,
			err: errNonNil,
		},
		// Test an empty string.
		TestCase{
			in:  "",
			out: 0,
			err: errNonNil,
		},
		// Test complete nonsense.
		TestCase{
			in:  "This isn't even remotely close to resembling a correct duration string.",
			out: 0,
			err: errNonNil,
		},
	}

	for _, tc := range tests {
		out, err := processDuration(tc.in)
		if tc.err == errNonNil {
			// Expect an error.
			if err == nil {
				t.Errorf("[%q] expected non-nil error", tc.in)
			}
		} else {
			if err != nil {
				t.Errorf("[%q] expected nil error, got %v", tc.in, err)
			} else if out != tc.out {
				t.Errorf("[%q] expected %v, got %v", tc.in, tc.out, out)
			}
		}
	}

}

func TestHandleDurationHint(t *testing.T) {
	type TestCase struct {
		in  string
		out time.Duration
	}
	tests := []TestCase{
		// Test various forms of the "Duration" stub.
		TestCase{
			in:  "<p>Duration: 2:00</p>",
			out: 2 * time.Hour,
		},
		TestCase{
			in:  "<p>Duration 2:00</p>",
			out: 2 * time.Hour,
		},
		TestCase{
			in:  "<p>duration: 2:00</p>",
			out: 2 * time.Hour,
		},
		TestCase{
			in:  "<p>duration 2:00</p>",
			out: 2 * time.Hour,
		},
		// Test the zero case.
		TestCase{
			in:  "<p>Duration: 0</p>",
			out: 0,
		},
		// Test false positive cases.
		TestCase{
			in:  "<p>Redwood is my favorite kind of tree.</p>",
			out: 0,
		},
		TestCase{
			in:  "<p><h6>Subsubsubsubsub Header</h6></p>",
			out: 0,
		},
	}

	for _, tc := range tests {
		ps := buildParserWithStep(tc.in)
		ps.advance()
		handleDurationHint(ps)
		if ps.currentStep.Duration != tc.out {
			t.Errorf("[%q] expected %v, got %v", tc.in, tc.out, ps.currentStep.Duration)
		}
	}
}

func TestComputeTotalDuration(t *testing.T) {
	type TestCase struct {
		in  []time.Duration
		out int
	}
	tests := []TestCase{
		TestCase{
			in:  []time.Duration{45 * time.Minute, 90 * time.Minute, 15 * time.Minute},
			out: 150,
		},
		TestCase{
			in:  []time.Duration{0, 0, 0, 0},
			out: 0,
		},
		TestCase{
			in:  nil,
			out: 0,
		},
	}

	for _, tc := range tests {
		ps := parserState{
			c: &types.Codelab{},
		}
		for _, v := range tc.in {
			ps.currentStep = ps.c.NewStep("")
			ps.currentStep.Duration = v
		}
		computeTotalDuration(ps.c)
		if tc.out != ps.c.Duration {
			t.Errorf("[%q] expected %v, got %v", tc.in, tc.out, ps.c.Duration)
		}
	}
}

func TestStandardSplit(t *testing.T) {
	type TestCase struct {
		in  string
		out []string
	}
	tests := []TestCase{
		TestCase{
			in:  "qwe,rty,ui,op",
			out: []string{"qwe", "rty", "ui", "op"},
		},
		TestCase{
			in:  "AsD,fGH,jkL;",
			out: []string{"asd", "fgh", "jkl;"},
		},
		TestCase{
			in:  "zxc, vb,\tnm",
			out: []string{"zxc", "vb", "nm"},
		},
		TestCase{
			in:  "QwE,   rt, YUIOP",
			out: []string{"qwe", "rt", "yuiop"},
		},
		TestCase{
			in:  "asdfghjkl;",
			out: []string{"asdfghjkl;"},
		},
	}
	for _, tc := range tests {
		out := standardSplit(tc.in)
		if len(out) != len(tc.out) {
			t.Errorf("[%q] expected %v, got %v", tc.in, tc.out, out)
		}
		for k, v := range out {
			if v != tc.out[k] {
				t.Errorf("[%q] expected %v, got %v", tc.in, tc.out, out)
			}
		}
	}
}

func TestAddMetadataToCodelab(t *testing.T) {
	tempStatus := types.LegacyStatus([]string{"draft"})
	type TestCase struct {
		in  map[string]string
		out types.Codelab
	}
	tests := []TestCase{
		TestCase{
			in: map[string]string{
				"summary":           "abcdefghij",
				"id":                "zyxwvut",
				"categories":        "not, really",
				"environments":      "web, kiosk",
				"status":            "draft",
				"feedback link":     "https://www.google.com",
				"analytics account": "12345",
			},
			out: types.Codelab{
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
	for _, tc := range tests {
		out := types.Codelab{}
		addMetadataToCodelab(tc.in, &out)
		if !reflect.DeepEqual(tc.out, out) {
			t.Errorf("[%q] expected %v, got %v", tc.in, tc.out, out)
		}
	}
}

func TestNewBreaklessTextNode(t *testing.T) {
	type TestCase struct {
		in  string
		out *types.TextNode
	}
	tests := []TestCase{
		TestCase{
			in:  "one\ntwo\nthree",
			out: types.NewTextNode("one two three"),
		},
		TestCase{
			in:  "four fivesix",
			out: types.NewTextNode("four fivesix"),
		},
	}
	for _, tc := range tests {
		out := newBreaklessTextNode(tc.in)
		if !reflect.DeepEqual(out, tc.out) {
			t.Errorf("[%q] expected %v, got %v", tc.in, tc.out, out)
		}
	}
}
