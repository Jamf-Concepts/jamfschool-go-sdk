// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"github.com/Jamf-Concepts/jamfschool-go-sdk/internal/client"
)

// ErrAuthentication is returned when the API responds with 401 Unauthorized or 403 Forbidden.
var ErrAuthentication = client.ErrAuthentication

// ErrNotFound is returned when the API responds with 404 Not Found.
var ErrNotFound = client.ErrNotFound

// ErrHTTP is returned for other HTTP error responses (4xx/5xx excluding 401, 403, 404).
var ErrHTTP = client.ErrHTTP
