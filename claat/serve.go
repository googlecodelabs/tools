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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
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
	err = fetchRepo(depsDir, "googlecodelabs/codelab-components")
	if err != nil {
		fatalf(err.Error())
	}
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
	if resp.StatusCode != 200 {
		return errors.New("Status code " + strconv.Itoa(resp.StatusCode) + " returned on http get")
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

// overrideSpec maps a given component name to desired spec, as described in fetchRepo
// doc comments.
var overrideSpec = map[string]string{
	"accessibility-developer-tools": "PolymerElements/accessibility-developer-tools#2.11.0",
	"async":                         "PolymerElements/async#1.5.2",
	"font-roboto":                   "PolymerElements/font-roboto#1.0.1",
	"google-apis":                   "PolymerElements/google-apis#1.1.7",
	"google-codelab-elements":       "GoogleCodelabs/google-codelab-elements#1.0.5",
	"iron-a11y-announcer":           "PolymerElements/iron-a11y-announcer#1.0.5",
	"iron-a11y-keys-behavior":       "PolymerElements/iron-a11y-keys-behavior#1.1.9",
	"iron-autogrow-textarea":        "PolymerElements/iron-autogrow-textarea#1.0.15",
	"iron-behaviors":                "PolymerElements/iron-behaviors#1.0.17",
	"iron-checked-element-behavior": "PolymerElements/iron-checked-element-behavior#1.0.5",
	"iron-collapse":                 "PolymerElements/iron-collapse#1.3.0",
	"iron-dropdown":                 "PolymerElements/iron-dropdown#1.5.5",
	"iron-fit-behavior":             "PolymerElements/iron-fit-behavior#1.2.6",
	"iron-flex-layout":              "PolymerElements/iron-flex-layout#1.3.2",
	"iron-form-element-behavior":    "PolymerElements/iron-form-element-behavior#1.0.6",
	"iron-icon":                     "PolymerElements/iron-icon#1.0.12",
	"iron-icons":                    "PolymerElements/iron-icons#1.2.0",
	"iron-iconset-svg":              "PolymerElements/iron-iconset-svg#1.1.0",
	"iron-input":                    "PolymerElements/iron-input#1.0.10",
	"iron-jsonp-library":            "PolymerElements/iron-jsonp-library#1.0.4",
	"iron-localstorage":             "PolymerElements/iron-localstorage#1.0.6",
	"iron-media-query":              "PolymerElements/iron-media-query#1.0.8",
	"iron-menu-behavior":            "PolymerElements/iron-menu-behavior#1.2.0",
	"iron-meta":                     "PolymerElements/iron-meta#1.1.2",
	"iron-overlay-behavior":         "PolymerElements/iron-overlay-behavior#1.10.3",
	"iron-pages":                    "PolymerElements/iron-pages#1.0.8",
	"iron-resizable-behavior":       "PolymerElements/iron-resizable-behavior#1.0.5",
	"iron-selector":                 "PolymerElements/iron-selector#1.5.2",
	"iron-validatable-behavior":     "PolymerElements/iron-validatable-behavior#1.1.1",
	"neon-animation":                "PolymerElements/neon-animation#1.2.4",
	"paper-behaviors":               "PolymerElements/paper-behaviors#1.0.12",
	"paper-button":                  "PolymerElements/paper-button#1.0.14",
	"paper-dialog-behavior":         "PolymerElements/paper-dialog-behavior#1.2.7",
	"paper-dialog":                  "PolymerElements/paper-dialog#1.1.0",
	"paper-drawer-panel":            "PolymerElements/paper-drawer-panel#1.0.11",
	"paper-dropdown-menu":           "PolymerElements/paper-dropdown-menu#1.5.0",
	"paper-fab":                     "PolymerElements/paper-fab#1.2.0",
	"paper-header-panel":            "PolymerElements/paper-header-panel#1.1.7",
	"paper-icon-button":             "PolymerElements/paper-icon-button#1.1.4",
	"paper-input":                   "PolymerElements/paper-input#1.1.23",
	"paper-item":                    "PolymerElements/paper-item#1.2.1",
	"paper-listbox":                 "PolymerElements/paper-listbox#1.1.2",
	"paper-material":                "PolymerElements/paper-material#1.0.6",
	"paper-menu-button":             "PolymerElements/paper-menu-button#1.5.2",
	"paper-menu":                    "PolymerElements/paper-menu#1.2.2",
	"paper-radio-button":            "PolymerElements/paper-radio-button#1.3.1",
	"paper-radio-group":             "PolymerElements/paper-radio-group#1.2.1",
	"paper-ripple":                  "PolymerElements/paper-ripple#1.0.9",
	"paper-scroll-header-panel":     "PolymerElements/paper-scroll-header-panel#1.0.16",
	"paper-styles":                  "PolymerElements/paper-styles#1.2.0",
	"paper-tabs":                    "PolymerElements/paper-tabs#1.8.0",
	"paper-toast":                   "PolymerElements/paper-toast#1.3.0",
	"paper-toolbar":                 "PolymerElements/paper-toolbar#1.1.7",
	"platinum-sw":                   "PolymerElements/platinum-sw#1.2.4",
	"polymer":                       "Polymer/polymer#1.4.0",
	"sinonjs":                       "PolymerElements/sinonjs#1.17.1",
	"stacky":                        "PolymerElements/stacky#1.3.2",
	"webcomponentsjs":               "WebComponents/webcomponentsjs#0.7.24",
}

// overrideComp maps a given spec, as described in fetchRepo doc comments to a
// desired component name.
var overrideComp = map[string]string{
	"googlecodelabs/codelab-components": "google-codelab-elements",
	"google/code-prettify":              "google-prettify",
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
	// Check exception map to see if we should override the component name.
	comp = repo
	if c, ok := overrideComp[path]; ok {
		comp = c
	}
	// if repo already exists locally, return immediately, we're done.
	if _, err := os.Stat(basedir + "/" + comp); err == nil {
		return nil
	}
	zipFile := basedir + "/" + comp + ".zip"
	url := "https://github.com/" + user + "/" + repo + "/archive/v" + ver + ".zip"
	err := downloadFile(zipFile, url)
	if err != nil {
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
	// Check exception map to see if we should overide spec.
	for c, s := range b.Dependencies {
		if spec, ok := overrideSpec[c]; ok {
			s = spec
		}
		err = fetchRepo(basedir, s)
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
