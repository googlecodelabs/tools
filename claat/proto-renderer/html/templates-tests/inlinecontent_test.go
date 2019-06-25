package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer/html"
	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
)

func TestRenderStylizedTextTemplateEscaping(t *testing.T) {
	tests := []*testingUtils.RendererTestingBatch{
		{
			{},
			"",
			true,
		},
	}
	testingUtils.CanonicalRenderingTestBatch(html.Render, tests, t)
}
