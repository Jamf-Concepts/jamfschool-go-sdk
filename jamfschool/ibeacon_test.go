// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetIBeacon(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/ibeacons/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"beacon":{"id":1,"UUID":"beacon-uuid","major":100,"minor":200,"name":"Test","description":"desc"}}`))
	})

	beacon, err := c.GetIBeacon(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if beacon.ID != 1 {
		t.Errorf("expected ID 1, got %d", beacon.ID)
	}
	if beacon.UUID != "beacon-uuid" {
		t.Errorf("expected UUID %q, got %q", "beacon-uuid", beacon.UUID)
	}
	if beacon.Major != 100 {
		t.Errorf("expected major 100, got %d", beacon.Major)
	}
	if beacon.Minor != 200 {
		t.Errorf("expected minor 200, got %d", beacon.Minor)
	}
	if beacon.Name != "Test" {
		t.Errorf("expected name %q, got %q", "Test", beacon.Name)
	}
	if beacon.Description != "desc" {
		t.Errorf("expected description %q, got %q", "desc", beacon.Description)
	}
}

func TestGetIBeacon_StringMajorMinor(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/ibeacons/2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"beacon":{"id":2,"UUID":"str-uuid","major":"100","minor":"200","name":"StringBeacon","description":"string values"}}`))
	})

	beacon, err := c.GetIBeacon(context.Background(), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if beacon.Major != 100 {
		t.Errorf("expected major 100 (parsed from string), got %d", beacon.Major)
	}
	if beacon.Minor != 200 {
		t.Errorf("expected minor 200 (parsed from string), got %d", beacon.Minor)
	}
}

func TestGetIBeacon_NotFound(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/ibeacons/999", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`Not Found`))
	})

	_, err := c.GetIBeacon(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetIBeacons(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/ibeacons", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"beacons":[{"id":1,"name":"Beacon A","major":10,"minor":20},{"id":2,"name":"Beacon B","major":30,"minor":40}]}`))
	})

	beacons, err := c.GetIBeacons(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(beacons) != 2 {
		t.Fatalf("expected 2 beacons, got %d", len(beacons))
	}
	if beacons[0].Name != "Beacon A" {
		t.Errorf("expected first beacon name %q, got %q", "Beacon A", beacons[0].Name)
	}
	if beacons[1].Name != "Beacon B" {
		t.Errorf("expected second beacon name %q, got %q", "Beacon B", beacons[1].Name)
	}
}

func TestCreateIBeacon(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/ibeacons", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"ok","beacon":{"id":99,"UUID":"new-uuid","major":10,"minor":20,"name":"New Beacon","description":"new"}}`))
	})

	major := int64(10)
	minor := int64(20)
	beacon, err := c.CreateIBeacon(context.Background(), IBeaconCreateInput{
		Name:  "New Beacon",
		UUID:  "new-uuid",
		Major: &major,
		Minor: &minor,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if beacon.ID != 99 {
		t.Errorf("expected ID 99, got %d", beacon.ID)
	}
	if beacon.Name != "New Beacon" {
		t.Errorf("expected name %q, got %q", "New Beacon", beacon.Name)
	}
}

func TestCreateIBeacon_InvalidJSON(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/ibeacons", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{not valid json`))
	})

	_, err := c.CreateIBeacon(context.Background(), IBeaconCreateInput{
		Name: "Test",
	})
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestUpdateIBeacon(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/ibeacons/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"ok","beacon":{"id":1,"UUID":"updated-uuid","major":50,"minor":60,"name":"Updated Beacon","description":"updated"}}`))
	})

	beacon, err := c.UpdateIBeacon(context.Background(), 1, IBeaconUpdateInput{
		Name: "Updated Beacon",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if beacon.ID != 1 {
		t.Errorf("expected ID 1, got %d", beacon.ID)
	}
	if beacon.Name != "Updated Beacon" {
		t.Errorf("expected name %q, got %q", "Updated Beacon", beacon.Name)
	}
	if beacon.Major != 50 {
		t.Errorf("expected major 50, got %d", beacon.Major)
	}
}

func TestDeleteIBeacon(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/ibeacons/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"ok"}`))
	})

	err := c.DeleteIBeacon(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// IBeacon UnmarshalJSON (custom unmarshaler)
// ---------------------------------------------------------------------------

func TestIBeaconUnmarshalJSON_NumericValues(t *testing.T) {
	t.Parallel()

	data := []byte(`{"id":1,"UUID":"u","major":100,"minor":200,"name":"n","description":"d"}`)
	var b IBeacon
	if err := json.Unmarshal(data, &b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.ID != 1 {
		t.Errorf("expected ID 1, got %d", b.ID)
	}
	if b.UUID != "u" {
		t.Errorf("expected UUID %q, got %q", "u", b.UUID)
	}
	if b.Major != 100 {
		t.Errorf("expected major 100, got %d", b.Major)
	}
	if b.Minor != 200 {
		t.Errorf("expected minor 200, got %d", b.Minor)
	}
	if b.Name != "n" {
		t.Errorf("expected name %q, got %q", "n", b.Name)
	}
	if b.Description != "d" {
		t.Errorf("expected description %q, got %q", "d", b.Description)
	}
}

func TestIBeaconUnmarshalJSON_StringValues(t *testing.T) {
	t.Parallel()

	data := []byte(`{"id":1,"UUID":"u","major":"100","minor":"200","name":"n","description":"d"}`)
	var b IBeacon
	if err := json.Unmarshal(data, &b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Major != 100 {
		t.Errorf("expected major 100, got %d", b.Major)
	}
	if b.Minor != 200 {
		t.Errorf("expected minor 200, got %d", b.Minor)
	}
}

func TestIBeaconUnmarshalJSON_ZeroValues(t *testing.T) {
	t.Parallel()

	data := []byte(`{"id":1,"name":"n"}`)
	var b IBeacon
	if err := json.Unmarshal(data, &b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Major != 0 {
		t.Errorf("expected major 0, got %d", b.Major)
	}
	if b.Minor != 0 {
		t.Errorf("expected minor 0, got %d", b.Minor)
	}
}

func TestIBeaconUnmarshalJSON_EmptyStringValues(t *testing.T) {
	t.Parallel()

	data := []byte(`{"id":1,"major":"","minor":"","name":"n"}`)
	var b IBeacon
	if err := json.Unmarshal(data, &b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Major != 0 {
		t.Errorf("expected major 0 for empty string, got %d", b.Major)
	}
	if b.Minor != 0 {
		t.Errorf("expected minor 0 for empty string, got %d", b.Minor)
	}
}

// ---------------------------------------------------------------------------
// parseFlexibleInt64
// ---------------------------------------------------------------------------

func TestParseFlexibleInt64_Number(t *testing.T) {
	t.Parallel()

	result := parseFlexibleInt64(json.RawMessage("42"))
	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}

func TestParseFlexibleInt64_String(t *testing.T) {
	t.Parallel()

	result := parseFlexibleInt64(json.RawMessage(`"42"`))
	if result != 42 {
		t.Errorf("expected 42, got %d", result)
	}
}

func TestParseFlexibleInt64_EmptyString(t *testing.T) {
	t.Parallel()

	result := parseFlexibleInt64(json.RawMessage(`""`))
	if result != 0 {
		t.Errorf("expected 0 for empty string, got %d", result)
	}
}

func TestParseFlexibleInt64_Null(t *testing.T) {
	t.Parallel()

	result := parseFlexibleInt64(json.RawMessage("null"))
	if result != 0 {
		t.Errorf("expected 0 for null, got %d", result)
	}
}
