// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"net/http"
	"testing"
)

func TestGetDEPDevice(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/dep/C02X9999", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code": 0,
			"placeholder": map[string]any{
				"id":              50,
				"userId":          10,
				"locationId":      1,
				"model":           "MacBook Air",
				"color":           "Silver",
				"serialNumber":    "C02X9999",
				"status":          "assigned",
				"deviceName":      "Lab Mac",
				"dateAdded":       1700000000,
				"datePushed":      1700100000,
				"profileName":     "Default",
				"placeholderName": "Placeholder 50",
			},
		})
	})

	dep, err := c.GetDEPDevice(context.Background(), "C02X9999")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dep.ID != 50 {
		t.Errorf("expected ID 50, got %d", dep.ID)
	}
	if dep.UserID != 10 {
		t.Errorf("expected userId 10, got %d", dep.UserID)
	}
	if dep.LocationID != 1 {
		t.Errorf("expected locationId 1, got %d", dep.LocationID)
	}
	if dep.Model != "MacBook Air" {
		t.Errorf("expected model %q, got %q", "MacBook Air", dep.Model)
	}
	if dep.Color != "Silver" {
		t.Errorf("expected color %q, got %q", "Silver", dep.Color)
	}
	if dep.SerialNumber != "C02X9999" {
		t.Errorf("expected serialNumber %q, got %q", "C02X9999", dep.SerialNumber)
	}
	if dep.Status != "assigned" {
		t.Errorf("expected status %q, got %q", "assigned", dep.Status)
	}
	if dep.DeviceName != "Lab Mac" {
		t.Errorf("expected deviceName %q, got %q", "Lab Mac", dep.DeviceName)
	}
	if dep.DateAdded != 1700000000 {
		t.Errorf("expected dateAdded 1700000000, got %d", dep.DateAdded)
	}
	if dep.DatePushed != 1700100000 {
		t.Errorf("expected datePushed 1700100000, got %d", dep.DatePushed)
	}
	if dep.ProfileName != "Default" {
		t.Errorf("expected profileName %q, got %q", "Default", dep.ProfileName)
	}
	if dep.PlaceholderName != "Placeholder 50" {
		t.Errorf("expected placeholderName %q, got %q", "Placeholder 50", dep.PlaceholderName)
	}
}

func TestGetDEPDevice_StringDates(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/dep/C02XSTR", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code": 0,
			"placeholder": map[string]any{
				"id":           60,
				"serialNumber": "C02XSTR",
				"dateAdded":    "1700000000",
				"datePushed":   "1700100000",
			},
		})
	})

	dep, err := c.GetDEPDevice(context.Background(), "C02XSTR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dep.DateAdded != 1700000000 {
		t.Errorf("expected dateAdded 1700000000, got %d", dep.DateAdded)
	}
	if dep.DatePushed != 1700100000 {
		t.Errorf("expected datePushed 1700100000, got %d", dep.DatePushed)
	}
}

func TestGetDEPDevice_NotFound(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/dep/NONEXISTENT", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`Not Found`))
	})

	_, err := c.GetDEPDevice(context.Background(), "NONEXISTENT")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetDEPDevices(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/dep", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code": 0,
			"placeholders": []map[string]any{
				{"id": 1, "serialNumber": "SN001"},
				{"id": 2, "serialNumber": "SN002"},
			},
		})
	})

	devices, err := c.GetDEPDevices(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(devices) != 2 {
		t.Fatalf("expected 2 DEP devices, got %d", len(devices))
	}
	if devices[0].SerialNumber != "SN001" {
		t.Errorf("expected first serial %q, got %q", "SN001", devices[0].SerialNumber)
	}
	if devices[1].SerialNumber != "SN002" {
		t.Errorf("expected second serial %q, got %q", "SN002", devices[1].SerialNumber)
	}
}
