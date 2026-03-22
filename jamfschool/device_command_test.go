// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"net/http"
	"testing"
)

func TestEraseDevice(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/UDID-123/wipe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "DeviceWipeScheduled",
		})
	})

	err := c.EraseDevice(context.Background(), "UDID-123", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEraseDevice_WithClearLock(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/UDID-123/wipe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "DeviceWipeScheduled",
		})
	})

	err := c.EraseDevice(context.Background(), "UDID-123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEraseDevice_NotFound(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/NONEXISTENT/wipe", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":404,"message":"DeviceNotFound"}`))
	})

	err := c.EraseDevice(context.Background(), "NONEXISTENT", false)
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestRestartDevice(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/UDID-123/restart", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "DeviceRestartScheduled",
		})
	})

	err := c.RestartDevice(context.Background(), "UDID-123", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRestartDevice_WithClearPasscode(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/UDID-123/restart", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "DeviceRestartScheduled",
		})
	})

	err := c.RestartDevice(context.Background(), "UDID-123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRestartDevice_NotFound(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/NONEXISTENT/restart", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":404,"message":"DeviceNotFound"}`))
	})

	err := c.RestartDevice(context.Background(), "NONEXISTENT", false)
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestRefreshDevice(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/UDID-123/refresh", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "CommandScheduled",
		})
	})

	err := c.RefreshDevice(context.Background(), "UDID-123", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRefreshDevice_WithClearErrors(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/UDID-123/refresh", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "CommandScheduled",
		})
	})

	err := c.RefreshDevice(context.Background(), "UDID-123", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnenrollDevice(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/UDID-123/unenroll", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "DeviceUnenrolled",
		})
	})

	err := c.UnenrollDevice(context.Background(), "UDID-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClearDeviceActivationLock(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/UDID-123/activationlock/clear", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "Unlocked",
		})
	})

	err := c.ClearDeviceActivationLock(context.Background(), "UDID-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTrashDevice(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/UDID-123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "DeviceTrashed",
		})
	})

	err := c.TrashDevice(context.Background(), "UDID-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRestoreDevice(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/UDID-123/restore", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "DeviceRestored",
		})
	})

	err := c.RestoreDevice(context.Background(), "UDID-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateDeviceESIM(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/devices/UDID-123/cellularPlan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "RefreshCellularPlanScheduled",
		})
	})

	err := c.UpdateDeviceESIM(context.Background(), "UDID-123", "https://carrier.example.com", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
