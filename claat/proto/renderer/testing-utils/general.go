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
	"bytes"
	"io"
)

// ReaderToString makes io.Reader more readable for errors
func ReaderToString(i io.Reader) string {
	if i == nil {
		return ""
	}
	var b bytes.Buffer
	b.ReadFrom(i)
	return b.String()
}

// renderingFunc is the function signature for output-format agnositic 'Render'
type renderingFunc func(interface{}) (io.Reader, error)
