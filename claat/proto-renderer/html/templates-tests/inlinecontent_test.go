package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer/html"
	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderInlineContentTemplate(t *testing.T) {
	tests := []*testingUtils.RendererTestingBatch{
		{
			&tutorial.InlineContent{},
			"",
			false,
		},
		{
			testingUtils.NewInlineContentPlain(""),
			"",
			true,
		},
		{
			testingUtils.NewInlineContentStrong("strong"),
			"<strong>strong</strong>",
			true,
		},
		{
			testingUtils.NewInlineContentEmphazied("emphasized"),
			"<em>emphasized</em>",
			true,
		},
		{
			testingUtils.NewInlineContentStrongAndEmphazied("strong & emphasized"),
			"<strong><em>strong &amp; emphasized</em></strong>",
			true,
		},
		{
			testingUtils.NewInlineContentCode(`~!@#$%^&*()_+-=[]{}\|'";:/?.><,`),
			`<code>~!@#$%^&amp;*()_+-=[]{}\\|&#39;&#34;;:/?.&gt;&lt;,</code>`,
			true,
		},
	}
	testingUtils.CanonicalRenderingTestBatch(html.Render, tests, t)
}

func TestRenderInlineContentTemplateIdentiy(t *testing.T) {
	tests := []*testingUtils.RendererTestingIdendityBatch{
		{
			testingUtils.NewInlineContentPlain(`<script>alert("you've been hacked!");</script>!`),
			testingUtils.NewStylizedTextPlain(`<script>alert("you've been hacked!");</script>!`),
			`&lt;script&gt;alert(&#34;you&#39;ve been hacked!&#34;);&lt;/script&gt;!`,
			true,
		},
		{
			testingUtils.NewInlineContentCode(`<script>alert("you've been hacked!");</script>!`),
			testingUtils.NewInlineCode(`<script>alert("you've been hacked!");</script>!`),
			`<code>&lt;script&gt;alert(&#34;you&#39;ve been hacked!&#34;);&lt;/script&gt;!</code>`,
			true,
		},
	}
	testingUtils.RenderingTestIdendityBatch(html.Render, tests, t)
}
