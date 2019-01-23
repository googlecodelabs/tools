package filesystem

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// FileSystemFetcher implements fetch.Fetcher. It retrieves resources from the file system.
type FileSystemFetcher struct {
	resPath string
}

// NewFileSystemFetcher returns a new, initialized FileSystemFetcher.
// The input string is the path to the file to read the resource from.
func NewFileSystemFetcher(resPath string) FileSystemFetcher {
	return FileSystemFetcher{
		resPath: resPath,
	}
}

// Fetch fetches the resource.
// Instead of holding a file descriptor, the entire file is eagerly read into memory.
func (fsf FileSystemFetcher) Fetch() (io.Reader, error) {
	f, err := os.Open(fsf.resPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open file: %s", err)
	}
	defer f.Close()

	res, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file: %s", err)
	}

	return bytes.NewReader(res), nil
}
