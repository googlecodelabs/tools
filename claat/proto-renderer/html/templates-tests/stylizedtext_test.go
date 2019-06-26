package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer/html"
	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderStylizedTextTemplateEscaping(t *testing.T) {
	tests := []*testingUtils.RendererTestingBatch{
		{
			testingUtils.NewStylizedTextPlain(`<script>alert("you've been hacked!");</script>!`),
			"&lt;script&gt;alert(&#34;you&#39;ve been hacked!&#34;);&lt;/script&gt;!",
			true,
		},
		{
			testingUtils.NewStylizedTextPlain("D@ ?òü ǝ$çâpæ? ^>^ '>__<' {&]"),
			"D@ ?òü ǝ$çâpæ? ^&gt;^ &#39;&gt;__&lt;&#39; {&amp;]",
			true,
		},
		{
			testingUtils.NewStylizedTextPlain("<h3>**__Markdown not ![esca](ped)__**</h3>"),
			"&lt;h3&gt;**__Markdown not ![esca](ped)__**&lt;/h3&gt;",
			true,
		},
	}
	testingUtils.CanonicalRenderingTestBatch(html.Render, tests, t)
}

func TestRenderStylizedTextTemplate(t *testing.T) {
	tests := []*testingUtils.RendererTestingBatch{
		{
			&tutorial.StylizedText{},
			"",
			true,
		},
		{
			testingUtils.NewStylizedTextPlain(""),
			"",
			true,
		},
		{
			testingUtils.NewStylizedTextPlain("hello!"),
			"hello!",
			true,
		},
		{
			testingUtils.NewStylizedTextStrong("hello!"),
			"<strong>hello!</strong>",
			true,
		},
		{
			testingUtils.NewStylizedTextEmphazied("hello!"),
			"<em>hello!</em>",
			true,
		},
		{
			testingUtils.NewStylizedTextStrongAndEmphazied("hello!"),
			"<strong><em>hello!</em></strong>",
			true,
		},
	}
	testingUtils.CanonicalRenderingTestBatch(html.Render, tests, t)
}
