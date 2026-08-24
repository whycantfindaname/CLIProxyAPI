---
name: error-handling
description: Error propagation, classification, API responses, and cleanup rules
paths:
  - "cmd/**/*.go"
  - "internal/**/*.go"
  - "sdk/**/*.go"
  - "test/**/*.go"
---

# Error Handling

## Propagate Context Without Losing Identity

Return errors from libraries and helpers. Wrap them with operation context and
`%w`, using method-specific variable names when needed:

```go
result, errLoad := loadConfig(path)
if errLoad != nil {
    return nil, fmt.Errorf("load config %s: %w", path, errLoad)
}
```

This keeps `errors.Is` and `errors.As` usable. The pattern is widespread in
`internal/store/gitstore.go`; `internal/clienterror/client_error_test.go`
specifically verifies wrapped cancellation, deadlines, and status-bearing
errors.

Do not call `log.Fatal` or `log.Fatalf`: they terminate the process and bypass
normal cleanup. Command entry points may log and return from `main`; packages
must return errors to their caller.

## Classify Before Retrying or Rotating Credentials

Use `internal/clienterror` for upstream HTTP classification.
`HTTPStatusFromError` gives explicit `StatusCode()` errors precedence, maps
`context.Canceled` to 499, and maps `context.DeadlineExceeded` to 504.
`IsRequestFault` distinguishes client input faults from credential, quota, and
transport failures so invalid input does not rotate or penalize credentials.

Do not classify by status alone when the upstream body carries a supported
structured code/type. Conversely, preserve 401, 402, and 429 as credential,
payment, or rate-limit failures even if the body also says
`invalid_request_error`. Add table-driven cases to
`internal/clienterror/client_error_test.go` whenever this taxonomy changes.

## HTTP Error Contracts

Protocol handlers under `sdk/api/handlers/` own protocol-shaped responses.
For OpenAI-compatible endpoints use `handlers.ErrorResponse` and
`handlers.BuildErrorResponseBody` from `sdk/api/handlers/handlers.go`; valid
upstream JSON is preserved, while plain text is wrapped with the appropriate
type and code.

Management handlers use their existing management contract, such as
`c.AbortWithStatusJSON(status, gin.H{"error": message})` in
`internal/api/handlers/management/handler.go`. Do not replace one API family's
wire shape with another's shared struct.

Once a streaming response has started, an error cannot be treated like a
normal pre-header JSON response. Follow the SSE/WebSocket termination pattern
in the matching protocol handler and executor, and test both pre-stream and
post-stream failures.

## Panic and Cleanup Rules

HTTP handlers should return meaningful status codes rather than panic.
`internal/logging.GinLogrusRecovery` is the last-resort boundary: it logs a
stack and returns 500, while deliberately re-panicking `http.ErrAbortHandler`
so `net/http` can abort the connection. Keep both behaviors covered by
`internal/logging/gin_logger_test.go`.

Handle cleanup failures when they affect correctness or diagnostics:

```go
defer func() {
    if errClose := resource.Close(); errClose != nil {
        log.WithError(errClose).Warn("failed to close resource")
    }
}()
```

If cleanup failure must be returned, preserve the primary failure and join or
wrap the cleanup error, as recovery paths in `internal/store/gitstore.go` do.

## Timeout Boundary

Timeouts are allowed during credential acquisition. After an upstream
connection is established, do not add request deadlines. The intentional
exceptions are listed in `AGENTS.md`: Codex WebSocket liveness, wsrelay session
deadlines, the management API call tool, and the Antigravity model-fetch
utility.

## Anti-Patterns

- Do not compare error strings when a sentinel, wrapped cause, status interface,
  or structured upstream body exists.
- Do not discard a wrapped cause with `%v` when callers need `errors.Is/As`.
- Do not rotate credentials for a known client request fault.
- Do not expose internal stack traces, auth data, or storage paths in client
  responses.
- Do not write a second OpenAI error envelope or use the OpenAI envelope for
  management endpoints.
- Do not add blanket HTTP timeouts after connection establishment.

## Verification

```bash
env -u http_proxy -u https_proxy -u HTTP_PROXY -u HTTPS_PROXY go test ./internal/clienterror
env -u http_proxy -u https_proxy -u HTTP_PROXY -u HTTPS_PROXY go test ./internal/logging -run 'TestGinLogrusRecovery'
env -u http_proxy -u https_proxy -u HTTP_PROXY -u HTTPS_PROXY go test ./sdk/api/handlers/...
rg -n "log\.Fatal|log\.Fatalf|fmt\.Errorf" cmd internal sdk
```
