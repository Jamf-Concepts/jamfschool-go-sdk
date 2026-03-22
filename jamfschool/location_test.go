// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"net/http"
	"testing"
)

func TestGetLocation(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/locations/3", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":3,"name":"Main Campus","isDistrict":false,"street":"123 Main St","streetNumber":"123","postalCode":"12345","city":"Springfield","source":"manual","asmIdentifier":"ASM-001","schoolNumber":"SCH001"}`))
	})

	loc, err := c.GetLocation(context.Background(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.ID != 3 {
		t.Errorf("expected ID 3, got %d", loc.ID)
	}
	if loc.Name != "Main Campus" {
		t.Errorf("expected name %q, got %q", "Main Campus", loc.Name)
	}
	if loc.IsDistrict {
		t.Error("expected isDistrict false, got true")
	}
	if loc.Source != "manual" {
		t.Errorf("expected source %q, got %q", "manual", loc.Source)
	}
	if loc.SchoolNumber != "SCH001" {
		t.Errorf("expected schoolNumber %q, got %q", "SCH001", loc.SchoolNumber)
	}
	if loc.Street == nil || *loc.Street != "123 Main St" {
		t.Errorf("expected street %q, got %v", "123 Main St", loc.Street)
	}
	if loc.StreetNumber == nil || *loc.StreetNumber != "123" {
		t.Errorf("expected streetNumber %q, got %v", "123", loc.StreetNumber)
	}
	if loc.PostalCode == nil || *loc.PostalCode != "12345" {
		t.Errorf("expected postalCode %q, got %v", "12345", loc.PostalCode)
	}
	if loc.City == nil || *loc.City != "Springfield" {
		t.Errorf("expected city %q, got %v", "Springfield", loc.City)
	}
	if loc.ASMIdentifier == nil || *loc.ASMIdentifier != "ASM-001" {
		t.Errorf("expected asmIdentifier %q, got %v", "ASM-001", loc.ASMIdentifier)
	}
}

func TestGetLocation_NotFound(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/locations/999", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`Not Found`))
	})

	_, err := c.GetLocation(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetLocations(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/locations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"locations": []map[string]any{
				{"id": 1, "name": "Location A"},
				{"id": 2, "name": "Location B"},
			},
		})
	})

	locations, err := c.GetLocations(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(locations) != 2 {
		t.Fatalf("expected 2 locations, got %d", len(locations))
	}
	if locations[0].Name != "Location A" {
		t.Errorf("expected first location name %q, got %q", "Location A", locations[0].Name)
	}
	if locations[1].Name != "Location B" {
		t.Errorf("expected second location name %q, got %q", "Location B", locations[1].Name)
	}
}
