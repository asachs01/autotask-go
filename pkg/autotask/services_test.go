package autotask

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

// Define a contextKey type for context values
type contextKey string

// setupMockServer creates a test server and configures the client to use it
func setupMockServer(handler http.Handler) (*httptest.Server, Client) {
	server := httptest.NewServer(handler)
	client := NewClient("username", "password", "integrationCode")

	// Use reflection to set the baseURL field
	clientValue := reflect.ValueOf(client).Elem()
	baseURLField := clientValue.FieldByName("baseURL")
	if baseURLField.IsValid() && baseURLField.CanSet() {
		baseURL, _ := url.Parse(server.URL)
		baseURLField.Set(reflect.ValueOf(baseURL))
	}

	// Set zoneInfo directly to bypass the API call
	zoneInfoField := clientValue.FieldByName("zoneInfo")
	if zoneInfoField.IsValid() && zoneInfoField.CanSet() {
		zoneInfo := &ZoneInfo{
			ZoneName: "Test Zone",
			URL:      server.URL,
			WebURL:   "https://example.com",
			CI:       123,
		}
		zoneInfoField.Set(reflect.ValueOf(zoneInfo))
	}

	return server, client
}

func TestCompaniesService(t *testing.T) {
	// Create a mock server
	mockServer := NewMockServer(t)
	defer mockServer.Close()

	// Add handlers for the Companies endpoints
	mockServer.AddHandler("/Companies/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			mockServer.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
				"item": map[string]interface{}{
					"id":   float64(1),
					"name": "Test Company",
				},
			})
		} else if r.Method == http.MethodPatch {
			mockServer.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
				"item": map[string]interface{}{
					"id":   float64(1),
					"name": "Updated Company",
				},
			})
		} else if r.Method == http.MethodDelete {
			mockServer.RespondWithJSON(w, http.StatusOK, nil)
		}
	})

	mockServer.AddHandler("/Companies/query", func(w http.ResponseWriter, r *http.Request) {
		items := []interface{}{
			map[string]interface{}{
				"id":   float64(1),
				"name": "Test Company",
			},
			map[string]interface{}{
				"id":   float64(2),
				"name": "Another Company",
			},
		}
		mockServer.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"items": items,
		})
	})

	mockServer.AddHandler("/Companies", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			mockServer.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
				"item": map[string]interface{}{
					"id":   float64(3),
					"name": "New Company",
				},
			})
		}
	})

	mockServer.AddHandler("/Companies/query/count", func(w http.ResponseWriter, r *http.Request) {
		mockServer.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"count": float64(2),
		})
	})

	// Create a client using the mock server
	client := mockServer.NewTestClient()

	// Create a companies service
	service := &companiesService{
		BaseEntityService: NewBaseEntityService(client, "Companies"),
	}

	ctx := context.Background()

	// Test Get
	t.Run("Get", func(t *testing.T) {
		company, err := service.Get(ctx, 1)
		if err != nil {
			t.Fatalf("Failed to get company: %v", err)
		}

		companyMap, ok := company.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected map[string]interface{}, got %T", company)
		}

		if id, ok := companyMap["id"].(float64); !ok || id != 1 {
			t.Errorf("Expected id 1, got %v", companyMap["id"])
		}

		if name, ok := companyMap["name"].(string); !ok || name != "Test Company" {
			t.Errorf("Expected name 'Test Company', got %v", companyMap["name"])
		}
	})

	// Test Query
	t.Run("Query", func(t *testing.T) {
		var response struct {
			Items []map[string]interface{} `json:"items"`
		}
		err := service.Query(ctx, "", &response)
		if err != nil {
			t.Fatalf("Failed to query companies: %v", err)
		}

		if len(response.Items) != 2 {
			t.Errorf("Expected 2 companies, got %d", len(response.Items))
		}
	})

	// Test Create
	t.Run("Create", func(t *testing.T) {
		newCompany := map[string]interface{}{
			"name": "New Company",
		}
		company, err := service.Create(ctx, newCompany)
		if err != nil {
			t.Fatalf("Failed to create company: %v", err)
		}

		companyMap, ok := company.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected map[string]interface{}, got %T", company)
		}

		if id, ok := companyMap["id"].(float64); !ok || id != 3 {
			t.Errorf("Expected id 3, got %v", companyMap["id"])
		}

		if name, ok := companyMap["name"].(string); !ok || name != "New Company" {
			t.Errorf("Expected name 'New Company', got %v", companyMap["name"])
		}
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		updatedCompany := map[string]interface{}{
			"name": "Updated Company",
		}
		company, err := service.Update(ctx, 1, updatedCompany)
		if err != nil {
			t.Fatalf("Failed to update company: %v", err)
		}

		companyMap, ok := company.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected map[string]interface{}, got %T", company)
		}

		if name, ok := companyMap["name"].(string); !ok || name != "Updated Company" {
			t.Errorf("Expected name 'Updated Company', got %v", companyMap["name"])
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		err := service.Delete(ctx, 1)
		if err != nil {
			t.Fatalf("Failed to delete company: %v", err)
		}
	})

	// Test Count
	t.Run("Count", func(t *testing.T) {
		count, err := service.Count(ctx, "")
		if err != nil {
			t.Fatalf("Failed to count companies: %v", err)
		}

		if count != 2 {
			t.Errorf("Expected count 2, got %d", count)
		}
	})
}

func TestWebhookService(t *testing.T) {
	// Create a mock server
	mockServer := NewMockServer(t)
	defer mockServer.Close()

	// Add handlers for the Webhooks endpoints
	mockServer.AddHandler("/Webhooks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			mockServer.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
				"id":     float64(1),
				"url":    "https://example.com/webhook",
				"events": []string{"ticket.created", "ticket.updated"},
			})
		}
	})

	mockServer.AddHandler("/Webhooks/query", func(w http.ResponseWriter, r *http.Request) {
		items := []interface{}{
			map[string]interface{}{
				"id":     float64(1),
				"url":    "https://example.com/webhook",
				"events": []string{"ticket.created", "ticket.updated"},
			},
			map[string]interface{}{
				"id":     float64(2),
				"url":    "https://example.com/webhook2",
				"events": []string{"company.created"},
			},
		}
		mockServer.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"items": items,
		})
	})

	mockServer.AddHandler("/Webhooks/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			mockServer.RespondWithJSON(w, http.StatusOK, nil)
		}
	})

	// Create a client using the mock server
	client := mockServer.NewTestClient()

	// Create a webhook service
	service := &webhookService{
		BaseEntityService: NewBaseEntityService(client, "Webhooks"),
		handlers:          make(map[string][]WebhookHandler),
	}

	ctx := context.Background()

	// Test RegisterHandler
	t.Run("RegisterHandler", func(t *testing.T) {
		handler := func(event *WebhookEvent) error {
			return nil
		}
		service.RegisterHandler("ticket.created", handler)

		if len(service.handlers["ticket.created"]) != 1 {
			t.Errorf("Expected 1 handler, got %d", len(service.handlers["ticket.created"]))
		}
	})

	// Test SetWebhookSecret
	t.Run("SetWebhookSecret", func(t *testing.T) {
		service.SetWebhookSecret("test-secret")
		if service.secret != "test-secret" {
			t.Errorf("Expected secret to be 'test-secret', got '%s'", service.secret)
		}
	})

	// Test CreateWebhook
	t.Run("CreateWebhook", func(t *testing.T) {
		err := service.CreateWebhook(ctx, "https://example.com/webhook", []string{"ticket.created", "ticket.updated"})
		if err != nil {
			t.Fatalf("Failed to create webhook: %v", err)
		}
	})

	// Test ListWebhooks
	t.Run("ListWebhooks", func(t *testing.T) {
		webhooks, err := service.ListWebhooks(ctx)
		if err != nil {
			t.Fatalf("Failed to list webhooks: %v", err)
		}

		if len(webhooks) != 2 {
			t.Errorf("Expected 2 webhooks, got %d", len(webhooks))
		}
	})

	// Test DeleteWebhook
	t.Run("DeleteWebhook", func(t *testing.T) {
		err := service.DeleteWebhook(ctx, 1)
		if err != nil {
			t.Fatalf("Failed to delete webhook: %v", err)
		}
	})

	// Test HandleWebhook
	t.Run("HandleWebhook", func(t *testing.T) {
		handlerCalled := false
		handler := func(event *WebhookEvent) error {
			handlerCalled = true
			return nil
		}
		service.RegisterHandler("ticket.created", handler)
		service.SetWebhookSecret("test-secret")

		// Create a webhook event
		event := map[string]interface{}{
			"eventType": "ticket.created",
			"entityId":  float64(123),
			"entity":    "Ticket",
			"timestamp": "2023-01-01T12:00:00Z",
			"data":      json.RawMessage(`{"id": 123, "title": "Test Ticket", "status": 1}`),
		}

		// Create a request with the event
		eventBytes, _ := json.Marshal(event)
		req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(eventBytes))

		// Calculate signature
		mac := hmac.New(sha256.New, []byte("test-secret"))
		mac.Write(eventBytes)
		signature := hex.EncodeToString(mac.Sum(nil))
		req.Header.Set("X-Autotask-Signature", signature)
		req.Header.Set("Content-Type", "application/json")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Handle the webhook
		service.HandleWebhook(rr, req)

		// Check if the handler was called
		if !handlerCalled {
			t.Errorf("Expected handler to be called")
		}
	})
}

func TestProjectsService(t *testing.T) {
	// Create a mock server
	mockServer := NewMockServer(t)
	defer mockServer.Close()

	// Add handlers for the Projects endpoints
	mockServer.AddHandler("/Projects/1", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			mockServer.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
				"item": map[string]interface{}{
					"id":   float64(1),
					"name": "Test Project",
				},
			})
		} else if r.Method == http.MethodPatch {
			mockServer.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
				"item": map[string]interface{}{
					"id":   float64(1),
					"name": "Updated Project",
				},
			})
		} else if r.Method == http.MethodDelete {
			mockServer.RespondWithJSON(w, http.StatusOK, nil)
		}
	})

	mockServer.AddHandler("/Projects/query", func(w http.ResponseWriter, r *http.Request) {
		items := []interface{}{
			map[string]interface{}{
				"id":   float64(1),
				"name": "Test Project",
			},
			map[string]interface{}{
				"id":   float64(2),
				"name": "Another Project",
			},
		}
		mockServer.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"items": items,
		})
	})

	mockServer.AddHandler("/Projects", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			mockServer.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
				"item": map[string]interface{}{
					"id":   float64(3),
					"name": "New Project",
				},
			})
		}
	})

	// Create a client using the mock server
	client := mockServer.NewTestClient()

	// Create a projects service
	service := &projectsService{
		BaseEntityService: NewBaseEntityService(client, "Projects"),
	}

	ctx := context.Background()

	// Test Get
	t.Run("Get", func(t *testing.T) {
		project, err := service.Get(ctx, 1)
		if err != nil {
			t.Fatalf("Failed to get project: %v", err)
		}

		projectMap, ok := project.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected map[string]interface{}, got %T", project)
		}

		if id, ok := projectMap["id"].(float64); !ok || id != 1 {
			t.Errorf("Expected id 1, got %v", projectMap["id"])
		}

		if name, ok := projectMap["name"].(string); !ok || name != "Test Project" {
			t.Errorf("Expected name 'Test Project', got %v", projectMap["name"])
		}
	})

	// Test Query
	t.Run("Query", func(t *testing.T) {
		var response struct {
			Items []map[string]interface{} `json:"items"`
		}
		err := service.Query(ctx, "", &response)
		if err != nil {
			t.Fatalf("Failed to query projects: %v", err)
		}

		if len(response.Items) != 2 {
			t.Errorf("Expected 2 projects, got %d", len(response.Items))
		}
	})

	// Test Create
	t.Run("Create", func(t *testing.T) {
		newProject := map[string]interface{}{
			"name": "New Project",
		}
		project, err := service.Create(ctx, newProject)
		if err != nil {
			t.Fatalf("Failed to create project: %v", err)
		}

		projectMap, ok := project.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected map[string]interface{}, got %T", project)
		}

		if id, ok := projectMap["id"].(float64); !ok || id != 3 {
			t.Errorf("Expected id 3, got %v", projectMap["id"])
		}

		if name, ok := projectMap["name"].(string); !ok || name != "New Project" {
			t.Errorf("Expected name 'New Project', got %v", projectMap["name"])
		}
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		updatedProject := map[string]interface{}{
			"name": "Updated Project",
		}
		project, err := service.Update(ctx, 1, updatedProject)
		if err != nil {
			t.Fatalf("Failed to update project: %v", err)
		}

		projectMap, ok := project.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected map[string]interface{}, got %T", project)
		}

		if name, ok := projectMap["name"].(string); !ok || name != "Updated Project" {
			t.Errorf("Expected name 'Updated Project', got %v", projectMap["name"])
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		err := service.Delete(ctx, 1)
		if err != nil {
			t.Fatalf("Failed to delete project: %v", err)
		}
	})
}
