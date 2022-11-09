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
package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"

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
	ProviderGoogle = "goog"
)

var (
	googleAuthConfig = oauth2.Config{
		ClientID:     googClient,
		ClientSecret: googSecret,
		Scopes:       []string{scopeDriveReadOnly},
		RedirectURL:  "http://localhost:8091",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
)

// The webserver waits for an oauth code in the three-legged auth flow.
func startWebServer() (code string, err error) {
	listener, err := net.Listen("tcp", "localhost:8091")
	if err != nil {
		return "", err
	}
	codeCh := make(chan string)

	go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		codeCh <- code // send code to OAuth flow
		listener.Close()
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Received oauth code\r\nYou can now safely close this browser window.")
	}))
        code = <- codeCh
	return code, nil
}

type authorizationHandler func(conf *oauth2.Config) (*oauth2.Token, error)

type internalOptions struct {
	authHandler authorizationHandler
}

type Helper struct {
	authToken string
	provider  string
	client    *http.Client
	opts      internalOptions
}

func NewHelper(at, p string, rt http.RoundTripper) (*Helper, error) {
	return newHelper(at, p, rt, internalOptions{
		authHandler: authorize,
	})
}

func newHelper(at, p string, rt http.RoundTripper, io internalOptions) (*Helper, error) {
	h := Helper{
		authToken: at,
		provider:  p,
		opts:      io,
	}

	var err error
	h.client, err = h.produceDriveClient(rt)
	if err != nil {
		return nil, err
	}

	return &h, nil
}

func (h *Helper) DriveClient() *http.Client {
	return h.client
}

func (h *Helper) produceDriveClient(rt http.RoundTripper) (*http.Client, error) {
	ts, err := h.tokenSource()
	if err != nil {
		return nil, err
	}

	if rt == nil {
		rt = http.DefaultTransport
	}

	return &http.Client{
		Transport: &oauth2.Transport{
			Source: ts,
			Base:   rt,
		},
	}, nil
}

// tokenSource creates a new oauth2.TokenSource backed by tokenRefresher,
// using previously stored user credentials if available.
// If authToken is not given at Helper init, we use the Google provider.
// Otherwise, we use the auth config for the given provider.
func (h *Helper) tokenSource() (oauth2.TokenSource, error) {
	// Create a static token source if we have an auth token.
	if h.authToken != "" {
		return oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: h.authToken,
		}), nil
	}

	// Otherwise, use the Google provider.
	t, err := readToken(h.provider)
	if err != nil {
		t, err = h.opts.authHandler(&googleAuthConfig)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to obtain access token for %q", h.provider)
	}
	cache := &cachedTokenSource{
		src:      googleAuthConfig.TokenSource(context.Background(), t),
		provider: h.provider,
		config:   &googleAuthConfig,
	}
	return oauth2.ReuseTokenSource(nil, cache), nil
}

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

// authorize performs user authorization flow, asking for permissions grant.
func authorize(conf *oauth2.Config) (*oauth2.Token, error) {
	aURL := conf.AuthCodeURL("unused", oauth2.AccessTypeOffline)
	fmt.Printf("Authorize me at following URL, please:\n\n%s\n", aURL)
	code, err := startWebServer()
	if err != nil {
		return nil, err
	}
	return conf.Exchange(context.Background(), code)
}
