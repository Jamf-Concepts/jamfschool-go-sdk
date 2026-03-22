# Contributing

Thank you for your interest in contributing to the Jamf School Go SDK.

## Prerequisites

- **Go** >= 1.26 (see `go.mod` for the exact version)
- **golangci-lint** for linting
- A Jamf School tenant with API credentials (for acceptance tests only)

## Getting Started

```bash
# Clone the repository
git clone https://github.com/Jamf-Concepts/jamfschool-go-sdk.git
cd jamfschool-go-sdk

# Build and lint
make

# Run unit tests
make test
```

## Development Workflow

1. Create a feature branch from `main`.
2. Make your changes.
3. Run tests and linting before committing:

   ```bash
   make fmt
   make lint
   make test
   ```

4. Open a pull request against `main`.

## Make Targets

| Target     | Description                                           |
| ---------- | ----------------------------------------------------- |
| `build`    | Build all packages                                    |
| `fmt`      | Format Go source files                                |
| `lint`     | Run golangci-lint                                     |
| `vet`      | Run go vet                                            |
| `test`     | Run unit tests                                        |
| `testacc`  | Run acceptance tests (requires environment variables) |
| `generate` | Run go generate                                       |

The default target runs: `fmt lint build`.

## Project Structure

| Directory          | Purpose                                              |
| ------------------ | ---------------------------------------------------- |
| `jamfschool/`      | Exported SDK package — typed client methods and types |
| `internal/client/` | Transport layer — HTTP Basic Auth, error handling     |

## Adding a New Resource

1. Add types and client methods in `jamfschool/<resource>.go`.
2. Add unit tests in `jamfschool/<resource>_test.go` using `testServer(t)`.
3. Add an acceptance test in `jamfschool/acc_<resource>_test.go` following the existing CRUD pattern.
4. Run tests and linting.

## Running Acceptance Tests

Acceptance tests run against a live Jamf School tenant. Set the following environment variables:

```bash
export JAMFSCHOOL_URL="https://myschool.jamfcloud.com"
export JAMFSCHOOL_NETWORK_ID="your-network-id"
export JAMFSCHOOL_API_KEY="your-api-key"
```

Then run:

```bash
make testacc
```

Or manually:

```bash
JAMFSCHOOL_ACC=1 go test -v -cover -count=1 -tags acceptance -timeout 120m -p=1 ./...
```

## Dependencies

This project uses native Go packages only. Do not introduce third-party dependencies without discussion.

## Commit Messages

Use [conventional commit](https://www.conventionalcommits.org/) style messages:

- `feat: add device_group resource support`
- `fix: handle nil response in GetUser`
- `test: add acceptance tests for iBeacon CRUD`
- `refactor: extract shared response parsing`
- `docs: update README with new usage examples`

## Pull Requests

- Keep PRs focused — one feature or fix per PR.
- Include unit tests for new code.
- Include acceptance tests for new resources.
- Linting must pass before merge.

## Reporting Issues

Open an issue on GitHub with:

- SDK version (Go module version or commit SHA).
- Relevant code snippet (redact credentials).
- Expected vs actual behaviour.
- Any error messages or logs.
