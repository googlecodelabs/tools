package genrenderer

import (
	"errors"
)

// AssertError handles type-assertions of panic-recovery messages.
// Used by output-format packages to delegate error-handling to callers
func AssertError(el interface{}) error {
	switch x := el.(type) {
	case string:
		return errors.New(x)
	case error:
		return x
	}
	return nil
}
