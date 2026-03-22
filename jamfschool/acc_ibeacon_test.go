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

func TestAccIBeaconCRUD(t *testing.T) {
	c := accClient(t)
	ctx := context.Background()

	name := fmt.Sprintf("tf-sdk-acc-beacon-%d", os.Getpid())
	major := int64(100)
	minor := int64(200)

	// Create
	beacon, err := c.CreateIBeacon(ctx, jamfschool.IBeaconCreateInput{
		Name:        name,
		Description: "sdk acc test",
		Major:       &major,
		Minor:       &minor,
	})
	if err != nil {
		t.Fatalf("CreateIBeacon: %v", err)
	}
	if beacon.ID == 0 {
		t.Fatal("CreateIBeacon returned id 0")
	}
	if beacon.Name != name {
		t.Errorf("expected name %q in create response, got %q", name, beacon.Name)
	}
	t.Cleanup(func() { _ = c.DeleteIBeacon(ctx, beacon.ID) })

	// Read
	b, err := c.GetIBeacon(ctx, beacon.ID)
	if err != nil {
		t.Fatalf("GetIBeacon: %v", err)
	}
	if b.Name != name {
		t.Errorf("expected name %q, got %q", name, b.Name)
	}
	if b.Major != 100 {
		t.Errorf("expected major 100, got %d", b.Major)
	}
	if b.Minor != 200 {
		t.Errorf("expected minor 200, got %d", b.Minor)
	}
	if b.UUID == "" {
		t.Error("expected non-empty UUID")
	}

	// Update
	newMajor := int64(300)
	updated, err := c.UpdateIBeacon(ctx, beacon.ID, jamfschool.IBeaconUpdateInput{
		Name:  name,
		Major: &newMajor,
	})
	if err != nil {
		t.Fatalf("UpdateIBeacon: %v", err)
	}
	if updated.Major != 300 {
		t.Errorf("expected major 300 in update response, got %d", updated.Major)
	}

	// Read back
	b, err = c.GetIBeacon(ctx, beacon.ID)
	if err != nil {
		t.Fatalf("GetIBeacon after update: %v", err)
	}
	if b.Major != 300 {
		t.Errorf("expected major 300, got %d", b.Major)
	}
	if b.Minor != 200 {
		t.Errorf("expected minor 200 unchanged, got %d", b.Minor)
	}

	// Delete
	err = c.DeleteIBeacon(ctx, beacon.ID)
	if err != nil {
		t.Fatalf("DeleteIBeacon: %v", err)
	}
}
