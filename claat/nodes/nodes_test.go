package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEmptyNodes(t *testing.T) {
	tests := []struct {
		name    string
		inNodes []Node
		out     bool
	}{
		{
			name: "Zero",
			out:  true,
		},
		{
			name: "OneEmpty",
			inNodes: []Node{
				NewTextNode(""),
			},
			out: true,
		},
		{
			name: "OneNonEmpty",
			inNodes: []Node{
				NewTextNode("foo"),
			},
		},
		{
			name: "MultipleEmpty",
			inNodes: []Node{
				NewTextNode(""),
				NewTextNode(""),
				NewTextNode(""),
			},
			out: true,
		},
		{
			name: "MultipleSomeNonEmpty",
			inNodes: []Node{
				NewTextNode("foo"),
				NewTextNode(""),
				NewTextNode("bar"),
			},
		},
		{
			name: "MultipleAllNonEmpty",
			inNodes: []Node{
				NewTextNode("foo"),
				NewTextNode("bar"),
				NewTextNode("baz"),
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := EmptyNodes(tc.inNodes)
			if out != tc.out {
				t.Errorf("EmptyNodes(%+v) = %t, want %t", tc.inNodes, out, tc.out)
				return
			}
		})
	}
}

func TestMutateEnv(t *testing.T) {
	tests := []struct {
		name  string
		inEnv []string
		out   []string
	}{
		{
			name:  "Sorted",
			inEnv: []string{"a", "b", "c"},
			out:   []string{"a", "b", "c"},
		},
		{
			name:  "Unsorted",
			inEnv: []string{"b", "c", "a"},
			out:   []string{"a", "b", "c"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewTextNode("foobar")
			n.MutateEnv(tc.inEnv)
			if diff := cmp.Diff(tc.out, n.env); diff != "" {
				t.Errorf("MutateEnv(%q) got diff (-want +got): %s", tc.inEnv, diff)
				return

			}

			// Also test that a copy was made.
			if &(tc.out) == &(n.env) {
				t.Errorf("MutateEnv(%q) did not copy input", tc.inEnv)
				return
			}
		})
	}
}
