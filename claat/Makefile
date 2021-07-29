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

OUTDIRNAME := bin
OUTDIR := ./$(OUTDIRNAME)
OUTRELEASE := gs://claat
RELEASES=$(OUTDIR)/claat-darwin-amd64 \
	 $(OUTDIR)/claat-linux-amd64 \
	 $(OUTDIR)/claat-linux-386 \
	 $(OUTDIR)/claat-windows-amd64.exe \
	 $(OUTDIR)/claat-windows-386.exe

ifeq ($(OS), Windows_NT)
		SRCS = $(shell dir *.go /b/s) 
	else
		SRCS = $(shell find . -name '*.go') 
endif

all: $(OUTDIR)/claat

$(OUTDIR)/claat: $(SRCS) VERSION
	go build -o $@ -ldflags "-X main.version=$(VERSION)"

serve: $(OUTDIR)/claat
	mkdir -p $(OUTDIR)/deps
	cd ../; \
	    bazel build :bundle; \
	    mv -f bazel-genfiles/bundle.zip claat/$(OUTDIRNAME)/deps
	unzip -o $(OUTDIR)/deps/bundle.zip -d $(OUTDIR)/deps/codelab-elements
	cd bin; ./claat serve

release: $(RELEASES)
	echo $(VERSION) > $(OUTDIR)/VERSION
	cd $(OUTDIR) && sha1sum claat* > sha1sum.txt

clean:
	rm -rf $(OUTDIR) render/tmpldata.go 

$(OUTDIR)/claat-%: GOOS=$(firstword $(subst -, ,$*))
$(OUTDIR)/claat-%: GOARCH=$(subst .exe,,$(word 2,$(subst -, ,$*)))
$(OUTDIR)/claat-%: $(SRCS) VERSION
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-buildmode exe \
		-ldflags "-X main.version=$(VERSION) -s -w" \
		-o $@

%/tmpldata.go: %/gen-tmpldata.go %/template.html %/template.md %/template-offline.html
	cd $* && go generate
