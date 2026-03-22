// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/internal/client"
)

// App represents a Jamf School app.
type App struct {
	ID         int64  `json:"id,omitempty"`
	LocationID int64  `json:"locationId,omitempty"`
	BundleID   string `json:"bundleId,omitempty"`
	AdamID     int64  `json:"adamId,omitempty"`
	Name       string `json:"name,omitempty"`
	Vendor     string `json:"vendor,omitempty"`
	Version    string `json:"version,omitempty"`
	Platform   string `json:"platform,omitempty"`
}

// AppCreateInput holds fields for creating an App Store app.
type AppCreateInput struct {
	AdamID      int64  `json:"adamId"`
	CountryCode string `json:"countryCode"`
}

type appCreateResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		MediaID    int64 `json:"mediaId"`
		LocationID int64 `json:"locationId"`
	} `json:"data"`
}

type appsResponse struct {
	Apps []App `json:"apps"`
}

// v4 is the protocol version required by app create/trash endpoints.
var v4 = client.WithProtocolVersion("4")

// CreateApp creates a new App Store app from its Adam ID and returns the app ID.
func (c *Client) CreateApp(ctx context.Context, input AppCreateInput) (int64, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodPost, "/apps", input, v4)
	if err != nil {
		return 0, err
	}
	r, err := decode[appCreateResponse](resp)
	if err != nil {
		return 0, err
	}
	return r.Data.MediaID, nil
}

// TrashApp moves an app to trash (soft delete).
func (c *Client) TrashApp(ctx context.Context, id int64) error {
	_, err := c.transport.DoRequest(ctx, http.MethodPost, fmt.Sprintf("/apps/%d/trash", id), nil, v4)
	return err
}

// GetApp retrieves an app by ID.
func (c *Client) GetApp(ctx context.Context, id int64) (*App, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, fmt.Sprintf("/apps/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return decode[App](resp)
}

// GetApps retrieves all apps.
func (c *Client) GetApps(ctx context.Context) ([]App, error) {
	resp, err := c.transport.DoRequest(ctx, http.MethodGet, "/apps", nil)
	if err != nil {
		return nil, err
	}
	r, err := decode[appsResponse](resp)
	if err != nil {
		return nil, err
	}
	return r.Apps, nil
}
