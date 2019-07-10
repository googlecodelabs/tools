package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto/constructor"
	"github.com/googlecodelabs/tools/claat/proto/renderer/html"
	"github.com/googlecodelabs/tools/claat/proto/renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderHeadingFailuresTemplate(t *testing.T) {
	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: &tutorial.Heading{},
			Out:     "",
			Ok:      false,
		},
		{
			InProto: protoconstructors.NewHeading(0),
			Out:     "",
			Ok:      false,
		},
	}
	testingutils.TestCanonicalRendererBatch(html.Render, tests, t)
}
func TestRenderHeadingTemplateFromFile(t *testing.T) {
	tests := []*testingutils.CanonicalFileRenderingBatch{
		{
			InProto: protoconstructors.NewHeading(
				0,
				protoconstructors.NewInlineContentTextPlain(""),
			),
			OutPath: "Heading/dummy.html",
			Ok:      true,
		},
	}
	testingutils.TestCanonicalFileRenderBatch("html", html.Render, tests, t)
}
