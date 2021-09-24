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

package render

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googlecodelabs/tools/claat/nodes"
)

func TestHTMLEnv(t *testing.T) {
	one := nodes.NewTextNode("one ")
	one.MutateEnv([]string{"one"})
	two := nodes.NewTextNode("two ")
	two.MutateEnv([]string{"two"})
	three := nodes.NewTextNode("three ")
	three.MutateEnv([]string{"one", "three"})

	tests := []struct {
		env    string
		output string
	}{
		{"", "one two three "},
		{"one", "one three "},
		{"two", "two "},
		{"three", "three "},
		{"four", ""},
	}
	for i, test := range tests {
		var ctx Context
		ctx.Env = test.env
		h, err := HTML(ctx, one, two, three)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}
		if v := string(h); v != test.output {
			t.Errorf("%d: v = %q; want %q", i, v, test.output)
		}
	}
}

// TODO: test HTML
// TODO: test writeHTML
// TODO: test ReplaceDoubleCurlyBracketsWithEntity
// TODO: test matchEnv
// TODO: test write
// TODO: test writeString
// TODO: test writeFmt
// TODO: test escape
// TODO: test writeEscape
// TODO: test text
// TODO: test image
// TODO: test url
// TODO: test button
// TODO: test code
// TODO: test list
// TODO: test onlyImages
// TODO: test itemsList
// TODO: test grid
// TODO: test infobox

func TestSurvey(t *testing.T) {
	tests := []struct {
		name   string
		inNode *nodes.SurveyNode
		out    string
	}{
		{
			name: "NoID",
			inNode: nodes.NewSurveyNode("",
				&nodes.SurveyGroup{
					Name:    "pick a number",
					Options: []string{"1", "2", "3"},
				},
				&nodes.SurveyGroup{
					Name:    "choose an answer",
					Options: []string{"yes", "no", "probably"},
				},
			),
			out: `<google-codelab-survey survey-id="">
<h4>pick a number</h4>
<paper-radio-group>
<paper-radio-button>1</paper-radio-button>
<paper-radio-button>2</paper-radio-button>
<paper-radio-button>3</paper-radio-button>
</paper-radio-group>
<h4>choose an answer</h4>
<paper-radio-group>
<paper-radio-button>yes</paper-radio-button>
<paper-radio-button>no</paper-radio-button>
<paper-radio-button>probably</paper-radio-button>
</paper-radio-group>
</google-codelab-survey>`,
		},
		{
			name:   "NoGroups",
			inNode: nodes.NewSurveyNode("foobar"),
			out: `<google-codelab-survey survey-id="foobar">
</google-codelab-survey>`,
		},
		{
			name: "Simple",
			inNode: nodes.NewSurveyNode("a simple example",
				&nodes.SurveyGroup{
					Name:    "pick a color",
					Options: []string{"red", "blue", "yellow"},
				}),
			out: `<google-codelab-survey survey-id="a simple example">
<h4>pick a color</h4>
<paper-radio-group>
<paper-radio-button>red</paper-radio-button>
<paper-radio-button>blue</paper-radio-button>
<paper-radio-button>yellow</paper-radio-button>
</paper-radio-group>
</google-codelab-survey>`,
		},
		{
			name: "MultipleGroups",
			inNode: nodes.NewSurveyNode("an example with multiple survey groups",
				&nodes.SurveyGroup{
					Name:    "a",
					Options: []string{"a", "aa", "aaa"},
				},
				&nodes.SurveyGroup{
					Name:    "b",
					Options: []string{"b", "bb", "bbb"},
				},
				&nodes.SurveyGroup{
					Name:    "c",
					Options: []string{"c", "cc", "ccc"},
				}),
			out: `<google-codelab-survey survey-id="an example with multiple survey groups">
<h4>a</h4>
<paper-radio-group>
<paper-radio-button>a</paper-radio-button>
<paper-radio-button>aa</paper-radio-button>
<paper-radio-button>aaa</paper-radio-button>
</paper-radio-group>
<h4>b</h4>
<paper-radio-group>
<paper-radio-button>b</paper-radio-button>
<paper-radio-button>bb</paper-radio-button>
<paper-radio-button>bbb</paper-radio-button>
</paper-radio-group>
<h4>c</h4>
<paper-radio-group>
<paper-radio-button>c</paper-radio-button>
<paper-radio-button>cc</paper-radio-button>
<paper-radio-button>ccc</paper-radio-button>
</paper-radio-group>
</google-codelab-survey>`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.survey(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.survey(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}

func TestHeader(t *testing.T) {
	a1 := nodes.NewTextNode("foo")
	a1.Italic = true
	a2 := nodes.NewTextNode("bar")
	a3 := nodes.NewTextNode("baz")
	a3.Code = true

	tests := []struct {
		name   string
		inNode *nodes.HeaderNode
		out    string
	}{
		{
			name:   "SimpleH1",
			inNode: nodes.NewHeaderNode(1, nodes.NewTextNode("foobar")),
			out:    "<h1 is-upgraded>foobar</h1>",
		},
		{
			name:   "LevelOutOfRange",
			inNode: nodes.NewHeaderNode(100, nodes.NewTextNode("foobar")),
			out:    "<h100 is-upgraded>foobar</h100>",
		},
		{
			name:   "EmptyContent",
			inNode: nodes.NewHeaderNode(2),
			out:    "<h2 is-upgraded></h2>",
		},
		{
			name:   "StyledText",
			inNode: nodes.NewHeaderNode(3, a1, a2, a3),
			out:    "<h3 is-upgraded><em>foo</em>bar<code>baz</code></h3>",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.header(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.header(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}

func TestYouTube(t *testing.T) {
	tests := []struct {
		name   string
		inNode *nodes.YouTubeNode
		out    string
	}{
		{
			name:   "NonEmpty",
			inNode: nodes.NewYouTubeNode("Mlk888FiI8A"),
			out:    `<iframe class="youtube-video" src="https://www.youtube.com/embed/Mlk888FiI8A?rel=0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`,
		},
		{
			name:   "Empty",
			inNode: nodes.NewYouTubeNode(""),
			out:    `<iframe class="youtube-video" src="https://www.youtube.com/embed/?rel=0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.youtube(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.youtube(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}

func TestIframe(t *testing.T) {
	tests := []struct {
		name   string
		inNode *nodes.IframeNode
		out    string
	}{
		{
			name:   "SomeText",
			inNode: nodes.NewIframeNode("maps.google.com"),
			out:    `<iframe class="embedded-iframe" src="maps.google.com"></iframe>`,
		},
		{
			name:   "Escape",
			inNode: nodes.NewIframeNode("ma ps.google.com"),
			out:    `<iframe class="embedded-iframe" src="ma ps.google.com"></iframe>`,
		},
		{
			name:   "Empty",
			inNode: nodes.NewIframeNode(""),
			out:    `<iframe class="embedded-iframe" src=""></iframe>`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			outBuffer := &bytes.Buffer{}
			hw := &htmlWriter{w: outBuffer}
			hw.iframe(tc.inNode)
			out := outBuffer.String()
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("hw.iframe(%+v) got diff (-want +got):\n%s", tc.inNode, diff)
			}
		})
	}
}
