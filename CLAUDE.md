# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this repo is

A **hands-on tutorial** for building Go mTLS over TCP, in three runnable steps under `tutorial/`:

| Step | Path | Topic |
|---|---|---|
| 01 | `tutorial/step01-tcp-echo/{server,client}/main.go` | Plaintext TCP echo (baseline) |
| 02 | `tutorial/step02-tls/{server,client}/main.go` | TLS with server-only auth (`Certificates`, `RootCAs`, `ServerName`) |
| 03 | `tutorial/step03-mtls/{server,client}/main.go` | mTLS (`ClientCAs`, `ClientAuth = RequireAndVerifyClientCert`) |

Each step is a **standalone main package** with its own README. Documentation/comments are in Japanese (matches the project README). Entry point for learners: [`docs/00-overview.md`](docs/00-overview.md). The Go module path is `github.com/gracefulm/go-template-project` (template origin) — use that for `goimports -local`, not the directory name.

`internal/`, `pkg/`, `api/`, `configs/`, `init/`, `web/` are empty placeholders inherited from the project-layout scaffold and currently unused.

## Common commands

All workflows go through the Makefile. Per Go 1.26 `tool` directives in `go.mod`, linters/formatters are invoked via `go tool` — do not assume globally installed binaries.

```bash
make init   # one-time: configures git commit.template -> ./.commit_template
make certs  # runs scripts/gen-certs.sh — generates CA/server/client certs into each step's certs/
make fmt    # go mod tidy + go fmt + goimports (-local github.com/gracefulm/go-template-project)
make lint   # go tool golangci-lint run
make sec    # go tool govulncheck ./...
make test   # runs lint + fmt, then: go test -race -cover ./...
```

Run any tutorial step from the repo root, e.g. `go run ./tutorial/step02-tls/server`. **The cert paths in each `main.go` are relative to the repo root** — running from inside a step directory will fail to find `certs/`.

Run a single test:

```bash
go test -race -run TestName ./path/to/pkg
```

Note: `make test` depends on `lint` and `fmt`, so it will reformat files and may fail on lint issues before any tests execute. Use `go test ./...` directly if you need to run tests without the formatting/lint gate.

## Layout convention

The directory structure follows [golang-standards/project-layout](https://github.com/golang-standards/project-layout) (community convention, not official Go layout). The README explicitly recommends keeping only directories you actually use rather than preserving empty scaffolding — feel free to delete unused top-level dirs as the project takes shape, and prefer a flat layout for small command-line projects.

`vender/` (sic — note spelling) is the vendor directory placeholder; it is not in use and `vendor/` is commented out in `.gitignore`.

## Commit messages

Commits use Japanese-style emoji prefixes defined in `.commit_template` (run `make init` to wire it into git). Format: `:prefix: <subject>` (e.g. `:add: コメント追加`, `:fix: 修正`). Common prefixes: `:fix:`, `:hotfix:`, `:add:`, `:feat:`, `:update:`, `:change:`, `:docs:`, `:remove:`, `:refactor:`, `:test:`, `:chore:`, `:style:`, `:upgrade:`, `:revert:`. See `.commit_template` for the full list.
