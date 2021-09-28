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
				NewTextNode(NewTextNodeOptions{Value: ""}),
			},
			out: true,
		},
		{
			name: "OneNonEmpty",
			inNodes: []Node{
				NewTextNode(NewTextNodeOptions{Value: "foo"}),
			},
		},
		{
			name: "MultipleEmpty",
			inNodes: []Node{
				NewTextNode(NewTextNodeOptions{Value: ""}),
				NewTextNode(NewTextNodeOptions{Value: ""}),
				NewTextNode(NewTextNodeOptions{Value: ""}),
			},
			out: true,
		},
		{
			name: "MultipleSomeNonEmpty",
			inNodes: []Node{
				NewTextNode(NewTextNodeOptions{Value: "foo"}),
				NewTextNode(NewTextNodeOptions{Value: ""}),
				NewTextNode(NewTextNodeOptions{Value: "bar"}),
			},
		},
		{
			name: "MultipleAllNonEmpty",
			inNodes: []Node{
				NewTextNode(NewTextNodeOptions{Value: "foo"}),
				NewTextNode(NewTextNodeOptions{Value: "bar"}),
				NewTextNode(NewTextNodeOptions{Value: "baz"}),
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
			n := NewTextNode(NewTextNodeOptions{Value: "foobar"})
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
