package nodes

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewSurveyNode(t *testing.T) {
	tests := []struct {
		name     string
		inID     string
		inGroups []*SurveyGroup
		out      *SurveyNode
	}{
		{
			name: "Empty",
			out: &SurveyNode{
				node: node{typ: NodeSurvey},
			},
		},
		{
			// TODO: should the absence of an ID be an error?
			name: "GroupsNoID",
			inGroups: []*SurveyGroup{
				&SurveyGroup{
					Name:    "pick a number",
					Options: []string{"1", "2", "3"},
				},
				&SurveyGroup{
					Name:    "choose an answer",
					Options: []string{"yes", "no", "probably"},
				},
			},
			out: &SurveyNode{
				node: node{typ: NodeSurvey},
				Groups: []*SurveyGroup{
					&SurveyGroup{
						Name:    "pick a number",
						Options: []string{"1", "2", "3"},
					},
					&SurveyGroup{
						Name:    "choose an answer",
						Options: []string{"yes", "no", "probably"},
					},
				},
			},
		},
		{
			name: "IDNoGroups",
			inID: "identifier",
			out: &SurveyNode{
				node: node{typ: NodeSurvey},
				ID:   "identifier",
			},
		},
		{
			name: "Simple",
			inID: "a simple example",
			inGroups: []*SurveyGroup{
				&SurveyGroup{
					Name:    "pick a color",
					Options: []string{"red", "blue", "yellow"},
				},
			},
			out: &SurveyNode{
				node: node{typ: NodeSurvey},
				ID:   "a simple example",
				Groups: []*SurveyGroup{
					&SurveyGroup{
						Name:    "pick a color",
						Options: []string{"red", "blue", "yellow"},
					},
				},
			},
		},
		{
			name: "Multiple",
			inID: "an example with multiple survey groups",
			inGroups: []*SurveyGroup{
				&SurveyGroup{
					Name:    "a",
					Options: []string{"a", "aa", "aaa"},
				},
				&SurveyGroup{
					Name:    "b",
					Options: []string{"b", "bb", "bbb"},
				},
				&SurveyGroup{
					Name:    "c",
					Options: []string{"c", "cc", "ccc"},
				},
			},
			out: &SurveyNode{
				node: node{typ: NodeSurvey},
				ID:   "an example with multiple survey groups",
				Groups: []*SurveyGroup{
					&SurveyGroup{
						Name:    "a",
						Options: []string{"a", "aa", "aaa"},
					},
					&SurveyGroup{
						Name:    "b",
						Options: []string{"b", "bb", "bbb"},
					},
					&SurveyGroup{
						Name:    "c",
						Options: []string{"c", "cc", "ccc"},
					},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NewSurveyNode(tc.inID, tc.inGroups...)
			if diff := cmp.Diff(tc.out, out, cmp.AllowUnexported(SurveyNode{}, node{})); diff != "" {
				t.Errorf("NewSurveyNode(%q, %v) got diff (-want +got): %s", tc.inID, tc.inGroups, diff)
				return
			}
		})
	}
}

func TestSurveyNodeEmpty(t *testing.T) {
	tests := []struct {
		name string
		in   *SurveyNode
		out  bool
	}{
		{
			name: "NoGroups",
			in:   NewSurveyNode("id"),
			out:  true,
		},
		{
			name: "OneGroupEmpty",
			in: NewSurveyNode("id",
				&SurveyGroup{
					Name: "one",
				}),
			out: true,
		},
		{
			name: "MultiGroupsEmpty",
			in: NewSurveyNode("id",
				&SurveyGroup{
					Name: "one",
				},
				&SurveyGroup{
					Name: "two",
				},
				&SurveyGroup{
					Name: "three",
				},
			),
			out: true,
		},
		{
			name: "OneGroupNonEmpty",
			in: NewSurveyNode("id",
				&SurveyGroup{
					Name:    "one",
					Options: []string{"two", "three"},
				}),
		},
		{
			name: "MultiGroupsNonEmpty",
			in: NewSurveyNode("id",
				&SurveyGroup{
					Name:    "one",
					Options: []string{"two", "three"},
				},
				&SurveyGroup{
					Name:    "four",
					Options: []string{"five", "six"},
				},
				&SurveyGroup{
					Name:    "seven",
					Options: []string{"eight", "nine"},
				},
			),
		},
		{
			name: "MultiGroupsNonEmptySomeNoOptions",
			in: NewSurveyNode("id",
				&SurveyGroup{
					Name:    "one",
					Options: []string{"two", "three"},
				},
				&SurveyGroup{
					Name:    "four",
					Options: []string{"five", "six"},
				},
				&SurveyGroup{
					Name: "seven",
				},
			),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := tc.in.Empty()
			if out != tc.out {
				t.Errorf("SurveyNode.Empty() = %t, want %t", out, tc.out)
				return
			}
		})
	}
}
