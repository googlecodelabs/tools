package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto/constructor"
	"github.com/googlecodelabs/tools/claat/proto/renderer/html"
	"github.com/googlecodelabs/tools/claat/proto/renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderInlineContentTemplateLinkFromFile(t *testing.T) {
	linkProto := protoconstructors.NewLink(
		"https://www.google.com/search?q=weather+in+nyc",
		protoconstructors.NewStylizedTextPlain("hey google,"),
		protoconstructors.NewStylizedTextStrong(" how's the"),
		protoconstructors.NewStylizedTextEmphasized(" weather in "),
		protoconstructors.NewStylizedTextStrongAndEmphasized("NYC today?"),
	)

	tests := []*testingutils.CanonicalFileRenderingBatch{
		{
			InProto: linkProto,
			OutPath: "Link/google_weather.html",
			Ok:      true,
		},
	}
	testingutils.TestCanonicalFileRenderBatch("html", html.Render, tests, t)
}

func TestRenderInlineContentTemplateIdentity(t *testing.T) {
	tests := []*testingutils.RendererIdendityBatch{
		{
			InProto:  protoconstructors.NewInlineContentTextPlain(`<script>alert("you've been hacked!");</script>!`),
			OutProto: protoconstructors.NewStylizedTextPlain(`<script>alert("you've been hacked!");</script>!`),
			Out:      `&lt;script&gt;alert(&#34;you&#39;ve been hacked!&#34;);&lt;/script&gt;!`,
			Ok:       true,
		},
		{
			InProto:  protoconstructors.NewInlineContentCode(`<script>alert("you've been hacked!");</script>!`),
			OutProto: protoconstructors.NewInlineCode(`<script>alert("you've been hacked!");</script>!`),
			Out:      `<code>&lt;script&gt;alert(&#34;you&#39;ve been hacked!&#34;);&lt;/script&gt;!</code>`,
			Ok:       true,
		},
	}
	testingutils.RenderingIdendityTestBatch(html.Render, tests, t)
}

func TestRenderInlineContentTemplate(t *testing.T) {
	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: &tutorial.InlineContent{},
			Out:     "",
			Ok:      false,
		},
		{
			InProto: protoconstructors.NewInlineContentTextPlain(""),
			Out:     "",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewInlineContentTextStrong("strong"),
			Out:     "<strong>strong</strong>",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewInlineContentTextEmphasized("emphasized"),
			Out:     "<em>emphasized</em>",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewInlineContentTextStrongAndEmphasized("strong & emphasized"),
			Out:     "<strong><em>strong &amp; emphasized</em></strong>",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewInlineContentCode(`~!@#$%^&*()_+-=[]{}\|'";:/?.><,`),
			Out:     `<code>~!@#$%^&amp;*()_+-=[]{}\\|&#39;&#34;;:/?.&gt;&lt;,</code>`,
			Ok:      true,
		},
	}
	testingutils.TestCanonicalRendererBatch(html.Render, tests, t)
}
