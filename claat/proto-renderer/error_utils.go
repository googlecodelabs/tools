package genrenderer

import (
	"errors"
)

func AssertError(el interface{}) error {
	var err error
	switch x := el.(type) {
	case string:
		err = errors.New(x)
	case error:
		err = x
	}
	return err
}
