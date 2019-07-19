package testingutils

import (
	"bytes"
	"io"
)

// ReaderToString makes io.Reader more readable for errors
func ReaderToString(i io.Reader) string {
	if i == nil {
		return ""
	}
	var b bytes.Buffer
	b.ReadFrom(i)
	return b.String()
}

// renderingFunc is the tunction signature for output-format agnositic 'Render'
type renderingFunc func(interface{}) (io.Reader, error)
