// Copyright 2016-2019 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cmd_test ensures the end to end works as intended.
package cmd_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/googlecodelabs/tools/claat/cmd"
)

func TestExportCodelabMemory(t *testing.T) {
	/*
		Test Plan: Ensure ExportCodelabMemory and ExportCodelab can generate identical
		artifact on valid cases with a few difference e.g. removal of "source"
		metadata field.
	*/
	tests := []struct {
		name     string
		filePath string
		knownBug string
	}{
		{
			name:     "Multiple Steps",
			filePath: "testdata/simple-2-steps.md",
			knownBug: "https://github.com/googlecodelabs/tools/issues/391",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.knownBug != "" {
				t.Skip("Skipping Tests with known issue", test.knownBug)
			}

			tmp, err := ioutil.TempDir(".", "TestExportCodelabMemory-*")
			if err != nil {
				t.Fatal(err)
			}

			defer os.RemoveAll(tmp)

			testFile, err := ioutil.ReadFile(test.filePath)
			if err != nil {
				t.Fatal(err)
			}

			testContent := ioutil.NopCloser(bytes.NewReader(testFile))
			gotBytes := bytes.NewBuffer([]byte{})
			opts := cmd.CmdExportOptions{
				Expenv:   "web",
				Output:   tmp,
				Tmplout:  "devsite",
				GlobalGA: "UA-99999999-99",
			}

			// Given the same markdown input, ExportCodelabMemory should have the same output content as ExportCodelab
			wantMeta, err := cmd.ExportCodelab(test.filePath, nil, opts)
			if err != nil {
				t.Fatal(err)
			}

			generatedFolder := path.Join(tmp, wantMeta.ID)
			files, err := ioutil.ReadDir(generatedFolder)
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("ExportCodelab generated files under %q", generatedFolder)
			for _, f := range files {
				t.Logf("Name: %s, IsDir: %v, Size: %d", f.Name(), f.IsDir(), f.Size())
			}

			wantBytes, err := ioutil.ReadFile(path.Join(tmp, wantMeta.ID, "index.html"))
			if err != nil {
				t.Fatal(err)
			}

			gotMeta, err := cmd.ExportCodelabMemory(testContent, gotBytes, opts)
			if err != nil {
				t.Errorf("ExportCodelabMemory got error %q, want nil", err)
			}

			// Because the In-Memory codelab doesn't have the source, when comparing, we remove Source
			wantMeta.Source = ""
			if !reflect.DeepEqual(wantMeta, gotMeta) {
				t.Errorf("ExportCodelabMemory returns metadata:\n%+v\nwant:\n%+v\n", gotMeta, wantMeta)
			}

			if bytes.Compare(wantBytes, gotBytes.Bytes()) != 0 {
				t.Errorf("ExportCodelabMemory returns diff: %s", cmp.Diff(string(wantBytes), string(gotBytes.Bytes())))
			}
		})
	}
}
