package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer/html"
	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderTestBlockTemplate(t *testing.T) {
	tests := []*testingUtils.RendererTestingBatch{
		{
			InProto: &tutorial.TextBlock{},
			Out:     "",
			Ok:      false,
		},
		{
			InProto: testingUtils.NewTextBlock(),
			Out:     "",
			Ok:      false,
		},
		{
			InProto: testingUtils.NewTextBlock(
				testingUtils.NewInlineContentTextPlain(`hello, `),
				testingUtils.NewInlineContentTextStrong(`world!`),
				testingUtils.NewInlineContentTextEmphazied(` goodbye`),
				testingUtils.NewInlineContentTextPlain(` `),
				testingUtils.NewInlineContentTextStrongAndEmphazied(`cruel `),
				testingUtils.NewInlineContentCode(`world!`),
			),
			Out: "hello, <strong>world!</strong><em> goodbye</em> <strong><em>cruel </em></strong><code>world!</code>",
			Ok:  true,
		},
	}
	testingUtils.CanonicalRenderingTestBatch(html.Render, tests, t)
}
