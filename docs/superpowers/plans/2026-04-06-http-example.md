# http-example Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go HTTP server that serves "I'm running on your machine", with GitHub Actions CI/CD that auto-releases cross-platform binaries on every push to main.

**Architecture:** Flat module with a `server` package exporting a pure `NewHandler()` function, and a `main.go` that wires up the handler with port-retry logic. GitHub Actions workflow with three jobs: test, build (matrix), release with auto-incrementing semver tags.

**Tech Stack:** Go 1.26.1, net/http stdlib, GitHub Actions, `gh` CLI for releases.

---

## File Map

| File | Responsibility |
|------|---------------|
| `go.mod` | Module definition: `github.com/antonioshadji/http-example` |
| `server/server.go` | Pure function `NewHandler()` returning `http.Handler` |
| `server/server_test.go` | Tests for handler: status, body, content-type |
| `main.go` | Wiring: handler setup, port retry loop 8080-8089, logging |
| `.gitignore` | Ignore Go binaries and OS artifacts |
| `.github/workflows/release.yml` | CI/CD: test, cross-compile, auto-versioned release |

---

### Task 1: Initialize Go module

**Files:**
- Create: `go.mod`

- [ ] **Step 1: Initialize the module**

```bash
cd /Users/ahadjigeorgalis/Code/gocode/http_example
go mod init github.com/antonioshadji/http-example
```

Expected: Creates `go.mod` with:
```
module github.com/antonioshadji/http-example

go 1.26.1
```

- [ ] **Step 2: Create .gitignore**

Create `.gitignore` with:
```
# Binaries
http-example
http-example-*
*.exe

# OS
.DS_Store
```

- [ ] **Step 3: Commit**

```bash
git add go.mod .gitignore
git commit -m "chore: initialize Go module and gitignore"
```

---

### Task 2: Write failing tests for the handler

**Files:**
- Create: `server/server_test.go`

- [ ] **Step 1: Create the test file**

```go
package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/antonioshadji/http-example/server"
)

func TestNewHandler_StatusOK(t *testing.T) {
	handler := server.NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestNewHandler_BodyContainsMessage(t *testing.T) {
	handler := server.NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "I'm running on your machine") {
		t.Errorf("expected body to contain %q, got %q", "I'm running on your machine", body)
	}
}

func TestNewHandler_ContentTypeHTML(t *testing.T) {
	handler := server.NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	ct := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		t.Errorf("expected Content-Type text/html, got %q", ct)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./server/...
```

Expected: FAIL — `server.NewHandler` does not exist yet.

- [ ] **Step 3: Commit the failing tests**

```bash
git add server/server_test.go
git commit -m "test: add failing tests for server handler"
```

---

### Task 3: Implement the handler to make tests pass

**Files:**
- Create: `server/server.go`

- [ ] **Step 1: Write the handler**

```go
package server

import (
	"fmt"
	"net/http"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<!DOCTYPE html><html><head><title>http-example</title></head><body><h1>I'm running on your machine</h1></body></html>")
	})
	return mux
}
```

- [ ] **Step 2: Run tests to verify they pass**

```bash
go test ./server/... -v
```

Expected: All 3 tests PASS.

- [ ] **Step 3: Commit**

```bash
git add server/server.go
git commit -m "feat: implement NewHandler serving home page"
```

---

### Task 4: Write main.go with port retry logic

**Files:**
- Create: `main.go`

- [ ] **Step 1: Write main.go**

```go
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/antonioshadji/http-example/server"
)

func main() {
	handler := server.NewHandler()

	const basePort = 8080
	const maxAttempts = 10

	for i := range maxAttempts {
		port := basePort + i
		addr := fmt.Sprintf(":%d", port)

		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Printf("port %d unavailable: %v", port, err)
			continue
		}

		log.Printf("listening on http://localhost:%d", port)
		if err := http.Serve(listener, handler); err != nil {
			log.Fatalf("server error: %v", err)
		}
		return
	}

	log.Fatalf("failed to bind to any port in range %d-%d", basePort, basePort+maxAttempts-1)
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go build -o /dev/null .
```

Expected: Exits 0, no errors.

- [ ] **Step 3: Smoke test manually**

```bash
go run . &
curl -s http://localhost:8080
kill %1
```

Expected: HTML output containing "I'm running on your machine".

- [ ] **Step 4: Commit**

```bash
git add main.go
git commit -m "feat: add main with port retry logic (8080-8089)"
```

---

### Task 5: Create GitHub repository

- [ ] **Step 1: Create the remote repo**

```bash
gh repo create antonioshadji/http-example --public --source=. --remote=origin
```

Expected: Creates the repo on GitHub and adds the `origin` remote.

- [ ] **Step 2: Push main branch**

```bash
git push -u origin main
```

Expected: All commits pushed to `origin/main`.

---

### Task 6: Add GitHub Actions workflow

**Files:**
- Create: `.github/workflows/release.yml`

- [ ] **Step 1: Create the workflow file**

```yaml
name: Build and Release

on:
  push:
    branches: [main]

permissions:
  contents: write

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.26.1"

      - name: Run tests
        run: go test ./... -v

  build:
    needs: test
    runs-on: ubuntu-latest
    strategy:
      matrix:
        include:
          - goos: linux
            goarch: amd64
          - goos: darwin
            goarch: arm64
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.26.1"

      - name: Build binary
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
        run: go build -o http-example-${{ matrix.goos }}-${{ matrix.goarch }} .

      - name: Upload artifact
        uses: actions/upload-artifact@v4
        with:
          name: http-example-${{ matrix.goos }}-${{ matrix.goarch }}
          path: http-example-${{ matrix.goos }}-${{ matrix.goarch }}

  release:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Download all artifacts
        uses: actions/download-artifact@v4
        with:
          path: artifacts

      - name: Determine next version
        id: version
        run: |
          LATEST=$(git tag --list 'v0.0.*' --sort=-version:refname | head -n 1)
          if [ -z "$LATEST" ]; then
            NEXT="v0.0.1"
          else
            PATCH=$(echo "$LATEST" | sed 's/v0\.0\.//')
            NEXT="v0.0.$((PATCH + 1))"
          fi
          echo "tag=$NEXT" >> "$GITHUB_OUTPUT"

      - name: Create release
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          gh release create "${{ steps.version.outputs.tag }}" \
            artifacts/**/* \
            --title "${{ steps.version.outputs.tag }}" \
            --generate-notes
```

- [ ] **Step 2: Commit and push**

```bash
git add .github/workflows/release.yml
git commit -m "ci: add build and release workflow"
git push
```

Expected: Push triggers the workflow. Check status:

```bash
gh run watch
```

- [ ] **Step 3: Verify the release was created**

```bash
gh release list
```

Expected: Shows `v0.0.1` with two binary assets.

- [ ] **Step 4: Verify go install works**

```bash
go install github.com/antonioshadji/http-example@v0.0.1
```

Expected: Binary installed to `$GOPATH/bin/http-example`.
