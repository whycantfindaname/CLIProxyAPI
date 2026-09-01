# CLIProxyAPI Managed Sync

This directory is the single in-repository entry for the CLIProxyAPI managed
workflow. The initial contract targets `lwj_dev`; Agent Infra owns repository
convergence, and the project workflow runs the complete Go test suite.

```bash
go test ./...
```

Configuration projection, auth material, service restart, and real-provider
acceptance remain platform-owned. The initial contract does not modify config
or runtime state. See [errors.md](errors.md) after a failure.
