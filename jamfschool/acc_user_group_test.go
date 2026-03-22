// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

//go:build acceptance

package jamfschool_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
)

func TestAccGroupCRUD(t *testing.T) {
	c := accClient(t)
	ctx := context.Background()

	name := fmt.Sprintf("tf-sdk-acc-grp-%d", os.Getpid())

	// Create
	id, err := c.CreateGroup(ctx, jamfschool.GroupCreateInput{
		Name:        name,
		Description: "sdk acc test",
		ACL: &jamfschool.ACL{
			Teacher: "deny",
			Parent:  "inherit",
		},
	})
	if err != nil {
		t.Fatalf("CreateGroup: %v", err)
	}
	if id == 0 {
		t.Fatal("CreateGroup returned id 0")
	}
	t.Cleanup(func() { _ = c.DeleteGroup(ctx, id) })

	// Read
	group, err := c.GetGroup(ctx, id)
	if err != nil {
		t.Fatalf("GetGroup: %v", err)
	}
	if group.Name != name {
		t.Errorf("expected name %q, got %q", name, group.Name)
	}
	if group.Description != "sdk acc test" {
		t.Errorf("expected description %q, got %q", "sdk acc test", group.Description)
	}
	if group.ACL.Teacher != "deny" {
		t.Errorf("expected acl.teacher %q, got %q", "deny", group.ACL.Teacher)
	}

	// Update
	err = c.UpdateGroup(ctx, id, jamfschool.GroupUpdateInput{
		Name: name,
		ACL: &jamfschool.ACL{
			Teacher: "allow",
			Parent:  "deny",
		},
	})
	if err != nil {
		t.Fatalf("UpdateGroup: %v", err)
	}

	// Read back
	group, err = c.GetGroup(ctx, id)
	if err != nil {
		t.Fatalf("GetGroup after update: %v", err)
	}
	if group.ACL.Teacher != "allow" {
		t.Errorf("expected acl.teacher %q after update, got %q", "allow", group.ACL.Teacher)
	}
	if group.ACL.Parent != "deny" {
		t.Errorf("expected acl.parent %q after update, got %q", "deny", group.ACL.Parent)
	}

	// Delete
	err = c.DeleteGroup(ctx, id)
	if err != nil {
		t.Fatalf("DeleteGroup: %v", err)
	}
}
