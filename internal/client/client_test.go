// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MIT

package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	c := NewClient("https://example.jamfcloud.com", "net123", "key456")

	if c.baseURL != "https://example.jamfcloud.com" {
		t.Errorf("expected baseURL %q, got %q", "https://example.jamfcloud.com", c.baseURL)
	}
	if c.networkID != "net123" {
		t.Errorf("expected networkID %q, got %q", "net123", c.networkID)
	}
	if c.apiKey != "key456" {
		t.Errorf("expected apiKey %q, got %q", "key456", c.apiKey)
	}
}

func TestNewClient_TrimsTrailingSlash(t *testing.T) {
	t.Parallel()

	c := NewClient("https://example.jamfcloud.com/", "net", "key")

	if c.baseURL != "https://example.jamfcloud.com" {
		t.Errorf("expected trailing slash trimmed, got %q", c.baseURL)
	}
}

func TestNewClient_DefaultUserAgent(t *testing.T) {
	t.Parallel()

	c := NewClient("https://example.jamfcloud.com", "net", "key")

	if c.userAgent != "jamfschool-go-sdk/dev" {
		t.Errorf("expected default user agent %q, got %q", "jamfschool-go-sdk/dev", c.userAgent)
	}
}

func TestNewClientWithUserAgent(t *testing.T) {
	t.Parallel()

	c := NewClientWithUserAgent("https://example.jamfcloud.com", "net", "key", "my-app/1.0")

	if c.userAgent != "my-app/1.0" {
		t.Errorf("expected user agent %q, got %q", "my-app/1.0", c.userAgent)
	}
}

func TestBaseURL(t *testing.T) {
	t.Parallel()

	c := NewClient("https://example.jamfcloud.com", "net", "key")

	if c.BaseURL() != "https://example.jamfcloud.com" {
		t.Errorf("expected %q, got %q", "https://example.jamfcloud.com", c.BaseURL())
	}
}

func TestDoRequest_Success(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "net123" || pass != "key456" {
			t.Errorf("unexpected auth: user=%q, pass=%q, ok=%v", user, pass, ok)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %q", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("expected Accept application/json, got %q", r.Header.Get("Accept"))
		}
		if r.Header.Get("X-Server-Protocol-Version") != "3" {
			t.Errorf("expected X-Server-Protocol-Version 3, got %q", r.Header.Get("X-Server-Protocol-Version"))
		}
		if r.URL.Path != "/api/users" {
			t.Errorf("expected path /api/users, got %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id": 1}`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "net123", "key456")
	body, err := c.DoRequest(context.Background(), http.MethodGet, "/users", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != `{"id": 1}` {
		t.Errorf("unexpected body: %s", string(body))
	}
}

func TestDoRequest_WithBody(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "net", "key")
	payload := map[string]string{"name": "test"}
	body, err := c.DoRequest(context.Background(), http.MethodPost, "/items", payload)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != `{"ok": true}` {
		t.Errorf("unexpected body: %s", string(body))
	}
}

func TestDoRequest_UserAgent(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("User-Agent"); got != "my-app/1.0" {
			t.Errorf("expected User-Agent %q, got %q", "my-app/1.0", got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := NewClientWithUserAgent(server.URL, "net", "key", "my-app/1.0")
	_, err := c.DoRequest(context.Background(), http.MethodGet, "/users", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDoRequest_Unauthorized(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`Unauthorized`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "bad", "creds")
	_, err := c.DoRequest(context.Background(), http.MethodGet, "/users", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrAuthentication) {
		t.Errorf("expected ErrAuthentication, got %v", err)
	}
}

func TestDoRequest_Forbidden(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`Forbidden`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "net", "key")
	_, err := c.DoRequest(context.Background(), http.MethodGet, "/users", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrAuthentication) {
		t.Errorf("expected ErrAuthentication, got %v", err)
	}
}

func TestDoRequest_NotFound(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`Not Found`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "net", "key")
	_, err := c.DoRequest(context.Background(), http.MethodGet, "/users/999", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestDoRequest_ServerError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`Internal Server Error`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "net", "key")
	_, err := c.DoRequest(context.Background(), http.MethodGet, "/users", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrHTTP) {
		t.Errorf("expected ErrHTTP, got %v", err)
	}
}

func TestDoRequest_BadRequest(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`Bad Request`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "net", "key")
	_, err := c.DoRequest(context.Background(), http.MethodPost, "/users", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrHTTP) {
		t.Errorf("expected ErrHTTP, got %v", err)
	}
}

func TestDoRequest_InvalidBody(t *testing.T) {
	t.Parallel()

	c := NewClient("https://example.com", "net", "key")
	_, err := c.DoRequest(context.Background(), http.MethodPost, "/users", make(chan int))

	if err == nil {
		t.Fatal("expected error for unmarshalable body, got nil")
	}
}

func TestDoRequest_CancelledContext(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := NewClient(server.URL, "net", "key")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.DoRequest(ctx, http.MethodGet, "/users", nil)

	if err == nil {
		t.Fatal("expected error for cancelled context, got nil")
	}
}

func TestRedactRequestHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		headers  http.Header
		wantAuth string
		wantNil  bool
	}{
		{
			name:    "nil headers",
			headers: nil,
			wantNil: true,
		},
		{
			name: "redacts authorization",
			headers: http.Header{
				"Authorization": []string{"Basic abc123"},
				"Content-Type":  []string{"application/json"},
			},
			wantAuth: "[REDACTED]",
		},
		{
			name: "no authorization header",
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
			wantAuth: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := redactRequestHeaders(tt.headers)
			if tt.wantNil {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
				return
			}
			if got.Get("Authorization") != tt.wantAuth {
				t.Errorf("expected Authorization %q, got %q", tt.wantAuth, got.Get("Authorization"))
			}
			if got.Get("Content-Type") != tt.headers.Get("Content-Type") {
				t.Errorf("expected Content-Type preserved, got %q", got.Get("Content-Type"))
			}
		})
	}
}

func TestRedactRequestHeaders_DoesNotMutateOriginal(t *testing.T) {
	t.Parallel()

	original := http.Header{"Authorization": []string{"Basic secret"}}
	_ = redactRequestHeaders(original)

	if original.Get("Authorization") != "Basic secret" {
		t.Errorf("original header was mutated: got %q", original.Get("Authorization"))
	}
}

func TestRedactRequestBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "redacts password",
			body: `{"username":"jsmith","password":"secret123"}`,
			want: `{"password":"[REDACTED]","username":"jsmith"}`,
		},
		{
			name: "redacts storePassword",
			body: `{"username":"jsmith","storePassword":true}`,
			want: `{"storePassword":"[REDACTED]","username":"jsmith"}`,
		},
		{
			name: "redacts both fields",
			body: `{"username":"jsmith","password":"secret","storePassword":true}`,
			want: `{"password":"[REDACTED]","storePassword":"[REDACTED]","username":"jsmith"}`,
		},
		{
			name: "no sensitive fields unchanged",
			body: `{"username":"jsmith","email":"j@example.com"}`,
			want: `{"username":"jsmith","email":"j@example.com"}`,
		},
		{
			name: "empty body unchanged",
			body: "",
			want: "",
		},
		{
			name: "non-json body unchanged",
			body: "not json",
			want: "not json",
		},
		{
			name: "json array unchanged",
			body: `[1,2,3]`,
			want: `[1,2,3]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := string(redactRequestBody([]byte(tt.body)))
			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestDoRequest_RedactsPasswordInLoggedBody(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id": 1}`))
	}))
	defer server.Close()

	logger := &capturingLogger{}
	c := NewClient(server.URL, "net", "key")
	c.SetLogger(logger)

	payload := map[string]any{
		"username": "jsmith",
		"password": "SuperSecret123!",
	}
	_, err := c.DoRequest(context.Background(), http.MethodPost, "/users", payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logged := string(logger.lastRequestBody)
	if contains(logged, "SuperSecret123!") {
		t.Errorf("password was not redacted in logged body: %s", logged)
	}
	if !contains(logged, "[REDACTED]") {
		t.Errorf("expected [REDACTED] in logged body: %s", logged)
	}
	if !contains(logged, "jsmith") {
		t.Errorf("expected non-sensitive field preserved in logged body: %s", logged)
	}
}

// capturingLogger captures the most recent LogRequest body for test assertions.
type capturingLogger struct {
	lastRequestBody []byte
}

func (l *capturingLogger) LogRequest(_ context.Context, _, _ string, _ http.Header, body []byte) {
	l.lastRequestBody = body
}

func (l *capturingLogger) LogResponse(_ context.Context, _ int, _ http.Header, _ []byte) {}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
