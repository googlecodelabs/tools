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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

// CmdBuild is the "claat build ..." subcommand.
// It returns a process exit code.
func CmdBuild() int {
	const depsDir = "bower_components"
	var codelabElem = []byte(`
<link rel="import" href="../bower_components/google-codelab-elements/google-codelab-elements.html">
`)

	err := os.MkdirAll(depsDir, 0755)
	if err != nil {
		log.Fatalf("%s: %v", depsDir, err)
	}
	// Go get the dependencies.
	if err := fetchRepo(depsDir, "google-codelab-elements", "googlecodelabs/codelab-components#2.0.2"); err != nil {
		log.Fatalf(err.Error())
	}
	if err := writeFile(filepath.Join("elements", "codelab.html"), codelabElem); err != nil {
		log.Fatalf(err.Error())
	}
	return 0
}

func writeFile(name string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(name, content, 0644)
}

// warnOnce makes sure we notify users about fetching dependencies only once.
var warnOnce sync.Once

// fetchRepo downloads a repo from github and unpacks it into basedir/name,
// calling itself recursively for every dependency found in the unpacked bower.json.
// If the basedir/name directory already exists, fetchRepo does nothing.
//
// The spec has the following format: "githubuser/repo#version".
// Any non-alphanumeric characters are stripped from the version.
// If version is missing, "master" is used instead, unless overridden in bowerVersionOverride.
// Instead of bower registry a local bowerSpecResolve map is used for naming resolution.
func fetchRepo(basedir, name, spec string) error {
	outdir := filepath.Join(basedir, name)
	if _, err := os.Stat(outdir); err == nil {
		return nil
	}

	warnOnce.Do(func() {
		log.Println("Fetching dependencies...")
	})
	tryURL, err := fromBowerSpec(name, spec)
	if err != nil {
		return err
	}
	zipFile := filepath.Join(basedir, name+".zip")
	var ok bool
	for _, u := range tryURL {
		if err := downloadFile(zipFile, u); err == nil {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("fetchRepo(%q, %q): could not find in any of %q", name, spec, tryURL)
	}
	if err := stripUnzip(outdir, zipFile); err != nil {
		return err
	}
	os.Remove(zipFile)

	// If unzipped archive contains a bower.json, parse it, and for each dependency therein,
	// recursively fetch the corresponding repo.
	bowerFile := filepath.Join(outdir, "bower.json")
	raw, err := ioutil.ReadFile(bowerFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var bower struct {
		Dependencies map[string]string
	}
	if err := json.Unmarshal(raw, &bower); err != nil {
		return fmt.Errorf("%s: %v", bowerFile, err)
	}

	for name, spec := range bower.Dependencies {
		if err := fetchRepo(basedir, name, spec); err != nil {
			return err
		}
	}
	return nil
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

// bowerSpecResolve resolves a bower component name to its github user/repo.
var bowerSpecResolve = map[string]string{
	"webcomponentsjs": "webcomponents/webcomponentsjs",
}

// bowerVersionOverride maps a githb user/repo to a particular fixed version to fetch
// when resolving and fetching dependencies.
var bowerVersionOverride = map[string]string{
	"polymer/polymer": "standard-layer",
	"polymerelements/accessibility-developer-tools": "2.11.0",
	"polymerelements/async":                         "1.5.2",
	"polymerelements/font-roboto":                   "1.0.1",
	"polymerelements/google-apis":                   "1.1.7",
	"polymerelements/iron-a11y-announcer":           "1.0.5",
	"polymerelements/iron-a11y-keys-behavior":       "1.1.9",
	"polymerelements/iron-autogrow-textarea":        "1.0.15",
	"polymerelements/iron-behaviors":                "1.0.17",
	"polymerelements/iron-checked-element-behavior": "1.0.5",
	"polymerelements/iron-collapse":                 "1.3.0",
	"polymerelements/iron-dropdown":                 "1.5.5",
	"polymerelements/iron-fit-behavior":             "1.2.6",
	"polymerelements/iron-flex-layout":              "1.3.2",
	"polymerelements/iron-form-element-behavior":    "1.0.6",
	"polymerelements/iron-icon":                     "1.0.12",
	"polymerelements/iron-icons":                    "1.2.0",
	"polymerelements/iron-iconset-svg":              "1.1.0",
	"polymerelements/iron-input":                    "1.0.10",
	"polymerelements/iron-jsonp-library":            "1.0.4",
	"polymerelements/iron-localstorage":             "1.0.6",
	"polymerelements/iron-media-query":              "1.0.8",
	"polymerelements/iron-menu-behavior":            "1.2.0",
	"polymerelements/iron-meta":                     "1.1.2",
	"polymerelements/iron-overlay-behavior":         "1.10.3",
	"polymerelements/iron-pages":                    "1.0.8",
	"polymerelements/iron-resizable-behavior":       "1.0.5",
	"polymerelements/iron-selector":                 "1.5.2",
	"polymerelements/iron-validatable-behavior":     "1.1.1",
	"polymerelements/neon-animation":                "1.2.4",
	"polymerelements/paper-behaviors":               "1.0.12",
	"polymerelements/paper-button":                  "1.0.14",
	"polymerelements/paper-dialog":                  "1.1.0",
	"polymerelements/paper-dialog-behavior":         "1.2.7",
	"polymerelements/paper-drawer-panel":            "1.0.11",
	"polymerelements/paper-dropdown-menu":           "1.5.0",
	"polymerelements/paper-fab":                     "1.2.0",
	"polymerelements/paper-header-panel":            "1.1.7",
	"polymerelements/paper-icon-button":             "1.1.4",
	"polymerelements/paper-input":                   "1.1.23",
	"polymerelements/paper-item":                    "1.2.1",
	"polymerelements/paper-listbox":                 "1.1.2",
	"polymerelements/paper-material":                "1.0.6",
	"polymerelements/paper-menu":                    "1.2.2",
	"polymerelements/paper-menu-button":             "1.5.2",
	"polymerelements/paper-radio-button":            "1.3.1",
	"polymerelements/paper-radio-group":             "1.2.1",
	"polymerelements/paper-ripple":                  "1.0.9",
	"polymerelements/paper-scroll-header-panel":     "1.0.16",
	"polymerelements/paper-styles":                  "1.2.0",
	"polymerelements/paper-tabs":                    "1.8.0",
	"polymerelements/paper-toast":                   "1.3.0",
	"polymerelements/paper-toolbar":                 "1.1.7",
	"polymerelements/platinum-sw":                   "1.2.4",
	"polymerelements/sinonjs":                       "1.17.1",
	"polymerelements/stacky":                        "1.3.2",
}

// fromBowerSpec returns possible locations of the bower package identified by the name/spec pair.
func fromBowerSpec(name, spec string) (urls []string, err error) {
	if strings.IndexByte(spec, '/') == -1 {
		s, ok := bowerSpecResolve[strings.ToLower(name)]
		if !ok {
			return nil, fmt.Errorf("unable to resolve %q", name)
		}
		spec = s + "#" + spec
	}
	parts := strings.Split(spec, "#") // {repo, version}
	if len(parts) > 2 {
		return nil, fmt.Errorf("invalid spec for %s: %q", name, spec)
	}
	repo := path.Clean(parts[0])
	ver := bowerVersionOverride[strings.ToLower(repo)]
	if ver == "" && len(parts) > 1 {
		ver = parts[1]
	}
	ver = strings.Trim(ver, "^~ ")
	if ver == "" {
		ver = "master"
	}
	u := []string{"https://github.com" + path.Join("/", repo, "archive", ver) + ".zip"}
	if ver != "master" && !strings.HasPrefix(ver, "v") {
		u = append(u, "https://github.com"+path.Join("/", repo, "archive", "v"+ver)+".zip")
	}
	return u, nil
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
