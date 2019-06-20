package html

import (
	"bytes"
	"io"
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer"
)

// unsupportedType is a dummy type used to showcase the failures of rendering
// non-proto custom types since we take in "any" type as rendering input.
type unsupportedType struct{}

func TestRender(t *testing.T) {
	// TODO: update to be proto dependent on next PR
	tests := []struct {
		in interface{}
		ok bool
	}{
		{nil, false},
		{"invalid input type", false},
		{unsupportedType{}, false},
		{genrenderer.NewSampleProtoTemplate("3"), true},
	}

	for _, tc := range tests {
		o, err := Render(tc.in)

		if err != nil && tc.ok {
			t.Errorf("\nRender(\n\t%#v\n)\nPanic: %v(false negative)", tc.in, err)
		}

		// plain want error, in != out verification is not in scope for 'Render'
		if err == nil && !tc.ok {
			rndrOut := readerToString(o)
			t.Errorf("\nRender(\n\t%#v\n) = %#v\nWant error\n(false positive)", tc.in, rndrOut)
		}
	}
}

// readerToString Making io.Reader more readable for errors
func readerToString(i io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(i)
	return b.String()
}
