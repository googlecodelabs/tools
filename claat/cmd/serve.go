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

package cmd

import (
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

// CmdServe is the "claat serve ..." subcommand.
// addr is the hostname and port to bind the web server to.
// It returns a process exit code.
func CmdServe(addr string) int {
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Printf("Serving codelabs on %s, opening browser tab now...", addr)
	ch := make(chan error, 1)
	go func() {
		ch <- http.ListenAndServe(addr, nil)
	}()
	openBrowser("http://" + addr)
	log.Fatalf("claat serve: %v", <-ch)
	return 0
}

// openBrowser tries to open the URL in a browser.
func openBrowser(url string) error {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start()
}
