// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package render

import (
	"bytes"
	"testing"

	"github.com/googlecodelabs/tools/claat/nodes"
	"github.com/googlecodelabs/tools/claat/types"
)

func TestExecuteBuiltin(t *testing.T) {
	step := &types.Step{
		Title:   "Test step",
		Content: nodes.NewListNode(nodes.NewTextNode(nodes.NewTextNodeOptions{Value: "text"})),
	}
	data := &struct {
		Context
	}{Context: Context{
		Meta:  &types.Meta{},
		Steps: []*types.Step{step},
	}}
	for _, f := range []string{"html", "md"} {
		var buf bytes.Buffer
		if err := Execute(&buf, f, data); err != nil {
			t.Errorf("%s: %v", f, err)
		}
	}
}
