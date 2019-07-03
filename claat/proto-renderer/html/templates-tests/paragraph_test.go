package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer/html"
	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderTestBlockTemplate(t *testing.T) {
	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: &tutorial.Paragraph{},
			Out:     "",
			Ok:      false,
		},
		{
			InProto: testingutils.NewParagraph(),
			Out:     "",
			Ok:      false,
		},
		{
			InProto: testingutils.NewParagraph(
				testingutils.NewInlineContentTextPlain(`hello, `),
				testingutils.NewInlineContentTextStrong(`world!`),
				testingutils.NewInlineContentTextEmphasized(` goodbye`),
				testingutils.NewInlineContentTextPlain(` `),
				testingutils.NewInlineContentTextStrongAndEmphasized(`cruel `),
				testingutils.NewInlineContentCode(`world!`),
			),
			Out: "<p>hello, <strong>world!</strong><em> goodbye</em> <strong><em>cruel </em></strong><code>world!</code></p>",
			Ok:  true,
		},
	}
	testingutils.CanonicalRenderTestBatch(html.Render, tests, t)
}
