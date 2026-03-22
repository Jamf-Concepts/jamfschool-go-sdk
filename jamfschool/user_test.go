// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGetUser(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/users/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code": 0,
			"user": map[string]any{
				"id":          1,
				"username":    "testuser",
				"email":       "test@example.com",
				"firstName":   "Test",
				"lastName":    "User",
				"status":      "Active",
				"deviceCount": 0,
				"groupIds":    []int{10, 20},
				"groups":      []string{"group-a", "group-b"},
			},
		})
	})

	result, err := c.GetUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}
	if result.Username != "testuser" {
		t.Errorf("expected username %q, got %q", "testuser", result.Username)
	}
	if result.Email != "test@example.com" {
		t.Errorf("expected email %q, got %q", "test@example.com", result.Email)
	}
	if result.FirstName != "Test" {
		t.Errorf("expected firstName %q, got %q", "Test", result.FirstName)
	}
	if result.LastName != "User" {
		t.Errorf("expected lastName %q, got %q", "User", result.LastName)
	}
	if result.Status != "Active" {
		t.Errorf("expected status %q, got %q", "Active", result.Status)
	}
	if result.DeviceCount != 0 {
		t.Errorf("expected deviceCount 0, got %d", result.DeviceCount)
	}
	if len(result.GroupIDs) != 2 || result.GroupIDs[0] != 10 || result.GroupIDs[1] != 20 {
		t.Errorf("expected groupIds [10,20], got %v", result.GroupIDs)
	}
	if len(result.Groups) != 2 || result.Groups[0] != "group-a" || result.Groups[1] != "group-b" {
		t.Errorf("expected groups [group-a,group-b], got %v", result.Groups)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/users/999", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`Not Found`))
	})

	_, err := c.GetUser(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetUsers(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":  0,
			"count": 2,
			"users": []map[string]any{
				{"id": 1, "username": "alice"},
				{"id": 2, "username": "bob"},
			},
		})
	})

	users, err := c.GetUsers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].Username != "alice" {
		t.Errorf("expected first user %q, got %q", "alice", users[0].Username)
	}
	if users[1].Username != "bob" {
		t.Errorf("expected second user %q, got %q", "bob", users[1].Username)
	}
}

func TestCreateUser(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
			"id":      42,
		})
	})

	id, err := c.CreateUser(context.Background(), UserCreateInput{
		Username:  "newuser",
		Password:  "secret",
		Email:     "new@example.com",
		FirstName: "New",
		LastName:  "User",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 42 {
		t.Errorf("expected ID 42, got %d", id)
	}
}

func TestCreateUser_WithMemberOf(t *testing.T) {
	t.Parallel()

	var gotBody map[string]any
	c, mux := testServer(t)
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
			"id":      43,
		})
	})

	id, err := c.CreateUser(context.Background(), UserCreateInput{
		Username:  "groupuser",
		Password:  "secret",
		Email:     "group@example.com",
		FirstName: "Group",
		LastName:  "User",
		MemberOf:  []any{int64(10), int64(20)},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 43 {
		t.Errorf("expected ID 43, got %d", id)
	}
	memberOf, ok := gotBody["memberOf"].([]any)
	if !ok {
		t.Fatalf("expected memberOf array in request body, got %v", gotBody["memberOf"])
	}
	if len(memberOf) != 2 {
		t.Errorf("expected 2 memberOf entries, got %d", len(memberOf))
	}
}

func TestUpdateUser_WithMemberOf(t *testing.T) {
	t.Parallel()

	var gotBody map[string]any
	c, mux := testServer(t)
	mux.HandleFunc("/api/users/1", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&gotBody)
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "UserDetailsSaved",
		})
	})

	memberOf := []any{int64(5)}
	err := c.UpdateUser(context.Background(), 1, UserUpdateInput{
		MemberOf: &memberOf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := gotBody["memberOf"].([]any)
	if !ok {
		t.Fatalf("expected memberOf array in request body, got %v", gotBody["memberOf"])
	}
	if len(got) != 1 {
		t.Errorf("expected 1 memberOf entry, got %d", len(got))
	}
}

func TestUpdateUser_WithEmptyMemberOf(t *testing.T) {
	t.Parallel()

	var rawBody []byte
	c, mux := testServer(t)
	mux.HandleFunc("/api/users/1", func(w http.ResponseWriter, r *http.Request) {
		rawBody, _ = io.ReadAll(r.Body)
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "UserDetailsSaved",
		})
	})

	memberOf := []any{}
	err := c.UpdateUser(context.Background(), 1, UserUpdateInput{
		MemberOf: &memberOf,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(rawBody), `"memberOf":[]`) {
		t.Errorf("expected request body to contain memberOf:[], got %s", string(rawBody))
	}
}

func TestUpdateUser_WithoutMemberOf(t *testing.T) {
	t.Parallel()

	var rawBody []byte
	c, mux := testServer(t)
	mux.HandleFunc("/api/users/1", func(w http.ResponseWriter, r *http.Request) {
		rawBody, _ = io.ReadAll(r.Body)
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "UserDetailsSaved",
		})
	})

	err := c.UpdateUser(context.Background(), 1, UserUpdateInput{
		FirstName: "Updated",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(string(rawBody), "memberOf") {
		t.Errorf("expected request body to NOT contain memberOf when nil, got %s", string(rawBody))
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/users/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
		})
	})

	err := c.UpdateUser(context.Background(), 1, UserUpdateInput{
		FirstName: "Updated",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/users/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
		})
	})

	err := c.DeleteUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMigrateUser(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/users/1/migrate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    200,
			"message": "UserMigrated",
		})
	})

	err := c.MigrateUser(context.Background(), 1, 5, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
