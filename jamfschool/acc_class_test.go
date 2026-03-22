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

func TestAccClassCRUD(t *testing.T) {
	c := accClient(t)
	ctx := context.Background()

	name := fmt.Sprintf("tf-sdk-acc-cls-%d", os.Getpid())

	// Create
	uuid, err := c.CreateClass(ctx, jamfschool.ClassCreateInput{
		Name:        name,
		Description: "sdk acc test",
	})
	if err != nil {
		t.Fatalf("CreateClass: %v", err)
	}
	if uuid == "" {
		t.Fatal("CreateClass returned empty uuid")
	}
	t.Cleanup(func() { _ = c.DeleteClass(ctx, uuid) })

	// Read
	class, err := c.GetClass(ctx, uuid)
	if err != nil {
		t.Fatalf("GetClass: %v", err)
	}
	if class.Name != name {
		t.Errorf("expected name %q, got %q", name, class.Name)
	}
	if class.UUID != uuid {
		t.Errorf("expected uuid %q, got %q", uuid, class.UUID)
	}
	if class.Description != "sdk acc test" {
		t.Errorf("expected description %q, got %q", "sdk acc test", class.Description)
	}
	if class.Source == "" {
		t.Error("expected non-empty source")
	}

	// Update
	err = c.UpdateClass(ctx, uuid, jamfschool.ClassUpdateInput{
		Description: "updated desc",
	})
	if err != nil {
		t.Fatalf("UpdateClass: %v", err)
	}

	// Read back
	class, err = c.GetClass(ctx, uuid)
	if err != nil {
		t.Fatalf("GetClass after update: %v", err)
	}
	if class.Description != "updated desc" {
		t.Errorf("expected description %q, got %q", "updated desc", class.Description)
	}

	// Delete
	err = c.DeleteClass(ctx, uuid)
	if err != nil {
		t.Fatalf("DeleteClass: %v", err)
	}
}
