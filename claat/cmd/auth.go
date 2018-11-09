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

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const (
	// auth scopes needed by the program
	scopeDriveReadOnly = "https://www.googleapis.com/auth/drive.readonly"

	// program credentials for installed apps
	googClient = "183908478743-e8rth9fbo7juk9eeivgp23asnt791g63.apps.googleusercontent.com"
	googSecret = "ljELuf5jUrzcOxZGL7OQfkIC"

	// token providers
	providerGoogle = "goog"
)

var (
	// OAuth2 configs for OOB flow
	authConfig = map[string]*oauth2.Config{
		providerGoogle: {
			ClientID:     googClient,
			ClientSecret: googSecret,
			Scopes:       []string{scopeDriveReadOnly},
			RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/auth",
				TokenURL: "https://accounts.google.com/o/oauth2/token",
			},
		},
	}

	// reusable HTTP clients
	clientsMu sync.Mutex // guards clients
	clients   map[string]*http.Client
)

func init() {
	clients = make(map[string]*http.Client)
}

// driveClient returns an HTTP client which knows how to perform authenticated
// requests to Google Drive API.
func driveClient(authToken string) (*http.Client, error) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	if hc, ok := clients[providerGoogle]; ok {
		return hc, nil
	}
	ts, err := tokenSource(providerGoogle, authToken)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Transport{
		Source: ts,
		Base:   http.DefaultTransport,
	}
	hc := &http.Client{Transport: t}
	clients[providerGoogle] = hc
	return hc, nil
}

// tokenSource creates a new oauth2.TokenSource backed by tokenRefresher,
// using previously stored user credentials if available.
// If authToken is given, we disregard the value of provider.
// Otherwise, we use the auth config for the given provider.
func tokenSource(provider, authToken string) (oauth2.TokenSource, error) {
	// Ignore provider if authToken is given.
	if authToken != "" {
		tok := &oauth2.Token{AccessToken: authToken}
		return oauth2.StaticTokenSource(tok), nil
	}

	conf := authConfig[provider]
	if conf == nil {
		return nil, fmt.Errorf("no auth config for %q", provider)
	}
	t, err := readToken(provider)
	if err != nil {
		t, err = authorize(conf)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to obtain access token for %q", provider)
	}
	cache := &cachedTokenSource{
		src:      conf.TokenSource(context.Background(), t),
		provider: provider,
		config:   conf,
	}
	return oauth2.ReuseTokenSource(nil, cache), nil
}

// authorize performs user authorization flow, asking for permissions grant.
func authorize(conf *oauth2.Config) (*oauth2.Token, error) {
	aurl := conf.AuthCodeURL("unused", oauth2.AccessTypeOffline)
	fmt.Printf("Authorize me at following URL, please:\n\n%s\n\nCode: ", aurl)
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, err
	}
	return conf.Exchange(context.Background(), code)
}

// cachedTokenSource stores tokens returned from src on local disk.
// It is usually combined with oauth2.ReuseTokenSource.
type cachedTokenSource struct {
	src      oauth2.TokenSource
	provider string
	config   *oauth2.Config
}

func (c *cachedTokenSource) Token() (*oauth2.Token, error) {
	t, err := c.src.Token()
	if err != nil {
		t, err = authorize(c.config)
	}
	if err != nil {
		return nil, err
	}
	writeToken(c.provider, t)
	return t, nil
}

// readToken deserializes token from local disk.
func readToken(provider string) (*oauth2.Token, error) {
	l, err := tokenLocation(provider)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(l)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	return t, json.Unmarshal(b, t)
}

// writeToken serializes token tok to local disk.
func writeToken(provider string, tok *oauth2.Token) error {
	l, err := tokenLocation(provider)
	if err != nil {
		return err
	}
	w, err := os.Create(l)
	if err != nil {
		return err
	}
	defer w.Close()
	b, err := json.MarshalIndent(tok, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// tokenLocation returns a local file path, suitable for storing user credentials.
func tokenLocation(provider string) (string, error) {
	d := homedir()
	if d == "" {
		log.Printf("WARNING: unable to identify user home dir")
	}
	d = path.Join(d, ".config", "claat")
	if err := os.MkdirAll(d, 0700); err != nil {
		return "", err
	}
	return path.Join(d, provider+"-cred.json"), nil
}

func homedir() string {
	if v := os.Getenv("HOME"); v != "" {
		return v
	}
	d, p := os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH")
	if d != "" && p != "" {
		return d + p
	}
	return os.Getenv("USERPROFILE")
}
