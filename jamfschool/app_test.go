// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"net/http"
	"testing"
)

func TestCreateApp(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/apps", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("X-Server-Protocol-Version") != "4" {
			t.Errorf("expected protocol version 4, got %s", r.Header.Get("X-Server-Protocol-Version"))
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "AppCreated",
			"data": map[string]any{
				"mediaId":    42,
				"locationId": 0,
			},
		})
	})

	id, err := c.CreateApp(context.Background(), AppCreateInput{
		AdamID:      361309726,
		CountryCode: "us",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 42 {
		t.Errorf("expected id 42, got %d", id)
	}
}

func TestCreateApp_InvalidJSON(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/apps", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{not valid json`))
	})

	_, err := c.CreateApp(context.Background(), AppCreateInput{
		AdamID:      12345,
		CountryCode: "us",
	})
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestTrashApp(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/apps/42/trash", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("X-Server-Protocol-Version") != "4" {
			t.Errorf("expected protocol version 4, got %s", r.Header.Get("X-Server-Protocol-Version"))
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "AppTrashed",
		})
	})

	err := c.TrashApp(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetApp(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/apps/10", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"id":       10,
			"bundleId": "com.example.app",
			"adamId":   12345,
			"name":     "Example App",
			"vendor":   "Example Inc",
			"version":  "2.1.0",
			"platform": "iOS",
		})
	})

	app, err := c.GetApp(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if app.ID != 10 {
		t.Errorf("expected ID 10, got %d", app.ID)
	}
	if app.BundleID != "com.example.app" {
		t.Errorf("expected bundleId %q, got %q", "com.example.app", app.BundleID)
	}
	if app.AdamID != 12345 {
		t.Errorf("expected adamId 12345, got %d", app.AdamID)
	}
	if app.Name != "Example App" {
		t.Errorf("expected name %q, got %q", "Example App", app.Name)
	}
	if app.Vendor != "Example Inc" {
		t.Errorf("expected vendor %q, got %q", "Example Inc", app.Vendor)
	}
	if app.Version != "2.1.0" {
		t.Errorf("expected version %q, got %q", "2.1.0", app.Version)
	}
	if app.Platform != "iOS" {
		t.Errorf("expected platform %q, got %q", "iOS", app.Platform)
	}
}

func TestGetApp_NotFound(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/apps/999", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`Not Found`))
	})

	_, err := c.GetApp(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetApps(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/apps", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"apps": []map[string]any{
				{"id": 1, "name": "App A"},
				{"id": 2, "name": "App B"},
			},
		})
	})

	apps, err := c.GetApps(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(apps) != 2 {
		t.Fatalf("expected 2 apps, got %d", len(apps))
	}
	if apps[0].Name != "App A" {
		t.Errorf("expected first app name %q, got %q", "App A", apps[0].Name)
	}
	if apps[1].Name != "App B" {
		t.Errorf("expected second app name %q, got %q", "App B", apps[1].Name)
	}
}
