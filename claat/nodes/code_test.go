package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewCodeNode(t *testing.T) {
	tests := []struct {
		name    string
		inValue string
		inTerm  bool
		inLang  string
		out     *CodeNode
	}{
		{
			name: "Empty",
			out: &CodeNode{
				node: node{typ: NodeCode},
			},
		},
		{
			name:    "Terminal",
			inValue: "sl",
			inTerm:  true,
			out: &CodeNode{
				node:  node{typ: NodeCode},
				Value: "sl",
				Term:  true,
			},
		},
		{
			name:    "SourceCode",
			inValue: "fmt.Println(\"foobar\")",
			inLang:  "go",
			out: &CodeNode{
				node:  node{typ: NodeCode},
				Value: "fmt.Println(\"foobar\")",
				Lang:  "go",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewCodeNode(tc.inValue, tc.inTerm, tc.inLang)
			if diff := cmp.Diff(tc.out, out, cmp.AllowUnexported(CodeNode{}, node{})); diff != "" {
				t.Errorf("NewCodeNode(%q, %t, %q) got diff (-want +got): %s", tc.inValue, tc.inTerm, tc.inLang, diff)
				return
			}
		})
	}
}

func TestCodeNodeEmpty(t *testing.T) {
	tests := []struct {
		name    string
		inValue string
		inTerm  bool
		inLang  string
		out     bool
	}{
		{
			name:   "EmptyTerminal",
			inTerm: true,
			out:    true,
		},
		{
			name:   "EmptySourceCode",
			inLang: "go",
			out:    true,
		},
		{
			name:    "NonEmptyTerminal",
			inTerm:  true,
			inValue: "sl",
		},
		{
			name:    "NonEmptySourceCode",
			inLang:  "go",
			inValue: "fmt.Println(\"foobar\")",
		},
		{
			name:    "EmptyWithSpacesTerminal",
			inTerm:  true,
			inValue: "\n \t",
			out:     true,
		},
		{
			name:    "EmptyWithSpacesSourceCode",
			inLang:  "go",
			inValue: "\n \t",
			out:     true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewCodeNode(tc.inValue, tc.inTerm, tc.inLang)
			out := n.Empty()
			if out != tc.out {
				t.Errorf("CodeNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
