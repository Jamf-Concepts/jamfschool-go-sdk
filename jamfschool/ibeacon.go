// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// IBeacon represents a Jamf School iBeacon.
type IBeacon struct {
	ID          int64  `json:"id,omitempty"`
	UUID        string `json:"UUID,omitempty"`
	Major       int64  `json:"major,omitempty"`
	Minor       int64  `json:"minor,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// UnmarshalJSON handles the API returning major/minor as either strings or numbers.
func (b *IBeacon) UnmarshalJSON(data []byte) error {
	type Alias IBeacon
	aux := &struct {
		Major json.RawMessage `json:"major,omitempty"`
		Minor json.RawMessage `json:"minor,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(b),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if len(aux.Major) > 0 {
		b.Major = parseFlexibleInt64(aux.Major)
	}
	if len(aux.Minor) > 0 {
		b.Minor = parseFlexibleInt64(aux.Minor)
	}
	return nil
}

// parseFlexibleInt64 parses a JSON value that may be a number or a quoted string.
func parseFlexibleInt64(raw json.RawMessage) int64 {
	var n json.Number
	if err := json.Unmarshal(raw, &n); err == nil {
		if v, err := n.Int64(); err == nil {
			return v
		}
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		v, _ := strconv.ParseInt(s, 10, 64)
		return v
	}
	return 0
}

// IBeaconCreateInput holds fields for creating an iBeacon.
type IBeaconCreateInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	UUID        string `json:"UUID,omitempty"`
	Major       *int64 `json:"major,omitempty"`
	Minor       *int64 `json:"minor,omitempty"`
}

// IBeaconUpdateInput holds fields for updating an iBeacon.
type IBeaconUpdateInput struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	UUID        string `json:"UUID,omitempty"`
	Major       *int64 `json:"major,omitempty"`
	Minor       *int64 `json:"minor,omitempty"`
}

type ibeaconResponse struct {
	Beacon IBeacon `json:"beacon"`
}

type ibeaconsResponse struct {
	Beacons []IBeacon `json:"beacons"`
}

type ibeaconCreateResponse struct {
	Message string  `json:"message"`
	Beacon  IBeacon `json:"beacon"`
}

type ibeaconUpdateResponse struct {
	Message string  `json:"message"`
	Beacon  IBeacon `json:"beacon"`
}

// GetIBeacon retrieves an iBeacon by ID.
func (c *Client) GetIBeacon(ctx context.Context, id int64) (*IBeacon, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/ibeacons/%d", id), nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[ibeaconResponse](resp)
	if err != nil {
		return nil, err
	}
	return &r.Beacon, nil
}

// GetIBeacons retrieves all iBeacons.
func (c *Client) GetIBeacons(ctx context.Context) ([]IBeacon, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, "/ibeacons", nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[ibeaconsResponse](resp)
	if err != nil {
		return nil, err
	}
	return r.Beacons, nil
}

// CreateIBeacon creates a new iBeacon.
func (c *Client) CreateIBeacon(ctx context.Context, input IBeaconCreateInput) (*IBeacon, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodPost, "/ibeacons", input)
	if err != nil {
		return nil, err
	}
	r, err := decode[ibeaconCreateResponse](resp)
	if err != nil {
		return nil, err
	}
	return &r.Beacon, nil
}

// UpdateIBeacon updates an iBeacon (uses POST per API spec).
func (c *Client) UpdateIBeacon(ctx context.Context, id int64, input IBeaconUpdateInput) (*IBeacon, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodPost, fmt.Sprintf("/ibeacons/%d", id), input)
	if err != nil {
		return nil, err
	}
	r, err := decode[ibeaconUpdateResponse](resp)
	if err != nil {
		return nil, err
	}
	return &r.Beacon, nil
}

// DeleteIBeacon deletes an iBeacon.
func (c *Client) DeleteIBeacon(ctx context.Context, id int64) error {
	_, err := c.transport.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/ibeacons/%d", id), nil)
	return err
}
