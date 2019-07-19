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
package remotehttpfetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFetch(t *testing.T) {
	type input struct {
		statusCode int
	}
	type output struct {
		ok bool
	}
	tests := []struct {
		in  input
		out output
	}{
		{
			in: input{
				statusCode: http.StatusOK,
			},
			out: output{
				ok: true,
			},
		},
		{
			in: input{
				statusCode: http.StatusNotFound,
			},
			out: output{},
		},
	}
	for _, tc := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.in.statusCode)
			fmt.Fprintf(w, "data")
		}))
		// Don't defer Close() here because this runs in a loop.

		rhf := New(server.URL, nil)
		res, err := rhf.Fetch()
		server.Close()

		if err == nil && !tc.out.ok {
			t.Errorf("Fetch() = %v, want err", res)
		}
		if err != nil && tc.out.ok {
			t.Errorf("Fetch() = %s, want []byte(\"data\")", err)
		}

		if err == nil {
			resBytes, err := ioutil.ReadAll(res)
			if err != nil {
				t.Errorf("could not read Fetch() result: %s", err)
				break
			}
			if !reflect.DeepEqual(resBytes, []byte("data")) {
				t.Errorf("Fetch() = %v, want []byte(\"data\")", resBytes)
			}
		}
	}
}
