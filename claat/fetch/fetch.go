package fetch

import "io"

type Fetcher interface {
	Fetch() (io.Reader, error)
}
