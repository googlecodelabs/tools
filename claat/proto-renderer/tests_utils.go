package genrenderer

import (
	"bytes"
	"io"
	"testing"

	"github.com/googlecodelabs/tools/third_party"
)

// unsupportedType is a dummy type used to showcase the failures of rendering
// non-proto custom types since we take in "any" type as rendering input.
type UnsupportedType struct{}

// Simple proto constructor
func NewDummyProto(in string) *devrel_tutorial.StylizedText {
	return &devrel_tutorial.StylizedText{
		Text: in,
	}
}

// Type for Render() handling output
type RenderedTemplate struct {
	txt io.Reader
	err error
}

// Type for canonical i != o template tests
type TestingBatch struct {
	// in  interface{}
	In  *RenderedTemplate
	Out string
	Ok  bool
}

// Converts output from Render() into a RenderedTemplate for batch processing
func AssertRenderedTemplate(r io.Reader, err error) *RenderedTemplate {
	return &RenderedTemplate{
		txt: r,
		err: err,
	}
}

// Helper for canonical i != o and !ok tests
func CanonicalTemplateTestBatch(tests []*TestingBatch, t *testing.T) {
	for _, tc := range tests {
		if tc.In.err != nil && tc.Ok {
			t.Errorf("boi")
			continue
		}

		if tc.Out != ReaderToString(tc.In.txt) {
			// change this
			t.Errorf("Expecting:\n\t%s, but got \n\t%#v", tc.Out, tc.In)
			continue
		}
	}
}

// readerToString Making io.Reader more readable for errors
func ReaderToString(i io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(i)
	return b.String()
}
