// Copyright 2019 Google LLC. All Rights Reserved.
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
	"testing"

	"golang.org/x/oauth2"
)

// Fake authorization handler to skip interactivity breaking the tests.
func fakeAuthorizationHandler(_ *oauth2.Config) (*oauth2.Token, error) {
	return &oauth2.Token{}, nil
}

func TestNewHelperWithAccessToken(t *testing.T) {
	accessToken := "accTok"
	provider := "pvdr"

	h, err := NewHelper(accessToken, provider, nil)

	if err != nil {
		t.Fatalf("NewHelper(%q, %q, nil) = err, want Helper", accessToken, provider)
	}
	// For our purposes, it's enough to check that we got a helper object back at all.
	if h == nil {
		t.Fatalf("NewHelper(%q, %q, nil) = nil, want non-nil", accessToken, provider)
	}
}

func TestNewHelperWithoutAccessToken(t *testing.T) {
	provider := "pvdr"

	// Test the private version so we can sub in our fake auth handler.
	h, err := newHelper("", provider, nil, internalOptions{
		authHandler: fakeAuthorizationHandler,
	})

	if err != nil {
		t.Fatalf("NewHelper(\"\", %q, nil) = err, want Helper", provider)
	}
	// For our purposes, it's enough to check that we got a helper object back at all.
	if h == nil {
		t.Fatalf("NewHelper(\"\", %q, nil) = nil, want non-nil", provider)
	}
}

func TestDriveClient(t *testing.T) {
	accessToken := "accTok"
	provider := "pvdr"

	h, err := NewHelper(accessToken, provider, nil)
	if err != nil {
		t.Fatalf("NewHelper(%q, %q, nil) = err, want Helper", accessToken, provider)
	}

	c1 := h.DriveClient()
	if c1 == nil {
		t.Fatalf("h.DriveClient() == nil, want non-nil")
	}
}
