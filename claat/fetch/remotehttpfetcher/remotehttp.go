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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// RemoteHTTPFetcher implements fetch.Fetcher. It retrieves resources over HTTP.
type RemoteHTTPFetcher struct {
	url string
	c   *http.Client
}

// New returns a new, initialized RemoteHTTPFetcher.
// The input string is the URL to read the resource from.
// If given, the http.Client will be used to make the request. Otherwise, the default client will be used.
func New(url string, c *http.Client) *RemoteHTTPFetcher {
	rhf := RemoteHTTPFetcher{
		url: url,
		c:   c,
	}
	if rhf.c == nil {
		rhf.c = http.DefaultClient
	}

	return &rhf
}

// Fetch fetches the resource.
func (rhf *RemoteHTTPFetcher) Fetch() (io.Reader, error) {
	res, err := rhf.c.Get(rhf.url)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error making HTTP request: status code %v and err %s", res.StatusCode, err)
	}
	defer res.Body.Close()

	resource, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %s", err)
	}

	return bytes.NewReader(resource), nil
}
