// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package jamfschool

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/internal/client"
)

// Client provides typed methods for all Jamf School API operations.
type Client struct {
	transport *client.Client
}

// NewClient creates a new Jamf School API client.
func NewClient(baseURL, networkID, apiKey string, opts ...Option) *Client {
	cfg := &clientConfig{
		userAgent: "jamfschool-go-sdk/dev",
	}
	for _, opt := range opts {
		opt(cfg)
	}

	var transportOpts []client.Option
	if cfg.httpClient != nil {
		transportOpts = append(transportOpts, client.WithHTTPClient(cfg.httpClient))
	}

	transport := client.NewClientWithUserAgent(baseURL, networkID, apiKey, cfg.userAgent, transportOpts...)
	if cfg.logger != nil {
		transport.SetLogger(cfg.logger)
	}

	return &Client{
		transport: transport,
	}
}

// BaseURL returns the base URL configured for the client.
func (c *Client) BaseURL() string {
	return c.transport.BaseURL()
}

// clientConfig holds configuration applied via Option functions.
type clientConfig struct {
	userAgent  string
	httpClient *http.Client
	logger     Logger
}

// Option configures a Client.
type Option func(*clientConfig)

// WithUserAgent sets a custom user agent string.
func WithUserAgent(userAgent string) Option {
	return func(cfg *clientConfig) {
		if userAgent != "" {
			cfg.userAgent = userAgent
		}
	}
}

// WithHTTPClient overrides the default HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(cfg *clientConfig) {
		cfg.httpClient = httpClient
	}
}

// WithLogger sets a logger for HTTP request/response logging.
func WithLogger(logger Logger) Option {
	return func(cfg *clientConfig) {
		cfg.logger = logger
	}
}

// createResponse is shared across multiple resource types.
type createResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	ID      int64  `json:"id"`
}

// decode unmarshals JSON response data into a value of type T.
func decode[T any](data []byte) (*T, error) {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &v, nil
}
