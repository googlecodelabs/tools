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
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func CmdBuild() {
	const depsDir = "deps"

	err := os.MkdirAll(depsDir, 0755)
	if err != nil {
		log.Fatalf("%s: %v", depsDir, err)
	}
	// Go get the dependencies.
	const filename = "codelab-elements"
	zippath := filepath.Join(depsDir, filename + ".zip")
	const bundleURL = "https://github.com/shawnbuso/codelab-elements/releases/download/0.1/bundle.zip"
	err = downloadFile(zippath, bundleURL)
	if err != nil {
		log.Fatalf("Error downloading deps: %v", err)
	}
	stripUnzip(filepath.Join(depsDir, filename), zippath)
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
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s: %s", url, resp.Status)
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

// stripUnzip unpacks src file to the dest directory, stripping parent dir from the zip file paths.
func stripUnzip(dest, src string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}
	extractAndWriteFile := func(f *zip.File) error {
		zpath := filepath.FromSlash(f.Name)
		if i := strings.IndexByte(zpath, filepath.Separator); i > 0 {
			zpath = zpath[i+1:]
		}
		zpath = filepath.Join(dest, zpath)
		if f.FileInfo().IsDir() {
			return os.MkdirAll(zpath, f.Mode())
		}

		zf, err := f.Open()
		if err != nil {
			return err
		}
		defer zf.Close()
		if err := os.MkdirAll(filepath.Dir(zpath), f.Mode()); err != nil {
			return err
		}
		w, err := os.OpenFile(zpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		// Catch both io.Copy and w.Close errors but prefer returning the former.
		_, err = io.Copy(w, zf)
		if err1 := w.Close(); err1 != nil && err == nil {
			err = err1
		}
		return err
	}
	for _, f := range r.File {
		if err := extractAndWriteFile(f); err != nil {
			return err
		}
	}
	return nil
}
