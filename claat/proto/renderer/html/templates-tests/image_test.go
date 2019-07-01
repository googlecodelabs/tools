package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto/renderer/html"
	"github.com/googlecodelabs/tools/claat/proto/renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderImageTemplate(t *testing.T) {
	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: &tutorial.Image{},
			Out:     `<img src="" alt="" style="height: 0px; width: 0px">`,
			Ok:      true,
		},
	}
	testingutils.TestCanonicalRendererBatch(html.Render, tests, t)
}
