// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto/constructor"
	"github.com/googlecodelabs/tools/claat/proto/renderer/html"
	"github.com/googlecodelabs/tools/claat/proto/renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

// TestRenderInlineContentTemplateIdentity proves tutorial.InlineContent
// constructors work properly
func TestRenderInlineContentTemplateIdentity(t *testing.T) {
	tests := []*testingutils.RendererIdendityBatch{
		{
			InFunc:  protoconstructors.NewInlineContentTextPlain,
			InProto: protoconstructors.NewStylizedTextPlain(""),
		},
		{
			InFunc:  protoconstructors.NewInlineContentTextPlain,
			InProto: protoconstructors.NewStylizedTextPlain(`<script>alert("you've been hacked!");</script>!`),
		},
		{
			InFunc:  protoconstructors.NewInlineContentTextStrong,
			InProto: protoconstructors.NewStylizedTextStrong("strong"),
		},
		{
			InFunc:  protoconstructors.NewInlineContentTextEmphasized,
			InProto: protoconstructors.NewStylizedTextEmphasized("emphasized"),
		},
		{
			InFunc:  protoconstructors.NewInlineContentTextStrongAndEmphasized,
			InProto: protoconstructors.NewStylizedTextStrongAndEmphasized("strong & emphasized"),
		},
		{
			InFunc:  protoconstructors.NewInlineContentCode,
			InProto: protoconstructors.NewInlineCode(`<script>alert("you've been hacked!");</script>!`),
		},
		{
			InFunc: protoconstructors.NewInlineContentLink,
			InProto: protoconstructors.NewLink(
				"https://www.google.com/search?q=weather+in+nyc",
				protoconstructors.NewStylizedTextPlain("hey google,"),
				protoconstructors.NewStylizedTextStrong(" how's the"),
				protoconstructors.NewStylizedTextEmphasized(" weather in "),
				protoconstructors.NewStylizedTextStrongAndEmphasized("NYC today?"),
			),
		},
		{
			InFunc: protoconstructors.NewInlineContentButton,
			InProto: protoconstructors.NewButtonPlain(
				protoconstructors.NewLink(
					"http://github.com/favicon.ico",
					protoconstructors.NewStylizedTextStrongAndEmphasized("Github's favicon"),
				)),
		},
		{
			InFunc: protoconstructors.NewInlineContentButton,
			InProto: protoconstructors.NewButtonDownload(
				protoconstructors.NewLink(
					"https://firebase.google.com/favicon.ico",
					protoconstructors.NewStylizedTextStrongAndEmphasized("Firebase's favicon"),
				)),
		},
		{
			InFunc:  protoconstructors.NewInlineContentInlineImage,
			InProto: protoconstructors.NewInlineImage("https://cloud.google.com/favicon.ico", "hello GCloud", 40, 40),
		},
	}
	testingutils.RenderingIdendityTestBatch(html.Render, tests, t)
}

func TestRenderInlineContentStylizedTextTemplate(t *testing.T) {
	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: &tutorial.InlineContent{},
			Out:     "",
			Ok:      false,
		},
		{
			InProto: tutorial.InlineContent{},
			Out:     "",
			Ok:      false,
		},
	}
	testingutils.TestCanonicalRendererBatch(html.Render, tests, t)
}
