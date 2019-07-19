package htmltests

import (
	"go/build"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/googlecodelabs/tools/claat/proto/constructor"
	"github.com/googlecodelabs/tools/claat/proto/renderer/html"
	"github.com/googlecodelabs/tools/claat/proto/renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

const linkFileRelDir = "src/github.com/googlecodelabs/tools/claat/proto/renderer/html/templates-tests/testdata/InlineContent/google_weather.txt"

func TestRenderLinkTemplate(t *testing.T) {
	linkFileAbsDir := filepath.Join(build.Default.GOPATH, linkFileRelDir)
	weatherLinkBytes, err := ioutil.ReadFile(linkFileAbsDir)
	if err != nil {
		t.Errorf("Reading %#v outputted %#v", linkFileAbsDir, err)
		return
	}
	weatherLinkOutput := string(weatherLinkBytes[:])

	linkProto := protoconstructors.NewLink(
		"https://www.google.com/search?q=weather+in+nyc",
		protoconstructors.NewStylizedTextPlain("hey google,"),
		protoconstructors.NewStylizedTextStrong(" how's the"),
		protoconstructors.NewStylizedTextEmphasized(" weather in "),
		protoconstructors.NewStylizedTextStrongAndEmphasized("NYC today?"),
	)

	canonicalTests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: &tutorial.Link{},
			Out:     "",
			Ok:      false,
		},
		{
			InProto: protoconstructors.NewLink("only://link.does.not/work?#ok"),
			Out:     "",
			Ok:      false,
		},
		{
			InProto: linkProto,
			Out:     weatherLinkOutput,
			Ok:      true,
		},
	}
	testingutils.TestCanonicalRendererBatch(html.Render, canonicalTests, t)
}
