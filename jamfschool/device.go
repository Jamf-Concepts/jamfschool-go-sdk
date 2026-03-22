// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"fmt"
	"net/http"
)

// DeviceModel represents model information for a device.
type DeviceModel struct {
	Name       string `json:"name,omitempty"`
	Identifier string `json:"identifier,omitempty"`
	Type       string `json:"type,omitempty"`
}

// DeviceOS represents OS information for a device.
type DeviceOS struct {
	Prefix  string `json:"prefix,omitempty"`
	Version string `json:"version,omitempty"`
}

// Device represents a Jamf School device.
type Device struct {
	UDID             string      `json:"UDID,omitempty"`
	LocationID       int64       `json:"locationId,omitempty"`
	SerialNumber     string      `json:"serialNumber,omitempty"`
	AssetTag         string      `json:"assetTag,omitempty"`
	Name             string      `json:"name,omitempty"`
	IsManaged        bool        `json:"isManaged,omitempty"`
	IsSupervised     bool        `json:"isSupervised,omitempty"`
	DeviceEnrollType string      `json:"deviceEnrollType,omitempty"`
	BatteryLevel     float64     `json:"batteryLevel,omitempty"`
	TotalCapacity    float64     `json:"totalCapacity,omitempty"`
	Notes            string      `json:"notes,omitempty"`
	LastCheckin      string      `json:"lastCheckin,omitempty"`
	Model            DeviceModel `json:"model"`
	OS               DeviceOS    `json:"os"`
	InTrash          bool        `json:"inTrash,omitempty"`
}

type deviceResponse struct {
	Code   int    `json:"code"`
	Device Device `json:"device"`
}

type devicesResponse struct {
	Code    int      `json:"code"`
	Count   int      `json:"count"`
	Devices []Device `json:"devices"`
}

// GetDevice retrieves a device by UDID.
func (c *Client) GetDevice(ctx context.Context, udid string) (*Device, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/devices/%s", udid), nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[deviceResponse](resp)
	if err != nil {
		return nil, err
	}
	return &r.Device, nil
}

// GetDevices retrieves all devices.
func (c *Client) GetDevices(ctx context.Context) ([]Device, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, "/devices", nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[devicesResponse](resp)
	if err != nil {
		return nil, err
	}
	return r.Devices, nil
}
