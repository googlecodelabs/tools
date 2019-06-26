package testingUtils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"testing"

	"github.com/googlecodelabs/tools/third_party"
)

// unsupportedType is a dummy type used to showcase the failures of rendering
// non-proto custom types since we take in "any" type as rendering input.
type UnsupportedType struct{}

// NewDummyProto is a simple proto constructor
func NewDummyProto() *tutorial.StylizedText {
	return &tutorial.StylizedText{
		Text: "dummy",
	}
}

// renderingFunc is the tunction signature for output-format agnositic 'Render'
type renderingFunc func(interface{}) (io.Reader, error)

// RendererTestingBatch type for canonical i != o and !ok rendering template tests
type RendererTestingBatch struct {
	InProto interface{}
	Out     string
	Ok      bool
}

// CanonicalRenderingTestBatch helper for canonical i != o and !ok tests
func CanonicalRenderingTestBatch(renderer renderingFunc, tests []*RendererTestingBatch, t *testing.T) {
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

// readerToString makes io.Reader more readable for errors
func ReaderToString(i io.Reader) string {
	if i == nil {
		return ""
	}
	var b bytes.Buffer
	b.ReadFrom(i)
	return b.String()
}

// RendererTestingIdendityBatch type for i != o and !ok rendering tests or oneof and their underlying proto
type RendererTestingIdendityBatch struct {
	InProto  interface{}
	OutProto interface{}
	Out      string
	Ok       bool
}

// RenderingTestIdendityBatch is a wrapper on 'CanonicalRenderingTestBatch' to prove that oneof types
// are equal to their underlying type rendering
func RenderingTestIdendityBatch(renderer renderingFunc, tests []*RendererTestingIdendityBatch, t *testing.T) {
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
		newTc := []*RendererTestingBatch{
			{
				tc.InProto,
				tc.Out,
				tc.Ok,
			},
		}
		CanonicalRenderingTestBatch(renderer, newTc, t)
	}
}

// runEncapsulatedRendering constrains the scope of panics for 'RenderingTestIdendityBatch'
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
