---
name: logging
description: CLIProxyAPI logrus, Gin, request ID, and sensitive-data conventions
paths:
  - "cmd/**/*.go"
  - "internal/**/*.go"
  - "sdk/**/*.go"
---

# Logging

## Shared Logging Stack

Use the repository-wide logrus logger (`log` alias) instead of creating a new
logger. `internal/logging/global_logger.go` configures caller reporting, the
custom formatter, stdout or rotating `main.log`, and Gin writers.
`internal/api/server.go` installs `GinLogrusLogger`, `GinLogrusRecovery`, and
CPA trace middleware before API routes.

The formatter emits timestamp, request ID, level, caller, message, and an
allowlisted set of structured fields. Add a new formatter field only when it
is useful across logs, then cover its rendered and redacted form in
`internal/logging/global_logger_test.go`.

## Levels and Context

- `Debug`: routing decisions, normalization details, and diagnostics useful
  only while investigating behavior. `internal/thinking/apply.go` uses fields
  such as provider and model for this purpose.
- `Info`: startup mode, enabled backend, lifecycle transition, or successful
  operation that an operator needs to understand.
- `Warn`: recoverable degradation, ignored optional input, cleanup failure, or
  4xx request outcome.
- `Error`: failed initialization, operation failure that changes the result,
  recovered panic, or 5xx request outcome.

Use `WithError` and `WithField(s)` for stable dimensions:

```go
log.WithFields(log.Fields{
    "provider": provider,
    "model":    model,
}).Debug("thinking provider selected")
```

Do not interpolate stable dimensions into a long message when they already
have an approved structured field.

## Request Correlation

`internal/logging/gin_logger.go` creates an eight-character request ID for AI
API path groups (`/v1`, `/v1beta`, `/openai/v1`, and `/backend-api/codex`) and
stores it in both Gin and request contexts through
`internal/logging/requestid.go`. Pass the request context through handlers,
auth selection, translators, and executors so downstream logs remain
correlated. Do not generate a second request ID inside a lower layer.

Gin request logs derive level from the final status: 5xx is error, 4xx is warn,
and success is info. `internal/api/middleware/request_logging.go` owns detailed
request/response capture and error-only body spooling; endpoint handlers should
not duplicate raw request logging.

## Sensitive Data

Mask query strings through existing utilities and avoid raw request bodies,
authorization headers, provider tokens, API keys, OAuth payloads, auth file
contents, management passwords, and storage credentials. A credential label or
non-secret ID is acceptable when it is needed for correlation.

The formatter deliberately omits generic path fields unless a plugin ID makes
them plugin metadata; `TestLogFormatterOmitsGenericPathField` protects against
leaking auth paths. Preserve that boundary.

## Anti-Patterns

- Do not use `fmt.Printf`, the standard `log` package, or a package-local
  logrus instance for service diagnostics.
- Do not log the same request body in handlers and request middleware.
- Do not put secrets into structured fields just because the formatter omits
  unknown keys; hooks or alternate formatters may still observe them.
- Do not log an expected client cancellation as an internal server failure.
- Do not use `Fatal` or `Panic` for recoverable runtime errors.
- Do not add high-volume per-chunk info logs to streaming paths.

## Verification

```bash
env -u http_proxy -u https_proxy -u HTTP_PROXY -u HTTPS_PROXY go test ./internal/logging ./internal/api/middleware
rg -n "fmt\.Print|log\.Fatal|log\.Panic|Authorization|api[_-]?key|access[_-]?token" cmd internal sdk
```

Review search matches in context; protocol constants and redaction tests are
valid, while emitted raw credential values are not.
