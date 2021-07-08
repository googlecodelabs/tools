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

func TestNormalizedSplit(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  []string
	}{
		{
			name: "none",
			in:   "",
			out:  []string{},
		},
		{
			name: "one",
			in:   "peach",
			out:  []string{"peach"},
		},
		{
			name: "whitespace",
			in:   "   ",
			out:  []string{},
		},
		{
			name: "commas",
			in:   ",,,",
			out:  []string{},
		},
		{
			name: "split",
			in:   "peach,pear",
			out:  []string{"peach", "pear"},
		},
		{
			name: "split and trim space",
			in:   "peach , pear",
			out:  []string{"peach", "pear"},
		},
		{
			name: "split and trim consequtive spaces",
			in:   "peach  ,  pear",
			out:  []string{"peach", "pear"},
		},
		{
			name: "split and collapse space",
			in:   "p e a c h,pear",
			out:  []string{"peach", "pear"},
		},
		{
			name: "split and remove duplicates",
			in:   "peach,pear,pear",
			out:  []string{"peach", "pear"},
		},
		{
			name: "split and lowercase",
			in:   "PEACH,pear",
			out:  []string{"peach", "pear"},
		},
		{
			name: "split and strip new lines",
			in:   "pea\nch,pear\n",
			out:  []string{"peach", "pear"},
		},
		{
			name: "split and strip tabs",
			in:   "pea\tch,pear\t",
			out:  []string{"peach", "pear"},
		},
		{
			name: "split and strip whitespace",
			in:   "pea\tch, pear\n",
			out:  []string{"peach", "pear"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := NormalizedSplit(tc.in)
			if diff := cmp.Diff(tc.out, out); diff != "" {
				t.Errorf("NormalizedSplit(%+v) got diff (-want, +got)=\n%s", tc.in, diff)
			}
		})
	}
}
