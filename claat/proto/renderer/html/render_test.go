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
package html

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto/renderer/testing-utils"
)

func TestRender(t *testing.T) {
	tests := []struct {
		in interface{}
		ok bool
	}{
		// invalid cases
		{nil, false},
		{"invalid input type", false},
		{testingutils.UnsupportedType{}, false},
		// valid cases
		{testingutils.NewDummyProto(), true},
	}

	for _, tc := range tests {
		o, err := Render(tc.in)

		if err != nil && tc.ok {
			t.Errorf("\nRender(\n\t%#v\n)\nPanic: %v(false negative)", tc.in, err)
		}

		// plain want error, in != out verification is not in scope for 'Render'
		if err == nil && !tc.ok {
			rndrOut := testingutils.ReaderToString(o)
			t.Errorf("\nRender(\n\t%#v\n) = %#v\nWant error\n(false positive)", tc.in, rndrOut)
		}
	}
}
