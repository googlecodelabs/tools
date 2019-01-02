// Copyright 2018 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// The webtest command runs a Google Closure based JS tests.
// It exits with non-zero code if any of the tests fail.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bazelbuild/rules_webtesting/go/webtest"
	"github.com/tebeka/selenium"
)

var (
	testURL = flag.String("test_url", "/demo/hello_test.html", "Initial URL to start the tests.")
	host    = flag.String("host", "localhost", "Host to bind static server to.")
	debug   = flag.Bool("debug", false, "Enable verbose logging.")
)

func main() {
	log.SetPrefix("[tools/webtest] ")
	flag.Parse()
	if *debug {
		log.Printf("arguments: \n\t%s", strings.Join(os.Args, "\n\t"))
		info, err := webtest.GetBrowserInfo()
		if err != nil {
			log.Fatalf("get browser info: %v", err)
		}
		log.Printf("webtest browser: %s", info.BrowserLabel)
		log.Printf("webtest environment: %s", info.Environment)
	}
	if *testURL == "" {
		log.Fatal("--test_url argument is required")
	}
	if !strings.HasPrefix(*testURL, "/") {
		*testURL = "/" + *testURL
	}

	cwdir, err := os.Getwd()
	if err != nil {
		log.Fatalf("getwd: %v", err)
	}
	addr, err := serve(cwdir, *host)
	if err != nil {
		log.Fatalf("serve: %v", err)
	}
	debugf("serving on %s; root dir: %s", addr, cwdir)

	// TODO: Consider making capabilities configurable.
	drv, err := webtest.NewWebDriverSession(selenium.Capabilities{})
	if err != nil {
		log.Fatalf("webdriver new session: %v", err)
	}
	runURL := fmt.Sprintf("http://%s%s", addr, *testURL)
	debugf("running tests on URL: %s", runURL)
	runErr := run(drv, runURL)
	if err := drv.Quit(); err != nil {
		log.Printf("webdriver quit: %v", err)
	}
	if runErr != nil {
		log.Fatal(runErr)
	}
}

func serve(root, host string) (addr string, err error) {
	l, err := net.Listen("tcp", host+":")
	if err != nil {
		return "", err
	}

	go func() {
		fs := http.FileServer(http.Dir(root))
		h := func(w http.ResponseWriter, r *http.Request) {
			debugf("%s %s", r.Method, r.URL)
			fs.ServeHTTP(w, r)
		}
		errServe := http.Serve(l, http.HandlerFunc(h))
		log.Fatalf("serve(%q): %v", root, errServe)
	}()

	return l.Addr().String(), nil
}

// TODO: Make run timeout and configurable based on test size (small, large, etc.)
func run(drv selenium.WebDriver, testURL string) error {
	if err := drv.Get(testURL); err != nil {
		return fmt.Errorf("webdriver GET %s: %v", testURL, err)
	}
	if err := waitTests(drv); err != nil {
		return fmt.Errorf("waitTests: %v", err)
	}
	v, err := drv.ExecuteScript(`return window.top.G_testRunner.getTestResultsAsJson();`, nil)
	if err != nil {
		return fmt.Errorf("G_testRunner.getTestResultsAsJson: %v", err)
	}
	var results map[string][]*struct{ Message, Stacktrace string }
	if err := json.Unmarshal([]byte(v.(string)), &results); err != nil {
		return fmt.Errorf("G_testRunner.getTestResultsAsJson unmarshal: %v", err)
	}
	var (
		nfail int
		fails strings.Builder
	)
	for name, res := range results {
		switch {
		default:
			debugf("%s: OK", name)
		case len(res) > 0:
			nfail++
			fmt.Fprintf(&fails, "%s: FAIL\n", name)
			for i, item := range res {
				fmt.Fprintf(&fails, "%d. %s\n%s\n", i+1, item.Message, item.Stacktrace)
			}
		}
	}
	if nfail > 0 {
		return fmt.Errorf("%d out of %d test(s) failed:\n%s", nfail, len(results), &fails)
	}
	return nil
}

func waitTests(drv selenium.WebDriver) error {
	f := func(drv selenium.WebDriver) (bool, error) {
		v, err := drv.ExecuteScript(`return window.top.G_testRunner.isFinished();`, nil)
		if err != nil {
			return false, fmt.Errorf("G_testRunner.isFinished: %v", err)
		}
		return v.(bool), nil
	}
	// TODO: make interval configurable based on test size (small, large, etc.)
	return drv.WaitWithTimeoutAndInterval(f, time.Minute, 100*time.Millisecond)
}

func debugf(format string, args ...interface{}) {
	if *debug {
		log.Printf(format, args...)
	}
}
