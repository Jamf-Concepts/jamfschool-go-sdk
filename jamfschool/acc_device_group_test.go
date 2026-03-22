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

func TestAccDeviceGroupCRUD(t *testing.T) {
	c := accClient(t)
	ctx := context.Background()

	name := fmt.Sprintf("tf-sdk-acc-dg-%d", os.Getpid())

	// Create
	id, err := c.CreateDeviceGroup(ctx, jamfschool.DeviceGroupCreateInput{
		Name:        name,
		Description: "sdk acc test",
	})
	if err != nil {
		t.Fatalf("CreateDeviceGroup: %v", err)
	}
	if id == 0 {
		t.Fatal("CreateDeviceGroup returned id 0")
	}
	t.Cleanup(func() { _ = c.DeleteDeviceGroup(ctx, id) })

	// Read
	dg, err := c.GetDeviceGroup(ctx, id)
	if err != nil {
		t.Fatalf("GetDeviceGroup: %v", err)
	}
	if dg.Name != name {
		t.Errorf("expected name %q, got %q", name, dg.Name)
	}
	if dg.Description != "sdk acc test" {
		t.Errorf("expected description %q, got %q", "sdk acc test", dg.Description)
	}
	if dg.IsSmartGroup {
		t.Error("expected IsSmartGroup false")
	}

	// Update
	err = c.UpdateDeviceGroup(ctx, id, jamfschool.DeviceGroupUpdateInput{
		Description: "updated desc",
	})
	if err != nil {
		t.Fatalf("UpdateDeviceGroup: %v", err)
	}

	// Read back
	dg, err = c.GetDeviceGroup(ctx, id)
	if err != nil {
		t.Fatalf("GetDeviceGroup after update: %v", err)
	}
	if dg.Description != "updated desc" {
		t.Errorf("expected description %q, got %q", "updated desc", dg.Description)
	}

	// Delete
	err = c.DeleteDeviceGroup(ctx, id)
	if err != nil {
		t.Fatalf("DeleteDeviceGroup: %v", err)
	}
}
