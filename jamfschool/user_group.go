// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"fmt"
	"net/http"
)

// Group represents a Jamf School user group.
type Group struct {
	ID          int64  `json:"id,omitempty"`
	LocationID  int64  `json:"locationId,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	UserCount   int64  `json:"userCount,omitempty"`
	ACL         ACL    `json:"acl"`
	Modified    string `json:"modified,omitempty"`
}

// ACL represents the access control list for a group.
type ACL struct {
	SelfService              string `json:"selfService,omitempty"`
	SelfServiceInfo          string `json:"selfServiceInfo,omitempty"`
	SelfServiceLocation      string `json:"selfServiceLocation,omitempty"`
	SelfServiceClearPasscode string `json:"selfServiceClearPasscode,omitempty"`
	SelfServiceLock          string `json:"selfServiceLock,omitempty"`
	SelfServiceWipe          string `json:"selfServiceWipe,omitempty"`
	SelfServiceUnenroll      string `json:"selfServiceUnenroll,omitempty"`
	Teacher                  string `json:"teacher,omitempty"`
	Parent                   string `json:"parent,omitempty"`
}

// GroupCreateInput holds fields for creating a group.
type GroupCreateInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	LocationID  *int64 `json:"locationId,omitempty"`
	ACL         *ACL   `json:"acl,omitempty"`
}

// GroupUpdateInput holds fields for updating a group.
type GroupUpdateInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ACL         *ACL   `json:"acl,omitempty"`
}

type groupResponse struct {
	Code  int   `json:"code"`
	Group Group `json:"group"`
}

type groupsResponse struct {
	Code   int     `json:"code"`
	Count  int     `json:"count"`
	Groups []Group `json:"groups"`
}

// GetGroup retrieves a group by ID.
func (c *Client) GetGroup(ctx context.Context, id int64) (*Group, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/users/groups/%d", id), nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[groupResponse](resp)
	if err != nil {
		return nil, err
	}
	return &r.Group, nil
}

// GetGroups retrieves all groups.
func (c *Client) GetGroups(ctx context.Context) ([]Group, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, "/users/groups", nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[groupsResponse](resp)
	if err != nil {
		return nil, err
	}
	return r.Groups, nil
}

// CreateGroup creates a new group.
func (c *Client) CreateGroup(ctx context.Context, input GroupCreateInput) (int64, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodPost, "/users/groups", input)
	if err != nil {
		return 0, err
	}
	r, err := decode[createResponse](resp)
	if err != nil {
		return 0, err
	}
	return r.ID, nil
}

// UpdateGroup updates a group.
func (c *Client) UpdateGroup(ctx context.Context, id int64, input GroupUpdateInput) error {
	_, err := c.transport.DoRequest(ctx, http.MethodPut, fmt.Sprintf("/users/groups/%d", id), input)
	return err
}

// DeleteGroup deletes a group.
func (c *Client) DeleteGroup(ctx context.Context, id int64) error {
	_, err := c.transport.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/users/groups/%d", id), nil)
	return err
}
