package htmltests

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto/constructor"
	"github.com/googlecodelabs/tools/claat/proto/renderer/html"
	"github.com/googlecodelabs/tools/claat/proto/renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderYoutubeVideoTemplateFromFile(t *testing.T) {
	tests := []*testingutils.CanonicalFileRenderingBatch{
		{
			InProto: &tutorial.YoutubeVideo{},
			OutPath: "YoutubeVideo/no-id.html",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewYTVideo(""),
			OutPath: "YoutubeVideo/no-id.html",
			Ok:      true,
		},
		{
			InProto: protoconstructors.NewYTVideo("PWL7uUWh-xcw"),
			OutPath: "YoutubeVideo/valid.html",
			Ok:      true,
		},
	}
	testingutils.TestCanonicalFileRenderBatch("html", html.Render, tests, t)
}
