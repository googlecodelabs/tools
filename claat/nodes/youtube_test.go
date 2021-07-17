package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewYouTubeNode(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  *YouTubeNode
	}{
		{
			name: "Empty",
			out: &YouTubeNode{
				node: node{typ: NodeYouTube},
			},
		},
		{
			name: "NonEmpty",
			in:   "Mlk888FiI8A",
			out: &YouTubeNode{
				node:    node{typ: NodeYouTube},
				VideoID: "Mlk888FiI8A",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewYouTubeNode(tc.in)
			if diff := cmp.Diff(tc.out, out, cmp.AllowUnexported(YouTubeNode{}, node{})); diff != "" {
				t.Errorf("NewYouTubeNode(%q) got diff (-want +got): %s", tc.in, diff)
				return
			}
		})
	}
}

func TestYouTubeNodeEmpty(t *testing.T) {
	tests := []struct {
		name string
		in   *YouTubeNode
		out  bool
	}{
		{
			name: "Empty",
			in: &YouTubeNode{
				node: node{typ: NodeYouTube},
			},
			out: true,
		},
		{
			name: "NonEmpty",
			in: &YouTubeNode{
				node:    node{typ: NodeYouTube},
				VideoID: "Mlk888FiI8A",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := tc.in.Empty()
			if out != tc.out {
				t.Errorf("YouTubeNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
