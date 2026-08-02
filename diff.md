# Local `lwj_dev` Differences from `main`

## Comparison state

- Repository: `router-for-me/CLIProxyAPI`
- Comparison baseline: local `main` and `origin/main` at
  `bc71c77f5cc42f3fbe1bf040cf14d4f166894835`
- Active branch: `lwj_dev`
- Commit `8d1abb420388b0067cd7a91b49dd23b35d314d74` is the existing local
  instruction-topology change on top of `main`.
- The branch commits the PR #4587 site patch and this document after that
  instruction-topology baseline.

## Macro summary

`lwj_dev` preserves this installation's local Codex/Claude instruction entry
topology and carries a narrow, two-file site patch from upstream PR #4587. The
site patch routes eligible streaming priority requests arriving through
HTTP/SSE (including Claude `/v1/messages`) onto the existing Codex upstream
WebSocket executor when the selected Codex auth explicitly enables
`websockets`.

No whole PR branch or unrelated upstream commits are merged. Local `main`
remains the upstream synchronization branch; future updates fetch and
fast-forward `main`, then merge `main` into `lwj_dev`.

## File-by-file differences

### `AGENTS.md`

- Committed in `8d1abb42`.
- Changed from the upstream regular file into a symbolic link to `CLAUDE.md`.
- Keeps Codex and Claude Code on one repository-local instruction source
  without a recursive `AGENTS.md`/`CLAUDE.md` reference.

### `CLAUDE.md`

- Committed in `8d1abb42`.
- Replaces the upstream one-line `@AGENTS.md` indirection with the preserved
  regular instruction document used by both entrypoints.

### `internal/runtime/executor/codex_websockets_executor.go`

- Committed exact source-file change from PR #4587
  (`0b4e8b54aa62ae8c9f0d72427d71cb06fd406cc3`).
- Adds a side-effect-free eligibility check that resolves translation,
  thinking settings, and payload rules before transport selection.
- Routes only streaming priority requests with `websockets=true` auths to
  `CodexWebsocketsExecutor`; standard requests, image requests, non-streaming
  execution, and existing downstream-WebSocket dispatch retain their prior
  paths.

### `internal/runtime/executor/codex_websockets_executor_priority_test.go`

- Committed exact test file from PR #4587
  (`0b4e8b54aa62ae8c9f0d72427d71cb06fd406cc3`).
- Covers direct priority payloads, fast-alias payload rules, Claude ingress,
  thinking-suffix rules, standard HTTP fallback, disabled WebSocket auths,
  downstream-WebSocket behavior, image exclusion, and non-streaming scope.

### `diff.md`

- Committed local documentation file.
- Records the committed and working-tree differences between `lwj_dev` and
  `main`, including this file itself.

## State outside the Git diff

Codex auth JSON state is runtime configuration outside this repository's
tracked diff. Setting `websockets: true`, its restricted backup, and any later
restore are operational state changes and are not represented as source
changes above.

The deployed binary under the workspace `.local/bin` directory, supervisor
PIDs, shadow instance state, benchmark evidence, and production rollout or
rollback are also repository-external runtime state. A successful source build
or working-tree diff does not imply that production has been published.

The supervisor control script under workspace `.config/cliproxyapi` is also
outside this Git repository. During rollout it was corrected to use exact tmux
session existence matching, preventing an acceptance session whose name starts
with `cliproxyapi-` from being mistaken for the production `cliproxyapi`
session. That operational script change is not part of the PR #4587 source
diff.
