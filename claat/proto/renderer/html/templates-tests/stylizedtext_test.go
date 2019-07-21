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

func TestRenderStylizedTextTemplateEscaping(t *testing.T) {
	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: protoconstructors.NewStylizedTextPlain(`<script>alert("you've been hacked!");</script>!`),
			Out:     "&lt;script&gt;alert(&#34;you&#39;ve been hacked!&#34;);&lt;/script&gt;!",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewStylizedTextPlain("D@ ?òü ǝ$çâpæ? ^>^ '>__<' {&]"),
			Out:     "D@ ?òü ǝ$çâpæ? ^&gt;^ &#39;&gt;__&lt;&#39; {&amp;]",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewStylizedTextPlain("<h3>**__Markdown not ![esca](ped)__**</h3>"),
			Out:     "&lt;h3&gt;**__Markdown not ![esca](ped)__**&lt;/h3&gt;",
			Ok:      true,
		},
	}
	testingutils.TestCanonicalRendererBatch(html.Render, tests, t)
}

func TestRenderStylizedTextTemplate(t *testing.T) {
	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: &tutorial.StylizedText{},
			Out:     "",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewStylizedTextPlain(""),
			Out:     "",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewStylizedTextPlain("hello!"),
			Out:     "hello!",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewStylizedTextStrong("hello!"),
			Out:     "<strong>hello!</strong>",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewStylizedTextEmphasized("hello!"),
			Out:     "<em>hello!</em>",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewStylizedTextStrongAndEmphasized("hello!"),
			Out:     "<strong><em>hello!</em></strong>",
			Ok:      true,
		},
	}
	testingutils.TestCanonicalRendererBatch(html.Render, tests, t)
}
