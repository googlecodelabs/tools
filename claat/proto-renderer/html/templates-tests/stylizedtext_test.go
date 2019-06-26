package tests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer/html"
	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderStylizedTextTemplateEscaping(t *testing.T) {
	tests := []*testingUtils.RendererTestingBatch{
		{
			InProto: testingUtils.NewStylizedTextPlain(`<script>alert("you've been hacked!");</script>!`),
			Out:     "&lt;script&gt;alert(&#34;you&#39;ve been hacked!&#34;);&lt;/script&gt;!",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewStylizedTextPlain("D@ ?òü ǝ$çâpæ? ^>^ '>__<' {&]"),
			Out:     "D@ ?òü ǝ$çâpæ? ^&gt;^ &#39;&gt;__&lt;&#39; {&amp;]",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewStylizedTextPlain("<h3>**__Markdown not ![esca](ped)__**</h3>"),
			Out:     "&lt;h3&gt;**__Markdown not ![esca](ped)__**&lt;/h3&gt;",
			Ok:      true,
		},
	}
	testingUtils.CanonicalRenderingTestBatch(html.Render, tests, t)
}

func TestRenderStylizedTextTemplate(t *testing.T) {
	tests := []*testingUtils.RendererTestingBatch{
		{
			InProto: &tutorial.StylizedText{},
			Out:     "",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewStylizedTextPlain(""),
			Out:     "",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewStylizedTextPlain("hello!"),
			Out:     "hello!",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewStylizedTextStrong("hello!"),
			Out:     "<strong>hello!</strong>",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewStylizedTextEmphasized("hello!"),
			Out:     "<em>hello!</em>",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewStylizedTextStrongAndEmphasized("hello!"),
			Out:     "<strong><em>hello!</em></strong>",
			Ok:      true,
		},
	}
	testingUtils.CanonicalRenderingTestBatch(html.Render, tests, t)
}
