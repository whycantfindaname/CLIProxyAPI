# Local `lwj_dev` Differences from `main`

## Comparison state

- Repository: `router-for-me/CLIProxyAPI`
- Upstream comparison baseline: `origin/main` at
  `81e1b5374f99c212f196f34956eeed964a46b8fa`.
- Active personal branch: `lwj_dev`.
- `origin` is the upstream read/synchronization remote; `fork` is the personal
  publication remote.
- This document describes source differences only. It does not claim that the
  branch is pushed, built into a deployed binary, activated, or live-verified.

## Branch-owned surfaces

### Agent workflow and platform metadata

- `AGENTS.md` is a symbolic link to the regular `CLAUDE.md` instruction source,
  keeping Codex and Claude Code on one repository-local rule document.
- `.agents/`, `.claude/`, `.codex/`, `.trellis/`, and `.zcode/` carry the
  branch-local Trellis workflow, agent definitions, hooks, Skills, specs, and
  platform adapters.
- `.gitattributes` preserves the repository treatment required by those
  workflow files.

### Codex priority WebSocket routing

- `internal/runtime/executor/codex_websockets_executor.go` routes eligible
  streaming requests whose resolved payload requests `service_tier=priority`
  to the existing Codex upstream WebSocket transport when the selected auth
  explicitly enables `websockets=true`.
- `internal/runtime/executor/codex_websockets_executor_priority_test.go`
  covers direct and rule-derived priority, Claude ingress, thinking suffixes,
  HTTP fallback, image exclusion, downstream WebSocket behavior, and
  non-streaming scope.

### Runtime management password file

- `internal/api/handlers/management/runtime_secret.go` resolves a runtime
  management password from `MANAGEMENT_PASSWORD` or
  `MANAGEMENT_PASSWORD_FILE`, with a non-empty direct value taking precedence.
- The management handler and API server wiring use the resolved value to
  authenticate and enable management routes.
- Tests cover unset, direct, file, missing-file, empty-file, precedence, and
  route-enablement behavior without committing any secret value.

### Prefixed force-mapped model aliases

- `sdk/cliproxy/auth/conductor_models.go` restores the full auth-prefixed alias
  in non-streaming and streaming responses when API-key force mapping is
  enabled, preventing the upstream model name from leaking into the
  client-visible response model.
- `sdk/cliproxy/auth/conductor_force_mapping_test.go` covers both response
  modes for the prefixed Claude API-key path.

### Claude fingerprint test baseline

- `internal/runtime/executor/claude_executor_test.go` preserves the personal
  branch's configured `MacOS`/`arm64` fingerprint assertions. This is a test
  baseline and does not add a runtime platform override.

### Upstream integration formatting

- `internal/pluginhost/host.go`,
  `internal/runtime/executor/claude_thinking_replay_test.go`, and
  `internal/runtime/executor/codex_stream_bootstrap_buffering_test.go` contain
  formatting-only `gofmt` normalization applied while integrating upstream
  `81e1b5374f99c212f196f34956eeed964a46b8fa`.

## Verification entry points

```bash
go test ./internal/runtime/executor -run 'CodexAutoExecutor|CodexPriority|TestApplyClaudeHeaders'
go test ./internal/api/handlers/management -run TestLoadRuntimeManagementSecret
go test ./internal/api -run TestManagementPasswordFile
go test ./sdk/cliproxy/auth -run TestManagerExecute_APIKeyPrefixedAliasForceMappingRestoresFullModel
go test ./...
go build -o test-output ./cmd/server
git diff --check
```

Use the explicit command environments required by `CLAUDE.md`. Build output is
temporary and must not be committed.

## State outside the Git diff

Auth JSON, management-password files, provider credentials, deployed binaries,
service supervisors, PIDs, logs, and runtime configuration are outside this
repository's tracked source. Source tests and builds do not establish runtime
activation, deployment, or provider connectivity.
