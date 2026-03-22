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
