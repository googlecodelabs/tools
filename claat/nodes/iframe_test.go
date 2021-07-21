package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewIframeNode(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  *IframeNode
	}{
		{
			name: "Empty",
			out: &IframeNode{
				node: node{typ: NodeIframe},
			},
		},
		{
			name: "Simple",
			in:   "google.com",
			out: &IframeNode{
				node: node{typ: NodeIframe},
				URL:  "google.com",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewIframeNode(tc.in)
			if diff := cmp.Diff(tc.out, out, cmp.AllowUnexported(IframeNode{}, node{})); diff != "" {
				t.Errorf("NewIframeNode(%q) got diff (-want +got): %s", tc.in, diff)
				return
			}
		})
	}
}

func TestIframeNodeEmpty(t *testing.T) {
	tests := []struct {
		name string
		in   *IframeNode
		out  bool
	}{
		{
			name: "Empty",
			in: &IframeNode{
				node: node{typ: NodeIframe},
			},
			out: true,
		},
		{
			name: "NonEmpty",
			in: &IframeNode{
				node: node{typ: NodeIframe},
				URL:  "google.com",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := tc.in.Empty()
			if out != tc.out {
				t.Errorf("IframeNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
