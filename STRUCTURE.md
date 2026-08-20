# CLIProxyAPI Structure and Maintenance Boundaries

This document describes the role of this personal fork in the Infra workspace, its code boundaries, and its normal verification entry points. It is a repository structure and maintenance contract, not a deployment record, runtime-status report, or release promise.

## Evidence basis

This description is based on the repository's Git inventory, `find` output, source/configuration inspection, and the branch comparison against the upstream baseline.

The observed Infra `manifests/companion-repositories.json` entry records:

- `workspace_path`: `Tools/CLIProxyAPI`
- `repository`: `https://github.com/whycantfindaname/CLIProxyAPI.git`
- `upstream_repository`: `https://github.com/router-for-me/CLIProxyAPI.git`
- `branch`: `lwj_dev`
- `component`: `cliproxyapi`
- `publish_status`: `published`
- `clone_remote`: `fork`
- `consumers`: `cliproxyapi-service`, `cpamp-collector`, `coding-agents`

The manifest's `commit` field is the immutable bootstrap pin. Read that value from the Infra manifest rather than duplicating it here: publishing this document advances the companion commit, so a copied SHA would be stale by construction. The pin is not a claim about deployment state or remote freshness.

If the branch, remotes, manifest pin, or directory responsibilities change, update this evidence section and the boundary descriptions below.

## Infra role and consumers

The Infra manifest declares this repository as the `cliproxyapi` companion repository for the CLIProxyAPI service and its reusable SDK.

Manifest-declared consumers are:

- `cliproxyapi-service`: builds the proxy service from `cmd/server`.
- `cpamp-collector`: the CPA-Manager-Plus collector consumer recorded by the manifest. The concrete interface contract is maintained by the consumer and Infra configuration; the manifest does not prove that a collector process is running.
- `coding-agents`: use the service as an OpenAI-, Gemini-, Claude-, or Codex-compatible endpoint.

The same Infra configuration declares `cpa-usage-monitoring` as a service relationship produced by `cliproxyapi` and consumed by `cpamp`, and selects `cliproxyapi` in the macOS, Oppo Linux, and Oppo Windows repository sets. These are workspace orchestration relationships. They do not mean that this repository contains CPA-Manager-Plus, or that deployment on any machine has completed.

## Personal fork and upstream boundaries

### Remote and branch responsibilities

- `origin` is the upstream read/synchronization boundary: `https://github.com/router-for-me/CLIProxyAPI.git`.
- `fork` is the personal repository boundary: `https://github.com/whycantfindaname/CLIProxyAPI.git`.
- The personal branch described by the Infra entry is `lwj_dev`.
- To identify personal changes, use `git diff origin/main...lwj_dev` together with the commit history. Do not infer the current difference from the older `diff.md` alone.
- The maintenance rule is to update the local `main` from upstream first, then merge the intended upstream state into `lwj_dev`. This document does not authorize fetch, merge, commit, or push, and does not claim that any of those actions occurred.

### Responsibility of the personal patch

Relative to the upstream comparison baseline, the personal branch's controlled changes are confined to the following responsibilities:

| Path | Fork-side responsibility | Main verification entry |
| --- | --- | --- |
| `AGENTS.md`, `CLAUDE.md` | Keep the repository-local Agent instructions on one source: `AGENTS.md` is a symbolic link to `CLAUDE.md`. This affects Agent entry behavior, not service runtime behavior. | `git ls-tree HEAD AGENTS.md CLAUDE.md` |
| `internal/runtime/executor/codex_websockets_executor.go` | On top of the existing Codex HTTP/WebSocket executors, detect streaming requests whose resolved payload has `service_tier=priority` and whose auth enables `websockets=true`, then select the upstream WebSocket. Standard HTTP, image requests, non-streaming requests, and existing downstream-WebSocket paths retain their existing boundaries. The probe reuses translation, thinking, and payload-rule resolution. | `go test ./internal/runtime/executor -run 'CodexAutoExecutor|CodexPriority'` |
| `internal/runtime/executor/codex_websockets_executor_priority_test.go` | Lock down the transport choice, payload overrides, Claude ingress, thinking suffix, HTTP fallback, image exclusion, and non-streaming boundaries described above. | Same as above |
| `internal/api/handlers/management/runtime_secret.go`, `internal/api/handlers/management/runtime_secret_test.go`, `internal/api/handlers/management/handler.go`, `internal/api/server.go`, `internal/api/server_reload.go`, `internal/api/server_test.go` | Add the runtime management-password file input `MANAGEMENT_PASSWORD_FILE`. The direct environment value `MANAGEMENT_PASSWORD` takes precedence; file contents are trimmed; the server resolves the runtime secret at startup and uses it to enable management routes and the authentication override. | `go test ./internal/api/handlers/management -run TestLoadRuntimeManagementSecret`; `go test ./internal/api -run TestManagementPasswordFile` |
| `sdk/cliproxy/auth/conductor_models.go` | When force mapping is enabled and the requested model carries an auth prefix, restore the full prefixed alias in the response so the upstream model name is not exposed as the client-visible model name. | `go test ./sdk/cliproxy/auth -run TestManagerExecute_APIKeyPrefixedAliasForceMappingRestoresFullModel` |
| `sdk/cliproxy/auth/conductor_force_mapping_test.go` | Cover non-streaming and streaming response rewriting for a prefixed Claude API-key alias. | Same as above |
| `internal/runtime/executor/claude_executor_test.go` | Pin the Claude fingerprint test baseline to the configured `MacOS`/`arm64` assertions. This is a test baseline, not a new runtime setting. | `go test ./internal/runtime/executor -run TestApplyClaudeHeaders` |
| `diff.md` | Preserve the existing human-maintained branch-difference note. Its comparison baseline can become stale; it does not replace the current Git diff. | `git diff origin/main...lwj_dev` |

Production code, SDK code, examples, documentation, and CI files outside this table are treated as upstream-inherited content by default. A fork-side edit to one of those areas must become an identifiable personal patch and be reflected in this document.

## Top-level directories

The following directories appear as top-level paths in the repository's current `git ls-files` inventory:

| Path | Content and boundary |
| --- | --- |
| `.github/` | Issue templates, workflows, and the model-catalog refresh script; CI/upstream process files are not runtime data. |
| `assets/` | Logos, sponsor images, and other static images used by the README; these are controlled static resources. |
| `auths/` | The repository intentionally retains only `auths/.gitkeep` as a directory placeholder; auth JSON is excluded by `.gitignore`. |
| `cmd/` | Executable entry points for the server, model fetchers, and the Codex catalog validator. |
| `docs/` | SDK usage, advanced, access-control, and watcher documentation in English and Chinese variants. |
| `examples/` | Custom-provider, HTTP-request, translator, realtime OpenAI, and plugin examples. |
| `internal/` | Private server implementation: configuration, authentication, routing, provider executors, translation, storage, plugins, logging, and runtime components. |
| `sdk/` | Embeddable CLIProxy service, auth, executor, pipeline, session, usage, plugin ABI/API, and translator interfaces. |
| `test/` | Cross-module protocol-compatibility, tool-call, thinking, parallel-call, and usage integration tests. |

## Controlled root files

Important controlled root files are grouped by responsibility below; ordinary files do not need individual implementation descriptions here.

- Project metadata: `.dockerignore`, `.gitignore`, `go.mod`, `go.sum`, `LICENSE`.
- Agent and project documentation: [AGENTS.md](AGENTS.md), [CLAUDE.md](CLAUDE.md), [README.md](README.md), [README_CN.md](README_CN.md), [README_JA.md](README_JA.md), [diff.md](diff.md), `STRUCTURE.md`.
- Configuration templates: [config.example.yaml](config.example.yaml), [.env.example](.env.example), [.env.cluster.example](.env.cluster.example). The real `config.yaml` and `.env` are not version-controlled.
- Build and container entry points: `Dockerfile`, `docker-compose.yml`, `docker-compose.cluster.yml`, `docker-build.sh`, `docker-build.ps1`.

The README describes public capabilities and SDK documentation entry points; configuration fields are defined by `config.example.yaml` and `internal/config`; repository-local Agent rules are defined by `CLAUDE.md`.

## Key second-level directories

### Commands, service, and tests

- `cmd/server/`: the `main` entry point; parses flags, loads `.env`/configuration, initializes auth, plugins, and the API server, and starts HTTP/TLS/TUI modes.
- `cmd/fetch_antigravity_models/` and `cmd/fetch_codex_models/`: use existing auth records to fetch upstream model catalogs and write JSON to a requested path. They are explicit network tools, not mandatory server-start paths.
- `cmd/validate_codex_models/`: validates the Codex client model catalog.
- `test/`: cross-package tests; provider or internal-component tests remain beside their implementation.

### Key `internal/` components

- `internal/api/`: Gin HTTP API, routes, middleware, management API, protocol multiplexing, hot reload, and server lifecycle. Management handlers are under `internal/api/handlers/management/`; request middleware is under `internal/api/middleware/`.
- `internal/auth/`: authentication implementations for Antigravity, Claude, Codex, Kimi, Vertex, and xAI.
- `internal/config/`: YAML configuration types, defaults, loading/saving, validation, cloning, and runtime normalization.
- `internal/runtime/executor/`: provider executors and their tests; helper implementations belong under `internal/runtime/executor/helps/`. The fork's Codex priority-WebSocket patch is here.
- `internal/translator/`: provider-protocol request/response translation and shared translators. It is not the default personal-fork modification area.
- `internal/registry/`: model registry, remote updater, and the embedded model catalogs under `models/`.
- `internal/store/`: Postgres, Git, and object-storage token/config backends.
- `internal/pluginhost/` and `internal/pluginstore/`: in-process plugin ABI/host, lifecycle, registry, installation, and authentication.
- `internal/thinking/`: thinking suffixes, configuration normalization, validation, and provider-specific application.
- `internal/logging/`, `internal/cache/`, `internal/watcher/`, `internal/wsrelay/`, and `internal/tui/`: logging, request caching, configuration/client watchers, WebSocket relay, and terminal UI.

### Key `sdk/` components

- `sdk/cliproxy/`: embeddable service builder, lifecycle, provider/executor registration, model management, and session management. Its main second-level components are `auth/`, `executor/`, `pipeline/`, `session/`, and `usage/`.
- `sdk/auth/`: file token store, provider auth, and refresh interfaces.
- `sdk/access/`: inbound access-authentication provider registration and validation.
- `sdk/api/`: management API adapters callable by host programs.
- `sdk/pluginabi/`, `sdk/pluginapi/`, `sdk/pluginhost/`, and `sdk/pluginstore/`: plugin ABI, plugin protocol, host bridge, and plugin-store interfaces.
- `sdk/translator/`: reusable protocol formats, translation pipeline, and plugin hooks.

### Examples, documentation, and model catalogs

- `examples/plugin/`: C, Go, and Rust plugin examples with independent Go modules; `examples/plugin/scripts/` contains example-generation scripts.
- `docs/sdk-*.md`: SDK usage, access, advanced, and watcher documentation; Chinese variants use the `_CN` suffix.
- `internal/registry/models/models.json` and `internal/registry/models/codex_client_models.json`: version-controlled embedded model catalogs, not auth or runtime request data.

## Source, generated, runtime, and credential boundaries

### Source files and controlled resources

Go source, tests, examples, SDK documentation, READMEs, CI workflows, Docker/Compose templates, and static images under `assets/` are repository content. Go changes must follow [CLAUDE.md](CLAUDE.md) for the Go version, English comments, `gofmt`, and explicit command environments.

### Generated or process-maintained content

- `internal/registry/models/*.json` are embedded model catalogs. `.github/scripts/refresh-model-catalogs.sh` fetches the catalogs from the model repository and calls `cmd/validate_codex_models` to validate the Codex catalog; do not replace them manually with unvalidated content.
- JSON generated by `cmd/fetch_*_models` is intended for inspection or offline use and should be written to a temporary or explicitly chosen path; do not treat an artifact containing runtime account context as a credential to commit.
- Build outputs such as `cli-proxy-api`, `*.exe`, `bin/`, plugin build artifacts, temporary test binaries, and runtime `static/` are protected by `.gitignore`.

### Runtime data and credentials

The following are machine or deployment state and must not be written into `STRUCTURE.md`, the README, or Git:

- `config.yaml`, `.env`, `auths/*.json`, and token/OAuth data under the default auth directory `~/.cli-proxy-api`.
- `MANAGEMENT_PASSWORD`, the file referenced by `MANAGEMENT_PASSWORD_FILE`, and credentials carried by Postgres/Git/object-store or plugin-store environment variables.
- Runtime or storage directories such as `logs/`, `conv/`, `temp/`, `refs/`, `pgstore/`, `gitstore/`, `objectstore/`, `plugins/`, and `static/`.
- Config/auth/log/plugin paths mounted by Docker Compose, service PIDs, supervisor/tmux state, deployment binaries, and rollback copies.

`.env.example` and `.env.cluster.example` contain placeholders only. `auths/.gitkeep` is the only intentionally retained auth-directory file; ignored local files such as `.DS_Store` must be preserved and must not be cleaned up as part of documentation work.

## Paths that may change, and paths maintained elsewhere

| Boundary | Maintenance rule |
| --- | --- |
| Explicit fork patch | Maintain implementation, tests, and documentation only within the personal-patch table above; update tests and this document when behavior or interfaces change. |
| Upstream-inherited code | For services, SDKs, providers, translators, examples, and CI files outside the personal patch, use upstream `main` interfaces and history as the baseline. If the fork changes one, make the change an identifiable patch instead of describing local experiments as synchronized upstream state. |
| Model catalogs | Refresh through `.github/scripts/refresh-model-catalogs.sh` and validate with `cmd/validate_codex_models`; update this document if the script or upstream model-repository responsibility changes. |
| Runtime directories and credentials | Maintained by startup flags, configuration, environment variables, external storage, or deployment processes; do not hand-commit them to this repository. |
| Infra manifest | Repository role, consumers, workspace path, branch, and remote relationship are maintained by Infra's `manifests/companion-repositories.json`; do not copy a drifting manifest into this repository. |

## Common development and verification entry points

The commands below describe entry points; they do not mean that these services or network operations were run while this document was written. For tests, builds, model fetches, and service operations, follow [CLAUDE.md](CLAUDE.md) and set proxy state, Go caches, configuration paths, and credential scope explicitly for each command.

### Code tests and builds

```
go test ./...
go test -v -run TestName ./path/to/pkg
go build -o cli-proxy-api ./cmd/server
go build -o test-output ./cmd/server && rm test-output
```

A documentation-only change does not require `gofmt`; Go changes must run `gofmt -w .` and the project's required build verification. To avoid leaving a build artifact in the worktree, use a controlled temporary output path when appropriate, while still verifying that `./cmd/server` compiles.

### Minimum regression for the fork patch

```
go test ./internal/runtime/executor -run 'CodexAutoExecutor|CodexPriority'
go test ./internal/api/handlers/management -run TestLoadRuntimeManagementSecret
go test ./internal/api -run TestManagementPasswordFile
go test ./sdk/cliproxy/auth -run TestManagerExecute_APIKeyPrefixedAliasForceMappingRestoresFullModel
go test ./internal/runtime/executor -run TestApplyClaudeHeaders
```

### Model catalogs, containers, and this document

```
go run ./cmd/validate_codex_models --file internal/registry/models/codex_client_models.json
./.github/scripts/refresh-model-catalogs.sh
docker compose build
docker compose up -d --remove-orphans
git diff --check
git ls-files --error-unmatch STRUCTURE.md
```

Model refresh and Compose startup access external systems or change local state; run them only when explicitly needed and the environment is prepared. Markdown-link checks must use the current filesystem; every repository-local link in this document must continue to exist.

## When to update this document

Update `STRUCTURE.md` in the same change whenever:

1. A top-level directory, important root file, or key second-level directory is added, removed, or renamed.
2. A personal `lwj_dev` patch changes its files, responsibility, interface, tests, or runtime boundary.
3. A fork/upstream remote, working branch, Infra manifest workspace path, consumer, or synchronization relationship changes.
4. The model-catalog generation script, validation entry point, embedded path, or submission policy changes.
5. The ignored/persisted boundary for configuration, auth, logs, plugins, storage, or build artifacts changes.
6. A repository-local file referenced here moves, is removed, or no longer has the described responsibility.

After updates, rerun `git ls-files`, `git diff origin/main...lwj_dev`, the relevant entry-point tests, and `git diff --check`. Report tests, deployment, push, and rollback separately; never infer runtime activation from source presence.
