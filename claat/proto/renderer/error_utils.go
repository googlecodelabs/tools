// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
