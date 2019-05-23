# Codelabs command line tool

The program takes an input in form of a resource location,
which can either be a Google Doc ID, local file path or an arbitrary URL.
It then converts the input into a codelab format, HTML by default.

For more info run `claat help`.

## Installation

The easiest way is to download pre-compiled binary.
The binaries, as well as their checksums are available at the
[Releases page](https://github.com/googlecodelabs/tools/releases/latest).

Alternatively, follow the instructions below:

**Prerequisites**
- Install [Go](https://golang.org/dl/) if you don't already have it.
- Install `protoc`. Follow instructions [here](http://google.github.io/proto-lens/installing-protoc.html) (MacOS and Linux), or grab the latest `protoc` executable binary (`protoc-<ver>-<os>-<arch>.zip/bin/protoc` under `Assets`) for your system [here](https://github.com/protocolbuffers/protobuf/releases) and place it under your `PATH`.
- Install Go's protobuf runtime
  - `go get -u github.com/golang/protobuf/protoc-gen-go`

**Tool installation**
- Download this repo and place it under your `$GOPATH`/`go env GOPATH`
  - `go get github.com/googlecodelabs/tools/claat`
- Install package dependencies
  - `cd $(go env GOPATH)/src/github.com/googlecodelabs/tools && go get ./...`
- Build the binary
  - `cd $(go env GOPATH)/src/github.com/googlecodelabs/tools/claat && make`
- `claat` should now be accessible directly from command line.


## Dev workflow

To build the binary, run `make`.

Testing is done with `make test` or `go test ./...` if preferred.

Don't forget to run `make lint` or `golint ./...` before creating a new CL.

To create cross-compiled versions for all supported OS/Arch, run `make release`.
It will place the output in `bin/claat-<os>-<arch>`.
