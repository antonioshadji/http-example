# Embed HTML with go:embed — Design Spec

## Overview

Refactor the HTTP server to serve an external HTML file embedded into the binary via Go's `embed` package, replacing the inline HTML string. The HTML page uses a dark terminal-inspired style with a teal accent.

## File Changes

### New: `content/templates/index.html`

Standalone HTML file with dark terminal styling:
- Dark background (`#1a1a2e` or similar)
- Teal/green accent color (`#00d4aa` or similar) for the heading
- Centered content, system font stack
- `<h1>` containing "I'm running on your machine"
- Responsive, minimal CSS embedded in a `<style>` tag

### New: `content/content.go`

Package `content`:
- `//go:embed templates/index.html` directive (path relative to `content/` directory)
- Exports `var IndexHTML []byte`

The `templates/` folder lives inside `content/` so that `//go:embed` can reference it directly (Go embed paths are relative to the source file's directory and cannot use `..`).

### Modified: `server/server.go`

- Signature changes from `NewHandler() http.Handler` to `NewHandler(html []byte) http.Handler`
- Handler writes the received bytes instead of an inline string
- Still sets `Content-Type: text/html; charset=utf-8`
- Remains a pure function — takes data as input, returns a handler

### Modified: `server/server_test.go`

- Tests pass test HTML bytes: `server.NewHandler([]byte("<html>I'm running on your machine</html>"))`
- Same three assertions: status 200, body contains message, content-type is text/html

### Modified: `main.go`

- Imports `github.com/antonioshadji/http-example/content`
- Passes `content.IndexHTML` to `server.NewHandler(content.IndexHTML)`
- No other changes to port retry logic

## Testing

- Existing tests adapted to new `NewHandler([]byte)` signature
- Tests use test HTML, not the real embedded file — keeps them fast and independent
- All three test cases preserved: status, body content, content-type

## No Changes

- `.github/workflows/release.yml` — no changes needed, `go build` automatically includes embedded files
- Port retry logic in `main.go` — unchanged
- `.gitignore` — unchanged
