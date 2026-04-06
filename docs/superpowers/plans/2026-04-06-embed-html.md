# Embed HTML with go:embed — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the inline HTML string with an embedded HTML file served via Go's `embed` package, styled with a dark terminal theme.

**Architecture:** HTML lives in `content/templates/index.html`, embedded by `content/content.go` which exports `IndexHTML []byte`. `main.go` imports `content` and passes the bytes to `NewHandler(html []byte)`. Pure functional — handler receives its data, doesn't fetch it.

**Tech Stack:** Go 1.26.1, `embed` stdlib package.

---

## File Map

| File | Action | Responsibility |
|------|--------|---------------|
| `content/templates/index.html` | Create | Styled HTML page with dark terminal theme |
| `content/content.go` | Create | Embeds `templates/index.html` as `IndexHTML []byte` |
| `server/server.go` | Modify | `NewHandler(html []byte)` — accepts HTML, serves it |
| `server/server_test.go` | Modify | Pass test HTML bytes to `NewHandler` |
| `main.go` | Modify | Import `content`, pass `content.IndexHTML` to `NewHandler` |

## Structure

```
http-example/
├── content/
│   ├── content.go              # //go:embed templates/index.html
│   └── templates/
│       └── index.html          # Dark terminal-themed HTML page
├── server/
│   ├── server.go               # NewHandler(html []byte) http.Handler
│   └── server_test.go          # Tests with test HTML bytes
├── main.go                     # Imports content, wires handler, port retry
├── go.mod
└── ...
```

---

### Task 1: Create the HTML template file

**Files:**
- Create: `content/templates/index.html`

- [ ] **Step 1: Create the HTML file**

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>http-example</title>
  <style>
    * {
      margin: 0;
      padding: 0;
      box-sizing: border-box;
    }
    body {
      font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      background-color: #1a1a2e;
      color: #e0e0e0;
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 100vh;
    }
    main {
      text-align: center;
      padding: 2rem;
    }
    h1 {
      font-size: 2.5rem;
      font-weight: 700;
      color: #00d4aa;
      margin-bottom: 0.5rem;
    }
    p {
      font-size: 1rem;
      color: #888;
    }
  </style>
</head>
<body>
  <main>
    <h1>I'm running on your machine</h1>
    <p>http-example</p>
  </main>
</body>
</html>
```

- [ ] **Step 2: Commit**

```bash
git add content/templates/index.html
git commit -m "feat: add styled HTML template with dark terminal theme"
```

---

### Task 2: Update tests for new NewHandler signature

**Files:**
- Modify: `server/server_test.go`

- [ ] **Step 1: Update the test file**

Replace the entire contents of `server/server_test.go` with:

```go
package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/antonioshadji/http-example/server"
)

var testHTML = []byte("<html><body>I'm running on your machine</body></html>")

func TestNewHandler_StatusOK(t *testing.T) {
	handler := server.NewHandler(testHTML)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestNewHandler_BodyContainsMessage(t *testing.T) {
	handler := server.NewHandler(testHTML)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	body := rec.Body.String()
	if !strings.Contains(body, "I'm running on your machine") {
		t.Errorf("expected body to contain %q, got %q", "I'm running on your machine", body)
	}
}

func TestNewHandler_ContentTypeHTML(t *testing.T) {
	handler := server.NewHandler(testHTML)
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
go test ./server/... -v
```

Expected: FAIL — `NewHandler` still takes no arguments, so `server.NewHandler(testHTML)` won't compile.

- [ ] **Step 3: Commit the failing tests**

```bash
git add server/server_test.go
git commit -m "test: update tests for NewHandler([]byte) signature"
```

---

### Task 3: Update server.go to accept html bytes

**Files:**
- Modify: `server/server.go`

- [ ] **Step 1: Update the handler**

Replace the entire contents of `server/server.go` with:

```go
package server

import (
	"net/http"
)

func NewHandler(html []byte) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(html)
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
git commit -m "refactor: NewHandler accepts html []byte parameter"
```

---

### Task 4: Create the content embed package and update main.go

**Files:**
- Create: `content/content.go`
- Modify: `main.go`

- [ ] **Step 1: Create content/content.go**

```go
package content

import _ "embed"

//go:embed templates/index.html
var IndexHTML []byte
```

- [ ] **Step 2: Update main.go**

Replace the entire contents of `main.go` with:

```go
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/antonioshadji/http-example/content"
	"github.com/antonioshadji/http-example/server"
)

func main() {
	handler := server.NewHandler(content.IndexHTML)

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

- [ ] **Step 3: Verify it compiles**

```bash
go build -o /dev/null .
```

Expected: Exits 0, no errors.

- [ ] **Step 4: Run all tests**

```bash
go test ./... -v
```

Expected: All tests PASS.

- [ ] **Step 5: Commit**

```bash
git add content/content.go main.go
git commit -m "feat: add content embed package and wire into main"
```

---

### Task 5: Verify end-to-end and push

- [ ] **Step 1: Smoke test**

```bash
go run . &
sleep 1
curl -s http://localhost:8080 | head -5
kill %1
```

Expected: First lines of the dark-themed HTML page containing the style and "I'm running on your machine".

- [ ] **Step 2: Push to trigger release**

```bash
git push
```

Expected: Push triggers workflow, creates next release version with embedded HTML.

- [ ] **Step 3: Verify workflow**

```bash
gh run list --limit 1
```

Expected: Shows a queued or in-progress run for the latest push.
