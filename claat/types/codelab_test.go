package types

import (
	"testing"
)

func TestNewCodelab(t *testing.T) {
	c := NewCodelab()
	if c.Extra == nil {
		t.Errorf("NewCodelab() failed to initialize Extra")
	}
}
