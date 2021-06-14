package util

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUnique(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		out  []string
	}{
		{
			name: "HasDuplicatesSequential",
			in:   []string{"peach", "peach", "apple", "orange", "orange", "orange", "pineapple", "pineapple", "nectarine"},
			out:  []string{"peach", "apple", "orange", "pineapple", "nectarine"},
		},
		{
			name: "HasDuplicatesNonsequential",
			in:   []string{"cantaloupe", "watermelon", "cantaloupe", "honeydew", "cantaloupe", "honeydew", "honeydew", "watermelon"},
			out:  []string{"cantaloupe", "watermelon", "honeydew"},
		},
		{
			name: "NoDuplicates",
			in:   []string{"strawberry", "blackberry", "blueberry"},
			out:  []string{"strawberry", "blackberry", "blueberry"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := Unique(tc.in)
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("Unique(%+v) got diff (-want, +got)=\n%s", tc.in, diff)
			}
		})
	}
}
