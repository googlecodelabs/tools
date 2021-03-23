// Copyright 2018-2019 Google LLC. All Rights Reserved.
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

// The claat command generates one or more codelabs from "source" documents,
// specified as either Google Doc IDs or local markdown files.
// The command also allows one to preview generated codelabs from local drive
// using "claat serve".
// See more details at https://github.com/googlecodelabs/tools.
package cmd

import (
	"path/filepath"

	"github.com/googlecodelabs/tools/claat/types"

	// allow parsers to register themselves
	_ "github.com/googlecodelabs/tools/claat/parser/gdoc"
	_ "github.com/googlecodelabs/tools/claat/parser/md"
)

const (
	// metaFilename is codelab metadata file.
	metaFilename = "codelab.json"
	// stdout is a special value for -o cli arg to identify stdout writer.
	stdout = "-"

	// log report formats
	reportErr = "err\t%s %v"
	reportOk  = "ok\t%s"
)

// isStdout reports whether filename is stdout.
func isStdout(filename string) bool {
	return filename == stdout
}

// codelabDir returns codelab root directory.
// The base argument is codelab parent directory.
func codelabDir(base string, m *types.Meta) string {
	return filepath.Join(base, m.ID)
}
