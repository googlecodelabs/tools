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
