---
name: directory-structure
description: Package ownership and placement rules for CLIProxyAPI
paths:
  - "cmd/**/*.go"
  - "internal/**/*.go"
  - "sdk/**/*.go"
  - "test/**/*.go"
  - "examples/**/*.go"
---

# Directory Structure

## Runtime Shape

The standalone path begins in `cmd/server/main.go`: it loads `.env` and YAML,
selects a persistence backend, registers the shared token store, builds runtime
services, and starts HTTP/TLS/TUI modes. The reusable path begins at
`sdk/cliproxy/builder.go` and `sdk/cliproxy/service.go`.

HTTP assembly belongs in `internal/api/`. `internal/api/server.go` constructs
the Gin engine and middleware; `internal/api/server_routes.go` owns public
route groups. Protocol handlers live under `sdk/api/handlers/`, while provider
network execution lives under `internal/runtime/executor/`.

Request translation and thinking are distinct boundaries:

```text
Gin route -> sdk/api handler -> auth manager -> executor
          -> translator/payload normalization -> upstream provider
```

`internal/thinking/apply.go` converts request or model-suffix input to the
canonical `ThinkingConfig` in `internal/thinking/types.go`, validates it in
`internal/thinking/validate.go`, then delegates to
provider packages such as `internal/thinking/provider/claude/` and
`internal/thinking/provider/codex/`. Keep that canonical-to-provider flow; do
not implement provider-specific thinking rules in handlers or executors.

## Ownership Map

| Path | Responsibility | Placement rule |
| --- | --- | --- |
| `cmd/` | Executable entry points | Keep flag parsing and process wiring here; reusable behavior belongs in `internal/` or `sdk/` |
| `internal/api/` | Gin server, route groups, middleware, management API, reload, protocol multiplexing | Add public routes in server route assembly and management endpoints under `handlers/management/` |
| `internal/config/` | YAML types, defaults, parsing, normalization, validation, and runtime clones | A new config field must update its type, example, normalization/defaulting, and tests together |
| `internal/runtime/executor/` | Provider executors and executor unit tests | Keep only executors and their tests here; shared executor helpers go in `internal/runtime/executor/helps/` |
| `internal/translator/` | Protocol request/response translation | Do not make a translator-only change unless the `AGENTS.md` permission rule is satisfied |
| `internal/store/` | Postgres, Git, and object-store persistence adapters | Keep backend-specific persistence here; depend on SDK contracts rather than API handlers |
| `internal/logging/` | Global formatter, Gin logging, request logging, request metadata | Reuse these helpers instead of creating package-local logging formats |
| `internal/pluginhost/`, `internal/pluginstore/` | Plugin loading, routing, lifecycle, registry, and installation | Keep host implementation private; shared contracts belong under `sdk/plugin*` |
| `sdk/cliproxy/` | Embeddable service, builder, auth scheduling, executor contracts, usage and sessions | Preserve public SDK interfaces; server-only wiring stays in `cmd/` or `internal/` |
| `sdk/auth/` | Provider login helpers and file token store | Do not duplicate token persistence in a command |
| `test/` | Cross-module protocol compatibility and translation tests | Package-local behavior stays beside its implementation as `*_test.go` |

`STRUCTURE.md` contains the fuller map for model catalogs, examples, docs,
generated artifacts, and runtime data.

## Adding or Changing Behavior

- Put a new endpoint beside its protocol family in `sdk/api/handlers/`, and
  register the route through `internal/api/server_routes.go` or the management
  route assembly. Test both handler output and route/middleware behavior when
  both contracts change.
- Put provider transport behavior in the matching executor. Reusable payload,
  cache, transport, and credential helpers belong in
  `internal/runtime/executor/helps/`.
- Put embeddable interfaces in `sdk/`; do not make public SDK code import an
  `internal/` package that an external consumer cannot access.
- Keep package names short and lowercase. File names use lowercase snake case,
  matching examples such as `server_routes.go`, `request_logging.go`, and
  `codex_websockets_executor.go`.

## Fork Boundary

This checkout is the `lwj_dev` personal branch. `STRUCTURE.md` identifies the
current fork patch surface and focused regressions. Code outside that table is
upstream-inherited by default. A new fork-specific change must be explicit and
must update `STRUCTURE.md`; never present an inherited area as locally owned.

## Anti-Patterns

- Do not place reusable service logic in `cmd/server/main.go` merely because
  startup calls it.
- Do not register routes from provider executors or make handlers own upstream
  transports.
- Do not bypass `ThinkingConfig` with one-off provider field mutations.
- Do not put general helpers directly in `internal/runtime/executor/`.
- Do not make standalone `internal/translator/` edits without following the
  repository permission rule.
- Do not write runtime config, OAuth tokens, logs, generated binaries, or
  storage workdirs into tracked source paths.

## Useful Inspection Commands

```bash
rg --files cmd internal sdk test | sort
rg -n "setupRoutes|NewServer|NewService|ApplyThinking|RegisterTokenStore" cmd internal sdk
git diff origin/main...lwj_dev -- AGENTS.md CLAUDE.md internal sdk STRUCTURE.md
```
