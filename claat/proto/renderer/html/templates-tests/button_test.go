package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto/constructors"
	"github.com/googlecodelabs/tools/claat/proto/renderer/html"
	"github.com/googlecodelabs/tools/claat/proto/renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderButtonTemplate(t *testing.T) {
	canonicalTests := []*testingutils.CanonicalFileRenderingBatch{
		{
			InProto: protoconstructors.NewButtonPlain(
				protoconstructors.NewLink(
					"www.cloud.io",
					protoconstructors.NewStylizedTextPlain("hosting"),
				),
			),
			OutPath: "Button/plain.txt",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewButtonDownload(
				protoconstructors.NewLink(
					"www.random.org",
					protoconstructors.NewStylizedTextPlain("FizzBuzz"),
				),
			),
			OutPath: "Button/download.txt",
			Ok:      true,
		},
	}
	testingutils.TestCanonicalFileRenderBatch("html", html.Render, canonicalTests, t)
}

func TestRenderButtonTemplateFailures(t *testing.T) {
	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: &tutorial.Button{},
			Out:     "",
			Ok:      false,
		},
		{
			InProto: protoconstructors.NewButtonPlain(nil),
			Out:     "",
			Ok:      false,
		},
		{
			InProto: protoconstructors.NewButtonDownload(nil),
			Out:     "",
			Ok:      false,
		},
	}
	testingutils.TestCanonicalRendererBatch(html.Render, tests, t)
}
