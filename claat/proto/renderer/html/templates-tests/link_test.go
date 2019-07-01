package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto/constructor"
	"github.com/googlecodelabs/tools/claat/proto/renderer/html"
	"github.com/googlecodelabs/tools/claat/proto/renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderLinkTemplate(t *testing.T) {
	tests := []*testingutils.CanonicalFileRenderingBatch{
		{
			InProto: protoconstructors.NewLink(
				"https://www.google.com/search?q=weather+in+nyc",
				protoconstructors.NewStylizedTextPlain("hey google,"),
				protoconstructors.NewStylizedTextStrong(" how's the"),
				protoconstructors.NewStylizedTextEmphasized(" weather in "),
				protoconstructors.NewStylizedTextStrongAndEmphasized("NYC today?"),
			),
			OutPath: "Link/google_weather.html",
			Ok:      true,
		},
	}
	testingutils.TestCanonicalFileRenderBatch("html", html.Render, tests, t)
}

func TestRenderLinkTemplateFailures(t *testing.T) {
	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: &tutorial.Link{},
			Out:     "",
			Ok:      false,
		},
		{
			InProto: protoconstructors.NewLink("only://link.does.not/work?#ok=false"),
			Out:     "",
			Ok:      false,
		},
	}
	testingutils.TestCanonicalRendererBatch(html.Render, tests, t)
}
