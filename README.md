# jamfschool-go-sdk

A Go client library for the [Jamf School](https://www.jamf.com/products/jamf-school/) API.

## Installation

```sh
go get github.com/Jamf-Concepts/jamfschool-go-sdk
```

Requires Go 1.26 or later.

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
)

func main() {
	client := jamfschool.NewClient(
		"https://myschool.jamfcloud.com",
		os.Getenv("JAMFSCHOOL_NETWORK_ID"),
		os.Getenv("JAMFSCHOOL_API_KEY"),
	)

	ctx := context.Background()

	users, err := client.GetUsers(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d users\n", len(users))
}
```

## Authentication

The client uses HTTP Basic Auth (Network ID + API Key) to authenticate with the Jamf School API.

- **Network ID** — found at Devices > Enroll Device(s) in Jamf School.
- **API Key** — generated at Organization > Settings > API in Jamf School.

## Configuration Options

```go
// Custom user agent
client := jamfschool.NewClient(baseURL, networkID, apiKey,
	jamfschool.WithUserAgent("my-app/1.0.0"),
)

// Custom HTTP client
client := jamfschool.NewClient(baseURL, networkID, apiKey,
	jamfschool.WithHTTPClient(myHTTPClient),
)

// Enable request/response logging
client := jamfschool.NewClient(baseURL, networkID, apiKey,
	jamfschool.WithLogger(myLogger),
)
```

## Error Handling

The SDK provides sentinel errors for common failure cases:

```go
import "errors"

user, err := client.GetUser(ctx, id)
if errors.Is(err, jamfschool.ErrNotFound) {
	// Resource does not exist
}
if errors.Is(err, jamfschool.ErrAuthentication) {
	// Invalid credentials
}
if errors.Is(err, jamfschool.ErrHTTP) {
	// Other HTTP error (4xx/5xx)
}
```

## API Documentation

Full API reference is available on [pkg.go.dev](https://pkg.go.dev/github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool).

## License

[MIT](LICENSE)
