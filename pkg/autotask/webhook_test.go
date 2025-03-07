package autotask

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
