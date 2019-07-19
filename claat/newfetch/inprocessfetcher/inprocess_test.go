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
package inprocessfetcher

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	r := strings.NewReader("this is a string")
	ipf := New(r)

	if ipf.source != r {
		t.Errorf("New(%v).source = %v, want %v", r, ipf.source, r)
	}
}

func TestFetch(t *testing.T) {
	r := strings.NewReader("this is also a string")
	ipf := New(r)

	out, err := ipf.Fetch()
	if err != nil {
		t.Errorf("Fetch() got err %v, want nil", err)
	}
	if out != r {
		t.Errorf("Fetch() = %v, want %v", out, r)
	}
}
