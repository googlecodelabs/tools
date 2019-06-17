package remotehttp

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

// NewRemoteHTTPFetcher returns a new, initialized RemoteHTTPFetcher.
// The input string is the URL to read the resource from.
// If given, the http.Client will be used to make the request. Otherwise, the default client will be used.
func NewRemoteHTTPFetcher(url string, c *http.Client) RemoteHTTPFetcher {
	rhf := RemoteHTTPFetcher{
		url: url,
		c:   c,
	}
	if rhf.c == nil {
		rhf.c = http.DefaultClient
	}

	return rhf
}

// Fetch fetches the resource.
func (rhf RemoteHTTPFetcher) Fetch() (io.Reader, error) {
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
