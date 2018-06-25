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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/googlecodelabs/tools/claat/render"
	"github.com/googlecodelabs/tools/claat/types"
)

// cmdExport is the "claat export ..." subcommand.
func cmdExport() {
	if flag.NArg() == 0 {
		fatalf("Need at least one source. Try '-h' for options.")
	}
	type result struct {
		src  string
		meta *types.Meta
		err  error
	}
	args := unique(flag.Args())
	ch := make(chan *result, len(args))
	for _, src := range args {
		go func(src string) {
			meta, err := exportCodelab(src)
			ch <- &result{src, meta, err}
		}(src)
	}
	for range args {
		res := <-ch
		if res.err != nil {
			errorf(reportErr, res.src, res.err)
		} else if !isStdout(*output) {
			printf(reportOk, res.meta.ID)
		}
	}
}

// exportCodelab fetches codelab src from either local disk or remote,
// parses and stores the results on disk, in a dir ancestored by *output.
//
// Stored results include codelab content formatted in *tmplout, its assets
// and metadata in JSON format.
//
// There's a special case where basedir has a value of "-", in which
// nothing is stored on disk and the only output, codelab formatted content,
// is printed to stdout.
func exportCodelab(src string) (*types.Meta, error) {
	clab, err := slurpCodelab(src)
	if err != nil {
		return nil, err
	}
	var client *http.Client // need for downloadImages
	if clab.typ == srcGoogleDoc {
		client, err = driveClient()
		if err != nil {
			return nil, err
		}
	}

	// codelab export context
	lastmod := types.ContextTime(clab.mod)
	meta := &clab.Meta
	ctx := &types.Context{
		Source:  src,
		Env:     *expenv,
		Format:  *tmplout,
		Prefix:  *prefix,
		MainGA:  *globalGA,
		Updated: &lastmod,
	}

	dir := *output // output dir or stdout
	if !isStdout(dir) {
		dir = codelabDir(dir, meta)
		// download or copy codelab assets to disk, and rewrite image URLs
		mdir := filepath.Join(dir, imgDirname)
		if _, err := slurpImages(client, src, mdir, clab.Steps); err != nil {
			return nil, err
		}
	}
	// write codelab and its metadata to disk
	return meta, writeCodelab(dir, clab.Codelab, ctx)
}

// writeCodelab stores codelab main content in ctx.Format and its metadata
// in JSON format on disk.
func writeCodelab(dir string, clab *types.Codelab, ctx *types.Context) error {
	// output to stdout does not include metadata
	if !isStdout(dir) {
		// make sure codelab dir exists
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		// codelab metadata
		cm := &types.ContextMeta{Context: *ctx, Meta: clab.Meta}
		f := filepath.Join(dir, metaFilename)
		if err := writeMeta(f, cm); err != nil {
			return err
		}
	}

	// main content file(s)
	data := &struct {
		render.Context
		Current *types.Step
		StepNum int
		Prev    bool
		Next    bool
	}{Context: render.Context{
		Env:      ctx.Env,
		Prefix:   ctx.Prefix,
		GlobalGA: ctx.MainGA,
		Meta:     &clab.Meta,
		Steps:    clab.Steps,
		Extra:    extraVars,
	}}
	if ctx.Format != "offline" {
		w := os.Stdout
		if !isStdout(dir) {
			f, err := os.Create(filepath.Join(dir, "index."+ctx.Format))
			if err != nil {
				return err
			}
			w = f
			defer f.Close()
		}
		return render.Execute(w, ctx.Format, data)
	}
	for i, step := range clab.Steps {
		data.Current = step
		data.StepNum = i + 1
		data.Prev = i > 0
		data.Next = i < len(clab.Steps)-1
		w := os.Stdout
		if !isStdout(dir) {
			name := "index.html"
			if i > 0 {
				name = fmt.Sprintf("step-%d.html", i+1)
			}
			f, err := os.Create(filepath.Join(dir, name))
			if err != nil {
				return err
			}
			w = f
			defer f.Close()
		}
		if err := render.Execute(w, ctx.Format, data); err != nil {
			return err
		}
	}
	return nil
}

func slurpImages(client *http.Client, src, dir string, steps []*types.Step) (map[string]string, error) {
	// make sure img dir exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	type res struct {
		url, file string
		err       error
	}

	ch := make(chan *res, 100)
	defer close(ch)
	var count int
	for _, st := range steps {
		nodes := imageNodes(st.Content.Nodes)
		count += len(nodes)
		for _, n := range nodes {
			go func(n *types.ImageNode) {
				url := n.Src
				file, err := slurpBytes(client, src, dir, url)
				if err == nil {
					n.Src = filepath.Join(imgDirname, file)
				}
				ch <- &res{url, file, err}
			}(n)
		}
	}

	var err error
	imap := make(map[string]string, count)
	for i := 0; i < count; i++ {
		r := <-ch
		imap[r.file] = r.url
		if r.err != nil && err == nil {
			// record first error
			err = fmt.Errorf("%s => %s: %v", r.url, r.file, r.err)
		}
	}

	return imap, err
}

// imageNodes filters out everything except types.NodeImage nodes, recursively.
func imageNodes(nodes []types.Node) []*types.ImageNode {
	var imgs []*types.ImageNode
	for _, n := range nodes {
		switch n := n.(type) {
		case *types.ImageNode:
			imgs = append(imgs, n)
		case *types.ListNode:
			imgs = append(imgs, imageNodes(n.Nodes)...)
		case *types.ItemsListNode:
			for _, i := range n.Items {
				imgs = append(imgs, imageNodes(i.Nodes)...)
			}
		case *types.HeaderNode:
			imgs = append(imgs, imageNodes(n.Content.Nodes)...)
		case *types.URLNode:
			imgs = append(imgs, imageNodes(n.Content.Nodes)...)
		case *types.ButtonNode:
			imgs = append(imgs, imageNodes(n.Content.Nodes)...)
		case *types.InfoboxNode:
			imgs = append(imgs, imageNodes(n.Content.Nodes)...)
		case *types.GridNode:
			for _, r := range n.Rows {
				for _, c := range r {
					imgs = append(imgs, imageNodes(c.Content.Nodes)...)
				}
			}
		}
	}
	return imgs
}

// importNodes filters out everything except types.NodeImport nodes, recursively.
func importNodes(nodes []types.Node) []*types.ImportNode {
	var imps []*types.ImportNode
	for _, n := range nodes {
		switch n := n.(type) {
		case *types.ImportNode:
			imps = append(imps, n)
		case *types.ListNode:
			imps = append(imps, importNodes(n.Nodes)...)
		case *types.InfoboxNode:
			imps = append(imps, importNodes(n.Content.Nodes)...)
		case *types.GridNode:
			for _, r := range n.Rows {
				for _, c := range r {
					imps = append(imps, importNodes(c.Content.Nodes)...)
				}
			}
		}
	}
	return imps
}

// writeMeta writes codelab metadata to a local disk location
// specified by path.
func writeMeta(path string, cm *types.ContextMeta) error {
	b, err := json.MarshalIndent(cm, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return ioutil.WriteFile(path, b, 0644)
}

// codelabDir returns codelab root directory.
// The base argument is codelab parent directory.
func codelabDir(base string, m *types.Meta) string {
	return filepath.Join(base, m.ID)
}

// unique de-dupes a.
// The argument a is not modified.
func unique(a []string) []string {
	seen := make(map[string]struct{}, len(a))
	res := make([]string, 0, len(a))
	for _, s := range a {
		if _, y := seen[s]; !y {
			res = append(res, s)
			seen[s] = struct{}{}
		}
	}
	return res
}
