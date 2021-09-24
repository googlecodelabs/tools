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
// TODO: test survey
// TODO: test header

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
