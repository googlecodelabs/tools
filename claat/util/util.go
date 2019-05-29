package util

import (
	"log"
	"testing"
)

// Unique de-dupes a.
// The argument a is not modified.
func Unique(a []string) []string {
	seen := make(map[string]struct{}, len(a))
	res := make([]string, 0, len(a))
	for _, s := range a {
		if _, y := seen[s]; !y {
			res = append(res, s)
			seen[s] = struct{}{}
		}
	}
	return res
}

func LogIfError(err error) {
  if err != nil {
    log.Print(err)
  }
}

type TestingBatch struct {
	i string
	o string
}

func testMdBatch(tests []TestingBatch, t *testing.T) {
	for _, test := range tests {
		r := string(test)
		if test.o != r {
			t.Errorf("Expecting:\n\t'%s', but got \n\t'%s'", test.o, r)
			continue
		}
	}
}