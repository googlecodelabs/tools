package html

import (
	"testing"

	"github.com/googlecodelabs/tools/claat/proto-renderer/testing-utils"
)

func TestRender(t *testing.T) {
	tests := []struct {
		in interface{}
		ok bool
	}{
		// invalid cases
		{nil, false},
		{"invalid input type", false},
		{testingUtils.UnsupportedType{}, false},
		// valid cases
		{testingUtils.NewDummyProto("3"), true},
	}

	for _, tc := range tests {
		o, err := Render(tc.in)

		if err != nil && tc.ok {
			t.Errorf("\nRender(\n\t%#v\n)\nPanic: %v(false negative)", tc.in, err)
		}

		// plain want error, in != out verification is not in scope for 'Render'
		if err == nil && !tc.ok {
			rndrOut := testingUtils.ReaderToString(o)
			t.Errorf("\nRender(\n\t%#v\n) = %#v\nWant error\n(false positive)", tc.in, rndrOut)
		}
	}
}
