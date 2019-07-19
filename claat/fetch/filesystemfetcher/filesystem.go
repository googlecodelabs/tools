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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// FileSystemFetcher implements fetch.Fetcher. It retrieves resources from the file system.
type FileSystemFetcher struct {
	resPath string
}

// New returns a new, initialized FileSystemFetcher.
// The input string is the path to the file to read the resource from.
func New(resPath string) *FileSystemFetcher {
	return &FileSystemFetcher{
		resPath: resPath,
	}
}

// Fetch fetches the resource.
// Instead of holding a file descriptor, the entire file is eagerly read into memory.
func (fsf *FileSystemFetcher) Fetch() (io.Reader, error) {
	f, err := os.Open(fsf.resPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open file: %s", err)
	}
	defer f.Close()

	res, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file: %s", err)
	}

	return bytes.NewReader(res), nil
}
