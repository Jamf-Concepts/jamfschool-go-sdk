// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"net/http"
	"testing"
)

func TestGetDeviceGroup(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/groups/5", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"deviceGroup":{"id":5,"name":"Lab Macs","description":"Computer lab","members":12,"isSmartGroup":false,"shared":true,"type":"computers"}}`))
	})

	dg, err := c.GetDeviceGroup(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dg.ID != 5 {
		t.Errorf("expected ID 5, got %d", dg.ID)
	}
	if dg.Name != "Lab Macs" {
		t.Errorf("expected name %q, got %q", "Lab Macs", dg.Name)
	}
	if dg.Description != "Computer lab" {
		t.Errorf("expected description %q, got %q", "Computer lab", dg.Description)
	}
	if dg.Members != 12 {
		t.Errorf("expected members 12, got %d", dg.Members)
	}
	if dg.IsSmartGroup {
		t.Error("expected isSmartGroup false, got true")
	}
	if !dg.Shared {
		t.Error("expected shared true, got false")
	}
	if dg.Type != "computers" {
		t.Errorf("expected type %q, got %q", "computers", dg.Type)
	}
}

func TestGetDeviceGroup_NotFound(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/groups/999", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`Not Found`))
	})

	_, err := c.GetDeviceGroup(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetDeviceGroups(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/groups", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"DeviceGroups":[{"id":1,"name":"iPads"},{"id":2,"name":"Macs"}]}`))
	})

	groups, err := c.GetDeviceGroups(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 2 {
		t.Fatalf("expected 2 device groups, got %d", len(groups))
	}
	if groups[0].Name != "iPads" {
		t.Errorf("expected first group %q, got %q", "iPads", groups[0].Name)
	}
	if groups[1].Name != "Macs" {
		t.Errorf("expected second group %q, got %q", "Macs", groups[1].Name)
	}
}

func TestCreateDeviceGroup(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/groups", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
			"id":      77,
		})
	})

	id, err := c.CreateDeviceGroup(context.Background(), DeviceGroupCreateInput{
		Name:        "New Device Group",
		Description: "A device group",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 77 {
		t.Errorf("expected ID 77, got %d", id)
	}
}

func TestCreateDeviceGroup_WithOptions(t *testing.T) {
	t.Parallel()

	var gotBody map[string]any
	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/groups", func(w http.ResponseWriter, r *http.Request) {
		readJSON(t, r, &gotBody)
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
			"id":      78,
		})
	})

	id, err := c.CreateDeviceGroup(context.Background(), DeviceGroupCreateInput{
		Name:           "Shared Group",
		CollectionType: "list",
		Shared:         true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 78 {
		t.Errorf("expected ID 78, got %d", id)
	}
	if gotBody["collectionType"] != "list" {
		t.Errorf("expected collectionType %q, got %v", "list", gotBody["collectionType"])
	}
}

func TestUpdateDeviceGroup(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/groups/5", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
		})
	})

	err := c.UpdateDeviceGroup(context.Background(), 5, DeviceGroupUpdateInput{
		Name: "Renamed Group",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteDeviceGroup(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/groups/5", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
		})
	})

	err := c.DeleteDeviceGroup(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddDevicesToGroup(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/groups/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":         200,
			"devicesAdded": 1,
		})
	})

	err := c.AddDevicesToGroup(context.Background(), 1, []string{"UDID-001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemoveDevicesFromGroup(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/groups/remove", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":           200,
			"devicesRemoved": 1,
		})
	})

	err := c.RemoveDevicesFromGroup(context.Background(), 1, []string{"UDID-001"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetDeviceGroupMembers(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("groups") != "1" {
			t.Errorf("expected query groups=1, got %s", r.URL.RawQuery)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":  200,
			"count": 2,
			"devices": []map[string]any{
				{"UDID": "UDID-001"},
				{"UDID": "UDID-002"},
			},
		})
	})

	udids, err := c.GetDeviceGroupMembers(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(udids) != 2 {
		t.Fatalf("expected 2 UDIDs, got %d", len(udids))
	}
	if udids[0] != "UDID-001" {
		t.Errorf("expected UDID-001, got %s", udids[0])
	}
	if udids[1] != "UDID-002" {
		t.Errorf("expected UDID-002, got %s", udids[1])
	}
}
