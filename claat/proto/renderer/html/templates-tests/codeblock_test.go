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

func TestRenderCodeBlockTemplate(t *testing.T) {
	tests := []*testingutils.CanonicalFileRenderingBatch{
		{
			InProto: &tutorial.CodeBlock{},
			OutPath: "CodeBlock/dummy.html",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewCodeBlockPlain(""),
			OutPath: "CodeBlock/dummy.html",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewCodeBlockHighlighted(""),
			OutPath: "CodeBlock/dummy_highlighted.html",
			Ok:      true,
		},
	}
	testingutils.TestCanonicalFileRenderBatch("html", html.Render, tests, t)
}
