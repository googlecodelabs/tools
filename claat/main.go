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
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/googlecodelabs/tools/claat/cmd"

	// allow parsers to register themselves
	_ "github.com/googlecodelabs/tools/claat/parser/gdoc"
	_ "github.com/googlecodelabs/tools/claat/parser/md"
)

var (
	version string // set by linker -X

	// Flags.
	addr      = flag.String("addr", "localhost:9090", "hostname and port to bind web server to")
	authToken = flag.String("auth", "", "OAuth2 Bearer token; alternative credentials override.")
	expenv    = flag.String("e", "web", "codelab environment")
	extra     = flag.String("extra", "", "Additional arguments to pass to format templates. JSON object of string,string key values.")
	globalGA  = flag.String("ga", "UA-49880327-14", "global Google Analytics account")
	output    = flag.String("o", ".", "output directory or '-' for stdout")
	prefix    = flag.String("prefix", "../../", "URL prefix for html format")
	tmplout   = flag.String("f", "html", "output format")
)

func main() {
	log.SetFlags(0)
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) == 1 {
		log.Fatalf("Need subcommand. Try '-h' for options.")
	}
	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		usage()
		return
	}

	flag.Usage = usage
	flag.CommandLine.Parse(os.Args[2:])
	cmd.ExtraVars = cmd.ParseExtraVars(*extra)

	switch os.Args[1] {
	case "export":
		cmd.CmdExport(cmd.CmdExportOptions{
			AuthToken: *authToken,
			Expenv:    *expenv,
			GlobalGA:  *globalGA,
			Output:    *output,
			Prefix:    *prefix,
			Tmplout:   *tmplout,
		})
	case "serve":
		cmd.CmdServe(*addr)
	case "build":
		cmd.CmdBuild()
	case "update":
		if ok := cmd.CmdUpdate(cmd.CmdUpdateOptions{
			AuthToken: *authToken,
			GlobalGA:  *globalGA,
			Prefix:    *prefix,
		}); !ok {
			os.Exit(1)
		}
	case "help":
		usage()
	case "version":
		fmt.Println(version)
	default:
		log.Fatalf("Unknown subcommand. Try '-h' for options.")
	}
	os.Exit(cmd.Exit)
}

// usage prints usageText and program arguments to stderr.
func usage() {
	fmt.Fprint(os.Stderr, usageText)
	flag.PrintDefaults()
}

const usageText = `Usage: claat <cmd> [options] src [src ...]

Available commands are: export, serve, update, version.

## Export command

Export takes one or more 'src' documents and converts them
to the format specified with -f option.

The following formats are built-in:

- html (Polymer-based app)
- md (Markdown)
- offline (plain HTML markup for offline consumption)

Note that the built-in templates of the formats are not guaranteed to be stable.
They can be found in https://github.com/googlecodelabs/tools/tree/master/claat/render.
Please avoid using default templates in production. Use your own copies.

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

## Build command

Install all the dependencies needed by the codelab unless already installed.
This is done automatically by the serve command.
To clean up and rebuild all the dependencies, remove bower_components
directory and run the build command again

## Serve command

Serve provides a simple web server for viewing exported codelabs.
It takes no arguments and presents the current directory contents.
Clicking on a directory representing an exported codelab will load
all the required dependencies and render the generated codelab as
it would appear in production.

The serve command always runs the build command first.

The serve command takes a -addr host:port option, to specify the
desired hostname or IP address and port number to bind to.

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
