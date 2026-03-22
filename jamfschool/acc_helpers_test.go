// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

//go:build acceptance

package jamfschool_test

import (
	"os"
	"testing"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
)

// accClient creates a real Client backed by the live API.
// It skips the test if the required environment variables are not set.
func accClient(t *testing.T) *jamfschool.Client {
	t.Helper()

	url := os.Getenv("JAMFSCHOOL_URL")
	networkID := os.Getenv("JAMFSCHOOL_NETWORK_ID")
	apiKey := os.Getenv("JAMFSCHOOL_API_KEY")

	if url == "" || networkID == "" || apiKey == "" {
		t.Skip("Set JAMFSCHOOL_URL, JAMFSCHOOL_NETWORK_ID and JAMFSCHOOL_API_KEY to run acceptance tests")
	}

	return jamfschool.NewClient(url, networkID, apiKey)
}
