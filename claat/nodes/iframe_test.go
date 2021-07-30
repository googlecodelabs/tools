package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewIframeNode(t *testing.T) {
	tests := []struct {
		name  string
		inURL string
		out   *IframeNode
	}{
		{
			name: "Empty",
			out: &IframeNode{
				node: node{typ: NodeIframe},
			},
		},
		{
			name:  "Simple",
			inURL: "google.com",
			out: &IframeNode{
				node: node{typ: NodeIframe},
				URL:  "google.com",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewIframeNode(tc.inURL)
			if diff := cmp.Diff(tc.out, out, cmp.AllowUnexported(IframeNode{}, node{})); diff != "" {
				t.Errorf("NewIframeNode(%q) got diff (-want +got): %s", tc.inURL, diff)
				return
			}
		})
	}
}

func TestIframeNodeEmpty(t *testing.T) {
	tests := []struct {
		name  string
		inURL string
		out   bool
	}{
		{
			name: "Empty",
			out:  true,
		},
		{
			name:  "NonEmpty",
			inURL: "google.com",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewIframeNode(tc.inURL)
			out := n.Empty()
			if out != tc.out {
				t.Errorf("IframeNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
