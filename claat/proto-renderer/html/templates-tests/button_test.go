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

func TestRenderButtonTemplate(t *testing.T) {
	// TODO: Generalize file reading logic
	var fOuts []string
	files := []string{
		"plain.html",
		"download.html",
	}
	for _, f := range files {
		linkFileRelDir := "src/github.com/googlecodelabs/tools/claat/proto-renderer/html/templates-tests/testdata/Button"
		linkFileAbsDir := filepath.Join(build.Default.GOPATH, linkFileRelDir, f)
		fileBytes, err := ioutil.ReadFile(linkFileAbsDir)
		if err != nil {
			t.Errorf("Reading %#v outputted %#v", linkFileAbsDir, err)
			continue
		}
		fOuts = append(fOuts, string(fileBytes[:]))
	}

	tests := []*testingutils.CanonicalRenderingBatch{
		{
			InProto: &tutorial.Button{},
			Out:     "",
			Ok:      false,
		},
		{
			InProto: testingutils.NewButtonPlain(
				testingutils.NewLink(
					"www.cloud.io",
					testingutils.NewStylizedTextPlain("hosting"),
				),
			),
			Out: fOuts[0],
			Ok:  true,
		},
		{
			InProto: testingutils.NewButtonDownload(
				testingutils.NewLink(
					"www.random.org",
					testingutils.NewStylizedTextPlain("FizzBuzz"),
				),
			),
			Out: fOuts[1],
			Ok:  true,
		},
	}
	testingutils.CanonicalRenderTestBatch(html.Render, tests, t)
}
