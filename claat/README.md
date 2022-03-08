# Codelabs command line tool

The program takes an input in form of a resource location,
which can either be a Google Doc ID, local file path or an arbitrary URL.
It then converts the input into a codelab format, HTML by default.

For more info run `claat help`.

## Install

The easiest way is to download pre-compiled binary.
The binaries, as well as their checksums are available at the
[Releases page](https://github.com/googlecodelabs/tools/releases/latest).

Alternatively, if you have [Go installed](https://golang.org/doc/install):

    go install github.com/googlecodelabs/tools/claat@latest

If none of the above works, compile the tool from source following Dev workflow
instructions below.

## Dev workflow

**Prerequisites**

1. Install [Go](https://golang.org/dl/) if you don't have it.
2. Make sure this directory is placed under
   `$GOPATH/src/github.com/googlecodelabs/tools`.
3. Install package dependencies with `go get ./...` from this directory.

To build the binary, run `make`.

Testing is done with `make test` or `go test ./...` if preferred.

Don't forget to run `make lint` or `golint ./...` before creating a new CL.

To create cross-compiled versions for all supported OS/Arch, run `make release`.
It will place the output in `bin/claat-<os>-<arch>`.
