package autotask

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestWebhookHandlerRegistration(t *testing.T) {
	service := &webhookService{}

	// Register a handler
	called := false
	handler := func(event *WebhookEvent) error {
		called = true
		return nil
	}

	service.RegisterHandler("test.event", handler)

	// Verify handler was registered
	if service.handlers == nil {
		t.Fatal("handlers map was not initialized")
	}

	handlers, exists := service.handlers["test.event"]
	if !exists {
		t.Fatal("handler was not registered for event type")
	}

	if len(handlers) != 1 {
		t.Fatalf("expected 1 handler, got %d", len(handlers))
	}

	// Create a test event
	event := &WebhookEvent{
		EventType: "test.event",
		Entity:    "Test",
		EntityID:  123,
		Timestamp: "2023-01-01T12:00:00Z",
		Data:      json.RawMessage(`{"id": 123, "name": "Test"}`),
	}

	// Call the handler
	err := handlers[0](event)
	if err != nil {
		t.Fatalf("handler returned error: %v", err)
	}

	if !called {
		t.Fatal("handler was not called")
	}
}

func TestWebhookHandling(t *testing.T) {
	// Create a mock client
	mockClient := &client{
		logger: New(LogLevelInfo, false),
	}

	// Create webhook service
	service := &webhookService{
		BaseEntityService: NewBaseEntityService(mockClient, "Webhooks"),
	}

	// Register a handler
	eventData := make(map[string]interface{})
	handler := func(event *WebhookEvent) error {
		if event.EventType != "test.event" {
			t.Fatalf("expected event type 'test.event', got '%s'", event.EventType)
		}

		if event.EntityID != 123 {
			t.Fatalf("expected entity ID 123, got %d", event.EntityID)
		}

		// Unmarshal the event data
		if err := json.Unmarshal(event.Data, &eventData); err != nil {
			t.Fatalf("failed to unmarshal event data: %v", err)
		}

		return nil
	}

	service.RegisterHandler("test.event", handler)

	// Create a test webhook event
	webhookEvent := WebhookEvent{
		EventType: "test.event",
		Entity:    "Test",
		EntityID:  123,
		Timestamp: "2023-01-01T12:00:00Z",
		Data:      json.RawMessage(`{"id": 123, "name": "Test Entity"}`),
	}

	// Convert to JSON
	eventJSON, err := json.Marshal(webhookEvent)
	if err != nil {
		t.Fatalf("failed to marshal webhook event: %v", err)
	}

	// Create a test request
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(eventJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Handle the webhook
	service.HandleWebhook(rr, req)

	// Check response
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Verify event data was processed
	if len(eventData) == 0 {
		t.Fatal("event data was not processed")
	}

	if id, ok := eventData["id"].(float64); !ok || id != 123 {
		t.Fatalf("expected id 123, got %v", eventData["id"])
	}

	if name, ok := eventData["name"].(string); !ok || name != "Test Entity" {
		t.Fatalf("expected name 'Test Entity', got %v", eventData["name"])
	}
}

func TestWebhookSignatureVerification(t *testing.T) {
	// Create a mock client
	mockClient := &client{
		logger: New(LogLevelInfo, false),
	}

	// Create webhook service with a secret
	service := &webhookService{
		BaseEntityService: NewBaseEntityService(mockClient, "Webhooks"),
		secret:            "test-secret",
	}

	// Create a test webhook event
	webhookEvent := WebhookEvent{
		EventType: "test.event",
		Entity:    "Test",
		EntityID:  123,
		Timestamp: "2023-01-01T12:00:00Z",
		Data:      json.RawMessage(`{"id": 123, "name": "Test Entity"}`),
	}

	// Convert to JSON
	eventJSON, err := json.Marshal(webhookEvent)
	if err != nil {
		t.Fatalf("failed to marshal webhook event: %v", err)
	}

	// Create a valid signature
	mac := hmac.New(sha256.New, []byte("test-secret"))
	mac.Write(eventJSON)
	validSignature := hex.EncodeToString(mac.Sum(nil))

	// Test with valid signature
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Autotask-Signature", validSignature)

	rr := httptest.NewRecorder()
	service.HandleWebhook(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status code %d with valid signature, got %d", http.StatusOK, rr.Code)
	}

	// Test with invalid signature
	req = httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Autotask-Signature", "invalid-signature")

	rr = httptest.NewRecorder()
	service.HandleWebhook(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status code %d with invalid signature, got %d", http.StatusUnauthorized, rr.Code)
	}

	// Test with missing signature
	req = httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(eventJSON))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	service.HandleWebhook(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status code %d with missing signature, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestSetWebhookSecret(t *testing.T) {
	service := &webhookService{}

	// Set the webhook secret
	service.SetWebhookSecret("test-secret")

	// Verify the secret was set
	if service.secret != "test-secret" {
		t.Fatalf("expected secret 'test-secret', got '%s'", service.secret)
	}

	// Change the secret
	service.SetWebhookSecret("new-secret")

	// Verify the secret was updated
	if service.secret != "new-secret" {
		t.Fatalf("expected secret 'new-secret', got '%s'", service.secret)
	}
}

func TestHandleWebhookWithMultipleHandlers(t *testing.T) {
	// Create a mock client
	mockClient := &client{
		logger: New(LogLevelInfo, false),
	}

	// Create webhook service
	service := &webhookService{
		BaseEntityService: NewBaseEntityService(mockClient, "Webhooks"),
	}

	// Track which handlers were called
	handler1Called := false
	handler2Called := false

	// Register multiple handlers for the same event type
	service.RegisterHandler("test.event", func(event *WebhookEvent) error {
		handler1Called = true
		return nil
	})

	service.RegisterHandler("test.event", func(event *WebhookEvent) error {
		handler2Called = true
		return nil
	})

	// Create a test webhook event
	webhookEvent := WebhookEvent{
		EventType: "test.event",
		Entity:    "Test",
		EntityID:  123,
		Timestamp: "2023-01-01T12:00:00Z",
		Data:      json.RawMessage(`{"id": 123, "name": "Test Entity"}`),
	}

	// Convert to JSON
	eventJSON, err := json.Marshal(webhookEvent)
	if err != nil {
		t.Fatalf("failed to marshal webhook event: %v", err)
	}

	// Create a test request
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(eventJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Handle the webhook
	service.HandleWebhook(rr, req)

	// Check response
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Verify both handlers were called
	if !handler1Called {
		t.Fatal("first handler was not called")
	}

	if !handler2Called {
		t.Fatal("second handler was not called")
	}
}

func TestHandleWebhookWithErrorHandler(t *testing.T) {
	// Create a mock client
	mockClient := &client{
		logger: New(LogLevelInfo, false),
	}

	// Create webhook service
	service := &webhookService{
		BaseEntityService: NewBaseEntityService(mockClient, "Webhooks"),
	}

	// Track which handlers were called
	handler1Called := false
	handler2Called := false

	// Register a handler that returns an error
	service.RegisterHandler("test.event", func(event *WebhookEvent) error {
		handler1Called = true
		return errors.New("test error")
	})

	// Register a second handler that should still be called despite the first error
	service.RegisterHandler("test.event", func(event *WebhookEvent) error {
		handler2Called = true
		return nil
	})

	// Create a test webhook event
	webhookEvent := WebhookEvent{
		EventType: "test.event",
		Entity:    "Test",
		EntityID:  123,
		Timestamp: "2023-01-01T12:00:00Z",
		Data:      json.RawMessage(`{"id": 123, "name": "Test Entity"}`),
	}

	// Convert to JSON
	eventJSON, err := json.Marshal(webhookEvent)
	if err != nil {
		t.Fatalf("failed to marshal webhook event: %v", err)
	}

	// Create a test request
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(eventJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Handle the webhook
	service.HandleWebhook(rr, req)

	// Check response - should still be OK even though a handler returned an error
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Verify both handlers were called
	if !handler1Called {
		t.Fatal("first handler was not called")
	}

	if !handler2Called {
		t.Fatal("second handler was not called despite error in first handler")
	}
}

func TestHandleWebhookWithInvalidMethod(t *testing.T) {
	// Create a mock client
	mockClient := &client{
		logger: New(LogLevelInfo, false),
	}

	// Create webhook service
	service := &webhookService{
		BaseEntityService: NewBaseEntityService(mockClient, "Webhooks"),
	}

	// Create a test request with GET method (should be POST)
	req := httptest.NewRequest(http.MethodGet, "/webhook", nil)

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Handle the webhook
	service.HandleWebhook(rr, req)

	// Check response - should be Method Not Allowed
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status code %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestHandleWebhookWithInvalidJSON(t *testing.T) {
	// Create a mock client
	mockClient := &client{
		logger: New(LogLevelInfo, false),
	}

	// Create webhook service
	service := &webhookService{
		BaseEntityService: NewBaseEntityService(mockClient, "Webhooks"),
	}

	// Create a test request with invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Handle the webhook
	service.HandleWebhook(rr, req)

	// Check response - should be Bad Request
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleWebhookWithNoHandlers(t *testing.T) {
	// Create a mock client
	mockClient := &client{
		logger: New(LogLevelInfo, false),
	}

	// Create webhook service
	service := &webhookService{
		BaseEntityService: NewBaseEntityService(mockClient, "Webhooks"),
	}

	// Create a test webhook event with an event type that has no handlers
	webhookEvent := WebhookEvent{
		EventType: "unhandled.event",
		Entity:    "Test",
		EntityID:  123,
		Timestamp: "2023-01-01T12:00:00Z",
		Data:      json.RawMessage(`{"id": 123, "name": "Test Entity"}`),
	}

	// Convert to JSON
	eventJSON, err := json.Marshal(webhookEvent)
	if err != nil {
		t.Fatalf("failed to marshal webhook event: %v", err)
	}

	// Create a test request
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewBuffer(eventJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Handle the webhook
	service.HandleWebhook(rr, req)

	// Check response - should still be OK even though there are no handlers
	if rr.Code != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}

// Mock HTTP client for testing API calls
type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// mockClient creates a client with a mock HTTP client for testing
func createMockClient(doFunc func(req *http.Request) (*http.Response, error)) *client {
	mockClient := &client{
		logger: New(LogLevelInfo, false),
		httpClient: &http.Client{
			Transport: &mockTransport{doFunc: doFunc},
		},
		rateLimiter: NewRateLimiter(60), // Initialize rate limiter
	}

	// Set a dummy base URL to avoid nil pointer dereference
	mockClient.baseURL, _ = url.Parse("https://example.com/atservicesrest/v1.0/")

	return mockClient
}

// mockTransport implements http.RoundTripper for testing
type mockTransport struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

func TestCreateWebhook(t *testing.T) {
	// Create a mock response
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(bytes.NewBufferString(`{
			"id": 123,
			"url": "https://example.com/webhook",
			"events": ["ticket.created", "ticket.updated"]
		}`)),
	}

	// Create a mock client
	mockClient := createMockClient(func(req *http.Request) (*http.Response, error) {
		// For the Create method, we expect a POST request
		if strings.Contains(req.URL.String(), "Create") {
			if req.Method != http.MethodPost {
				t.Fatalf("expected method %s for Create, got %s", http.MethodPost, req.Method)
			}

			// Read the request body
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("failed to read request body: %v", err)
			}

			// Verify the request body
			var webhook struct {
				URL    string   `json:"url"`
				Events []string `json:"events"`
			}
			if err := json.Unmarshal(body, &webhook); err != nil {
				t.Fatalf("failed to unmarshal request body: %v", err)
			}

			if webhook.URL != "https://example.com/webhook" {
				t.Fatalf("expected URL 'https://example.com/webhook', got '%s'", webhook.URL)
			}

			if len(webhook.Events) != 2 {
				t.Fatalf("expected 2 events, got %d", len(webhook.Events))
			}
		}

		return mockResp, nil
	})

	// Create webhook service
	service := &webhookService{
		BaseEntityService: NewBaseEntityService(mockClient, "Webhooks"),
	}

	// Create a webhook
	ctx := context.Background()
	err := service.CreateWebhook(ctx, "https://example.com/webhook", []string{"ticket.created", "ticket.updated"})

	// Verify no error
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteWebhook(t *testing.T) {
	// Create a mock response
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
	}

	// Create a mock client
	mockClient := createMockClient(func(req *http.Request) (*http.Response, error) {
		// For the Delete method, we expect a DELETE request
		if strings.Contains(req.URL.String(), "Delete") {
			if req.Method != http.MethodDelete {
				t.Fatalf("expected method %s for Delete, got %s", http.MethodDelete, req.Method)
			}

			// Verify the URL contains the webhook ID
			if !strings.Contains(req.URL.String(), "123") {
				t.Fatalf("expected URL to contain webhook ID '123', got '%s'", req.URL.String())
			}
		}

		return mockResp, nil
	})

	// Create webhook service
	service := &webhookService{
		BaseEntityService: NewBaseEntityService(mockClient, "Webhooks"),
	}

	// Delete a webhook
	ctx := context.Background()
	err := service.DeleteWebhook(ctx, 123)

	// Verify no error
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestListWebhooks(t *testing.T) {
	// Create a mock response with the expected structure
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(bytes.NewBufferString(`{
			"items": [
				{
					"id": 123,
					"url": "https://example.com/webhook1",
					"events": ["ticket.created"]
				},
				{
					"id": 456,
					"url": "https://example.com/webhook2",
					"events": ["ticket.updated", "ticket.deleted"]
				}
			],
			"pageDetails": {
				"count": 2,
				"pageNumber": 1,
				"pageSize": 10
			}
		}`)),
	}

	// Create a mock client
	mockClient := createMockClient(func(req *http.Request) (*http.Response, error) {
		// For the Query method, we expect a GET request
		return mockResp, nil
	})

	// Create webhook service
	service := &webhookService{
		BaseEntityService: NewBaseEntityService(mockClient, "Webhooks"),
	}

	// List webhooks
	ctx := context.Background()
	webhooks, err := service.ListWebhooks(ctx)

	// Verify no error
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify webhooks
	if len(webhooks) != 2 {
		t.Fatalf("expected 2 webhooks, got %d", len(webhooks))
	}
}
