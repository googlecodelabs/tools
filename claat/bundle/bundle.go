package bundle

import (
	"errors"
	"io"

	"github.com/googlecodelabs/tools/claat/types"

	"golang.org/x/net/context"
)

const assetsDirName = "img"

type ContentWriter interface {
	WriteAsset(ctx context.Context, clab, name string, body io.Reader) error
	WriteMarkup(ctx context.Context, clab string, body []byte) error
	WriteMeta(ctx context.Context, cmeta *types.ContextMeta) error
	// TODO: move these to a ContentLister
	//ListMeta(ctx context.Context) ([]*types.ContextMeta, error)
	//Assets(ctx context.Context, id string) ([]string, error)
}

type ContentBundler struct {
	Source ContentWriter
	Target ContentWriter
}

func (cb *ContentBundler) Sync(ctx context.Context, ids ...string) error {
	return errors.New("not implemented")
}

type View interface {
	//ListMeta(ctx context.Context) ([]*types.ViewMeta, error)
	//Assets(ctx context.Context, view string) ([]string, error)
	//WriteMeta(ctx context.Context, view string, meta *types.ViewMeta) error
	//WriteAsset(ctx context.Context, view, name string, body []byte) error
}

type ViewBundler struct {
	Source View
	Target View
}

func (vb *ViewBundler) Sync(ctx context.Context, ids ...string) error {
	return errors.New("not implemented")
}

type MultiError interface {
	Errors() []error
}

type multiError []error

func (me multiError) Errors() []error {
	return []error(me)
}

func (me multiError) Error() string {
	if len(me) == 0 {
		return "no errors"
	}
	if len(me) == 1 {
		return me[0].Error()
	}
	const sep = "\n"
	n := len(sep) * (len(me) - 1)
	for i := 0; i < len(me); i++ {
		n += len(me[i].Error())
	}

	b := make([]byte, n)
	bp := copy(b, me[0].Error())
	for _, s := range me[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], s.Error())
	}
	return string(b)
}
