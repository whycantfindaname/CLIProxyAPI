# Prune unused session ID wrapper

## Goal

Remove the unused private session ID wrapper approved in S2 of
`/Users/jasonliao/Desktop/code/Artifacts/infra-dead-code-audit-20261003/DELETION_PLAN.md`,
preserving the existing payload identity behavior and public SDK interface.

## Background and authorization

- The approved plan identifies `sdk/cliproxy/auth/selector.go:1906` and
  `sdk/cliproxy/auth/selector_test.go:895`. The unexported `extractSessionID`
  only forwards to `ExtractSessionID(nil, payload, nil)`; repository search
  found its definition and one test-only call, with no production callers.
- `TestExtractSessionID` retains seven real payload cases: Claude Code user
  ID, JSON user ID with and without session ID, ordinary user ID,
  conversation ID, missing metadata, and empty payload.
- The user approved the deletion plan and implementation with “可以，请继续”,
  then separately confirmed creating and starting affected Trellis tasks with
  “确认”. Reuse the approved plan without a new design or interview.
- Scope is source-only and uncommitted on the existing `lwj_dev` checkout.
  This preparation task remains `planning` for root review and activation.

## Requirements

- R1: Remove only the unexported `extractSessionID` forwarding wrapper and its
  obsolete compatibility/deprecation comment from `selector.go`.
- R2: Update its existing test caller in `selector_test.go` to
  `ExtractSessionID(nil, []byte(tt.payload), nil)`, retaining every payload and
  expected result. Update the assertion label to match the public API.
- R3: Preserve current and concurrent edits; reread the two files before
  editing. Preserve public `ExtractSessionID`, other session helpers, and all
  SDK compatibility interfaces.
- R4: Keep the original affected files in the permission-700 artifact backup
  directory and retain exact verification commands, actual exit codes, raw
  logs, Git baseline, scope paths, and final diff in
  `/Users/jasonliao/Desktop/code/Artifacts/infra-dead-code-cleanup-20261003/cliproxy/`.

## Acceptance Criteria

- [ ] AC1 (R1): No definition or call of the exact private symbol
  `extractSessionID` remains; plural `extractSessionIDs` is preserved.
- [ ] AC2 (R2): All seven existing payload cases and expected identities remain
  unchanged and pass through `ExtractSessionID(nil, payload, nil)`.
- [ ] AC3 (R3): Product changes are limited to the two named Go files, plus
  one required documentation-only fork-patch row in `STRUCTURE.md`; public
  SDK behavior and interfaces remain unchanged.
- [ ] AC4 (R1-R3): Changed Go files are gofmt-clean; auth package tests,
  repository-required `go test ./...`, server compilation, and
  `git diff --check` succeed with captured exit status. Server output is
  written to Artifacts. Use explicit toolchain/cache environment and clear
  all lowercase/uppercase proxy variables without logging credentials.
- [ ] AC5 (R4): Backups, baseline, exact commands, raw test/build logs, exit
  codes, and final patch are saved. Any toolchain/dependency failure is
  diagnosed and reported as a verification gap instead of blindly retried.

## Out of scope

- Public SDK and tree compatibility deletion or deprecation.
- Adjacent dead-code candidates, service restart, deployment, fetch, reset,
  stage, commit, push, and machine-wide installation.
- Runtime `legacy/` directories or restoration of previously deleted content.

## Key decisions and limits

This is a lightweight PRD-only cleanup. Existing behavior tests provide the
regression coverage; no new features, abstractions, or tests are required.
Tests and compilation establish source verification only. Runtime activation
and provider connectivity are outside acceptance. No product, compatibility,
or scope decision remains unresolved; toolchain availability will be checked
during implementation validation.
