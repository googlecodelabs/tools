package testingUtils

import (
	"bytes"
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

// Simple proto constructor
func NewDummyProto(in string) *tutorial.StylizedText {
	return &tutorial.StylizedText{
		Text: in,
	}
}

// Function signature for output-format agnositic 'Render'
type renderingFunc func(interface{}) (io.Reader, error)

// Type for canonical i != o and !ok rendering template tests
type RendererTestingBatch struct {
	InProto interface{}
	Out     string
	Ok      bool
}

// Helper for canonical i != o and !ok tests
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

// readerToString Making io.Reader more readable for errors
func ReaderToString(i io.Reader) string {
	if i == nil {
		return ""
	}
	var b bytes.Buffer
	b.ReadFrom(i)
	return b.String()
}
