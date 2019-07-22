// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package testingutils

import (
	"go/build"
	"io/ioutil"
	"path/filepath"
	"testing"
)

const (
	filesTestDataDir = "src/github.com/googlecodelabs/tools/claat/proto/renderer/"
	outPkgDataDir    = "templates-tests/testdata"
)

// CanonicalFileRenderingBatch type for canonical i != o and !ok rendering tests
// of file-based outputs
type CanonicalFileRenderingBatch struct {
	InProto interface{}
	OutPath string
	Ok      bool
}

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
