package gdoc

import (
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func makeBlinkNode() *html.Node {
	return &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Blink,
		Data:     "blink",
	}
}

func TestIsHeader(t *testing.T) {
	tests := []struct {
		name string
		in   *html.Node
		out  bool
	}{
		{
			name: "StepTitle",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H1,
				Data:     "h1",
			},
		},
		{
			name: "FirstLevel",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H2,
				Data:     "h2",
			},
			out: true,
		},
		{
			name: "SecondLevel",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H3,
				Data:     "h3",
			},
			out: true,
		},
		{
			name: "ThirdLevel",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H4,
				Data:     "h4",
			},
			out: true,
		},
		{
			name: "FourthLevel",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H5,
				Data:     "h5",
			},
			out: true,
		},
		{
			name: "FifthLevel",
			in: &html.Node{
				Type:     html.ElementNode,
				DataAtom: atom.H6,
				Data:     "h6",
			},
			out: true,
		},
		{
			name: "NotAHeader",
			in:   makeBlinkNode(),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if out := isHeader(tc.in); out != tc.out {
				t.Errorf("isHeader(%v) = %t, want %t", tc.in, out, tc.out)
			}
		})
	}
}
