package htmltests

import (
	"go/build"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer/html"
	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
	"github.com/googlecodelabs/tools/third_party"
)

func TestRenderLinkTemplate(t *testing.T) {
	linkFileRelDir := "src/github.com/googlecodelabs/tools/claat/proto-renderer/html/templates-tests/testdata/InlineContent/google_weather.txt"
	linkFileAbsDir := filepath.Join(build.Default.GOPATH, linkFileRelDir)
	weatherLinkBytes, err := ioutil.ReadFile(linkFileAbsDir)
	if err != nil {
		t.Errorf("Reading %#v outputted %#v", linkFileAbsDir, err)
		continue
	}
	weatherLinkOutput := string(weatherLinkBytes[:])

	linkProto := testingutils.NewLink(
		"https://www.google.com/search?q=weather+in+nyc",
		testingutils.NewStylizedTextPlain("hey google,"),
		testingutils.NewStylizedTextStrong(" how's the"),
		testingutils.NewStylizedTextEmphasized(" weather in "),
		testingutils.NewStylizedTextStrongAndEmphasized("NYC today?"),
	)

	canonicalTests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: &tutorial.Link{},
			Out:     "",
			Ok:      false,
		},
		{
			InProto: testingutils.NewLink("only://link.does.not/work?#ok"),
			Out:     "",
			Ok:      false,
		},
		{
			InProto: linkProto,
			Out:     weatherLinkOutput,
			Ok:      true,
		},
	}
	testingutils.CanonicalRenderTestBatch(html.Render, canonicalTests, t)
}
