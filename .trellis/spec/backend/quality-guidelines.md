---
name: quality-and-testing
description: Go formatting, regression, build, review, and fork-quality rules
paths:
  - "**/*.go"
  - "go.mod"
  - "go.sum"
  - ".github/workflows/**/*.yml"
---

# Quality and Testing

## Required Code Shape

- Keep changes small and within the owning package.
- Format Go changes with `gofmt`; imports must remain goimports-style.
- Write comments in English. Keep user-visible text in the language already
  used by the endpoint or file.
- Use method-specific error names such as `errLoad`, `errClose`, and
  `errStart` when multiple operations share a scope or a deferred closure.
- Return errors instead of terminating or panicking. Wrap operational context
  with `%w` when callers need the cause.
- Preserve the no-post-connect-timeout rule and its documented exceptions in
  `AGENTS.md`.

## Test Placement and Style

Package behavior is tested beside implementation with `*_test.go`. Cross-module
protocol guarantees belong in `test/`, including thinking conversion, Claude
Code compatibility, tool-call translation, parallel calls, and usage logging.

Prefer table-driven tests with named cases for parsers, classifiers, and
provider matrices. Use `t.TempDir`, `httptest.NewServer`, and in-memory fakes so
tests do not depend on developer credentials, a running service, or persistent
filesystem state. Register cleanup with `t.Cleanup`; report cleanup errors when
they can invalidate the test.

Examples to follow:

- `internal/clienterror/client_error_test.go` for status/error taxonomy.
- `internal/config/clone_test.go` for deep-copy and mutation isolation.
- `internal/store/gitstore_test.go` for temporary repository and recovery
  behavior.
- `internal/logging/global_logger_test.go` for exact formatted output and
  secret-path suppression.
- `internal/runtime/executor/codex_websockets_executor_priority_test.go` for
  fork-specific transport boundaries.

Regression tests must be capable of failing on the old behavior. Avoid tests
that only call the implementation to compute both the expected and actual
value.

## Verification by Change Type

Use explicit environment prefixes. Add `GOPROXY` only when dependency access
is actually required and selected deliberately.

For a focused package change:

```bash
env -u http_proxy -u https_proxy -u HTTP_PROXY -u HTTPS_PROXY go test ./internal/clienterror -run TestHTTPStatusFromError
```

For all Go changes before handoff:

```bash
gofmt -w .
env -u http_proxy -u https_proxy -u HTTP_PROXY -u HTTPS_PROXY go test ./...
env -u http_proxy -u https_proxy -u HTTP_PROXY -u HTTPS_PROXY go build -o test-output ./cmd/server
rm test-output
git diff --check
```

`.github/workflows/pr-test-build.yml` builds `./cmd/server`; the repository
rules also require this compile check locally after Go changes. A local build
or test proves source behavior, not deployment, service health, or provider
connectivity.

For the current `lwj_dev` fork patch, `STRUCTURE.md` lists the minimum focused
regressions for Codex priority WebSocket selection, runtime management secrets,
prefixed model aliases, and the Claude fingerprint baseline. Run the relevant
subset whenever those files or shared dependencies change.

## Review Boundaries

- Check the complete data path when changing a protocol field: handler,
  translation, auth selection, executor, streaming/non-streaming response, and
  tests.
- Check config type, example, defaults, validation, runtime clone, reload, and
  SDK wiring when changing configuration.
- Check cancellation, credential rotation, cooldown, and client-visible error
  shape when changing retry or failure handling.
- Keep generated model catalogs under `internal/registry/models/` synchronized
  only through `.github/scripts/refresh-model-catalogs.sh` and validate the
  Codex catalog with `cmd/validate_codex_models`.
- Preserve unrelated dirty-worktree changes and runtime files.

## Anti-Patterns

- Do not make broad mechanical edits to unrelated providers for a focused fix.
- Do not add standalone `internal/translator/` changes without satisfying the
  repository permission rule.
- Do not put support files directly in `internal/runtime/executor/`; use
  `internal/runtime/executor/helps/`.
- Do not weaken a failing regression to make a build green.
- Do not use live credentials or external provider availability in unit tests.
- Do not treat tests, builds, saved config, or a local screen as deployment
  proof.
