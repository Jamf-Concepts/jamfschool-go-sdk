// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package client

import "errors"

// Sentinel errors returned by the client.
var (
	ErrAuthentication = errors.New("jamfschool: authentication failed")
	ErrNotFound       = errors.New("jamfschool: resource not found")
	ErrHTTP           = errors.New("jamfschool: HTTP request failed")
)
