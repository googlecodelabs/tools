package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer/html"
	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderInlineContentTemplate(t *testing.T) {
	tests := []*testingutils.RendererTestingBatch{
		{
			InProto: &tutorial.InlineContent{},
			Out:     "",
			Ok:      false,
		},
		{
			InProto: testingutils.NewInlineContentTextPlain(""),
			Out:     "",
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineContentTextStrong("strong"),
			Out:     "<strong>strong</strong>",
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineContentTextEmphasized("emphasized"),
			Out:     "<em>emphasized</em>",
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineContentTextStrongAndEmphasized("strong & emphasized"),
			Out:     "<strong><em>strong &amp; emphasized</em></strong>",
			Ok:      true,
		},
		{
			InProto: testingutils.NewInlineContentCode(`~!@#$%^&*()_+-=[]{}\|'";:/?.><,`),
			Out:     `<code>~!@#$%^&amp;*()_+-=[]{}\\|&#39;&#34;;:/?.&gt;&lt;,</code>`,
			Ok:      true,
		},
	}
	testingutils.CanonicalRenderingTestBatch(html.Render, tests, t)
}

func TestRenderInlineContentTemplateIdentiy(t *testing.T) {
	tests := []*testingutils.RendererTestingIdendityBatch{
		{
			InProto:  testingutils.NewInlineContentTextPlain(`<script>alert("you've been hacked!");</script>!`),
			OutProto: testingutils.NewStylizedTextPlain(`<script>alert("you've been hacked!");</script>!`),
			Out:      `&lt;script&gt;alert(&#34;you&#39;ve been hacked!&#34;);&lt;/script&gt;!`,
			Ok:       true,
		},
		{
			InProto:  testingutils.NewInlineContentCode(`<script>alert("you've been hacked!");</script>!`),
			OutProto: testingutils.NewInlineCode(`<script>alert("you've been hacked!");</script>!`),
			Out:      `<code>&lt;script&gt;alert(&#34;you&#39;ve been hacked!&#34;);&lt;/script&gt;!</code>`,
			Ok:       true,
		},
	}
	testingutils.RenderingTestIdendityBatch(html.Render, tests, t)
}
