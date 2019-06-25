package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer/html"
	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderStylizedTextTemplate(t *testing.T) {
	tests := []*testingUtils.RendererTestingBatch{
		{
			{
				testingUtils.NewStylizedTextPlain(`<script>alert("you've been hacked!");</script>!`),
			},

			// testingUtils.NewStylizedTextStrong
			// testingUtils.NewStylizedTextEmphazied
			// testingUtils.NewStylizedTextStrongAndEmphazied
			// testingUtils.NewInlineCode
			"",
			true,
		},
	}
	testingUtils.CanonicalRenderingTestBatch(html.Render, tests, t)
}
