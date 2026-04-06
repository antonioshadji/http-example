# http-example Design Spec

## Overview

A Go HTTP server that displays "I'm running on your machine" on the home page. Published to GitHub with CI/CD that builds cross-platform binaries and creates a versioned release on every push to `main`. Installable via `go install github.com/antonioshadji/http-example@latest`.

## Project Structure

```
http-example/
├── main.go                        # Wires up server, starts with port retry
├── server/server.go               # Pure function: NewHandler() returns http.Handler
├── server/server_test.go          # Tests handler via httptest
├── go.mod                         # module github.com/antonioshadji/http-example, go 1.26.1
├── .github/workflows/release.yml  # CI: build + release on push to main
└── .gitignore
```

## Application Code

### server/server.go

Exports a single pure function:

- `NewHandler() http.Handler` — returns a handler that responds to `GET /` with an HTML page containing "I'm running on your machine". Returns `text/html` content type with status 200.

No global state. No side effects. Pure function returning a handler.

### main.go

Wiring only:

1. Calls `server.NewHandler()` to get the handler.
2. Attempts to bind to port 8080.
3. If the port is already bound, increments by 1 and retries.
4. Tries ports 8080 through 8089 (10 attempts max).
5. If all 10 ports fail, logs the error and exits with a non-zero code.
6. On success, logs which port the server is listening on.

### server/server_test.go

Uses `httptest.NewRecorder` to verify:

- `GET /` returns HTTP 200.
- Response body contains "I'm running on your machine".
- Content-Type header is `text/html`.

## CI/CD — GitHub Actions

### Workflow: release.yml

**Trigger:** Push to `main` branch.

**Job 1: Test**
- Checkout code.
- Set up Go 1.26.1.
- Run `go test ./...`.

**Job 2: Build**
- Depends on Test passing.
- Matrix: `linux/amd64`, `darwin/arm64`.
- Cross-compiles with `GOOS`/`GOARCH` environment variables.
- Produces binaries named `http-example-<os>-<arch>`.
- Uploads binaries as workflow artifacts.

**Job 3: Release**
- Depends on Build passing.
- Fetches the latest `v0.0.x` tag from the repo.
- Increments the patch number (starts at `v0.0.1` if no tags exist).
- Creates a Git tag and pushes it.
- Creates a GitHub release using `gh release create`.
- Attaches both platform binaries to the release.

### Versioning

- Sequential patch bumps: `v0.0.1`, `v0.0.2`, `v0.0.3`, etc.
- Each push to `main` creates exactly one new version.
- Tags are pushed so `go install github.com/antonioshadji/http-example@v0.0.x` works.

## GitHub Repository

- **Name:** http-example
- **Owner:** antonioshadji
- **Visibility:** Public
- **Module path:** `github.com/antonioshadji/http-example`

## Target Platforms

| OS      | Arch  |
|---------|-------|
| linux   | amd64 |
| darwin  | arm64 |

## Error Handling

- Port binding failure: retry ports 8080-8089, log each attempt, exit non-zero after 10 failures.
- `http.ListenAndServe` error: log and exit non-zero.
- No retry logic beyond port binding. No graceful shutdown.
