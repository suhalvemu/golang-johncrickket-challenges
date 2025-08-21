# ccwc

Lightweight Go implementation of the Unix `wc` (word/line/byte/char counter).  
This README explains how to build, test, run and produce versioned binaries reproducibly (using Make and Go build ldflags).

## Prerequisites

- Go (the project Go version is pinned in `go.mod` — run `go version` to confirm)
- git
- (optional) make

Ensure module files are present:
```bash
go mod tidy
```

## Quickstart

Build a local binary:
```bash
# from repo root
make build VERSION=0.0.0
# or without make
go build -ldflags "-X 'cmd.version=0.0.0'" -o build/ccwc ./wc/cmd
```

Run:
```bash
# show help
./build/ccwc --help

# count a file
./build/ccwc path/to/file.txt

# read from stdin
echo "hello world" | ./build/ccwc
```

## Versioning and embedding build info

The code exposes a `version` variable in package `cmd` that can be set at build time via `-ldflags`. Example (used by the Makefile below):

- set version: `-X 'cmd.version=1.2.3'`
- you can also embed commit/date by adding variables to `cmd` and setting them with `-X`.

Example manual build embedding version and commit:
```bash
VERSION=1.2.3
COMMIT=$(git rev-parse --short HEAD)
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

go build -ldflags "-X 'cmd.version=${VERSION}' -X 'cmd.commit=${COMMIT}' -X 'cmd.date=${DATE}'" -o build/ccwc-${VERSION}-$(uname -s)-$(uname -m) ./wc/cmd
```

## Makefile (recommended)

Use `make` to produce consistent artifacts. Minimal example Makefile targets you can adopt:

```makefile
# Example Makefile excerpts (add to your Makefile)
BUILD_DIR := build
DIST_DIR := dist
APP := ccwc

VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

GOFLAGS := -ldflags "-X 'cmd.version=$(VERSION)' -X 'cmd.commit=$(COMMIT)' -X 'cmd.date=$(DATE)'"

.PHONY: build test clean dist

build:
    mkdir -p $(BUILD_DIR)
    GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GOFLAGS) -o $(BUILD_DIR)/$(APP)-$(VERSION)-$(shell uname -s)-$(shell uname -m) ./wc/cmd

test:
    go test ./... -v

dist: build
    mkdir -p $(DIST_DIR)
    tar -C $(BUILD_DIR) -czf $(DIST_DIR)/$(APP)-$(VERSION)-$(shell uname -s)-$(shell uname -m).tgz .

clean:
    rm -rf $(BUILD_DIR) $(DIST_DIR)
```

Invoke:
```bash
# reproducible-ish local build (set VERSION explicitly)
make build VERSION=1.2.3

# run tests
make test

# create tarball
make dist VERSION=1.2.3
```

## Reproducibility best practices

- Keep `go.mod` and `go.sum` committed. Run `go mod tidy` to keep them accurate.
- Pin and document the Go toolchain used (e.g. `go 1.20` in `go.mod`) and use CI images or container builds that use the same Go version.
- Embed build metadata (version, commit, date) at build-time via `-ldflags` so binaries are traceable.
- Use `make` (or CI scripts) to centralize build commands and avoid ad-hoc local differences.
- For very large assets or binaries, use Git LFS or an artifacts server; avoid committing big blobs to the repo.
- Consider using goreleaser or a dedicated CI/CD pipeline for reproducible release artifacts across OS/arch.

## Continuous Integration

In CI:
- Use the same Go version as local dev.
- Run `go mod download` before `go build`.
- Use the Makefile to standardize commands used by CI runners.
- Publish built artifacts into a consistent `dist/` folder and attach them to releases.

## Troubleshooting

- If push fails with large objects, check for files >100MB and remove them from history or use Git LFS.
- If builds differ across machines, ensure identical `GOOS`, `GOARCH`, and Go toolchain versions.
