// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"fmt"
	"net/http"
)

// Location represents a Jamf School location.
type Location struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name,omitempty"`
	IsDistrict    bool    `json:"isDistrict,omitempty"`
	Street        *string `json:"street"`
	StreetNumber  *string `json:"streetNumber"`
	PostalCode    *string `json:"postalCode"`
	City          *string `json:"city"`
	Source        string  `json:"source,omitempty"`
	ASMIdentifier *string `json:"asmIdentifier"`
	SchoolNumber  string  `json:"schoolNumber,omitempty"`
}

type locationsResponse struct {
	Locations []Location `json:"locations"`
}

// GetLocation retrieves a location by ID.
func (c *Client) GetLocation(ctx context.Context, id int64) (*Location, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/locations/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return decode[Location](resp)
}

// GetLocations retrieves all locations.
func (c *Client) GetLocations(ctx context.Context) ([]Location, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, "/locations", nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[locationsResponse](resp)
	if err != nil {
		return nil, err
	}
	return r.Locations, nil
}
