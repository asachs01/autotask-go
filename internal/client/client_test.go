package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	// Create a new client
	client := NewClient("test-user", "test-secret", "test-integration-code")

	// Verify client is not nil
	if client == nil {
		t.Fatal("client should not be nil")
	}

	// Verify client properties
	if client.username != "test-user" {
		t.Errorf("expected username to be 'test-user', got '%s'", client.username)
	}
	if client.secret != "test-secret" {
		t.Errorf("expected secret to be 'test-secret', got '%s'", client.secret)
	}
	if client.integrationCode != "test-integration-code" {
		t.Errorf("expected integration code to be 'test-integration-code', got '%s'", client.integrationCode)
	}
	if client.UserAgent != "Autotask Go Client" {
		t.Errorf("expected user agent to be 'Autotask Go Client', got '%s'", client.UserAgent)
	}
	if client.client == nil {
		t.Error("HTTP client should not be nil")
	}
	if client.rateLimiter == nil {
		t.Error("rate limiter should not be nil")
	}
}

func TestNewRequest(t *testing.T) {
	// Create a new client
	client := NewClient("test-user", "test-secret", "test-integration-code")

	// Set base URL
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	baseURL, _ := url.Parse(server.URL)
	client.baseURL = baseURL

	// Test with nil body
	ctx := context.Background()
	req, err := client.newRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("newRequest returned error: %v", err)
	}

	// Verify request
	if req.Method != http.MethodGet {
		t.Errorf("expected method to be GET, got %s", req.Method)
	}
	if req.URL.Path != "/test" {
		t.Errorf("expected path to be /test, got %s", req.URL.Path)
	}
	if req.Header.Get("User-Agent") != client.UserAgent {
		t.Errorf("expected User-Agent header to be %s, got %s", client.UserAgent, req.Header.Get("User-Agent"))
	}

	// For GET requests without a body, Content-Type might not be set
	if req.Body != nil && req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type header to be application/json, got %s", req.Header.Get("Content-Type"))
	}

	// Test with body
	body := map[string]string{"test": "value"}
	req, err = client.newRequest(ctx, http.MethodPost, "/test", body)
	if err != nil {
		t.Fatalf("newRequest returned error: %v", err)
	}

	// Verify request
	if req.Method != http.MethodPost {
		t.Errorf("expected method to be POST, got %s", req.Method)
	}
	if req.URL.Path != "/test" {
		t.Errorf("expected path to be /test, got %s", req.URL.Path)
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type header to be application/json, got %s", req.Header.Get("Content-Type"))
	}

	// Decode body
	var decodedBody map[string]string
	err = json.NewDecoder(req.Body).Decode(&decodedBody)
	if err != nil {
		t.Fatalf("error decoding request body: %v", err)
	}
	if decodedBody["test"] != "value" {
		t.Errorf("expected body to contain test=value, got %v", decodedBody)
	}
}

func TestDo(t *testing.T) {
	// Create a new client
	client := NewClient("test-user", "test-secret", "test-integration-code")

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request headers
		if r.Header.Get("User-Agent") != client.UserAgent {
			t.Errorf("expected User-Agent header to be %s, got %s", client.UserAgent, r.Header.Get("User-Agent"))
		}

		// Return a test response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"test": "response"}`)); err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Set base URL
	baseURL, _ := client.baseURL.Parse(server.URL)
	client.baseURL = baseURL

	// Create a request
	ctx := context.Background()
	req, err := client.newRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("newRequest returned error: %v", err)
	}

	// Execute the request
	var result map[string]string
	err = client.do(req, &result)
	if err != nil {
		t.Fatalf("do returned error: %v", err)
	}

	// Verify response
	if result["test"] != "response" {
		t.Errorf("expected response to contain test=response, got %v", result)
	}

	// Test error response
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(`{"error": "test error"}`)); err != nil {
			t.Errorf("Failed to write error response: %v", err)
		}
	}))
	defer errorServer.Close()

	// Set base URL
	baseURL, _ = client.baseURL.Parse(errorServer.URL)
	client.baseURL = baseURL

	// Create a request
	req, err = client.newRequest(ctx, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("newRequest returned error: %v", err)
	}

	// Execute the request
	err = client.do(req, &result)
	if err == nil {
		t.Fatal("expected do to return an error, got nil")
	}
}
