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

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	// allow parsers to register themselves
	_ "github.com/googlecodelabs/tools/claat/parser/gdoc"
)

var (
	output   = flag.String("o", ".", "output directory or '-' for stdout")
	expenv   = flag.String("e", "web", "codelab environment")
	tmplout  = flag.String("f", "html", "output format")
	prefix   = flag.String("prefix", "../../", "URL prefix for html format")
	globalGA = flag.String("ga", "UA-49880327-14", "global Google Analytics account")
	tmpldata = flag.String("data", "", "Additional arguments to pass to format templates. JSON object of string,string key values.")

	version string // set by linker -X
)

const (
	// imgDirname is where a codelab images are stored,
	// relative to the codelab dir.
	imgDirname = "img"
	// contentFilename is the name of file for codelab content output,
	// without the format extension.
	contentFilename = "index"
	// metaFilename is codelab metadata file.
	metaFilename = "codelab.json"
	// stdout is a special value for -o cli arg to identify stdout writer.
	stdout = "-"

	// log report formats
	reportErr = "err\t%s %v"
	reportOk  = "ok\t%s"
)

var (
	// commands contains all valid subcommands, e.g. "claat export".
	commands = map[string]func(){
		"export":  cmdExport,
		"update":  cmdUpdate,
		"help":    usage,
		"version": func() { fmt.Println(version) },
	}

	exitMu sync.Mutex // guards exit
	exit   int        // program exit code
)

// isStdout reports whether filename is stdout.
func isStdout(filename string) bool {
	return filename == stdout
}

// contentFile returns codelab main output file given the specified format.
func contentFile(format string) string {
	return contentFilename + "." + format
}

// printf prints formatted string fmt with args to stderr.
func printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// errorf calls printf with fmt and args, and sets non-zero exit code.
func errorf(format string, args ...interface{}) {
	printf(format, args...)
	exitMu.Lock()
	exit = 1
	exitMu.Unlock()
}

// fatalf calls printf and exits immediatly with non-zero code.
func fatalf(format string, args ...interface{}) {
	printf(format, args...)
	os.Exit(1)
}

func main() {
	log.SetFlags(0)
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) == 1 {
		fatalf("Need subcommand. Try '-h' for options.")
	}
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		usage()
		return
	}

	cmd := commands[os.Args[1]]
	if cmd == nil {
		fatalf("Unknown subcommand. Try '-h' for options.")
	}
	flag.Usage = usage
	flag.CommandLine.Parse(os.Args[2:])

	cmd()
	os.Exit(exit)
}

// usage prints usageText and program arguments to stderr.
func usage() {
	fmt.Fprint(os.Stderr, usageText)
	flag.PrintDefaults()
}

const usageText = `Usage: claat <cmd> [export flags] src [src ...]

Available commands are: export, update, version.

## Export command

Export takes one or more 'src' documents and converts them
to the format specified with -f option.

The following formats are built-in:
- html (Polymer-based app)
- md (Markdown)
To use a custom format, specify a local file path to a Go template file.
More info on Go templates: https://golang.org/pkg/text/template/.

Each 'src' can be either a remote HTTP resource or a local file.
Source formats currently supported are:
- Google Doc (Codelab Format, go/codelab-guide)
- Markdown

When 'src' is a Google Doc, it must be specified as a doc ID,
omitting https://docs.google.com/... part.

Instead of writing to an output directory, use "-o -" to specify
stdout. In this case images and metadata are not exported.
When writing to a directory, existing files will be overwritten.

The program exits with non-zero code if at least one src could not be exported.

## Update command

Update scans one or more 'src' local directories for codelab.json metadata
files, recursively. A directory containing the metadata file is expected
to be a codelab previously created with the export command.

Current directory is assumed if no 'src' argument is given.

Each found codelab is then re-exported using parameters from the metadata file.
Unused codelab assets will be deleted, as well as the entire codelab directory,
if codelab ID has changed since last update or export.

In the latter case, where codelab ID has changed, the new directory
will be placed alongside the old one. In other words, it will have the same ancestor
as the old one.

While -prefix and -ga can override existing codelab metadata, the other
arguments have no effect during update.

The program does not follow symbolic links and exits with non-zero code
if no metadata found or at least one src could not be updated.

## Flags

`
