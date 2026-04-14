// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// DEPDevice represents a DEP placeholder.
type DEPDevice struct {
	ID              int64  `json:"id,omitempty"`
	UserID          int64  `json:"userId,omitempty"`
	LocationID      int64  `json:"locationId,omitempty"`
	Model           string `json:"model,omitempty"`
	Color           string `json:"color,omitempty"`
	SerialNumber    string `json:"serialNumber,omitempty"`
	Status          string `json:"status,omitempty"`
	DeviceName      string `json:"deviceName,omitempty"`
	DateAdded       int64  `json:"dateAdded,omitempty"`
	DatePushed      int64  `json:"datePushed,omitempty"`
	ProfileName     string `json:"profileName,omitempty"`
	PlaceholderName string `json:"placeholderName,omitempty"`
}

// UnmarshalJSON handles the API returning dateAdded/datePushed as either
// numbers or strings depending on the endpoint.
func (d *DEPDevice) UnmarshalJSON(data []byte) error {
	type Alias DEPDevice
	aux := &struct {
		DateAdded  json.RawMessage `json:"dateAdded,omitempty"`
		DatePushed json.RawMessage `json:"datePushed,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if len(aux.DateAdded) > 0 {
		d.DateAdded = parseFlexibleInt64(aux.DateAdded)
	}
	if len(aux.DatePushed) > 0 {
		d.DatePushed = parseFlexibleInt64(aux.DatePushed)
	}
	return nil
}

type depDeviceResponse struct {
	Code        int       `json:"code"`
	Placeholder DEPDevice `json:"placeholder"`
}

type depDevicesResponse struct {
	Code         int         `json:"code"`
	Placeholders []DEPDevice `json:"placeholders"`
}

// GetDEPDevice retrieves a DEP device by serial number.
func (c *Client) GetDEPDevice(ctx context.Context, serial string) (*DEPDevice, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/dep/%s", serial), nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[depDeviceResponse](resp)
	if err != nil {
		return nil, err
	}
	return &r.Placeholder, nil
}

// GetDEPDevices retrieves all DEP devices.
func (c *Client) GetDEPDevices(ctx context.Context) ([]DEPDevice, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, "/dep", nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[depDevicesResponse](resp)
	if err != nil {
		return nil, err
	}
	return r.Placeholders, nil
}
