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
			InProto: &tutorial.InlineContent{},
			Out:     "",
			Ok:      false,
		},
		{
			InProto: testingUtils.NewInlineContentTextPlain(""),
			Out:     "",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewInlineContentTextStrong("strong"),
			Out:     "<strong>strong</strong>",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewInlineContentTextEmphazied("emphasized"),
			Out:     "<em>emphasized</em>",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewInlineContentTextStrongAndEmphazied("strong & emphasized"),
			Out:     "<strong><em>strong &amp; emphasized</em></strong>",
			Ok:      true,
		},
		{
			InProto: testingUtils.NewInlineContentCode(`~!@#$%^&*()_+-=[]{}\|'";:/?.><,`),
			Out:     `<code>~!@#$%^&amp;*()_+-=[]{}\\|&#39;&#34;;:/?.&gt;&lt;,</code>`,
			Ok:      true,
		},
	}
	testingUtils.CanonicalRenderingTestBatch(html.Render, tests, t)
}

func TestRenderInlineContentTemplateIdentiy(t *testing.T) {
	tests := []*testingUtils.RendererTestingIdendityBatch{
		{
			InProto:  testingUtils.NewInlineContentTextPlain(`<script>alert("you've been hacked!");</script>!`),
			OutProto: testingUtils.NewStylizedTextPlain(`<script>alert("you've been hacked!");</script>!`),
			Out:      `&lt;script&gt;alert(&#34;you&#39;ve been hacked!&#34;);&lt;/script&gt;!`,
			Ok:       true,
		},
		{
			InProto:  testingUtils.NewInlineContentCode(`<script>alert("you've been hacked!");</script>!`),
			OutProto: testingUtils.NewInlineCode(`<script>alert("you've been hacked!");</script>!`),
			Out:      `<code>&lt;script&gt;alert(&#34;you&#39;ve been hacked!&#34;);&lt;/script&gt;!</code>`,
			Ok:       true,
		},
	}
	testingUtils.RenderingTestIdendityBatch(html.Render, tests, t)
}
