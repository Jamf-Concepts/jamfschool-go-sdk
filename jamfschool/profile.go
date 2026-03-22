// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"fmt"
	"net/http"
)

// Profile represents a Jamf School configuration profile.
type Profile struct {
	ID          int64  `json:"id,omitempty"`
	LocationID  int64  `json:"locationId,omitempty"`
	Identifier  string `json:"identifier,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Platform    string `json:"platform,omitempty"`
}

type profilesResponse struct {
	Profiles []Profile `json:"profiles"`
}

// GetProfile retrieves a profile by ID.
func (c *Client) GetProfile(ctx context.Context, id int64) (*Profile, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/profiles/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return decode[Profile](resp)
}

// GetProfiles retrieves all profiles.
func (c *Client) GetProfiles(ctx context.Context) ([]Profile, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, "/profiles", nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[profilesResponse](resp)
	if err != nil {
		return nil, err
	}
	return r.Profiles, nil
}
