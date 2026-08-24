---
name: configuration-and-storage
description: Configuration, credentials, hot reload, and persistence conventions
paths:
  - "config.example.yaml"
  - ".env*.example"
  - "cmd/server/**/*.go"
  - "internal/config/**/*.go"
  - "internal/store/**/*.go"
  - "internal/watcher/**/*.go"
  - "sdk/auth/**/*.go"
  - "sdk/cliproxy/**/*.go"
---

# Configuration and Storage

## Configuration Sources and Ownership

`config.example.yaml` is the user-facing field reference. Runtime types,
defaults, parsing, normalization, and validation live in `internal/config/`.
`internal/config/config_load.go` loads files, while
`internal/config/parse.go` applies the same in-memory rules to byte payloads.
Keep these paths behaviorally aligned.

`cmd/server/main.go` resolves process-level inputs such as `PGSTORE_*`,
`GITSTORE_*`, and `OBJECTSTORE_*`, chooses one persistence backend, then calls
`sdk/auth.RegisterTokenStore` once. SDK construction follows the builder and
service lifecycle in `sdk/cliproxy/`.

For mutable runtime configuration, take an independent snapshot with
`Config.CloneForRuntime()` from `internal/config/clone.go`. The regression in
`internal/config/clone_test.go` proves slices, maps, nested structs, interface
values, and YAML nodes do not share mutable references.

## Adding a Configuration Field

Update the complete contract in one change:

1. Add the field to the owning type in `internal/config/` or `sdk/config/`.
2. Add the documented key and safe example value to `config.example.yaml`.
3. Add defaults, normalization, and validation to the existing load and parse
   flow. `LoadConfigOptional` and `ParseConfigBytes` must agree.
4. Confirm hot reload and `CloneForRuntime` preserve the field when it contains
   mutable references.
5. Add focused parse/default/validation tests beside `internal/config/` and a
   server or SDK test when wiring changes.

Follow test shapes such as `internal/config/request_retry_test.go`,
`internal/config/weight_test.go`, `internal/config/plugin_config_test.go`, and
`internal/config/clone_test.go`. Test omitted, valid, and invalid values when
the field has defaults or constraints.

## Persistence Backends

The default credential backend is `sdk/auth.FileTokenStore`. Optional adapters
live in `internal/store/`:

- `postgresstore.go` stores config and auth records in Postgres and also owns
  cooldown persistence.
- `gitstore.go` synchronizes config/auth files through a configured repository
  and validates the paths included in commits.
- `objectstore.go` persists token/config data in S3-compatible object storage.

Keep backend-specific setup and recovery in the adapter. Callers should use the
shared auth/store contract and must not switch on concrete backend types for
ordinary reads and writes.

There is no general ORM or migration framework. SQL schema creation and
backend evolution belong with the Postgres adapter. Preserve transaction,
version, and last-write semantics demonstrated by
`internal/store/postgres_cooldown_store.go` and its tests. Git recovery must
preserve unrelated local changes; `internal/store/gitstore_test.go` covers
branch selection, recovery, deletion, and path-safety behavior.

## Runtime and Secret Boundary

Never commit `config.yaml`, `.env`, `auths/*.json`, management passwords,
provider tokens, storage credentials, logs, plugin binaries, or backend work
directories. Keep examples secret-free. Logs and errors must not print raw
tokens, full auth records, or secret file contents.

Do not infer a live deployment from source or example configuration. Saved
config, active process config, service health, and a real provider request are
separate evidence.

## Anti-Patterns

- Do not add a YAML key only to `config.example.yaml` or only to a Go struct.
- Do not let `LoadConfigOptional` and `ParseConfigBytes` apply different
  defaults or validation.
- Do not retain pointers, maps, slices, or YAML nodes from a mutable config in
  a runtime snapshot.
- Do not instantiate separate token stores in unrelated components after the
  shared store has been registered.
- Do not write backend credentials or auth payloads to logs or tests.
- Do not add an ORM, migration layer, or storage abstraction for a field that
  is owned by an existing adapter.

## Verification

```bash
env -u http_proxy -u https_proxy -u HTTP_PROXY -u HTTPS_PROXY go test ./internal/config
env -u http_proxy -u https_proxy -u HTTP_PROXY -u HTTPS_PROXY go test ./internal/store
env -u http_proxy -u https_proxy -u HTTP_PROXY -u HTTPS_PROXY go test ./sdk/auth ./sdk/cliproxy/...
rg -n "request-retry|max-retry-interval|RequestRetry|MaxRetryInterval" config.example.yaml internal/config sdk/config
```

The final command demonstrates the required two-sided search with existing
retry fields. For backend wiring changes, add the relevant `cmd/server` or
`sdk/cliproxy` package test rather than starting a real service with credentials.
