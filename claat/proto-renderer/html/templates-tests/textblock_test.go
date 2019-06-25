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
			testingUtils.NewStylizedTextPlain(`<script>alert("you've been hacked!");</script>!`),
			"&lt;script&gt;alert(&#34;you&#39;ve been hacked!&#34;);&lt;/script&gt;!",
			true,
		},
	}
	testingUtils.CanonicalRenderingTestBatch(html.Render, tests, t)
}
