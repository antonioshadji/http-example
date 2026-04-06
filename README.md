# http-example

A simple Go HTTP server that displays "I'm running on your machine".

## Install

```bash
go install github.com/antonioshadji/http-example@latest
```

## Usage

```bash
http-example
```

The server starts on port 8080. If the port is unavailable, it tries ports 8081–8089.

## Build from source

```bash
git clone https://github.com/antonioshadji/http-example.git
cd http-example
go build -o http-example .
./http-example
```

## Releases

Pre-built binaries for Linux (amd64) and macOS (arm64) are available on the [releases page](https://github.com/antonioshadji/http-example/releases). A new release is created automatically on every push to `main`.
