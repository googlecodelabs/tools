package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto/constructors"
	"github.com/googlecodelabs/tools/claat/proto/renderer/html"
	"github.com/googlecodelabs/tools/claat/proto/renderer/testing-utils"
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
			InProto: protoconstructors.NewParagraph(),
			Out:     "",
			Ok:      false,
		},
		{
			InProto: protoconstructors.NewParagraph(
				protoconstructors.NewInlineContentTextPlain(`hello, `),
				protoconstructors.NewInlineContentTextStrong(`world!`),
				protoconstructors.NewInlineContentTextEmphasized(` goodbye`),
				protoconstructors.NewInlineContentTextPlain(` `),
				protoconstructors.NewInlineContentTextStrongAndEmphasized(`cruel `),
				protoconstructors.NewInlineContentCode(`world!`),
			),
			Out: "<p>hello, <strong>world!</strong><em> goodbye</em> <strong><em>cruel </em></strong><code>world!</code></p>",
			Ok:  true,
		},
	}
	testingutils.TestCanonicalRendererBatch(html.Render, tests, t)
}
