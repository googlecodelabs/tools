// Copyright 2018 Google Inc. All Rights Reserved.
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
	"encoding/json"
	"flag"
	"log"
	"os"
	"sync"

	// allow parsers to register themselves
	_ "github.com/googlecodelabs/tools/claat/parser/gdoc"
	_ "github.com/googlecodelabs/tools/claat/parser/md"
)

var (
	authToken = flag.String("auth", "", "OAuth2 Bearer token; alternative credentials override.")
	output    = flag.String("o", ".", "output directory or '-' for stdout")
	expenv    = flag.String("e", "web", "codelab environment")
	tmplout   = flag.String("f", "html", "output format")
	prefix    = flag.String("prefix", "../../", "URL prefix for html format")
	globalGA  = flag.String("ga", "UA-49880327-14", "global Google Analytics account")
	extra     = flag.String("extra", "", "Additional arguments to pass to format templates. JSON object of string,string key values.")
	addr      = flag.String("addr", "localhost:9090", "hostname and port to bind web server to")
)

const (
	// imgDirname is where a codelab images are stored,
	// relative to the codelab dir.
	imgDirname = "img"
	// metaFilename is codelab metadata file.
	metaFilename = "codelab.json"
	// stdout is a special value for -o cli arg to identify stdout writer.
	stdout = "-"

	// log report formats
	reportErr = "err\t%s %v"
	reportOk  = "ok\t%s"
)

var (
	Exit      int               // program exit code
	exitMu    sync.Mutex        // guards exit
	ExtraVars map[string]string // Extra template variables passed on the command line.
)

// isStdout reports whether filename is stdout.
func isStdout(filename string) bool {
	return filename == stdout
}

// printf prints formatted string fmt with args to stderr.
func printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// errorf calls printf with fmt and args, and sets non-zero exit code.
func errorf(format string, args ...interface{}) {
	printf(format, args...)
	exitMu.Lock()
	Exit = 1
	exitMu.Unlock()
}

// Fatalf calls printf and exits immediately with non-zero code.
func Fatalf(format string, args ...interface{}) {
	printf(format, args...)
	os.Exit(1)
}

// ParseExtraVars parses extra template variables from command line.
func ParseExtraVars() map[string]string {
	vars := make(map[string]string)
	if *extra == "" {
		return vars
	}
	b := []byte(*extra)
	err := json.Unmarshal(b, &vars)
	if err != nil {
		errorf("Error parsing additional template data: %v", err)
	}
	return vars
}
