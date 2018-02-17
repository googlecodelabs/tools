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
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Global for tracking whether we've warned user about delay.
var warned = false

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
		fatalf(err.Error())
	}
	// Go get the dependencies.
	err = fetchRepo(depsDir, "googlecodelabs/codelab-components#1.0.5")
	if err != nil {
		fatalf(err.Error())
	}
	os.Rename(depsDir+"/code-prettify", depsDir+"/google-prettify")
	err = os.MkdirAll(elemDir, 0755)
	if err != nil {
		fatalf(err.Error())
	}
	if _, err := os.Stat(elemFile); os.IsNotExist(err) {
		f, err := os.Create(elemFile)
		if err != nil {
			fatalf(err.Error())
		}
		defer f.Close()
		f.WriteString(codelabElem)
	}
	http.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println("Serving on " + *addr + ", opening browser tab now...")
	ch := make(chan error, 1)
	go func() {
		ch <- http.ListenAndServe(*addr, nil)
	}()
	openBrowser("http://" + *addr)
	fatalf("claat serve: %v", <-ch)
}

// downloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func downloadFile(filepath string, url string) error {
	if !warned {
		fmt.Println("Fetching dependencies...")
		warned = true
	}
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
	return err
}

// bowerComp maps non-conforming spec as described in fetchRepo doc comments.
// It is keyed by the name found in a bower.json dependency list,
// with values corresponding to a valid githubuser/repo spec.
var bowerComp = map[string]string{
	"googlecodelabs/codelab-components": "google-codelab-elements",
	"webcomponents/webcomponentsjs":     "webcomponentsjs",
	"google/google-prettify":            "code-prettify",
}

// fetchRepo downloads a repo from github and unpacks it into basedir/dest,
// calling itself recursively for every dependency found in the unpacked bower.json.
// If the basedir/dest directory already exists, fetchRepo does nothing.
//
// The spec has the following format: "githubuser/repo#version".
// Any non-alphanumeric characters are stripped from the version.
// If version is missing, "master" is used instead.
// The bower registry is not used for naming resolution.
func fetchRepo(basedir, spec string) error {
	var user, repo, path, comp, ver string
	s := strings.Split(spec, "#")
	path = s[0]
	if len(s) > 1 {
		ver = s[1]
	}
	ver = strings.Trim(ver, "^~ ")
	if ver == "" {
		ver = "master"
	}
	s = strings.Split(path, "/")
	user = s[0]
	if len(s) > 1 {
		repo = s[1]
	}
	// Check exception map to see if we should use a special component name.
	comp = repo
	if v, ok := bowerComp[path]; ok {
		comp = v
	}
	// if repo already exists locally, return immediately, we're done.
	if _, err := os.Stat(basedir + "/" + comp); err == nil {
		return nil
	}
	zipFile := basedir + "/" + comp + ".zip"
	url := "https://github.com/" + user + "/" + repo + "/archive/v" + ver + ".zip"
	err := downloadFile(zipFile, url)
	if err != nil {
		return err
	}
	// If get fails, it will download a file containing only "404: Not Found".
	// We check for that case by looking for an unusually small file.
	var st os.FileInfo
	if st, err = os.Stat(zipFile); err != nil {
		return err
	}
	if st.Size() < 20 {
		os.Remove(zipFile)
		url = "https://github.com/" + user + "/" + repo + "/archive/" + ver + ".zip"
		err = downloadFile(zipFile, url)
		if err != nil {
			return err
		}
	}
	err = unzip(zipFile, basedir)
	if err != nil {
		return err
	}
	os.Remove(zipFile)
	os.Rename(basedir+"/"+repo+"-"+ver, basedir+"/"+comp)

	// if unzipped archive contains a bower.json, parse it, and for each dependency therein,
	// recursively fetch the corresponding repo.
	bowerFile := basedir + "/" + comp + "/bower.json"
	if _, err := os.Stat(bowerFile); os.IsNotExist(err) {
		return nil
	}
	raw, err := ioutil.ReadFile(bowerFile)
	if err != nil {
		return err
	}
	var b struct {
		Dependencies map[string]string
	}
	err = json.Unmarshal(raw, &b)
	if err != nil {
		return err
	}
	for _, v := range b.Dependencies {
		err = fetchRepo(basedir, v)
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
