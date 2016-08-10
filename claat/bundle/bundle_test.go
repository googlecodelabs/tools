package bundle

import (
	"errors"
	"testing"
)

func TestMultiError(t *testing.T) {
	table := []struct {
		err error
		out string
	}{
		{multiError{}, "no errors"},
		{multiError{errors.New("one")}, "one"},
		{multiError{errors.New("one"), errors.New("two")}, "one\ntwo"},
	}
	for i, test := range table {
		if out := test.err.Error(); out != test.out {
			t.Errorf("%d: out = %q; want %q", i, out, test.out)
		}
	}
}
