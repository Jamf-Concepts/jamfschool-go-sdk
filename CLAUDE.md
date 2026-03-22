# Repository Guidelines

## Overview

This is the Go SDK for the [Jamf School](https://www.jamf.com/products/jamf-school/) REST API. The module path is `github.com/Jamf-Concepts/jamfschool-go-sdk`.

## Architecture

The SDK follows a two-layer architecture:

### Layer 1: `internal/client/` — HTTP Transport

- Handles HTTP Basic Auth authentication
- JSON request/response marshaling
- Error handling with sentinel errors
- Logging interface (optional)
- Per-request options (`WithProtocolVersion`, `WithQuery`)
- **No resource-specific types** belong here

### Layer 2: `jamfschool/` — Public API

- All resource types (User, Group, DeviceGroup, Class, IBeacon, Device, App, Profile, Location, DEPDevice)
- Methods call `c.transport.DoRequest()` directly
- One file per API domain
- Re-exports sentinel errors from internal/client
- Options pattern for client configuration (`WithLogger`, `WithHTTPClient`, `WithUserAgent`)

## Tooling

- Go >= 1.26
- Use `make` for build, lint, and test targets.

### Available make targets

| Target    | Description                                         |
| --------- | --------------------------------------------------- |
| `test`    | Run unit tests                                      |
| `testacc` | Run acceptance tests (requires environment variables)|
| `lint`    | Run golangci-lint                                   |

## Code Style Guidelines

- Follow Go conventions and idiomatic patterns.
- Favor clear and descriptive naming for variables, functions, and types.
- Always ensure constants, functions, variable sets and types have a short comment describing their purpose.
- Do not add comments inside type definitions or function bodies.

## Client Pattern

```go
c := jamfschool.NewClient(baseURL, networkID, apiKey,
    jamfschool.WithLogger(logger),
    jamfschool.WithUserAgent("my-app/1.0"),
)
user, err := c.GetUser(ctx, 1)
```

### Authentication

- HTTP Basic Auth — Network ID as username, API Key as password
- Headers: `Content-Type: application/json`, `Accept: application/json`, `X-Server-Protocol-Version: 3`

### Error Handling

Three sentinel errors:
- `ErrAuthentication` — HTTP 401/403
- `ErrNotFound` — HTTP 404
- `ErrHTTP` — HTTP 4xx/5xx (excluding 401/403/404)

Use `errors.Is(err, jamfschool.ErrNotFound)` to check.

### Per-Request Options (internal/client)

- `WithProtocolVersion(v)` — Override protocol version header (e.g. "4" for app endpoints)
- `WithQuery(q)` — Append raw query string to URL

### Logger Interface

```go
type Logger interface {
    LogRequest(ctx context.Context, method, url string, body []byte)
    LogResponse(ctx context.Context, statusCode int, headers http.Header, body []byte)
}
```

## API Resource Pattern

Each resource has a file per API domain with:

- Types: exported structs for API responses and inputs
- Methods: on `*Client` receiver, calling `c.transport.DoRequest()`
- Response wrappers: unexported structs for JSON unmarshaling

### Naming Convention

- Response types: `TypeName` (e.g. `User`, `Group`, `DeviceGroup`)
- Create inputs: `TypeNameCreateInput`
- Update inputs: `TypeNameUpdateInput`

### Error Wrapping Convention

```go
return nil, fmt.Errorf("decoding user response: %w", err)
```

## Testing

### Unit Tests

- Per-domain test files (`user_test.go`, `device_group_test.go`, etc.)
- Use `testServer(t)` helper returning `(*Client, *http.ServeMux)`
- Register handlers on mux with `/api/` prefix
- Use `writeJSON(t, w, status, v)` and `readJSON(t, r, v)` helpers

### Acceptance Tests

- Build tag: `//go:build acceptance`
- Require `JAMFSCHOOL_URL`, `JAMFSCHOOL_NETWORK_ID`, `JAMFSCHOOL_API_KEY` environment variables
- Use `accClient(t)` helper
- Use `t.Cleanup()` for resource teardown
- Run with: `make testacc`

## Environment Variables

- `JAMFSCHOOL_URL` — Base URL of the Jamf School instance
- `JAMFSCHOOL_NETWORK_ID` — Network ID for HTTP Basic Auth
- `JAMFSCHOOL_API_KEY` — API key for HTTP Basic Auth
