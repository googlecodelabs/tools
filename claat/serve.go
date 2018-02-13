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
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var deps = map[string]string{
	"webcomponentsjs":     "webcomponents/webcomponentsjs",
	"codelab-components":  "googlecodelabs/codelab-components#^1.0.0",
	"google-youtube":      "GoogleWebComponents/google-youtube#^1.2.1",
	"iron-flex-layout":    "PolymerElements/iron-flex-layout#^1.0.3",
	"iron-pages":          "PolymerElements/iron-pages#^1.0.4",
	"iron-selector":       "PolymerElements/iron-selector#^1.0.8",
	"iron-ajax":           "PolymerElements/iron-ajax#^1.0.8",
	"paper-button":        "PolymerElements/paper-button#^1.0.3",
	"paper-drawer-panel":  "PolymerElements/paper-drawer-panel#^1.0.3",
	"paper-header-panel":  "PolymerElements/paper-header-panel#^1.0.4",
	"paper-icon-button":   "PolymerElements/paper-icon-button#^1.0.3",
	"paper-input":         "PolymerElements/paper-input#^1.1.1",
	"paper-item":          "PolymerElements/paper-item#^1.0.2",
	"paper-listbox":       "PolymerElements/paper-listbox#^1.0.0",
	"paper-menu-button":   "PolymerElements/paper-menu-button#^1.0.3",
	"paper-dropdown-menu": "PolymerElements/paper-dropdown-menu#^1.0.4",
	"paper-styles":        "PolymerElements/paper-styles#^1.0.11",
	"paper-tabs":          "PolymerElements/paper-tabs#^1.1.0",
	"paper-toast":         "PolymerElements/paper-toast#^1.0.0",
	"paper-toolbar":       "PolymerElements/paper-toolbar#^1.0.4",
	"platinum-sw":         "PolymerElements/platinum-sw#^1.2.0",
	"polymer":             "Polymer/polymer#standard-layer",
	"url-search-params":   "WebReflection/url-search-params#^0.5.0",
}

var depsDir = "bower_components"
var elemDir = "elements"
var elemFile = elemDir + "/codelab.html"

const codelabElem = `
<link rel="import" href="../bower_components/google-codelab-elements/google-codelab-elements.html">
<link rel="import" href="../bower_components/google-youtube/google-youtube.html">
<link rel="import" href="../bower_components/iron-ajax/iron-ajax.html">
`

// cmdServe is the "claat serve ..." subcommand.
func cmdServe() {
	err := mkdirIfNec(depsDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	// If we don't have the required dependencies, go get 'em.
	for k, v := range deps {
		if _, err := os.Stat(depsDir + "/" + k); os.IsNotExist(err) {
			s := strings.Split(v, "#")
			url := "https://github.com/" + s[0] + "/archive/master.zip"
			archFile := depsDir + "/master.zip"
			destDir := depsDir + "/" + k
			fmt.Println("Downloading " + url)
			err := DownloadFile(archFile, url)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = Unzip(archFile, depsDir)
			os.Rename(destDir+"-master", destDir)
			os.Remove(archFile)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
	err = mkdirIfNec(elemDir)
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

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
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

// Unzip a file to a dest directory.
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()
	os.MkdirAll(dest, 0755)
	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
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

// Make a directory if it doesn't already exist.
func mkdirIfNec(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
