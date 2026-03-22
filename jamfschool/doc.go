// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

// Package jamfschool provides a Go client for the Jamf School API.
//
// Create a client with [NewClient] and use the typed methods to manage
// Jamf School resources such as users, groups, classes, devices, and iBeacons.
//
//	c := jamfschool.NewClient(
//		"https://myschool.jamfcloud.com",
//		os.Getenv("JAMFSCHOOL_NETWORK_ID"),
//		os.Getenv("JAMFSCHOOL_API_KEY"),
//	)
//
//	users, err := c.GetUsers(ctx)
//
// The client authenticates via HTTP Basic Auth using the Network ID and API Key.
//
// Sentinel errors [ErrAuthentication], [ErrNotFound], and [ErrHTTP]
// can be used with [errors.Is] for error handling:
//
//	user, err := c.GetUser(ctx, id)
//	if errors.Is(err, jamfschool.ErrNotFound) {
//		// handle missing resource
//	}
package jamfschool
