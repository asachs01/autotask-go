package autotask

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
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
	// Setup a mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle zone information request
		if strings.Contains(r.URL.Path, "ZoneInformation") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"zoneName": "Test Zone",
				"url": "https://api.example.com",
				"webUrl": "https://example.com",
				"ci": 123
			}`))
			return
		}

		switch r.URL.Path {
		case "/Companies/1":
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"item": {"id": 1, "name": "Test Company"}}`))
			} else if r.Method == http.MethodPatch {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"item": {"id": 1, "name": "Updated Company"}}`))
			} else if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
			}
		case "/Companies":
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"items": [{"id": 1, "name": "Test Company"}, {"id": 2, "name": "Another Company"}]}`))
			} else if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{"item": {"id": 3, "name": "New Company"}}`))
			}
		case "/Companies/count":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"count": 2}`))
		}
	})

	server, client := setupMockServer(handler)
	defer server.Close()

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
		var companies []map[string]interface{}
		err := service.Query(ctx, "filter=name contains 'Company'", &companies)
		if err != nil {
			t.Fatalf("Failed to query companies: %v", err)
		}

		if len(companies) != 2 {
			t.Errorf("Expected 2 companies, got %d", len(companies))
		}
	})

	// Test Create
	t.Run("Create", func(t *testing.T) {
		newCompany := map[string]interface{}{
			"name": "New Company",
		}

		created, err := service.Create(ctx, newCompany)
		if err != nil {
			t.Fatalf("Failed to create company: %v", err)
		}

		createdMap, ok := created.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected map[string]interface{}, got %T", created)
		}

		if id, ok := createdMap["id"].(float64); !ok || id != 3 {
			t.Errorf("Expected id 3, got %v", createdMap["id"])
		}

		if name, ok := createdMap["name"].(string); !ok || name != "New Company" {
			t.Errorf("Expected name 'New Company', got %v", createdMap["name"])
		}
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		updateCompany := map[string]interface{}{
			"name": "Updated Company",
		}

		updated, err := service.Update(ctx, 1, updateCompany)
		if err != nil {
			t.Fatalf("Failed to update company: %v", err)
		}

		updatedMap, ok := updated.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected map[string]interface{}, got %T", updated)
		}

		if name, ok := updatedMap["name"].(string); !ok || name != "Updated Company" {
			t.Errorf("Expected name 'Updated Company', got %v", updatedMap["name"])
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
		count, err := service.Count(ctx, "filter=name contains 'Company'")
		if err != nil {
			t.Fatalf("Failed to count companies: %v", err)
		}

		if count != 2 {
			t.Errorf("Expected count 2, got %d", count)
		}
	})
}

func TestWebhookService(t *testing.T) {
	// Setup a mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle zone information request
		if strings.Contains(r.URL.Path, "ZoneInformation") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"zoneName": "Test Zone",
				"url": "https://api.example.com",
				"webUrl": "https://example.com",
				"ci": 123
			}`))
			return
		}

		switch r.URL.Path {
		case "/Webhooks":
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"items": [{"id": 1, "url": "https://example.com/webhook", "events": ["ticket.created"]}]}`))
			} else if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{"item": {"id": 2, "url": "https://example.com/new-webhook", "events": ["ticket.updated"]}}`))
			}
		case "/Webhooks/1":
			if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
			}
		}
	})

	server, client := setupMockServer(handler)
	defer server.Close()

	// Create a webhook service
	service := &webhookService{
		BaseEntityService: NewBaseEntityService(client, "Webhooks"),
		handlers:          make(map[string][]WebhookHandler),
		secret:            "test-secret",
	}

	ctx := context.Background()

	// Test RegisterHandler
	t.Run("RegisterHandler", func(t *testing.T) {
		handler := func(event *WebhookEvent) error {
			return nil
		}

		service.RegisterHandler("ticket.created", handler)

		if len(service.handlers["ticket.created"]) != 1 {
			t.Errorf("Expected 1 handler for ticket.created, got %d", len(service.handlers["ticket.created"]))
		}
	})

	// Test SetWebhookSecret
	t.Run("SetWebhookSecret", func(t *testing.T) {
		service.SetWebhookSecret("new-secret")
		if service.secret != "new-secret" {
			t.Errorf("Expected secret to be 'new-secret', got '%s'", service.secret)
		}
	})

	// Test CreateWebhook
	t.Run("CreateWebhook", func(t *testing.T) {
		err := service.CreateWebhook(ctx, "https://example.com/new-webhook", []string{"ticket.updated"})
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

		if len(webhooks) != 1 {
			t.Errorf("Expected 1 webhook, got %d", len(webhooks))
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
		// Create a test request
		payload := map[string]interface{}{
			"event": "ticket.created",
			"data":  map[string]interface{}{"id": 123, "title": "Test Ticket"},
		}
		payloadBytes, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/webhook", nil)
		req.Body = http.NoBody
		req = req.WithContext(context.WithValue(req.Context(), contextKey("payload"), payloadBytes))

		// Create a response recorder
		w := httptest.NewRecorder()

		// Register a handler
		handlerCalled := false
		service.RegisterHandler("ticket.created", func(event *WebhookEvent) error {
			handlerCalled = true
			return nil
		})

		// Call HandleWebhook
		service.HandleWebhook(w, req)

		// Check that the handler was called
		if !handlerCalled {
			t.Error("Expected handler to be called")
		}
	})
}

func TestProjectsService(t *testing.T) {
	// Setup a mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle zone information request
		if strings.Contains(r.URL.Path, "ZoneInformation") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"zoneName": "Test Zone",
				"url": "https://api.example.com",
				"webUrl": "https://example.com",
				"ci": 123
			}`))
			return
		}

		switch r.URL.Path {
		case "/Projects/1":
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"item": {"id": 1, "name": "Test Project"}}`))
			} else if r.Method == http.MethodPatch {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"item": {"id": 1, "name": "Updated Project"}}`))
			} else if r.Method == http.MethodDelete {
				w.WriteHeader(http.StatusNoContent)
			}
		case "/Projects":
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"items": [{"id": 1, "name": "Test Project"}, {"id": 2, "name": "Another Project"}]}`))
			} else if r.Method == http.MethodPost {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{"item": {"id": 3, "name": "New Project"}}`))
			}
		}
	})

	server, client := setupMockServer(handler)
	defer server.Close()

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
		var projects []map[string]interface{}
		err := service.Query(ctx, "filter=name contains 'Project'", &projects)
		if err != nil {
			t.Fatalf("Failed to query projects: %v", err)
		}

		if len(projects) != 2 {
			t.Errorf("Expected 2 projects, got %d", len(projects))
		}
	})

	// Test Create
	t.Run("Create", func(t *testing.T) {
		newProject := map[string]interface{}{
			"name": "New Project",
		}

		created, err := service.Create(ctx, newProject)
		if err != nil {
			t.Fatalf("Failed to create project: %v", err)
		}

		createdMap, ok := created.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected map[string]interface{}, got %T", created)
		}

		if id, ok := createdMap["id"].(float64); !ok || id != 3 {
			t.Errorf("Expected id 3, got %v", createdMap["id"])
		}

		if name, ok := createdMap["name"].(string); !ok || name != "New Project" {
			t.Errorf("Expected name 'New Project', got %v", createdMap["name"])
		}
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		updateProject := map[string]interface{}{
			"name": "Updated Project",
		}

		updated, err := service.Update(ctx, 1, updateProject)
		if err != nil {
			t.Fatalf("Failed to update project: %v", err)
		}

		updatedMap, ok := updated.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected map[string]interface{}, got %T", updated)
		}

		if name, ok := updatedMap["name"].(string); !ok || name != "Updated Project" {
			t.Errorf("Expected name 'Updated Project', got %v", updatedMap["name"])
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
