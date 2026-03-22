// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"net/http"
	"testing"
)

func TestGetGroup(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/users/groups/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":0,"group":{"id":1,"name":"Test Group","description":"desc","locationId":0,"userCount":5,"acl":{"selfService":"allow","selfServiceInfo":"allow","selfServiceLocation":"deny","selfServiceClearPasscode":"deny","selfServiceLock":"deny","selfServiceWipe":"deny","selfServiceUnenroll":"deny","teacher":"allow","parent":"deny"}}}`))
	})

	group, err := c.GetGroup(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if group.ID != 1 {
		t.Errorf("expected ID 1, got %d", group.ID)
	}
	if group.Name != "Test Group" {
		t.Errorf("expected name %q, got %q", "Test Group", group.Name)
	}
	if group.Description != "desc" {
		t.Errorf("expected description %q, got %q", "desc", group.Description)
	}
	if group.UserCount != 5 {
		t.Errorf("expected userCount 5, got %d", group.UserCount)
	}
	if group.ACL.SelfService != "allow" {
		t.Errorf("expected acl.selfService %q, got %q", "allow", group.ACL.SelfService)
	}
	if group.ACL.SelfServiceInfo != "allow" {
		t.Errorf("expected acl.selfServiceInfo %q, got %q", "allow", group.ACL.SelfServiceInfo)
	}
	if group.ACL.SelfServiceLocation != "deny" {
		t.Errorf("expected acl.selfServiceLocation %q, got %q", "deny", group.ACL.SelfServiceLocation)
	}
	if group.ACL.SelfServiceClearPasscode != "deny" {
		t.Errorf("expected acl.selfServiceClearPasscode %q, got %q", "deny", group.ACL.SelfServiceClearPasscode)
	}
	if group.ACL.SelfServiceLock != "deny" {
		t.Errorf("expected acl.selfServiceLock %q, got %q", "deny", group.ACL.SelfServiceLock)
	}
	if group.ACL.SelfServiceWipe != "deny" {
		t.Errorf("expected acl.selfServiceWipe %q, got %q", "deny", group.ACL.SelfServiceWipe)
	}
	if group.ACL.SelfServiceUnenroll != "deny" {
		t.Errorf("expected acl.selfServiceUnenroll %q, got %q", "deny", group.ACL.SelfServiceUnenroll)
	}
	if group.ACL.Teacher != "allow" {
		t.Errorf("expected acl.teacher %q, got %q", "allow", group.ACL.Teacher)
	}
	if group.ACL.Parent != "deny" {
		t.Errorf("expected acl.parent %q, got %q", "deny", group.ACL.Parent)
	}
}

func TestGetGroup_NotFound(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/users/groups/999", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`Not Found`))
	})

	_, err := c.GetGroup(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
}

func TestGetGroups(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/users/groups", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":  0,
			"count": 2,
			"groups": []map[string]any{
				{"id": 1, "name": "Group A"},
				{"id": 2, "name": "Group B"},
			},
		})
	})

	groups, err := c.GetGroups(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Name != "Group A" {
		t.Errorf("expected first group name %q, got %q", "Group A", groups[0].Name)
	}
	if groups[1].Name != "Group B" {
		t.Errorf("expected second group name %q, got %q", "Group B", groups[1].Name)
	}
}

func TestCreateGroup(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/users/groups", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
			"id":      10,
		})
	})

	id, err := c.CreateGroup(context.Background(), GroupCreateInput{
		Name:        "New Group",
		Description: "A new group",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 10 {
		t.Errorf("expected ID 10, got %d", id)
	}
}

func TestCreateGroup_WithACL(t *testing.T) {
	t.Parallel()

	var gotBody map[string]any
	c, mux := testServer(t)
	mux.HandleFunc("/api/users/groups", func(w http.ResponseWriter, r *http.Request) {
		readJSON(t, r, &gotBody)
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
			"id":      11,
		})
	})

	id, err := c.CreateGroup(context.Background(), GroupCreateInput{
		Name: "ACL Group",
		ACL: &ACL{
			Teacher: "allow",
			Parent:  "deny",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 11 {
		t.Errorf("expected ID 11, got %d", id)
	}
	acl, ok := gotBody["acl"].(map[string]any)
	if !ok {
		t.Fatalf("expected acl object in request body, got %v", gotBody["acl"])
	}
	if acl["teacher"] != "allow" {
		t.Errorf("expected acl.teacher %q, got %v", "allow", acl["teacher"])
	}
}

func TestUpdateGroup(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/users/groups/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
		})
	})

	err := c.UpdateGroup(context.Background(), 1, GroupUpdateInput{
		Name: "Updated Group",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteGroup(t *testing.T) {
	t.Parallel()

	c, mux := testServer(t)
	mux.HandleFunc("/api/users/groups/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		writeJSON(t, w, http.StatusOK, map[string]any{
			"code":    0,
			"message": "ok",
		})
	})

	err := c.DeleteGroup(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
