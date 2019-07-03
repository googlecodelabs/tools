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
