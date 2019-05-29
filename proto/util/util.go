package util

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

type TestingBatch struct {
	i string
	o string
}