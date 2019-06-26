package genrenderer

import (
	"errors"
	"fmt"
)

// AssertError handles type-assertions of panic-recovery messages.
// Used by output-format packages to delegate error-handling to callers
func AssertError(el interface{}) error {
	switch x := el.(type) {
	case string:
		return errors.New(el.(string))
	case error:
		return x
	}
	return nil
}

// TypeNotSupported allows for debug-friendly nested-error panic messages
func TypeNotSupported(funcName string, el interface{}) string {
	return fmt.Sprintf("%s: type not supported: %T %#v", funcName, el, el)
}
