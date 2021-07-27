package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewYouTubeNode(t *testing.T) {
	tests := []struct {
		name      string
		inVideoID string
		out       *YouTubeNode
	}{
		{
			name: "Empty",
			out: &YouTubeNode{
				node: node{typ: NodeYouTube},
			},
		},
		{
			name:      "NonEmpty",
			inVideoID: "Mlk888FiI8A",
			out: &YouTubeNode{
				node:    node{typ: NodeYouTube},
				VideoID: "Mlk888FiI8A",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewYouTubeNode(tc.inVideoID)
			if diff := cmp.Diff(tc.out, out, cmp.AllowUnexported(YouTubeNode{}, node{})); diff != "" {
				t.Errorf("NewYouTubeNode(%q) got diff (-want +got): %s", tc.inVideoID, diff)
				return
			}
		})
	}
}

func TestYouTubeNodeEmpty(t *testing.T) {
	tests := []struct {
		name      string
		inVideoID string
		out       bool
	}{
		{
			name: "Empty",
			out:  true,
		},
		{
			name:      "NonEmpty",
			inVideoID: "Mlk888FiI8A",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewYouTubeNode(tc.inVideoID)
			out := n.Empty()
			if out != tc.out {
				t.Errorf("YouTubeNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
