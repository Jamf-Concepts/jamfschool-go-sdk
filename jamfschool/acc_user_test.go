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

func TestAccUserCRUD(t *testing.T) {
	c := accClient(t)
	ctx := context.Background()

	// Create
	input := jamfschool.UserCreateInput{
		Username:  fmt.Sprintf("tf-sdk-acc-%d", os.Getpid()),
		Password:  "TestPass123!",
		Email:     fmt.Sprintf("tf-sdk-acc-%d@example.com", os.Getpid()),
		FirstName: "SDK",
		LastName:  "Test",
	}
	id, err := c.CreateUser(ctx, input)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if id == 0 {
		t.Fatal("CreateUser returned id 0")
	}
	t.Cleanup(func() { _ = c.DeleteUser(ctx, id) })

	// Read
	user, err := c.GetUser(ctx, id)
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if user.Username != input.Username {
		t.Errorf("expected username %q, got %q", input.Username, user.Username)
	}
	if user.Email != input.Email {
		t.Errorf("expected email %q, got %q", input.Email, user.Email)
	}
	if user.FirstName != input.FirstName {
		t.Errorf("expected firstName %q, got %q", input.FirstName, user.FirstName)
	}
	if user.LastName != input.LastName {
		t.Errorf("expected lastName %q, got %q", input.LastName, user.LastName)
	}

	// Update
	err = c.UpdateUser(ctx, id, jamfschool.UserUpdateInput{
		Email: "tf-sdk-acc-updated@example.com",
	})
	if err != nil {
		t.Fatalf("UpdateUser: %v", err)
	}

	// Read back
	user, err = c.GetUser(ctx, id)
	if err != nil {
		t.Fatalf("GetUser after update: %v", err)
	}
	if user.Email != "tf-sdk-acc-updated@example.com" {
		t.Errorf("expected updated email, got %q", user.Email)
	}
	if user.Username != input.Username {
		t.Errorf("expected username unchanged %q, got %q", input.Username, user.Username)
	}

	// Delete
	err = c.DeleteUser(ctx, id)
	if err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
}
