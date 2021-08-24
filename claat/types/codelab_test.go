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

func TestNewStep(t *testing.T) {
	c := NewCodelab()
	s := c.NewStep("foobar")

	if c.Steps[len(c.Steps)-1] != s {
		t.Errorf(`Codelab.NewStep("foobar") did not return added step`)
	}
	if s.Title != "foobar" {
		t.Errorf(`Codelab.NewStep("foobar") got title %q, want "foobar"`, s.Title)
	}
	if s.Content == nil {
		t.Errorf(`Codelab.NewStep("foobar") did not initialize s.Content`)
	}
}
