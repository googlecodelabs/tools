package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto/renderer/html"
	"github.com/googlecodelabs/tools/claat/proto/renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderListTemplateFailures(t *testing.T) {
	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: &tutorial.List{},
			Out:     "",
			Ok:      false,
		},
	}
	testingutils.TestCanonicalRendererBatch(html.Render, tests, t)
}
