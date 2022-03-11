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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
)

const (
	// auth scopes needed by the program
	scopeDriveReadOnly = "https://www.googleapis.com/auth/drive.readonly"

	// program credentials for installed apps
	googClient = "183908478743-e8rth9fbo7juk9eeivgp23asnt791g63.apps.googleusercontent.com"
	googSecret = "ljELuf5jUrzcOxZGL7OQfkIC"

	// token providers
	ProviderGoogle = "goog"

	// service account creds (plan B is to use a SA like GitWhisperer)
	qwiklabsServiceAccountCreds = "qwiklabs-services-prod-e569bfaea3cd.json"
)

var (
	googleAuthConfig = oauth2.Config{
		ClientID:     googClient,
		ClientSecret: googSecret,
		Scopes:       []string{scopeDriveReadOnly},
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
)

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
	// ts, err := h.tokenSourceServiceAccount()
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

// Plan B
// Creates a new oauth2.TokenSource from SA
// Use the same SA as GitWhisperer
func (h *Helper) tokenSourceServiceAccount() (oauth2.TokenSource, error) {
  l, err := tokenServiceAccountLocation()
  if err != nil {
	  return nil, err
  }

  b, err := ioutil.ReadFile(l)
  if err != nil {
	  return nil, err
  }

  var c = struct {
        Email      string `json:"client_email"`
        PrivateKey string `json:"private_key"`
        TokenUri   string `json:"token_uri"`
  }{}
  json.Unmarshal(b, &c)

  config := &jwt.Config{
        Email:      c.Email,
        PrivateKey: []byte(c.PrivateKey),
	Scopes:     []string{scopeDriveReadOnly},
        TokenURL:   c.TokenUri,
  }

  t := config.TokenSource(oauth2.NoContext)

  log.Printf("Plan B: attrs\n%s, %s, %s", c.Email, c.PrivateKey, c.TokenUri)

  token, err := t.Token()
  log.Printf("Plan B: token\n%s", token)

  return t, nil
}

func tokenServiceAccountLocation() (string, error) {
	d := homedir()
	if d == "" {
		log.Printf("WARNING: unable to identify user home dir")
	}
	d = path.Join(d, ".config", "claat")
	if err := os.MkdirAll(d, 0700); err != nil {
		return "", err
	}
	return path.Join(d, qwiklabsServiceAccountCreds), nil
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
		return nil, fmt.Errorf("unable to obtain access token for %q: %#v", h.provider, err)
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
// dtc: 3/11/2022
//   copy cred file to /tmp so AppEngine Flex can read/write
func tokenLocation(provider string) (string, error) {
	filename := provider+"-cred.json"

	// only copy if /tmp/goog-cred.json does not exist
	if exists, _ := FileExists(path.Join("/tmp", filename)); !exists {
	  d := homedir()
	  if d == "" {
		  log.Printf("WARNING: unable to identify user home dir")
	  }
	  d = path.Join(d, ".config", "claat")
	  if err := os.MkdirAll(d, 0700); err != nil {
		  return "", err
	  }

	  if exists, _ = FileExists(path.Join(d, filename)); exists {
	    // copy file to /tmp
	    bytesRead, err := ioutil.ReadFile(path.Join(d, filename))
	    if err != nil {
	      log.Fatal(err)
	    }
	    err = ioutil.WriteFile(path.Join("/tmp", filename), bytesRead, 0666)
	    if err != nil {
	      log.Fatal(err)
	    }
	    log.Printf("tokenLocation: config file copied!")
	  }
	}

	return path.Join("/tmp", filename), nil
}

func FileExists(name string) (bool, error) {
    _, err := os.Stat(name)
    if err == nil {
        return true, nil
    }
    if errors.Is(err, os.ErrNotExist) {
        return false, nil
    }
    return false, err
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
	err = writeToken(c.provider, t)
	if err != nil {
		// AppEngine debug
	        log.Printf("Can't writeToken [%#v]", err)
	}

	return t, nil
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
