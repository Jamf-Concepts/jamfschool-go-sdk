// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testServer creates an httptest.Server and returns a Client pointed at it.
// Tests register additional handlers on the returned mux.
func testServer(t *testing.T) (*Client, *http.ServeMux) {
	t.Helper()
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	c := NewClient(srv.URL, "test-net", "test-key")
	return c, mux
}

// writeJSON is a test helper that writes a JSON response with the given status code.
func writeJSON(t *testing.T, w http.ResponseWriter, status int, v any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			t.Fatalf("writeJSON: %v", err)
		}
	}
}

// readJSON is a test helper that decodes a JSON request body.
func readJSON(t *testing.T, r *http.Request, v any) {
	t.Helper()
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		t.Fatalf("readJSON: %v", err)
	}
}
