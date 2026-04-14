// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"net/http"
	"testing"
)

func TestGetDevice(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/UDID-001", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"code": 0,
			"device": {
				"UDID": "UDID-001",
				"serialNumber": "C02X1234",
				"name": "Test MacBook",
				"isManaged": true,
				"isSupervised": true,
				"batteryLevel": 0.85,
				"totalCapacity": 500.0,
				"notes": "Lab device",
				"model": {"name": "MacBook Pro", "identifier": "MacBookPro16,1", "type": "computer"},
				"os": {"prefix": "macOS", "version": "14.3"}
			}
		}`))
	})

	device, err := c.GetDevice(context.Background(), "UDID-001")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if device.UDID != "UDID-001" {
		t.Errorf("expected UDID %q, got %q", "UDID-001", device.UDID)
	}
	if device.SerialNumber != "C02X1234" {
		t.Errorf("expected serialNumber %q, got %q", "C02X1234", device.SerialNumber)
	}
	if device.Name != "Test MacBook" {
		t.Errorf("expected name %q, got %q", "Test MacBook", device.Name)
	}
	if !device.IsManaged {
		t.Error("expected isManaged true, got false")
	}
	if !device.IsSupervised {
		t.Error("expected isSupervised true, got false")
	}
	if device.Notes != "Lab device" {
		t.Errorf("expected notes %q, got %q", "Lab device", device.Notes)
	}
	if device.Model.Name != "MacBook Pro" {
		t.Errorf("expected model name %q, got %q", "MacBook Pro", device.Model.Name)
	}
	if device.Model.Identifier != "MacBookPro16,1" {
		t.Errorf("expected model identifier %q, got %q", "MacBookPro16,1", device.Model.Identifier)
	}
	if device.Model.Type != "computer" {
		t.Errorf("expected model type %q, got %q", "computer", device.Model.Type)
	}
	if device.OS.Prefix != "macOS" {
		t.Errorf("expected OS prefix %q, got %q", "macOS", device.OS.Prefix)
	}
	if device.OS.Version != "14.3" {
		t.Errorf("expected OS version %q, got %q", "14.3", device.OS.Version)
	}
}

func TestGetDevice_ObjectFields(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/UDID-002", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// The single-device endpoint returns some fields as objects
		// instead of strings: model.type, lastCheckin, deviceEnrollType, etc.
		_, _ = w.Write([]byte(`{
			"code": 0,
			"device": {
				"UDID": "UDID-002",
				"name": "iPad",
				"lastCheckin": {"date": "2026-03-18 20:48:54", "epoch": 1774064934},
				"deviceEnrollType": {"value": "dep"},
				"assetTag": {"value": "ASSET-001"},
				"notes": {"text": "Lab device"},
				"model": {"name": "iPad Air", "identifier": "iPad13,2", "type": {"name": "iPad", "type": "iOS"}},
				"os": {"prefix": "iPadOS", "version": "17.4"}
			}
		}`))
	})

	device, err := c.GetDevice(context.Background(), "UDID-002")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if device.Model.Name != "iPad Air" {
		t.Errorf("expected model name %q, got %q", "iPad Air", device.Model.Name)
	}
	if device.Model.Type != "iOS" {
		t.Errorf("expected model type %q, got %q", "iOS", device.Model.Type)
	}
	if device.LastCheckin != "2026-03-18 20:48:54" {
		t.Errorf("expected lastCheckin %q, got %q", "2026-03-18 20:48:54", device.LastCheckin)
	}
	if device.DeviceEnrollType != "dep" {
		t.Errorf("expected deviceEnrollType %q, got %q", "dep", device.DeviceEnrollType)
	}
	if device.AssetTag != "ASSET-001" {
		t.Errorf("expected assetTag %q, got %q", "ASSET-001", device.AssetTag)
	}
	if device.Notes != "Lab device" {
		t.Errorf("expected notes %q, got %q", "Lab device", device.Notes)
	}
}

func TestGetDevice_NotFound(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/NONEXISTENT", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`Not Found`))
	})

	_, err := c.GetDevice(context.Background(), "NONEXISTENT")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetDevices(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":  0,
			"count": 2,
			"devices": []map[string]any{
				{"UDID": "UDID-001", "name": "iPad 1"},
				{"UDID": "UDID-002", "name": "iPad 2"},
			},
		})
	})

	devices, err := c.GetDevices(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(devices) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(devices))
	}
	if devices[0].UDID != "UDID-001" {
		t.Errorf("expected first UDID %q, got %q", "UDID-001", devices[0].UDID)
	}
	if devices[1].Name != "iPad 2" {
		t.Errorf("expected second name %q, got %q", "iPad 2", devices[1].Name)
	}
}
