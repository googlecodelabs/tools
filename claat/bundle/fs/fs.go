package fs

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/googlecodelabs/tools/claat/types"
	"golang.org/x/net/context"
)

// Dir implements ContentWriter for a local file system
type Dir string

// WriteAsset reads from r and writes to a file at d/clab/name.
func (d Dir) WriteAsset(ctx context.Context, clab, name string, r io.Reader) error {
	p := cleanPath(string(d), clab, name)
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// WriteMarkup writes codelab markup body to a file at d/clab/index.html.
func (d Dir) WriteMarkup(ctx context.Context, clab string, body []byte) error {
	p := cleanPath(string(d), clab, "index.html")
	return ioutil.WriteFile(p, body, 0644)
}

// WriteMeta writes codelab metadata to a file at r/meta.ID/codelab.json in JSON format.
func (d Dir) WriteMeta(ctx context.Context, meta *types.ContextMeta) error {
	b, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	p := filepath.Join(string(d), meta.ID, "codelab.json")
	return ioutil.WriteFile(p, b, 0644)
}

func cleanPath(root string, p ...string) string {
	a := make([]string, len(p))
	for i, v := range p {
		a[i] = filepath.Clean(string(filepath.Separator) + v)
	}
	a = append([]string{root}, a...)
	return filepath.Join(a...)
}
