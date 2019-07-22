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

func TestRenderInlineCodeTemplateEscaping(t *testing.T) {
	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: tutorial.InlineCode{},
			Out:     "<code></code>",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewInlineCode("< less-than"),
			Out:     "<code>&lt; less-than</code>",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewInlineCode("> greater-than"),
			Out:     "<code>&gt; greater-than</code>",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewInlineCode("/ backslash"),
			Out:     "<code>/ backslash</code>",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewInlineCode(`\ forwardslash`),
			Out:     `<code>\\ forwardslash</code>`,
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewInlineCode("& ampersand"),
			Out:     "<code>&amp; ampersand</code>",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewInlineCode(`" quotation`),
			Out:     "<code>&#34; quotation</code>",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewInlineCode("' apostrophe"),
			Out:     "<code>&#39; apostrophe</code>",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewInlineCode("TODO"),
			Out:     "<code>TODO</code>",
			Ok:      true,
		},
	}
	testingutils.TestCanonicalRendererBatch(html.Render, tests, t)
}
