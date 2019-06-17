package remotehttp

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

		rhf := NewRemoteHTTPFetcher(server.URL, nil)
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
