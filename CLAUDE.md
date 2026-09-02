# AGENTS.md

Go 1.26+ proxy server providing OpenAI/Gemini/Claude/Codex compatible APIs with OAuth and round-robin load balancing.

## Managed repository workflow

- Delivery status: `delivery_pending`; service activation and secret ownership remain platform-owned.
- Read order: `AGENTS.md`/`CLAUDE.md` -> project documentation.
- No project-local delivery contract is active. Add a v2 contract only after shared delivery ownership is modeled and the Agent Infra registry is changed to `delivery_contract`.

## Repository
- GitHub: https://github.com/router-for-me/CLIProxyAPI

## Commands
```bash
gofmt -w . # Format (required after Go changes)
go build -o cli-proxy-api ./cmd/server # Build
go run ./cmd/server # Run dev server
go test ./... # Run all tests
go test -v -run TestName ./path/to/pkg # Run single test
go build -o test-output ./cmd/server && rm test-output # Verify compile (REQUIRED after changes)
```
- Common flags: `--config <path>`, `--tui`, `--standalone`, `--local-model`, `--no-browser`, `--oauth-callback-port <port>`

## Command Environment
- Do not rely on inherited terminal environment when running build, test, server, model-fetch, or Codex-related commands. Set the required environment explicitly for each invocation.
- Prefer command-local environment prefixes, for example `env -u http_proxy -u https_proxy GOPROXY=... go test ./...`, or script-local exports that document proxy state, cache roots, config paths, and credentials scope.
- Clear proxy variables before domestic mirror or local-network operations. Only enable proxy variables explicitly for upstream services that require them, and never let provider tokens or auth paths leak into logs.

## Config
- Default config: `config.yaml` (template: `config.example.yaml`)
- `.env` is auto-loaded from the working directory
- Auth material defaults under `auths/`
- Storage backends: file-based default; optional Postgres/git/object store (`PGSTORE_*`, `GITSTORE_*`, `OBJECTSTORE_*`)

## Architecture
- `cmd/server/` — Server entrypoint
- `internal/api/` — Gin HTTP API (routes, middleware, modules)
- `internal/api/modules/amp/` — Amp integration (Amp-style routes + reverse proxy)
- `internal/thinking/` — Main thinking/reasoning pipeline. `ApplyThinking()` (apply.go) parses suffixes (`suffix.go`, suffix overrides body), normalizes config to canonical `ThinkingConfig` (`types.go`), normalizes and validates centrally (`validate.go`/`convert.go`), then applies provider-specific output via `ProviderApplier`. Do not break this "canonical representation → per-provider translation" architecture.
- `internal/runtime/executor/` — Per-provider runtime executors (incl. Codex WebSocket)
- `internal/translator/` — Provider protocol translators (and shared `common`)
- `internal/registry/` — Model registry + remote updater (`StartModelsUpdater`); `--local-model` disables remote updates
- `internal/store/` — Storage implementations and secret resolution
- `internal/managementasset/` — Config snapshots and management assets
- `internal/cache/` — Request signature caching
- `internal/watcher/` — Config hot-reload and watchers
- `internal/wsrelay/` — WebSocket relay sessions
- `internal/usage/` — Usage and token accounting
- `internal/tui/` — Bubbletea terminal UI (`--tui`, `--standalone`)
- `sdk/cliproxy/` — Embeddable SDK entry (service/builder/watchers/pipeline)
- `test/` — Cross-module integration tests

## Code Conventions
- Keep changes small and simple (KISS)
- Comments in English only
- If editing code that already contains non-English comments, translate them to English (don’t add new non-English comments)
- For user-visible strings, keep the existing language used in that file/area
- New Markdown docs should be in English unless the file is explicitly language-specific (e.g. `README_CN.md`)
- As a rule, do not make standalone changes to `internal/translator/`. You may modify it only as part of broader changes elsewhere.
- If a task requires changing only `internal/translator/`, run `gh repo view --json viewerPermission -q .viewerPermission` to confirm you have `WRITE`, `MAINTAIN`, or `ADMIN`. If you do, you may proceed; otherwise, file a GitHub issue including the goal, rationale, and the intended implementation code, then stop further work.
- `internal/runtime/executor/` should contain executors and their unit tests only. Place any helper/supporting files under `internal/runtime/executor/helps/`.
- Follow `gofmt`; keep imports goimports-style; wrap errors with context where helpful
- Do not use `log.Fatal`/`log.Fatalf` (terminates the process); prefer returning errors and logging via logrus
- Shadowed variables: use method suffix (`errStart := server.Start()`)
- Wrap defer errors: `defer func() { if err := f.Close(); err != nil { log.Errorf(...) } }()`
- Use logrus structured logging; avoid leaking secrets/tokens in logs
- Avoid panics in HTTP handlers; prefer logged errors and meaningful HTTP status codes
- Timeouts are allowed only during credential acquisition; after an upstream connection is established, do not set timeouts for any subsequent network behavior. Intentional exceptions that must remain allowed are the Codex websocket liveness deadlines in `internal/runtime/executor/codex_websockets_executor.go`, the wsrelay session deadlines in `internal/wsrelay/session.go`, the management APICall timeout in `internal/api/handlers/management/api_tools.go`, and the `cmd/fetch_antigravity_models` utility timeouts

## Managed Repository Context

- Registry ID: `cliproxyapi` (Agent Infra companion manifest `manifests/companion-repositories.json`)
- Managed branch: `lwj_dev` (upstream mirror baseline: `main`)
- Repository convergence authority: Agent Infra registry and sync contract (fetch, classify, safe fast-forward)
- Owner workflow + product/runtime authority: this repository's own source, `CLAUDE.md`/`AGENTS.md` (symlinked), `README.md` and `docs/`
- Workflow status: `registered` (`project_workflow=not_migrated`; no managed-project contract yet)
- Read order: `AGENTS.md`/`CLAUDE.md` -> `README.md` (getting started / management API / usage statistics) -> `docs/` (SDK docs)
- Update triggers: managed branch or remote change; build/release chain change; platform activation change; service/config/secret ownership change; new stable error class; a completed reusable major update flow
- Registered clause: the next substantive update task hitting a trigger must either promote the verified workflow into a managed-project contract (`.agent-infra/managed-project.json`) plus human guide and current error catalog, or record a concrete no-op reason
- Do not invent workflow: until migration, follow only the docs above; do not guess build, service restart, or activation steps
