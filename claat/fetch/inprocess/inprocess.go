package inprocess

import "io"

// InProcessFetcher implements fetch.Fetcher. It retrieves resources from an in-process source via io.Reader.
type InProcessFetcher struct {
	source io.Reader
}

// NewInProcessFetcher returns a new, initialized InProcessFetcher.
// The input io.Reader is a reader over the resource bytes.
func NewInProcessFetcher(source io.Reader) InProcessFetcher {
	return InProcessFetcher{
		source: source,
	}
}

// Fetch fetches the resource.
// This doesn't really do anything.
func (ipf InProcessFetcher) Fetch() (io.Reader, error) {
	return ipf.source, nil
}
