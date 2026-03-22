// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/internal/client"
)

// DeviceGroup represents a Jamf School device group.
type DeviceGroup struct {
	ID           int64  `json:"id,omitempty"`
	LocationID   int64  `json:"locationId,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Information  string `json:"information,omitempty"`
	IsSmartGroup bool   `json:"isSmartGroup,omitempty"`
	Members      int64  `json:"members,omitempty"`
	Shared       bool   `json:"shared,omitempty"`
	ImageURL     string `json:"imageUrl,omitempty"`
	Type         string `json:"type,omitempty"`
}

// DeviceGroupCreateInput holds fields for creating a device group.
type DeviceGroupCreateInput struct {
	Name           string `json:"name"`
	LocationID     *int64 `json:"locationId,omitempty"`
	Description    string `json:"description,omitempty"`
	Information    string `json:"information,omitempty"`
	CollectionType string `json:"collectionType,omitempty"`
	Shared         bool   `json:"shared,omitempty"`
}

// DeviceGroupUpdateInput holds fields for updating a device group.
type DeviceGroupUpdateInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Shared      *bool  `json:"shared,omitempty"`
}

type deviceGroupResponse struct {
	Code        int         `json:"code"`
	DeviceGroup DeviceGroup `json:"deviceGroup"`
}

type deviceGroupsResponse struct {
	Code         int           `json:"code"`
	DeviceGroups []DeviceGroup `json:"DeviceGroups"`
}

// GetDeviceGroup retrieves a device group by ID.
func (c *Client) GetDeviceGroup(ctx context.Context, id int64) (*DeviceGroup, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/devices/groups/%d", id), nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[deviceGroupResponse](resp)
	if err != nil {
		return nil, err
	}
	return &r.DeviceGroup, nil
}

// GetDeviceGroups retrieves all device groups.
func (c *Client) GetDeviceGroups(ctx context.Context) ([]DeviceGroup, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, "/devices/groups", nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[deviceGroupsResponse](resp)
	if err != nil {
		return nil, err
	}
	return r.DeviceGroups, nil
}

// CreateDeviceGroup creates a new device group.
func (c *Client) CreateDeviceGroup(ctx context.Context, input DeviceGroupCreateInput) (int64, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodPost, "/devices/groups", input)
	if err != nil {
		return 0, err
	}
	r, err := decode[createResponse](resp)
	if err != nil {
		return 0, err
	}
	return r.ID, nil
}

// UpdateDeviceGroup updates a device group.
func (c *Client) UpdateDeviceGroup(ctx context.Context, id int64, input DeviceGroupUpdateInput) error {
	_, err := c.transport.DoRequest(ctx, http.MethodPut, fmt.Sprintf("/devices/groups/%d", id), input)
	return err
}

// DeleteDeviceGroup deletes a device group.
func (c *Client) DeleteDeviceGroup(ctx context.Context, id int64) error {
	_, err := c.transport.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/devices/groups/%d", id), nil)
	return err
}

// AddDevicesToGroup adds devices to a static device group.
func (c *Client) AddDevicesToGroup(ctx context.Context, groupID int64, udids []string) error {
	body := map[string]any{
		"groupId": groupID,
		"udids":   udids,
	}
	_, err := c.transport.DoRequest(ctx, http.MethodPost, "/devices/groups/add", body)
	return err
}

// RemoveDevicesFromGroup removes devices from a static device group.
func (c *Client) RemoveDevicesFromGroup(ctx context.Context, groupID int64, udids []string) error {
	body := map[string]any{
		"groupId": groupID,
		"udids":   udids,
	}
	_, err := c.transport.DoRequest(ctx, http.MethodPost, "/devices/groups/remove", body)
	return err
}

// GetDeviceGroupMembers retrieves the UDIDs of devices in a group.
func (c *Client) GetDeviceGroupMembers(ctx context.Context, groupID int64) ([]string, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, "/devices", nil, client.WithQuery(fmt.Sprintf("groups=%d", groupID)))
	if err != nil {
		return nil, err
	}
	var r struct {
		Devices []struct {
			UDID string `json:"UDID"`
		} `json:"devices"`
	}
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	udids := make([]string, len(r.Devices))
	for i, d := range r.Devices {
		udids[i] = d.UDID
	}
	return udids, nil
}
