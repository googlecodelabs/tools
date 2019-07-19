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
