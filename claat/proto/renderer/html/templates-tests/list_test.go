package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto/renderer/html"
	"github.com/googlecodelabs/tools/claat/proto/renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderListTemplateLinkFromFile(t *testing.T) {
	tests := []*testingutils.CanonicalFileRenderingBatch{
		{
			InProto: &tutorial.List{},
			OutFile: "List/...html",
			Ok:      false,
		},
	}
	testingutils.TestCanonicalFileRenderBatch(html.Render, tests, t)
}

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
