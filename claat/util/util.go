// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package util

import (
	"strings"
	"unicode"
)

// ImgDirname is where a codelab images are stored,
// relative to the codelab dir.
const ImgDirname = "img"

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

// NormalizedSplit takes a string, removes spaces, splits it along a comma delimiter, then on each fragment, trims Unicode spaces
// from both ends and converts them to lowercase. It returns a slice of the unique processed strings.
func NormalizedSplit(s string) []string {
	s = stripSpaces(s)
	s = strings.Trim(s, ",")
	if s == "" {
		return []string{}
	}
	strs := strings.Split(s, ",")
	for k, v := range strs {
		strs[k] = strings.ToLower(v)
	}
	return Unique(strs)
}

func stripSpaces(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, ch := range s {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}
