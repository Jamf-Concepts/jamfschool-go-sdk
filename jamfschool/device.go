// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// DeviceModel represents model information for a device.
type DeviceModel struct {
	Name       string `json:"name,omitempty"`
	Identifier string `json:"identifier,omitempty"`
	Type       string `json:"type,omitempty"`
}

// UnmarshalJSON handles the API returning model.type as either a string
// (list endpoint) or an object (single-device endpoint).
func (m *DeviceModel) UnmarshalJSON(data []byte) error {
	type Alias DeviceModel
	aux := &struct {
		Type json.RawMessage `json:"type,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if len(aux.Type) > 0 {
		// Try string first (list endpoint returns e.g. "computer").
		var s string
		if json.Unmarshal(aux.Type, &s) == nil {
			m.Type = s
			return nil
		}
		// Fall back to object with a "type" or "name" field
		// (single-device endpoint returns e.g. {"name":"iPad","type":"iOS"}).
		var obj map[string]any
		if json.Unmarshal(aux.Type, &obj) == nil {
			if t, ok := obj["type"].(string); ok {
				m.Type = t
			} else if n, ok := obj["name"].(string); ok {
				m.Type = n
			}
		}
	}
	return nil
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

// UnmarshalJSON handles the single-device API endpoint returning some string
// fields (lastCheckin, deviceEnrollType, assetTag, notes) as objects instead
// of plain strings. When the value is an object, the "value" or "date" key
// is extracted; otherwise the raw string is used.
func (d *Device) UnmarshalJSON(data []byte) error {
	type Alias Device
	aux := &struct {
		LastCheckin      json.RawMessage `json:"lastCheckin,omitempty"`
		DeviceEnrollType json.RawMessage `json:"deviceEnrollType,omitempty"`
		AssetTag         json.RawMessage `json:"assetTag,omitempty"`
		Notes            json.RawMessage `json:"notes,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	d.LastCheckin = flexString(aux.LastCheckin)
	d.DeviceEnrollType = flexString(aux.DeviceEnrollType)
	d.AssetTag = flexString(aux.AssetTag)
	d.Notes = flexString(aux.Notes)
	return nil
}

// flexString extracts a string from a JSON value that may be a plain string
// or an object. When it's an object, it tries "value", "date", and "text"
// keys in order. Returns "" for null or empty input.
func flexString(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	var obj map[string]any
	if json.Unmarshal(raw, &obj) == nil {
		for _, key := range []string{"value", "date", "text", "name"} {
			if v, ok := obj[key].(string); ok {
				return v
			}
		}
		// Last resort: if there's a single string field, use it.
		for _, v := range obj {
			if vs, ok := v.(string); ok {
				return vs
			}
		}
	}
	return ""
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
