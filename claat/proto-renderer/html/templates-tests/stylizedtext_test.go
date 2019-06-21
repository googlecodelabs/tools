package tests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer"
	"github.com/googlecodelabs/tools/claat/proto-renderer/html"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderStylizedTextTemplate(t *testing.T) {
	// h3 := &devrel_tutorial.Heading{
	// 	Level: 3,
	// 	Text:  "<script>some _very_ bad code!</script>",
	// 	Text:  "D@ ?òü ǝ$çâpæ urlquery? '>__<' {&]",
	// 	Text:  "**__Markdown not ![esca](ped)__**",
	// }

	tests := []*genrenderer.TestingBatch{
		{
			genrenderer.AssertRenderedTemplate(
				html.Render(&devrel_tutorial.StylizedText{
					Text: "hello!",
				}),
			),
			"hello!",
			true,
		},
		{
			genrenderer.AssertRenderedTemplate(
				html.Render(&devrel_tutorial.StylizedText{
					Text:     "hello!",
					IsStrong: true,
				}),
			),
			"<strong>hello!</strong>",
			true,
		},
		{
			genrenderer.AssertRenderedTemplate(
				html.Render(&devrel_tutorial.StylizedText{
					Text:         "hello!",
					IsEmphasized: true,
				}),
			),
			"<em>hello!</em>",
			true,
		},
		{
			genrenderer.AssertRenderedTemplate(
				html.Render(&devrel_tutorial.StylizedText{
					Text:         "hello!",
					IsStrong:     true,
					IsEmphasized: true,
				}),
			),
			"<strong><em>hello!</em></strong>",
			true,
		},
	}
	genrenderer.CanonicalTemplateTestBatch(tests, t)
}
