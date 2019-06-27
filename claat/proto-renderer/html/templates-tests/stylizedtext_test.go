package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer/html"
	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderStylizedTextTemplateEscaping(t *testing.T) {
	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: testingutils.NewStylizedTextPlain(`<script>alert("you've been hacked!");</script>!`),
			Out:     "&lt;script&gt;alert(&#34;you&#39;ve been hacked!&#34;);&lt;/script&gt;!",
			Ok:      true,
		},
		{
			InProto: testingutils.NewStylizedTextPlain("D@ ?òü ǝ$çâpæ? ^>^ '>__<' {&]"),
			Out:     "D@ ?òü ǝ$çâpæ? ^&gt;^ &#39;&gt;__&lt;&#39; {&amp;]",
			Ok:      true,
		},
		{
			InProto: testingutils.NewStylizedTextPlain("<h3>**__Markdown not ![esca](ped)__**</h3>"),
			Out:     "&lt;h3&gt;**__Markdown not ![esca](ped)__**&lt;/h3&gt;",
			Ok:      true,
		},
	}
	testingutils.CanonicalRenderTestBatch(html.Render, tests, t)
}

func TestRenderStylizedTextTemplate(t *testing.T) {
	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: &tutorial.StylizedText{},
			Out:     "",
			Ok:      true,
		},
		{
			InProto: testingutils.NewStylizedTextPlain(""),
			Out:     "",
			Ok:      true,
		},
		{
			InProto: testingutils.NewStylizedTextPlain("hello!"),
			Out:     "hello!",
			Ok:      true,
		},
		{
			InProto: testingutils.NewStylizedTextStrong("hello!"),
			Out:     "<strong>hello!</strong>",
			Ok:      true,
		},
		{
			InProto: testingutils.NewStylizedTextEmphasized("hello!"),
			Out:     "<em>hello!</em>",
			Ok:      true,
		},
		{
			InProto: testingutils.NewStylizedTextStrongAndEmphasized("hello!"),
			Out:     "<strong><em>hello!</em></strong>",
			Ok:      true,
		},
	}
	testingutils.CanonicalRenderTestBatch(html.Render, tests, t)
}
