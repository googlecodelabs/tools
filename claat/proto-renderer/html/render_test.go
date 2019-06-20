package html

import (
	"bytes"
	"io"
	"testing"
)

type encapsulatedTest struct {
	in  interface{}
	out string
	ok  bool
}

type unsupportedProto struct{}

func TestRender(t *testing.T) {
	// TODO: update tests to be proto dependent on next PR
	tests := []encapsulatedTest{
		{nil, "", false},
		{"invalid", "", false},
		{unsupportedProto{}, "", false},
	}

	for _, tc := range tests {
		runEncapsulatedTest(tc, t)
	}
}

func runEncapsulatedTest(test encapsulatedTest, t *testing.T) {
	// Check wheather template failed to render by checking for panic
	defer func(test encapsulatedTest) {
		err := recover()
		if err != nil && test.ok {
			t.Errorf("\nRender(\n\t%#v\n) = %#v\nPanic occured:\n\t%#v\n(false negative)", test, test.out, err)
		}

		if err == nil && !test.ok {
			t.Errorf("\nRender(\n\t%#v\n) = %#v\nWant panic\n(false positive)", test, test.out)
		}
	}(test)

	o, _ := Render(test.in)
	// never gets here if above panicked
	rndrOut := readerToString(o)
	if test.out != rndrOut {
		t.Errorf("Expecting:\n\t'%s'\nBut got:\n\t'%s'", test.out, rndrOut)
	}
}

func readerToString(i io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(i)
	return b.String()
}
