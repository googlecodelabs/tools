package filesystem

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNewFileSystemFetcher(t *testing.T) {
	p := "this/is/a/file/path"
	fsf := NewFileSystemFetcher(p)

	if fsf.resPath != p {
		t.Errorf("NewFileSystemFetcher(%v).resPath = %v, want %v", p, fsf.resPath, p)
	}
}

func TestFetch(t *testing.T) {
	// Make temporary file for testing purposes.
	contents := []byte("file contents!")
	f, err := ioutil.TempFile("", "filesystem_test_file")
	if err != nil {
		t.Errorf("error creating temp file: %s", err)
	}
	fname := f.Name()
	defer os.Remove(fname)

	// Write some bytes to the file.
	_, err = f.Write(contents)
	if err != nil {
		// Be sure to close the file if we're going to stop here.
		f.Close()
		t.Errorf("error writing to temp file: %s", err)
	}

	f.Close()

	fsf := NewFileSystemFetcher(fname)
	r, err := fsf.Fetch()
	if err != nil {
		t.Errorf("Fetch() = got err %v, want nil", err)
	}

	// Get the bytes out of the reader and compare.
	res, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("error reading from Fetch() result: %s", err)
	}
	if !reflect.DeepEqual(res, contents) {
		t.Errorf("Fetch() reader got %v, want %v", res, contents)
	}
}
