# CLIProxyAPI Managed Sync Errors

## CLIPROXYAPI_SOURCE_VERIFY_FAILED

- Stage: project workflow.
- Meaning: `go test ./...` failed or timed out.
- Inspect: identify the first failing package and test, then verify the explicit
  Go and proxy environment used by the run.
- Action: fix the owning source or test and rerun the complete command.
- Stop condition: do not project configuration or restart the service until the
  test suite passes.
