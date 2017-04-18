# Copyright 2016 IBM Corporation
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

.DEFAULT_GOAL	:= build

#------------------------------------------------------------------------------
# Variables
#------------------------------------------------------------------------------

SHELL 		:= /bin/bash
BINDIR		:= bin
BUILDDIR    := build
DOCKERDIR	:= docker
RELEASEDIR  := release

TARGET_OS 	:= linux windows darwin
GOHOSTOS 	:= $(shell go env GOHOSTOS)

ifndef GOOS
    GOOS := $GOHOSTOS
endif

ifndef GOARCH
	GOARCH := $(shell go env GOHOSTARCH)
endif

GOFILES		= $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GODIRS		= $(shell go list -f '{{.Dir}}' ./... | grep -vFf <(go list -f '{{.Dir}}' ./vendor/...))
GOPKGS		= $(shell go list ./... | grep -vFf <(go list ./vendor/...))

APP_VER		:= $(shell git describe 2> /dev/null || echo "unknown")

#REGISTRY_APP_NAME		:= a8registry

#REGISTRY_IMAGE_NAME		:= amalgam8/a8-registry:latest

REGISTRY_DOCKERFILE		:= $(DOCKERDIR)/Dockerfile.registry

REGISTRY_RELEASE_NAME	:= $(REGISTRY_APP_NAME)-$(APP_VER)-$(GOOS)-$(GOARCH)

EXAMPLES_RELEASE_NAME	:= a8examples-$(APP_VER)

# build flags
BUILDFLAGS	:=

# linker flags
LDFLAGS     :=

# These do not work on Mac.
ifeq ($(GOHOSTOS),linux)
	# install pkgs to speed up compilation
	BUILDFLAGS	+= -i

	# linker flags to strip symbol tables and debug information
	LDFLAGS     += -s -w
	LDFLAGS     += -linkmode external
endif

# linker flags to enable static linking
LDFLAGS 	+= -extldflags -static

# linker flags to set build info variables
#BUILD_SYM	:= github.com/amalgam8/amalgam8/pkg/version
LDFLAGS		+= -X $(BUILD_SYM).version=$(APP_VER)
LDFLAGS		+= -X $(BUILD_SYM).gitRevision=$(shell git rev-parse --short HEAD 2> /dev/null  || echo unknown)
LDFLAGS		+= -X $(BUILD_SYM).branch=$(shell git rev-parse --abbrev-ref HEAD 2> /dev/null  || echo unknown)
LDFLAGS		+= -X $(BUILD_SYM).buildUser=$(shell whoami || echo nobody)@$(shell hostname -f || echo builder)
LDFLAGS		+= -X $(BUILD_SYM).buildDate=$(shell date +%Y-%m-%dT%H:%M:%S%:z)
LDFLAGS		+= -X $(BUILD_SYM).goVersion=$(word 3,$(shell go version))

#--------------
#-- high-level
#--------------
.PHONY: verify precommit

# to be run by CI to verify validity of code changes
verify: check build test

# to be run by developer before checking-in code changes
precommit: format verify

#---------
#-- build
#---------
.PHONY: build build.registry build.controller build.sidecar build.k8srules build.cli compile clean

build: build.main

build.main:
	@echo "--> building main"
	# @go build $(BUILDFLAGS) -ldflags '$(LDFLAGS)' -o 

compile:
	@echo "--> compiling packages"
	@go build $(GOPKGS)

clean:
	@echo "--> cleaning compiled objects and binaries"
	@go clean -tags netgo -i $(GOPKGS)
	@rm -rf $(BINDIR)/*
	@rm -rf $(BUILDDIR)/*
	@rm -rf $(RELEASEDIR)/*

#--------
#-- test
#--------
.PHONY: test test.long test.integration

test:
	@echo "--> running unit tests, excluding long tests"
	@go test -v $(GOPKGS) -short

test.long:
	@echo "--> running unit tests, including long tests"
	@go test -v $(GOPKGS)

test.integration:
	@echo "--> running integration tests"
	@build_and_run.sh

#---------------
#-- checks
#---------------
.PHONY: check format format.check vet lint

check: format.check vet lint

format: tools.goimports
	@echo "--> formatting code with 'goimports' tool"
	@goimports -w -l $(GOFILES)

format.check: tools.goimports
	@echo "--> checking code formatting with 'goimports' tool"
	@goimports -l $(GOFILES) | sed -e "s/^/\?\t/" | tee >(test -z)

vet: tools.govet
	@echo "--> checking code correctness with 'go vet' tool"
	@go vet $(GOPKGS)

lint: tools.golint
	@echo "--> checking code style with 'golint' tool"
	@echo $(GODIRS) | xargs -n 1 golint

#------------------
#-- dependencies
#------------------
.PHONY: depend.update depend.install

depend.update: tools.glide
	@echo "--> updating dependencies from glide.yaml"
	@glide update --strip-vendor

depend.install:	tools.glide
	@echo "--> installing dependencies from glide.lock "
	@glide install --strip-vendor

#---------------
#-- dockerize
#---------------
.PHONY: dockerize dockerize.registry

dockerize: dockerize.registry dockerize.controller dockerize.sidecar

dockerize.registry:
	@echo "--> building registry docker image"
	@docker build -t $(REGISTRY_IMAGE_NAME) -f $(REGISTRY_DOCKERFILE) .

#---------------
#-- release
#---------------

.PHONY: release release.registry

release: release.registry


compress: COMPRESSED_FILE :=
compress:
	@upx -qqt $(COMPRESSED_FILE); RESULT=$$?; if [ $$RESULT -eq 2 ]; then \
		echo "--> compressing $(COMPRESSED_FILE)"; \
		upx -qq --best $(COMPRESSED_FILE); \
	elif [ $$RESULT -eq 1 ]; then \
		false; \
	fi

compress.registry: tools.upx
	@make --no-print-directory compress COMPRESSED_FILE=$(BINDIR)/$(REGISTRY_APP_NAME)


release.registry:
	@echo "--> packaging registry for release"
	@mkdir -p $(RELEASEDIR)
	@tar -czf $(RELEASEDIR)/$(REGISTRY_RELEASE_NAME).tar.gz --transform 's:^.*/::' $(BINDIR)/$(REGISTRY_APP_NAME) README.md LICENSE


#---------------
#-- tools
#---------------
.PHONY: tools tools.goimports tools.golint tools.govet tools.glide tools.upx

tools: tools.goimports tools.golint tools.govet tools.glide tools.upx

tools.goimports:
	@command -v goimports >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing goimports"; \
		go get golang.org/x/tools/cmd/goimports; \
    fi

tools.govet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		echo "--> installing govet"; \
		go get golang.org/x/tools/cmd/vet; \
	fi

tools.golint:
	@command -v golint >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing golint"; \
		go get github.com/golang/lint/golint; \
    fi

tools.glide:
	@command -v glide >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing glide"; \
		GLIDE_VERSION="v0.12.3"; \
		GLIDE_ARCH="$(GOOS)-$(GOARCH)"; \
		GLIDE_RELEASE="glide-$$GLIDE_VERSION-$$GLIDE_ARCH"; \
		mkdir -p /tmp/$$GLIDE_RELEASE; \
		wget -qO- https://github.com/Masterminds/glide/releases/download/$$GLIDE_VERSION/$$GLIDE_RELEASE.tar.gz | tar xz -C /tmp/$$GLIDE_RELEASE; \
		cp /tmp/$$GLIDE_RELEASE/$$GLIDE_ARCH/glide ~/bin; \
    fi

tools.upx:
	@command -v upx >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing upx"; \
		UPX_VERSION="3.91"; \
		UPX_ARCH="$(GOARCH)_$(GOOS)" # only linux (amd64|i386) are supported; \
		UPX_RELEASE="upx-$$UPX_VERSION-$$UPX_ARCH"; \
		mkdir -p /tmp/$$UPX_RELEASE; \
		wget -qO- https://github.com/upx/upx/releases/download/v$$UPX_VERSION/$$UPX_RELEASE.tar.bz2 | tar xj -C /tmp/$$UPX_RELEASE; \
		cp /tmp/$$UPX_RELEASE/$$UPX_RELEASE/upx ~/bin; \
	fi

# This package converts any file into managable Go source code
tools.go-bindata:
	@command -v go-bindata >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing go-bindata"; \
		go get -u github.com/jteeuwen/go-bindata/...; \
	fi
