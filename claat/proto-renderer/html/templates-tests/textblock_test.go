package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer/html"
	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderTestBlockTemplate(t *testing.T) {
	tests := []*testingutils.RendererTestingBatch{
		{
			InProto: &tutorial.TextBlock{},
			Out:     "",
			Ok:      false,
		},
		{
			InProto: testingutils.NewTextBlock(),
			Out:     "",
			Ok:      false,
		},
		{
			InProto: testingutils.NewTextBlock(
				testingutils.NewInlineContentTextPlain(`hello, `),
				testingutils.NewInlineContentTextStrong(`world!`),
				testingutils.NewInlineContentTextEmphasized(` goodbye`),
				testingutils.NewInlineContentTextPlain(` `),
				testingutils.NewInlineContentTextStrongAndEmphasized(`cruel `),
				testingutils.NewInlineContentCode(`world!`),
			),
			Out: "hello, <strong>world!</strong><em> goodbye</em> <strong><em>cruel </em></strong><code>world!</code>",
			Ok:  true,
		},
	}
	testingutils.CanonicalRenderingTestBatch(html.Render, tests, t)
}
