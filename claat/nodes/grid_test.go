package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var cmpOptGrid = cmp.AllowUnexported(GridNode{}, node{}, ListNode{}, TextNode{})

func TestNewGridNode(t *testing.T) {
	tests := []struct {
		name   string
		inRows [][]*GridCell
		out    *GridNode
	}{
		{
			name: "Empty",
			out: &GridNode{
				node: node{typ: NodeGrid},
			},
		},
		{
			name: "OneRow",
			inRows: [][]*GridCell{
				[]*GridCell{
					&GridCell{
						Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "aaa"})),
					},
					&GridCell{
						Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "bbb"})),
					},
				},
			},
			out: &GridNode{
				node: node{typ: NodeGrid},
				Rows: [][]*GridCell{
					[]*GridCell{
						&GridCell{
							Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "aaa"})),
						},
						&GridCell{
							Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "bbb"})),
						},
					},
				},
			},
		},
		{
			name: "MultipleRows",
			inRows: [][]*GridCell{
				[]*GridCell{
					&GridCell{
						Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "aaa"})),
					},
					&GridCell{
						Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "bbb"})),
					},
				},
				[]*GridCell{
					&GridCell{
						Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "ccc"})),
					},
					&GridCell{
						Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "ddd"})),
					},
				},
				[]*GridCell{
					&GridCell{
						Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "eee"})),
					},
					&GridCell{
						Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "fff"})),
					},
				},
			},
			out: &GridNode{
				node: node{typ: NodeGrid},
				Rows: [][]*GridCell{
					[]*GridCell{
						&GridCell{
							Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "aaa"})),
						},
						&GridCell{
							Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "bbb"})),
						},
					},
					[]*GridCell{
						&GridCell{
							Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "ccc"})),
						},
						&GridCell{
							Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "ddd"})),
						},
					},
					[]*GridCell{
						&GridCell{
							Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "eee"})),
						},
						&GridCell{
							Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "fff"})),
						},
					},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewGridNode(tc.inRows...)
			if diff := cmp.Diff(tc.out, out, cmpOptGrid); diff != "" {
				t.Errorf("NewGridNode(%v) got diff (-want +got): %s", tc.inRows, diff)
				return
			}
		})
	}
}

func TestGridNodeEmpty(t *testing.T) {
	tests := []struct {
		name   string
		inRows [][]*GridCell
		out    bool
	}{
		{
			name: "Empty",
			out:  true,
		},
		{
			name: "NonEmpty",
			inRows: [][]*GridCell{
				[]*GridCell{
					&GridCell{
						Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "aaa"})),
					},
					&GridCell{
						Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "bbb"})),
					},
				},
				[]*GridCell{
					&GridCell{
						Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "ccc"})),
					},
					&GridCell{
						Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "ddd"})),
					},
				},
				[]*GridCell{
					&GridCell{
						Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "eee"})),
					},
					&GridCell{
						Content: NewListNode(NewTextNode(NewTextNodeOptions{Value: "fff"})),
					},
				},
			},
		},
		{
			name: "EmptyWithRows",
			inRows: [][]*GridCell{
				[]*GridCell{
					&GridCell{
						Content: NewListNode(),
					},
					&GridCell{
						Content: NewListNode(),
					},
				},
			},
			out: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			n := NewGridNode(tc.inRows...)
			out := n.Empty()
			if out != tc.out {
				t.Errorf("GridNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
