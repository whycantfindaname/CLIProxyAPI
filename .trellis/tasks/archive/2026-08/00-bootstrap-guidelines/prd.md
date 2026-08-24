# Bootstrap CLIProxyAPI Development Guidelines

## Goal

Replace the initial Trellis scaffolding with project-specific guidance for the
CLIProxyAPI Go service and SDK, using current source, tests, repository rules,
and architecture documentation as evidence.

## Scope

- Specifications: `.trellis/spec/backend/`
- Source inspected: `cmd/server/`, `internal/api/`, `internal/config/`,
  `internal/store/`, `internal/clienterror/`, `internal/logging/`,
  `internal/thinking/`, `internal/runtime/executor/`, `sdk/`, and `test/`
- Documentation inspected: `AGENTS.md`, `CLAUDE.md`, `STRUCTURE.md`,
  `README.md`, `docs/`, `config.example.yaml`, and `go.mod`
- Excluded: product code, repository rules, generated runtime files,
  manifests, `.gitattributes`, dependencies, services, and Smart Search

## Completed Work

- [x] Replace backend guideline scaffolding with real repository conventions
- [x] Add source-backed paths, examples, anti-patterns, and verification commands
- [x] Replace the inapplicable ORM/migrations guide with configuration and storage guidance
- [x] Remove generic thinking guides that describe Trellis/TypeScript rather than CLIProxyAPI
- [x] Keep the backend index synchronized with the final file set
- [x] Verify Codex and Claude Trellis initialization artifacts
- [x] Verify developer identity resolves to `jasonliao`
- [x] Verify local links, context generation, archive status, and Git scope

## Final Specification Set

- `.trellis/spec/backend/index.md`
- `.trellis/spec/backend/directory-structure.md`
- `.trellis/spec/backend/configuration-and-storage.md`
- `.trellis/spec/backend/error-handling.md`
- `.trellis/spec/backend/logging-guidelines.md`
- `.trellis/spec/backend/quality-guidelines.md`

## Acceptance Criteria

- [x] Guidance describes the current Go service and SDK architecture.
- [x] Important rules cite current source, tests, or project documentation.
- [x] The index links exactly the retained guideline files.
- [x] No generic scaffold text or empty sections remain under `.trellis/spec/`.
- [x] The completed task is archived under `.trellis/tasks/archive/2026-08/`.
