// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"fmt"
	"net/http"
)

// User represents a Jamf School user.
type User struct {
	ID          int64    `json:"id,omitempty"`
	LocationID  int64    `json:"locationId,omitempty"`
	Status      string   `json:"status,omitempty"`
	DeviceCount int64    `json:"deviceCount,omitempty"`
	Email       string   `json:"email,omitempty"`
	Username    string   `json:"username,omitempty"`
	Domain      string   `json:"domain,omitempty"`
	FirstName   string   `json:"firstName,omitempty"`
	LastName    string   `json:"lastName,omitempty"`
	GroupIDs    []int64  `json:"groupIds,omitempty"`
	Groups      []string `json:"groups,omitempty"`
	Notes       string   `json:"notes,omitempty"`
	Exclude     bool     `json:"exclude,omitempty"`
	InTrash     bool     `json:"inTrash,omitempty"`
}

// UserCreateInput holds fields for creating a user.
type UserCreateInput struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Domain        string `json:"domain,omitempty"`
	Notes         string `json:"notes,omitempty"`
	Exclude       bool   `json:"exclude,omitempty"`
	StorePassword bool   `json:"storePassword,omitempty"`
	LocationID    *int64 `json:"locationId,omitempty"`
	MemberOf      []any  `json:"memberOf,omitempty"`
}

// UserUpdateInput holds fields for updating a user.
type UserUpdateInput struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Email         string `json:"email,omitempty"`
	FirstName     string `json:"firstName,omitempty"`
	LastName      string `json:"lastName,omitempty"`
	Domain        string `json:"domain,omitempty"`
	Notes         string `json:"notes,omitempty"`
	Exclude       *bool  `json:"exclude,omitempty"`
	StorePassword *bool  `json:"storePassword,omitempty"`
	MemberOf      *[]any `json:"memberOf,omitempty"`
}

type userResponse struct {
	Code int  `json:"code"`
	User User `json:"user"`
}

type usersResponse struct {
	Code  int    `json:"code"`
	Count int    `json:"count"`
	Users []User `json:"users"`
}

// GetUser retrieves a user by ID.
func (c *Client) GetUser(ctx context.Context, id int64) (*User, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/users/%d", id), nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[userResponse](resp)
	if err != nil {
		return nil, err
	}
	return &r.User, nil
}

// GetUsers retrieves all users.
func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, "/users", nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[usersResponse](resp)
	if err != nil {
		return nil, err
	}
	return r.Users, nil
}

// CreateUser creates a new user and returns the ID.
func (c *Client) CreateUser(ctx context.Context, input UserCreateInput) (int64, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodPost, "/users", input)
	if err != nil {
		return 0, err
	}
	r, err := decode[createResponse](resp)
	if err != nil {
		return 0, err
	}
	return r.ID, nil
}

// UpdateUser updates an existing user.
func (c *Client) UpdateUser(ctx context.Context, id int64, input UserUpdateInput) error {
	_, err := c.transport.DoRequest(ctx, http.MethodPut, fmt.Sprintf("/users/%d", id), input)
	return err
}

// MigrateUser moves a user to a different location.
func (c *Client) MigrateUser(ctx context.Context, id, locationID int64, onlyUser bool) error {
	body := map[string]any{
		"locationId": locationID,
		"onlyUser":   onlyUser,
	}
	_, err := c.transport.DoRequest(ctx, http.MethodPut, fmt.Sprintf("/users/%d/migrate", id), body)
	return err
}

// DeleteUser deletes a user by ID.
func (c *Client) DeleteUser(ctx context.Context, id int64) error {
	_, err := c.transport.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/users/%d", id), nil)
	return err
}
