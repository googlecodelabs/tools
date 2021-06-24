// Copyright 2016-2019 Google LLC. All Rights Reserved.
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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/googlecodelabs/tools/claat/fetch"
	"github.com/googlecodelabs/tools/claat/types"
	"github.com/googlecodelabs/tools/claat/util"
)

// Options type to make the CmdUpdate signature succinct.
type CmdUpdateOptions struct {
	// AuthToken is the token to use for the Drive API.
	AuthToken string
	// ExtraVars is extra template variables.
	ExtraVars map[string]string
	// GlobalGA is the global Google Analytics account to use.
	GlobalGA string
	// PassMetadata are the extra metadata fields to pass along.
	PassMetadata map[string]bool
	// Prefix is a URL prefix to prepend when using HTML format.
	Prefix string
}

// CmdUpdate is the "claat update ..." subcommand.
// It returns a process exit code.
func CmdUpdate(opts CmdUpdateOptions) int {
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	dirs, err := scanPaths(roots)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if len(dirs) == 0 {
		log.Fatalf("no codelabs found in %s", strings.Join(roots, ", "))
	}

	type result struct {
		dir  string
		meta *types.Meta
		err  error
	}
	ch := make(chan *result, len(dirs))
	for _, d := range dirs {
		go func(d string) {
			// random sleep up to 1 sec
			// to reduce number of rate limit errors
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			meta, err := updateCodelab(d, opts)
			ch <- &result{d, meta, err}
		}(d)
	}

	var exitCode int
	for range dirs {
		res := <-ch
		if res.err != nil {
			exitCode = 1
			log.Printf(reportErr, res.dir, res.err)
		} else {
			log.Printf(reportOk, res.meta.ID)
		}
	}
	return exitCode
}

// updateCodelab reads metadata from a dir/codelab.json file,
// re-exports the codelab just like it normally would in exportCodelab,
// and removes assets (images) which are not longer in use.
func updateCodelab(dir string, opts CmdUpdateOptions) (*types.Meta, error) {
	// get stored codelab metadata and fail early if we can't
	meta, err := readMeta(filepath.Join(dir, metaFilename))
	if err != nil {
		return nil, err
	}
	// override allowed options from cli
	if opts.Prefix != "" {
		meta.Prefix = opts.Prefix
	}
	if opts.GlobalGA != "" {
		meta.MainGA = opts.GlobalGA
	}

	// fetch and parse codelab source
	f, err := fetch.NewFetcher(opts.AuthToken, opts.PassMetadata, nil)
	if err != nil {
		return nil, err
	}
	basedir := filepath.Join(dir, "..")
	clab, err := f.SlurpCodelab(meta.Source, basedir)
	if err != nil {
		return nil, err
	}
	updated := types.ContextTime(clab.Mod)
	meta.Context.Updated = &updated

	newdir := codelabDir(basedir, &clab.Meta)
	imgdir := filepath.Join(newdir, util.ImgDirname)

	// write codelab and its metadata
	if err := writeCodelab(newdir, clab.Codelab, opts.ExtraVars, &meta.Context); err != nil {
		return nil, err
	}

	// cleanup:
	// - remove original dir if codelab ID has changed and so has the output dir
	// - otherwise, remove images which are not in imgs
	old := codelabDir(basedir, &meta.Meta)
	if old != newdir {
		return &meta.Meta, os.RemoveAll(old)
	}
	visit := func(p string, fi os.FileInfo, err error) error {
		if err != nil || p == imgdir {
			return err
		}
		if fi.IsDir() {
			return filepath.SkipDir
		}
		if _, ok := clab.Imgs[filepath.Base(p)]; !ok {
			return os.Remove(p)
		}
		return nil
	}
	return &meta.Meta, filepath.Walk(imgdir, visit)
}

// scanPaths looks for codelab metadata files in roots, recursively.
// The roots argument can contain overlapping directories as the return
// value is always de-duped.
func scanPaths(roots []string) ([]string, error) {
	type result struct {
		root string
		dirs []string
		err  error
	}
	ch := make(chan *result, len(roots))
	for _, r := range roots {
		go func(r string) {
			dirs, err := walkPath(r)
			ch <- &result{r, dirs, err}
		}(r)
	}
	var dirs []string
	for range roots {
		res := <-ch
		if res.err != nil {
			return nil, fmt.Errorf("%s: %v", res.root, res.err)
		}
		dirs = append(dirs, res.dirs...)
	}
	return util.Unique(dirs), nil
}

// walkPath walks root dir recursively, looking for metaFilename files.
func walkPath(root string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(root, func(p string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return err
		}
		if filepath.Base(p) == metaFilename {
			dirs = append(dirs, filepath.Dir(p))
		}
		return nil
	})
	return dirs, err
}

// readMeta reads codelab metadata from file.
// It will convert legacy fields to the actual.
func readMeta(file string) (*types.ContextMeta, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var cm types.ContextMeta
	if err := json.Unmarshal(b, &cm); err != nil {
		return nil, err
	}
	if cm.Format == "" {
		cm.Format = "html"
	}
	return &cm, nil
}
