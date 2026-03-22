// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"net/http"
	"testing"
)

func TestGetClass(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/classes/abc-123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code": 0,
			"class": map[string]any{
				"uuid":         "abc-123",
				"name":         "Math 101",
				"description":  "Intro to math",
				"locationId":   1,
				"source":       "manual",
				"studentCount": 25,
				"teacherCount": 2,
			},
		})
	})

	class, err := c.GetClass(context.Background(), "abc-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if class.UUID != "abc-123" {
		t.Errorf("expected UUID %q, got %q", "abc-123", class.UUID)
	}
	if class.Name != "Math 101" {
		t.Errorf("expected name %q, got %q", "Math 101", class.Name)
	}
	if class.Description != "Intro to math" {
		t.Errorf("expected description %q, got %q", "Intro to math", class.Description)
	}
	if class.LocationID != 1 {
		t.Errorf("expected locationId 1, got %d", class.LocationID)
	}
	if class.Source != "manual" {
		t.Errorf("expected source %q, got %q", "manual", class.Source)
	}
	if class.StudentCount != 25 {
		t.Errorf("expected studentCount 25, got %d", class.StudentCount)
	}
	if class.TeacherCount != 2 {
		t.Errorf("expected teacherCount 2, got %d", class.TeacherCount)
	}
}

func TestGetClass_NotFound(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/classes/nonexistent", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`Not Found`))
	})

	_, err := c.GetClass(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetClasses(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/classes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code": 0,
			"classes": []map[string]any{
				{"uuid": "u1", "name": "Class A"},
				{"uuid": "u2", "name": "Class B"},
			},
		})
	})

	classes, err := c.GetClasses(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(classes) != 2 {
		t.Fatalf("expected 2 classes, got %d", len(classes))
	}
	if classes[0].UUID != "u1" {
		t.Errorf("expected first UUID %q, got %q", "u1", classes[0].UUID)
	}
	if classes[1].Name != "Class B" {
		t.Errorf("expected second name %q, got %q", "Class B", classes[1].Name)
	}
}

func TestCreateClass(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/classes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
			"uuid":    "abc-123",
		})
	})

	uuid, err := c.CreateClass(context.Background(), ClassCreateInput{
		Name:        "New Class",
		Description: "A new class",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if uuid != "abc-123" {
		t.Errorf("expected UUID %q, got %q", "abc-123", uuid)
	}
}

func TestCreateClass_WithStudents(t *testing.T) {
	t.Parallel()

	var gotBody map[string]any
	c, mux := testServer(t)
	mux.HandleFunc("/api/classes", func(w http.ResponseWriter, r *http.Request) {
		readJSON(t, r, &gotBody)
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
			"uuid":    "def-456",
		})
	})

	uuid, err := c.CreateClass(context.Background(), ClassCreateInput{
		Name:     "Class With Students",
		Students: []int64{1, 2},
		Teachers: []int64{3},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if uuid != "def-456" {
		t.Errorf("expected UUID %q, got %q", "def-456", uuid)
	}
	students, ok := gotBody["students"].([]any)
	if !ok {
		t.Fatalf("expected students array in request body, got %v", gotBody["students"])
	}
	if len(students) != 2 {
		t.Errorf("expected 2 students, got %d", len(students))
	}
}

func TestUpdateClass(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/classes/abc-123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
		})
	})

	err := c.UpdateClass(context.Background(), "abc-123", ClassUpdateInput{
		Name: "Updated Class",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteClass(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/classes/abc-123", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
		})
	})

	err := c.DeleteClass(context.Background(), "abc-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAssignClassUsers(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/classes/test-uuid/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "ClassSaved",
		})
	})

	err := c.AssignClassUsers(context.Background(), "test-uuid", []int64{1, 2}, []int64{3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetClassDevices(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/classes/abc-123/devices", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":  200,
			"count": 2,
			"devices": []map[string]any{
				{"UDID": "UDID-001", "serialNumber": "SN001", "name": "iPad 1"},
				{"UDID": "UDID-002", "serialNumber": "SN002", "name": "iPad 2"},
			},
		})
	})

	devices, err := c.GetClassDevices(context.Background(), "abc-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(devices) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(devices))
	}
	if devices[0].UDID != "UDID-001" {
		t.Errorf("expected first UDID %q, got %q", "UDID-001", devices[0].UDID)
	}
	if devices[0].SerialNumber != "SN001" {
		t.Errorf("expected first serial %q, got %q", "SN001", devices[0].SerialNumber)
	}
	if devices[1].Name != "iPad 2" {
		t.Errorf("expected second name %q, got %q", "iPad 2", devices[1].Name)
	}
}

func TestGetClassDevices_Empty(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/classes/empty-class/devices", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":  200,
			"count": 0,
		})
	})

	devices, err := c.GetClassDevices(context.Background(), "empty-class")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(devices) != 0 {
		t.Errorf("expected 0 devices, got %d", len(devices))
	}
}
