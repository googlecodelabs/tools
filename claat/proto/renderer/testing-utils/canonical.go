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
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

// CanonicalRenderingBatch type for canonical i != o and !ok rendering tests
type CanonicalRenderingBatch struct {
	InProto interface{}
	Out     string
	Ok      bool
}

// TestCanonicalRendererBatch is the helper for canonical i != o and !ok rendering tests
func TestCanonicalRendererBatch(renderer renderingFunc, tests []*CanonicalRenderingBatch, t *testing.T) {
	for _, tc := range tests {
		funcName := runtime.FuncForPC(reflect.ValueOf(renderer).Pointer()).Name()
		reader, err := renderer(tc.InProto)

		cmd := fmt.Sprintf("\n%s(\n\t%#v\n)", funcName, tc.InProto)

		if err != nil && tc.Ok {
			t.Errorf("%s\nError: %v(false negative)\nWant: %#v", cmd, err, tc.Out)
		}

		if err == nil && !tc.Ok {
			t.Errorf("%s\n = %#v\nWant Error\n(false positive)", cmd, reader)
		}

		rndrout := ReaderToString(reader)
		if tc.Out != rndrout {
			t.Errorf("%s = %#v\nBut want: \n%#v", cmd, rndrout, tc.Out)
		}
	}
}
