// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"net/http"
	"testing"
)

func TestGetProfile(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/profiles/7", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"id":          7,
			"locationId":  1,
			"identifier":  "com.example.profile",
			"name":        "WiFi Profile",
			"description": "Corporate WiFi",
			"platform":    "iOS",
		})
	})

	profile, err := c.GetProfile(context.Background(), 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profile.ID != 7 {
		t.Errorf("expected ID 7, got %d", profile.ID)
	}
	if profile.LocationID != 1 {
		t.Errorf("expected locationId 1, got %d", profile.LocationID)
	}
	if profile.Identifier != "com.example.profile" {
		t.Errorf("expected identifier %q, got %q", "com.example.profile", profile.Identifier)
	}
	if profile.Name != "WiFi Profile" {
		t.Errorf("expected name %q, got %q", "WiFi Profile", profile.Name)
	}
	if profile.Description != "Corporate WiFi" {
		t.Errorf("expected description %q, got %q", "Corporate WiFi", profile.Description)
	}
	if profile.Platform != "iOS" {
		t.Errorf("expected platform %q, got %q", "iOS", profile.Platform)
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/profiles/999", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`Not Found`))
	})

	_, err := c.GetProfile(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetProfiles(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/profiles", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"profiles": []map[string]any{
				{"id": 1, "name": "Profile A"},
				{"id": 2, "name": "Profile B"},
			},
		})
	})

	profiles, err := c.GetProfiles(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(profiles))
	}
	if profiles[0].Name != "Profile A" {
		t.Errorf("expected first profile name %q, got %q", "Profile A", profiles[0].Name)
	}
	if profiles[1].Name != "Profile B" {
		t.Errorf("expected second profile name %q, got %q", "Profile B", profiles[1].Name)
	}
}
