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

package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// cmdServe is the "claat serve ..." subcommand.
func cmdServe() {
	var depsDir = "bower_components"
	var elemDir = "elements"
	var elemFile = elemDir + "/codelab.html"
	const codelabElem = `
<link rel="import" href="../bower_components/google-codelab-elements/google-codelab-elements.html">
`
	err := os.MkdirAll(depsDir, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Go get the dependencies.
	err = fetchRepo(depsDir, "googlecodelabs/codelab-components#^1.0.0")
	if err != nil {
		fmt.Println(err)
		return
	}
	//os.Rename(depsDir+"/codelab-components", depsDir+"/google-codelab-elements")
	//os.Rename(depsDir+"/code-prettify", depsDir+"/google-prettify")
	err = os.MkdirAll(elemDir, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err := os.Stat(elemFile); os.IsNotExist(err) {
		f, err := os.Create(elemFile)
		if err != nil {
			log.Fatal("Cannot create file", err)
		}
		defer f.Close()
		fmt.Fprintf(f, codelabElem)
	}
	fmt.Println("Dependencies installed.")

	port := "9090"
	http.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println("Serving on localhost:" + port + ", opening browser tab now...")
	openBrowser("http://127.0.0.1:" + port)
	err = http.ListenAndServe(":"+port, nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// downloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func fetchRepo(depsDir string, repoVers string) error {
	// parse url into repoName and version
	s := strings.Split(repoVers, "#^")
	vers := ""
	if len(s) > 1 {
		vers = s[1]
	}
	s = strings.Split(s[0], "/")
	user := s[0]
	repo := s[1]

	// if repo already exists locally, return immediately, we're done.
	if _, err := os.Stat(depsDir + "/" + repo); os.IsExist(err) {
		return nil
	}
	url := "https://github.com/" + user + "/" + repo + "/archive/" + vers + ".zip"
	zipFile := depsDir + "/" + repo + ".zip"
	fmt.Println("Downloading " + url)
	err := downloadFile(zipFile, url)
	if err != nil {
		return err
	}
	err = unzip(zipFile, depsDir)
	if err != nil {
		return err
	}
	os.Remove(zipFile)
	os.Rename(depsDir+"/"+repo+"-"+vers, depsDir+"/"+repo)

	// if unzipped archive contains a bower.json, parse it, and for each dependency therein,
	// recursively fetch the corresponding repo.
	bowerFile := depsDir + "/" + repo + "/bower.json"
	if _, err := os.Stat(bowerFile); os.IsNotExist(err) {
		return nil
	}
	raw, err := ioutil.ReadFile(bowerFile)
	if err != nil {
		return err
	}
	type Dep struct {
		Name     string
		RepoVers string
	}
	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(raw, &objmap)
	if err != nil {
		return err
	}
	var deps map[string]string
	err = json.Unmarshal(*objmap["dependencies"], &deps)
	if err != nil {
		return err
	}
	for _, v := range deps {
		err = fetchRepo(depsDir, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// unzip a file to a dest directory.
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	os.MkdirAll(dest, 0755)
	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			return os.MkdirAll(path, f.Mode())
		}
		os.MkdirAll(filepath.Dir(path), f.Mode())
		f2, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(f2, rc)
		// Catch both io.Copy and file.Close errors but prefer returning the former.
		if err1 := f2.Close(); err1 != nil && err == nil {
			err = err1
		}
		return err
	}
	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}
	return nil
}

// openBrowser tries to open the URL in a browser,
// and returns whether it succeed in doing so.
func openBrowser(url string) bool {
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
	return cmd.Start() == nil
}
