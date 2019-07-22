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
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

// RendererIdendityBatch type for i != o and !ok rendering tests or oneof and their underlying proto
type RendererIdendityBatch struct {
	InProto  interface{}
	OutProto interface{}
	Out      string
	Ok       bool
}

// RenderingIdendityTestBatch is a wrapper on 'TestCanonicalRendererBatch' to prove that oneof types
// are equal to their underlying type rendering
func RenderingIdendityTestBatch(renderer renderingFunc, tests []*RendererIdendityBatch, t *testing.T) {
	for _, tc := range tests {
		rndrout, underlyingTypeErr := runEncapsulatedRendering(tc.OutProto, renderer, t)

		// ignore the normal set of error checks if the underlying rendering panicked
		if underlyingTypeErr != nil {
			funcName := runtime.FuncForPC(reflect.ValueOf(renderer).Pointer()).Name()
			cmd := fmt.Sprintf("\n%s(\n\t%#v\n)", funcName, tc.OutProto)
			t.Errorf("%s\nUnderlying rendering error: %v(false negative)\nWant: %#v", cmd, underlyingTypeErr, tc.Out)
			continue
		}

		// ignore the normal set of error checks if rendered OutProto != Out
		if tc.Out != rndrout {
			funcName := runtime.FuncForPC(reflect.ValueOf(renderer).Pointer()).Name()
			cmd := fmt.Sprintf("\n%s(\n\t%#v\n)", funcName, tc.OutProto)
			t.Errorf("%s = %#v\nBut want: \n%#v", cmd, rndrout, tc.Out)
			continue
		}

		// Create cannonical test from the output from the underlying type
		newTc := []*CanonicalRenderingBatch{
			{
				InProto: tc.InProto,
				Out:     tc.Out,
				Ok:      tc.Ok,
			},
		}
		TestCanonicalRendererBatch(renderer, newTc, t)
	}
}

// runEncapsulatedRendering constrains the scope of panics for 'RenderingIdendityTestBatch'
// otherwise we cannot iterate through consecutive panic-causing test-cases
func runEncapsulatedRendering(el interface{}, renderer renderingFunc, t *testing.T) (output interface{}, err error) {
	defer func() {
		r := recover()
		if r != nil {
			output = ""
			// not reusing genrenderer.AssertError due to import cycle
			switch r.(type) {
			case string:
				err = errors.New(r.(string))
			case error:
				err = r.(error)
			}
		}
	}()

	reader, err := renderer(el)
	output = ReaderToString(reader)
	return output, nil
}
