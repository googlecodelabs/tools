# Copyright 2016 Google Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHA     := $(shell git rev-parse --short HEAD)
DATE    := $(shell TZ=UTC date +%FT%T)Z
VERSION := $(shell cat VERSION)-$(DATE)-$(SHA)

OUTDIR := ./bin
OUTRELEASE := gs://claat
RELEASES=$(OUTDIR)/claat-darwin-amd64 \
	 $(OUTDIR)/claat-linux-amd64 \
	 $(OUTDIR)/claat-linux-386 \
	 $(OUTDIR)/claat-windows-amd64.exe \
	 $(OUTDIR)/claat-windows-386.exe

SRCS = $(shell find . -name '*.go') render/tmpldata.go
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

install: $(OUTDIR)/claat
	cp $< $(GOPATH)/bin

$(OUTDIR)/claat: $(SRCS) VERSION
	go build -o $@ $(LDFLAGS)

$(OUTDIR)/claat-%: GOOS=$(firstword $(subst -, ,$*))
$(OUTDIR)/claat-%: GOARCH=$(subst .exe,,$(word 2,$(subst -, ,$*)))
$(OUTDIR)/claat-%: $(SRCS) VERSION
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $@ $(LDFLAGS)

%/tmpldata.go: %/gen-tmpldata.go %/template.html %/template.md
	cd $* && go generate

release: $(RELEASES)
	echo $(VERSION) > $(OUTDIR)/VERSION

test: .test
.test: $(SRCS)
	go test ./... && touch .test

lint: .lint
.lint: $(SRCS)
	golint ./... && touch .lint

clean:
	rm -rf $(OUTDIR) render/tmpldata.go .test .lint
