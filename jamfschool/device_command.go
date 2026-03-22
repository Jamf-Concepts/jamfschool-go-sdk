// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"context"
	"fmt"
	"net/http"
)

// EraseDevice schedules a wipe on a device.
func (c *Client) EraseDevice(ctx context.Context, udid string, clearActivationLock bool) error {
	var body any
	if clearActivationLock {
		body = map[string]any{"clearActivationLock": "true"}
	}
	_, err := c.transport.DoRequest(ctx, http.MethodPost, fmt.Sprintf("/devices/%s/wipe", udid), body)
	return err
}

// RestartDevice schedules a restart on a device.
func (c *Client) RestartDevice(ctx context.Context, udid string, clearPasscode bool) error {
	var body any
	if clearPasscode {
		body = map[string]any{"clearPasscode": "true"}
	}
	_, err := c.transport.DoRequest(ctx, http.MethodPost, fmt.Sprintf("/devices/%s/restart", udid), body)
	return err
}

// RefreshDevice schedules a full inventory refresh for a device.
func (c *Client) RefreshDevice(ctx context.Context, udid string, clearErrors bool) error {
	var body any
	if clearErrors {
		body = map[string]any{"clearErrors": true}
	}
	_, err := c.transport.DoRequest(ctx, http.MethodPost, fmt.Sprintf("/devices/%s/refresh", udid), body)
	return err
}

// UnenrollDevice schedules removal of the management profile from a device.
func (c *Client) UnenrollDevice(ctx context.Context, udid string) error {
	_, err := c.transport.DoRequest(ctx, http.MethodPost, fmt.Sprintf("/devices/%s/unenroll", udid), nil)
	return err
}

// ClearDeviceActivationLock clears the activation lock from a device.
func (c *Client) ClearDeviceActivationLock(ctx context.Context, udid string) error {
	_, err := c.transport.DoRequest(ctx, http.MethodPost, fmt.Sprintf("/devices/%s/activationlock/clear", udid), nil)
	return err
}

// TrashDevice moves a device to trash.
func (c *Client) TrashDevice(ctx context.Context, udid string) error {
	_, err := c.transport.DoRequest(ctx, http.MethodDelete, fmt.Sprintf("/devices/%s", udid), nil)
	return err
}

// RestoreDevice restores a device from trash.
func (c *Client) RestoreDevice(ctx context.Context, udid string) error {
	_, err := c.transport.DoRequest(ctx, http.MethodPost, fmt.Sprintf("/devices/%s/restore", udid), nil)
	return err
}

// UpdateDeviceESIM queries a carrier URL for active eSIM cellular-plan profiles.
func (c *Client) UpdateDeviceESIM(ctx context.Context, udid string, serverURL string, requiresNetworkTether bool) error {
	body := map[string]any{
		"serverUrl":             serverURL,
		"requiresNetworkTether": requiresNetworkTether,
	}
	_, err := c.transport.DoRequest(ctx, http.MethodPost, fmt.Sprintf("/devices/%s/cellularPlan", udid), body)
	return err
}
