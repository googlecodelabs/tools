// Copyright 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The server command starts a simple static file server using current work dir
// as the root directory.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var addr = flag.String("addr", "localhost:8080", "Server address to bind to.")

func main() {
	flag.Parse()
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	h := http.FileServer(http.Dir(dir))
	log.Printf("serving from %q on http://%s", dir, *addr)
	log.Fatal(http.ListenAndServe(*addr, h))
}
