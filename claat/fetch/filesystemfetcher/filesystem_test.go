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
package filesystemfetcher

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	p := "this/is/a/file/path"
	fsf := New(p)

	if fsf.resPath != p {
		t.Errorf("New(%v).resPath = %v, want %v", p, fsf.resPath, p)
	}
}

func TestFetch(t *testing.T) {
	// Make temporary file for testing purposes.
	contents := []byte("file contents!")
	f, err := ioutil.TempFile("", "filesystem_test_file")
	if err != nil {
		t.Errorf("error creating temp file: %s", err)
	}
	fname := f.Name()
	defer os.Remove(fname)

	// Write some bytes to the file.
	_, err = f.Write(contents)
	defer f.Close()
	if err != nil {
		t.Errorf("error writing to temp file: %s", err)
	}

	fsf := New(fname)
	r, err := fsf.Fetch()
	if err != nil {
		t.Errorf("Fetch() = got err %v, want nil", err)
	}

	// Get the bytes out of the reader and compare.
	res, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("error reading from Fetch() result: %s", err)
	}
	if !reflect.DeepEqual(res, contents) {
		t.Errorf("Fetch() reader got %v, want %v", res, contents)
	}
}
