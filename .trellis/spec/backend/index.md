---
name: backend-guidelines
description: Navigation for CLIProxyAPI Go service and SDK conventions
paths:
  - "**/*.go"
  - "go.mod"
  - "go.sum"
  - "config.example.yaml"
---

# CLIProxyAPI Development Guidelines

CLIProxyAPI is a Go 1.26 proxy server and embeddable SDK. It exposes OpenAI,
Gemini, Claude, and Codex-compatible protocols, selects credentials, translates
requests, executes provider calls, and optionally persists configuration and
auth state through file, Postgres, Git, or object storage.

Read `AGENTS.md` before changing code. `STRUCTURE.md` is the detailed source,
runtime-data, upstream, and personal-fork boundary map.

## Guides

| Guide | Use it for |
| --- | --- |
| [Directory Structure](./directory-structure.md) | Choosing the owning package and preserving API, SDK, executor, translator, plugin, and fork boundaries |
| [Configuration and Storage](./configuration-and-storage.md) | Config parsing, defaults, validation, hot reload, credentials, and persistence backends |
| [Error Handling](./error-handling.md) | Wrapped errors, client-fault classification, HTTP responses, streaming failures, and cleanup errors |
| [Logging](./logging-guidelines.md) | logrus fields, request IDs, Gin logging, secret redaction, and file logging |
| [Quality and Testing](./quality-guidelines.md) | Formatting, focused regressions, integration tests, required build checks, and review boundaries |

## Verification Baseline

Use an explicit command environment as required by `AGENTS.md`; do not inherit
proxy, credential, or config state accidentally.

```bash
env -u http_proxy -u https_proxy -u HTTP_PROXY -u HTTPS_PROXY go test ./...
env -u http_proxy -u https_proxy -u HTTP_PROXY -u HTTPS_PROXY go build -o test-output ./cmd/server
rm test-output
git diff --check
```

Run only the focused package tests while iterating. Run the required server
build after Go changes. Network model refreshes, service startup, and Docker
operations are separate state-changing checks and are not substitutes for unit
tests or build verification.
