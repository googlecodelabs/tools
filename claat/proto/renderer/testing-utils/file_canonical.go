package testingutils

import (
	"go/build"
	"io/ioutil"
	"path/filepath"
	"testing"
)

type CanonicalFileRenderingBatch struct {
	InProto interface{}
	OutPath string
	Ok      bool
}

const (
	filesTestDataDir = "src/github.com/googlecodelabs/tools/claat/proto/renderer/"
	outPkgDataDir    = "templates-tests/testdata"
)

// TestCanonicalFileRenderBatch is the helper for canonical i != o and !ok rendering tests
// of file-based outputs
func TestCanonicalFileRenderBatch(
	outpkg string, renderer renderingFunc,
	tests []*CanonicalFileRenderingBatch, t *testing.T) {
	for _, tc := range tests {

		OutPathAbsDir := filepath.Join(
			build.Default.GOPATH, filesTestDataDir, outpkg, outPkgDataDir, tc.OutPath)
		fileBytes, err := ioutil.ReadFile(OutPathAbsDir)
		if err != nil {
			t.Errorf("Reading %#v outputted %#v", OutPathAbsDir, err)
			continue
		}

		// Create cannonical test from the output from the underlying type
		newTc := []*CanonicalRenderingBatch{
			{
				InProto: tc.InProto,
				Out:     string(fileBytes[:]),
				Ok:      tc.Ok,
			},
		}
		TestCanonicalRendererBatch(renderer, newTc, t)
	}
}
