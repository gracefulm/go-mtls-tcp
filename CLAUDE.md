# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this repo is

A **hands-on tutorial** for building Go mTLS over TCP, in four runnable steps under `tutorial/`:

| Step | Path | Topic |
|---|---|---|
| 01 | `tutorial/step01-tcp-echo/{server,client}/main.go` | Plaintext TCP echo (baseline) |
| 02 | `tutorial/step02-tls/{server,client}/main.go` | TLS with server-only auth (`Certificates`, `RootCAs`, `ServerName`) |
| 03 | `tutorial/step03-mtls/{server,client}/main.go` | mTLS (`ClientCAs`, `ClientAuth = RequireAndVerifyClientCert`) |
| 04 | `tutorial/step04-envoy-mtls/{client,envoy}/` + `compose.yaml` | Egress-sidecar pattern: `docker compose` brings up Envoy + a containerized Step 03 server (vanilla `golang:1.26-alpine` + read-only repo bind-mount, no Dockerfile). Envoy reaches the server via compose DNS as `server:9444`. Plaintext Go client runs on the host. Reuses Step 03's `certs/` |

Each step is a **standalone main package** with its own README (Step 04 has only a client; the server is Step 03's, unmodified). Documentation/comments are in Japanese (matches the project README). Entry point for learners: [`docs/00-overview.md`](docs/00-overview.md). The `docs/` directory holds educational reference material (tcpdump, X.509/CA, openssl, certificate flow, TLS 1.2 vs 1.3, Envoy/sidecar concepts) and Step 02/03 each include a captured packet dump (`tls-dump.txt`, `mtls-dump.txt`) for handshake observation.

The Go module path is `github.com/gracefulm/go-template-project` (template origin) — use that for `goimports -local`, not the directory name.

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

Run any tutorial step from the repo root, e.g. `go run ./tutorial/step02-tls/server`. **The cert paths in each `main.go` are hard-coded relative to the repo root** (e.g. `tutorial/step02-tls/certs/server.crt`) — running from inside a step directory will fail to find `certs/`. Step 04 additionally requires Docker + `docker compose` v2; bring up Envoy + server with `docker compose -f tutorial/step04-envoy-mtls/compose.yaml up`, then run the host-side client with `go run ./tutorial/step04-envoy-mtls/client`.

Run a single test:

```bash
go test -race -run TestName ./path/to/pkg
```

Note: `make test` depends on `lint` and `fmt`, so it will reformat files and may fail on lint issues before any tests execute. Use `go test ./...` directly if you need to run tests without the formatting/lint gate.

## Cert generation

`scripts/gen-certs.sh` is idempotent and overwrites each step's `certs/` directory. It generates a fresh root CA every run, so any server already running with the old CA must be restarted. Server cert is `CN=localhost`, SAN `localhost,127.0.0.1`, EKU `serverAuth`. Client cert is `CN=alice@example.com`, EKU `clientAuth`. Step 02 receives `ca.crt`, `server.{crt,key}`; Step 03 additionally receives `client.{crt,key}`.

## Commit messages

Commits use Japanese-style emoji prefixes defined in `.commit_template` (run `make init` to wire it into git). Format: `:prefix: <subject>` (e.g. `:add: コメント追加`, `:fix: 修正`). Common prefixes: `:fix:`, `:hotfix:`, `:add:`, `:feat:`, `:update:`, `:change:`, `:docs:`, `:remove:`, `:refactor:`, `:test:`, `:chore:`, `:style:`, `:upgrade:`, `:revert:`. See `.commit_template` for the full list.
