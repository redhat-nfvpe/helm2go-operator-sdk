# kernel-style V=1 build verbosity
ifeq ("$(origin V)", "command line")
       BUILD_VERBOSE = $(V)
endif

ifeq ($(BUILD_VERBOSE),1)
       Q =
else
       Q = @
endif

VERSION = $(shell git describe --dirty --tags --always)
GIT_COMMIT = $(shell git rev-parse HEAD)
REPO = github.com/redhat-nfvpe/helm2go-operator-sdk
BUILD_PATH = $(REPO)
PKGS = $(shell go list ./... | grep -v /vendor/)
SOURCES = $(shell find . -name '*.go' -not -path "*/vendor/*")
export RELEASE_VERSION=v0.9.0
dependency:

	$(Q)wget  https://github.com/operator-framework/operator-sdk/releases/download/${RELEASE_VERSION}/operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu

	$(Q)chmod +x operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu

	$(Q)sudo cp operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu /usr/local/bin/operator-sdk

	$(Q)rm operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu

	$(Q)wget https://storage.googleapis.com/kubernetes-helm/helm-v2.10.0-linux-amd64.tar.gz

	$(Q)tar -xvzf helm-v2.10.0-linux-amd64.tar.gz

	$(Q)sudo mv linux-amd64/helm /usr/local/bin/helm
	
	$(Q)rm helm-v2.10.0-linux-amd64.tar.gz

	helm init --client-only


export CGO_ENABLED:=0
export GO111MODULE:=on
export GOPROXY?=https://proxy.golang.org/

all: format test build/helm2go-operator-sdk

format:
	$(Q)go fmt $(PKGS)

tidy:
	$(Q)go mod tidy -v

clean:
	$(Q)rm -rf build

.PHONY: all test format tidy clean

install:
	$(Q)go install \
		-gcflags "all=-trimpath=${GOPATH}" \
		-asmflags "all=-trimpath=${GOPATH}" \
		-ldflags " \
			-X '${REPO}/version.GitVersion=${VERSION}' \
			-X '${REPO}/version.GitCommit=${GIT_COMMIT}' \
		" \
		$(BUILD_PATH)

ci-build: build/helm2go-operator-sdk-$(VERSION)-x86_64-linux-gnu ci-install

ci-install:
	mv build/helm2go-operator-sdk-$(VERSION)-x86_64-linux-gnu build/helm2go-operator-sdk

release_x86_64 := \
	build/helm2go-operator-sdk-$(VERSION)-x86_64-linux-gnu \
	build/helm2go-operator-sdk-$(VERSION)-x86_64-apple-darwin

release: clean $(release_x86_64) $(release_x86_64:=.asc)

build/helm2go-operator-sdk-%-x86_64-linux-gnu: GOARGS = GOOS=linux GOARCH=amd64
build/helm2go-operator-sdk-%-x86_64-apple-darwin: GOARGS = GOOS=darwin GOARCH=amd64

build/%: $(SOURCES)
	$(Q)$(GOARGS) go build \
		-gcflags "all=-trimpath=${GOPATH}" \
		-asmflags "all=-trimpath=${GOPATH}" \
		-ldflags " \
			-X '${REPO}/version.GitVersion=${VERSION}' \
			-X '${REPO}/version.GitCommit=${GIT_COMMIT}' \
		" \
		-o $@ $(BUILD_PATH)

build/%.asc:
	$(Q){ \
	default_key=$$(gpgconf --list-options gpg | awk -F: '$$1 == "default-key" { gsub(/"/,""); print toupper($$10)}'); \
	git_key=$$(git config --get user.signingkey | awk '{ print toupper($$0) }'); \
	if [ "$${default_key}" = "$${git_key}" ]; then \
		gpg --output $@ --detach-sig build/$*; \
		gpg --verify $@ build/$*; \
	else \
		echo "git and/or gpg are not configured to have default signing key $${default_key}"; \
		exit 1; \
	fi; \
	}

.PHONY: install release_x86_64 release

test: test/unit

test/unit:
	$(Q)go test -count=1 -short ./cmd/...
	$(Q)go test -count=1 -short ./pkg/...
	$(Q)go test -count=1 -short ./internal/...

