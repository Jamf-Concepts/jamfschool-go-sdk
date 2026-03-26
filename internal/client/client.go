// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

// Package client provides the HTTP transport layer for the Jamf School API.
//
// This package handles authentication, request/response processing, error handling,
// and logging. It does not contain any resource-specific types or methods;
// those belong in the jamfschool package.

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Logger is an interface for logging HTTP requests and responses.
type Logger interface {
	LogRequest(ctx context.Context, method, url string, headers http.Header, body []byte)
	LogResponse(ctx context.Context, statusCode int, headers http.Header, body []byte)
}

// RequestOption configures a single DoRequest call.
type RequestOption func(*requestOptions)

type requestOptions struct {
	protocolVersion string
	query           string
}

// WithProtocolVersion overrides the X-Server-Protocol-Version header for a
// single request. Some endpoints (e.g. app create/trash) require version "4".
func WithProtocolVersion(v string) RequestOption {
	return func(o *requestOptions) {
		o.protocolVersion = v
	}
}

// WithQuery appends a raw query string to the request URL.
func WithQuery(q string) RequestOption {
	return func(o *requestOptions) {
		o.query = q
	}
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient overrides the default HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// Client wraps the Jamf School REST API with HTTP Basic Auth.
type Client struct {
	baseURL    string
	networkID  string
	apiKey     string
	httpClient *http.Client
	logger     Logger
	userAgent  string
}

// NewClient creates a new Jamf School API client.
func NewClient(baseURL, networkID, apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL:   strings.TrimRight(baseURL, "/"),
		networkID: networkID,
		apiKey:    apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: "jamfschool-go-sdk/dev",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// NewClientWithUserAgent creates a new Jamf School API client with a custom user agent string.
func NewClientWithUserAgent(baseURL, networkID, apiKey, userAgent string, opts ...Option) *Client {
	opts = append([]Option{}, opts...)
	c := NewClient(baseURL, networkID, apiKey, opts...)
	if userAgent != "" {
		c.userAgent = userAgent
	}
	return c
}

// BaseURL returns the base URL configured for the client.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// SetLogger sets the logger for the client.
func (c *Client) SetLogger(logger Logger) {
	c.logger = logger
}

// DoRequest executes an HTTP request against the Jamf School API.
func (c *Client) DoRequest(ctx context.Context, method, path string, body any, opts ...RequestOption) ([]byte, error) {
	o := requestOptions{protocolVersion: "3"}
	for _, opt := range opts {
		opt(&o)
	}

	reqURL, err := url.JoinPath(c.baseURL, "api", path)
	if err != nil {
		return nil, fmt.Errorf("building request URL: %w", err)
	}
	if o.query != "" {
		reqURL += "?" + o.query
	}

	var reqBody io.Reader
	var reqBytes []byte
	if body != nil {
		reqBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshalling request body: %w", err)
		}
		reqBody = bytes.NewReader(reqBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.SetBasicAuth(c.networkID, c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Server-Protocol-Version", o.protocolVersion)
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	if c.logger != nil {
		c.logger.LogRequest(ctx, method, reqURL, redactRequestHeaders(req.Header), redactRequestBody(reqBytes))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrHTTP, err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if c.logger != nil {
		c.logger.LogResponse(ctx, resp.StatusCode, resp.Header, respBody)
	}

	switch {
	case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
		return nil, fmt.Errorf("%w: %s", ErrAuthentication, string(respBody))
	case resp.StatusCode == http.StatusNotFound:
		return nil, fmt.Errorf("%w: %s", ErrNotFound, string(respBody))
	case resp.StatusCode >= 400:
		return nil, fmt.Errorf("%w: status %d: %s", ErrHTTP, resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// sensitiveBodyKeys lists JSON field names whose values must be redacted before logging.
var sensitiveBodyKeys = map[string]bool{
	"password":      true,
	"storePassword": true,
}

// redactRequestBody returns a copy of body with sensitive JSON field values replaced.
func redactRequestBody(body []byte) []byte {
	if len(body) == 0 {
		return body
	}

	var parsed map[string]any
	if err := json.Unmarshal(body, &parsed); err != nil {
		return body
	}

	redacted := false
	for key := range sensitiveBodyKeys {
		if _, ok := parsed[key]; ok {
			parsed[key] = "[REDACTED]"
			redacted = true
		}
	}
	if !redacted {
		return body
	}

	out, err := json.Marshal(parsed)
	if err != nil {
		return body
	}
	return out
}

// redactRequestHeaders returns a copy of headers with sensitive values redacted.
func redactRequestHeaders(headers http.Header) http.Header {
	if headers == nil {
		return nil
	}
	clone := headers.Clone()
	if clone.Get("Authorization") != "" {
		clone.Set("Authorization", "[REDACTED]")
	}
	return clone
}
