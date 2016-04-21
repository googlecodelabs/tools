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
	"bytes"
	"crypto/md5"
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
	for _ = range args {
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

	// rewritten image urls
	var imap map[string]string

	dir := *output // output dir or stdout
	if !isStdout(dir) {
		dir = codelabDir(dir, meta)
		imap = rewriteImages(clab.Steps)
	}
	// write codelab and its metadata to disk
	if err := writeCodelab(dir, clab.Codelab, ctx); err != nil {
		return nil, err
	}
	// slurp codelab assets to disk, if any
	mdir := filepath.Join(dir, imgDirname)
	return meta, downloadImages(client, mdir, imap)
}

// writeCodelab stores codelab main content in ctx.Format and its metadata
// in JSON format on disk.
func writeCodelab(dir string, clab *types.Codelab, ctx *types.Context) error {
	// render main codelab content to a tmp buffer,
	// which will also verify output format is valid,
	// and avoid creating empty files in case this goes wrong
	data := &render.Context{
		Env:      ctx.Env,
		Prefix:   ctx.Prefix,
		GlobalGA: ctx.MainGA,
		Meta:     &clab.Meta,
		Steps:    clab.Steps,
		Extra:    extraVars,
	}
	var buf bytes.Buffer
	if err := render.Execute(&buf, ctx.Format, data); err != nil {
		return err
	}
	// output to stdout does not include metadata
	w := os.Stdout
	if !isStdout(dir) {
		// make sure codelab dir exists
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		// codelab metadata
		cm := &types.ContextMeta{*ctx, clab.Meta}
		f := filepath.Join(dir, metaFilename)
		if err := writeMeta(f, cm); err != nil {
			return err
		}
		// main content file
		f = filepath.Join(dir, contentFile(ctx.Format))
		var err error
		if w, err = os.Create(f); err != nil {
			return err
		}
		defer w.Close()
	}
	_, err := w.Write(buf.Bytes())
	return err
}

// rewriteImages returns a mapping of local codelab asset file
// to its original URL.
// The local filename is MD5 hash of the original URL.
func rewriteImages(steps []*types.Step) map[string]string {
	var imap = make(map[string]string)
	for _, st := range steps {
		nodes := imageNodes(st.Content.Nodes)
		for _, n := range nodes {
			file := fmt.Sprintf("%x.png", md5.Sum([]byte(n.Src)))
			imap[file] = n.Src
			n.Src = filepath.Join(imgDirname, file)
		}
	}
	return imap
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
