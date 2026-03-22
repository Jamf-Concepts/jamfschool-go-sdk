// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// ClassMember represents a user assigned to a class.
type ClassMember struct {
	ID int64 `json:"id"`
}

// Class represents a Jamf School class.
type Class struct {
	UUID          string        `json:"uuid,omitempty"`
	Name          string        `json:"name,omitempty"`
	Description   string        `json:"description,omitempty"`
	LocationID    int64         `json:"locationId,omitempty"`
	Source        string        `json:"source,omitempty"`
	ImageURL      string        `json:"image,omitempty"`
	UserGroupID   int64         `json:"userGroupId,omitempty"`
	StudentCount  int64         `json:"studentCount,omitempty"`
	TeacherCount  int64         `json:"teacherCount,omitempty"`
	DeviceGroupID int64         `json:"deviceGroupId,omitempty"`
	DeviceCount   int64         `json:"deviceCount,omitempty"`
	Students      []ClassMember `json:"students,omitempty"`
	Teachers      []ClassMember `json:"teachers,omitempty"`
}

// ClassCreateInput holds fields for creating a class.
type ClassCreateInput struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	LocationID  *int64  `json:"locationId,omitempty"`
	Students    []int64 `json:"students,omitempty"`
	Teachers    []int64 `json:"teachers,omitempty"`
}

// ClassUpdateInput holds fields for updating a class.
type ClassUpdateInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type classResponse struct {
	Code  int   `json:"code"`
	Class Class `json:"class"`
}

type classesResponse struct {
	Code    int     `json:"code"`
	Classes []Class `json:"classes"`
}

type classCreateResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	UUID    string `json:"uuid"`
}

// GetClass retrieves a class by UUID.
func (c *Client) GetClass(ctx context.Context, uuid string) (*Class, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/classes/%s", uuid), nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[classResponse](resp)
	if err != nil {
		return nil, err
	}
	return &r.Class, nil
}

// GetClasses retrieves all classes.
func (c *Client) GetClasses(ctx context.Context) ([]Class, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, "/classes", nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[classesResponse](resp)
	if err != nil {
		return nil, err
	}
	return r.Classes, nil
}

// CreateClass creates a new class.
func (c *Client) CreateClass(ctx context.Context, input ClassCreateInput) (string, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodPost, "/classes", input)
	if err != nil {
		return "", err
	}
	r, err := decode[classCreateResponse](resp)
	if err != nil {
		return "", err
	}
	return r.UUID, nil
}

// UpdateClass updates a class.
func (c *Client) UpdateClass(ctx context.Context, uuid string, input ClassUpdateInput) error {
	_, err := c.transport.DoRequest(ctx, http.MethodPut, fmt.Sprintf("/classes/%s", uuid), input)
	return err
}

// DeleteClass deletes a class.
func (c *Client) DeleteClass(ctx context.Context, uuid string) error {
	_, err := c.transport.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/classes/%s", uuid), nil)
	return err
}

// AssignClassUsers sets the students and teachers for a class.
func (c *Client) AssignClassUsers(ctx context.Context, uuid string, studentIDs, teacherIDs []int64) error {
	body := map[string]any{
		"students": int64sToStrings(studentIDs),
		"teachers": int64sToStrings(teacherIDs),
	}
	_, err := c.transport.DoRequest(ctx, http.MethodPut, fmt.Sprintf("/classes/%s/users", uuid), body)
	return err
}

// ClassDevice represents a device assigned to a class.
type ClassDevice struct {
	UDID         string `json:"UDID"`
	SerialNumber string `json:"serialNumber"`
	Name         string `json:"name"`
}

type classDevicesResponse struct {
	Code    int           `json:"code"`
	Count   int           `json:"count"`
	Devices []ClassDevice `json:"devices"`
}

// GetClassDevices retrieves the devices assigned to a class.
func (c *Client) GetClassDevices(ctx context.Context, uuid string) ([]ClassDevice, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/classes/%s/devices", uuid), nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[classDevicesResponse](resp)
	if err != nil {
		return nil, err
	}
	return r.Devices, nil
}

// int64sToStrings converts int64 IDs to string IDs as required by the API.
func int64sToStrings(ids []int64) []string {
	s := make([]string, len(ids))
	for i, id := range ids {
		s[i] = strconv.FormatInt(id, 10)
	}
	return s
}
