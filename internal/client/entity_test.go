package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/asachs01/autotask-go/pkg/autotask"
)

// setupTestServer creates a test server and configures the client to use it
func setupTestServer(handler http.Handler) (*httptest.Server, *Client, *EntityService) {
	server := httptest.NewServer(handler)

	client := NewClient("test-user", "test-secret", "test-integration-code")
	baseURL, _ := url.Parse(server.URL)
	client.baseURL = baseURL

	entityService := &EntityService{
		client:     client,
		entityName: "TestEntities",
	}

	return server, client, entityService
}

func TestEntityGet(t *testing.T) {
	// Setup test server
	server, _, entityService := setupTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodGet {
			t.Errorf("expected method to be GET, got %s", r.Method)
		}
		if r.URL.Path != "/TestEntities/123" {
			t.Errorf("expected path to be /TestEntities/123, got %s", r.URL.Path)
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 123, "name": "Test Entity"}`))
	}))
	defer server.Close()

	// Call Get
	ctx := context.Background()
	result, err := entityService.Get(ctx, 123)

	// Verify result
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	// Convert result to map
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected result to be map[string]interface{}, got %T", result)
	}

	// Verify result values
	if resultMap["id"].(float64) != 123 {
		t.Errorf("expected id to be 123, got %v", resultMap["id"])
	}
	if resultMap["name"].(string) != "Test Entity" {
		t.Errorf("expected name to be 'Test Entity', got %v", resultMap["name"])
	}
}

func TestEntityQuery(t *testing.T) {
	// Setup test server
	server, _, entityService := setupTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodGet {
			t.Errorf("expected method to be GET, got %s", r.Method)
		}
		if r.URL.Path != "/TestEntities/query" {
			t.Errorf("expected path to be /TestEntities/query, got %s", r.URL.Path)
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"items": [{"id": 123, "name": "Test Entity"}]}`))
	}))
	defer server.Close()

	// Call Query
	ctx := context.Background()
	var result map[string]interface{}
	err := entityService.Query(ctx, "name='Test'", &result)

	// Verify result
	if err != nil {
		t.Fatalf("Query returned error: %v", err)
	}

	// Verify result values
	items, ok := result["items"].([]interface{})
	if !ok {
		t.Fatalf("expected result to contain items array, got %v", result)
	}
	if len(items) != 1 {
		t.Fatalf("expected items to contain 1 item, got %d", len(items))
	}

	item, ok := items[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected item to be map[string]interface{}, got %T", items[0])
	}

	if item["id"].(float64) != 123 {
		t.Errorf("expected id to be 123, got %v", item["id"])
	}
	if item["name"].(string) != "Test Entity" {
		t.Errorf("expected name to be 'Test Entity', got %v", item["name"])
	}
}

func TestEntityCreate(t *testing.T) {
	// Setup test server
	server, _, entityService := setupTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodPost {
			t.Errorf("expected method to be POST, got %s", r.Method)
		}
		if r.URL.Path != "/TestEntities" {
			t.Errorf("expected path to be /TestEntities, got %s", r.URL.Path)
		}

		// Decode request body
		var requestBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			t.Fatalf("error decoding request body: %v", err)
		}

		// Verify request body
		if requestBody["name"] != "New Entity" {
			t.Errorf("expected name to be 'New Entity', got %v", requestBody["name"])
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 456, "name": "New Entity"}`))
	}))
	defer server.Close()

	// Call Create
	ctx := context.Background()
	entity := map[string]interface{}{
		"name": "New Entity",
	}
	result, err := entityService.Create(ctx, entity)

	// Verify result
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	// Convert result to map
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected result to be map[string]interface{}, got %T", result)
	}

	// Verify result values
	if resultMap["id"].(float64) != 456 {
		t.Errorf("expected id to be 456, got %v", resultMap["id"])
	}
	if resultMap["name"].(string) != "New Entity" {
		t.Errorf("expected name to be 'New Entity', got %v", resultMap["name"])
	}
}

func TestEntityUpdate(t *testing.T) {
	// Setup test server
	server, _, entityService := setupTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodPatch {
			t.Errorf("expected method to be PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/TestEntities/123" {
			t.Errorf("expected path to be /TestEntities/123, got %s", r.URL.Path)
		}

		// Decode request body
		var requestBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			t.Fatalf("error decoding request body: %v", err)
		}

		// Verify request body
		if requestBody["name"] != "Updated Entity" {
			t.Errorf("expected name to be 'Updated Entity', got %v", requestBody["name"])
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 123, "name": "Updated Entity"}`))
	}))
	defer server.Close()

	// Call Update
	ctx := context.Background()
	entity := map[string]interface{}{
		"name": "Updated Entity",
	}
	result, err := entityService.Update(ctx, 123, entity)

	// Verify result
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	// Convert result to map
	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected result to be map[string]interface{}, got %T", result)
	}

	// Verify result values
	if resultMap["id"].(float64) != 123 {
		t.Errorf("expected id to be 123, got %v", resultMap["id"])
	}
	if resultMap["name"].(string) != "Updated Entity" {
		t.Errorf("expected name to be 'Updated Entity', got %v", resultMap["name"])
	}
}

func TestEntityDelete(t *testing.T) {
	// Setup test server
	server, _, entityService := setupTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodDelete {
			t.Errorf("expected method to be DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/TestEntities/123" {
			t.Errorf("expected path to be /TestEntities/123, got %s", r.URL.Path)
		}

		// Return response
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Call Delete
	ctx := context.Background()
	err := entityService.Delete(ctx, 123)

	// Verify result
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
}

func TestEntityCount(t *testing.T) {
	// Setup test server
	server, _, entityService := setupTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodGet {
			t.Errorf("expected method to be GET, got %s", r.Method)
		}
		if r.URL.Path != "/TestEntities/count" {
			t.Errorf("expected path to be /TestEntities/count, got %s", r.URL.Path)
		}

		// Return response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"count": 42}`))
	}))
	defer server.Close()

	// Call Count
	ctx := context.Background()
	count, err := entityService.Count(ctx, "name='Test'")

	// Verify result
	if err != nil {
		t.Fatalf("Count returned error: %v", err)
	}
	if count != 42 {
		t.Errorf("expected count to be 42, got %d", count)
	}
}

func TestEntityBatchOperations(t *testing.T) {
	// Test BatchCreate
	t.Run("BatchCreate", func(t *testing.T) {
		// Setup test server
		server, _, entityService := setupTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request
			if r.Method != http.MethodPost {
				t.Errorf("expected method to be POST, got %s", r.Method)
			}
			if r.URL.Path != "/TestEntities/batch" {
				t.Errorf("expected path to be /TestEntities/batch, got %s", r.URL.Path)
			}

			// Decode request body
			var requestBody []interface{}
			err := json.NewDecoder(r.Body).Decode(&requestBody)
			if err != nil {
				t.Fatalf("error decoding request body: %v", err)
			}

			// Verify request body
			if len(requestBody) != 2 {
				t.Fatalf("expected 2 entities, got %d", len(requestBody))
			}

			// Return response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"items": [{"id": 1, "name": "Entity 1"}, {"id": 2, "name": "Entity 2"}]}`))
		}))
		defer server.Close()

		// Call BatchCreate
		ctx := context.Background()
		entities := []interface{}{
			map[string]interface{}{"name": "Entity 1"},
			map[string]interface{}{"name": "Entity 2"},
		}
		var result map[string]interface{}
		err := entityService.BatchCreate(ctx, entities, &result)

		// Verify result
		if err != nil {
			t.Fatalf("BatchCreate returned error: %v", err)
		}

		// Verify result values
		items, ok := result["items"].([]interface{})
		if !ok {
			t.Fatalf("expected result to contain items array, got %v", result)
		}
		if len(items) != 2 {
			t.Fatalf("expected items to contain 2 items, got %d", len(items))
		}
	})

	// Test BatchUpdate
	t.Run("BatchUpdate", func(t *testing.T) {
		// Setup test server
		server, _, entityService := setupTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request
			if r.Method != http.MethodPatch {
				t.Errorf("expected method to be PATCH, got %s", r.Method)
			}
			if r.URL.Path != "/TestEntities/batch" {
				t.Errorf("expected path to be /TestEntities/batch, got %s", r.URL.Path)
			}

			// Decode request body
			var requestBody []interface{}
			err := json.NewDecoder(r.Body).Decode(&requestBody)
			if err != nil {
				t.Fatalf("error decoding request body: %v", err)
			}

			// Verify request body
			if len(requestBody) != 2 {
				t.Fatalf("expected 2 entities, got %d", len(requestBody))
			}

			// Return response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"items": [{"id": 1, "name": "Updated 1"}, {"id": 2, "name": "Updated 2"}]}`))
		}))
		defer server.Close()

		// Call BatchUpdate
		ctx := context.Background()
		entities := []interface{}{
			map[string]interface{}{"id": 1, "name": "Updated 1"},
			map[string]interface{}{"id": 2, "name": "Updated 2"},
		}
		var result map[string]interface{}
		err := entityService.BatchUpdate(ctx, entities, &result)

		// Verify result
		if err != nil {
			t.Fatalf("BatchUpdate returned error: %v", err)
		}

		// Verify result values
		items, ok := result["items"].([]interface{})
		if !ok {
			t.Fatalf("expected result to contain items array, got %v", result)
		}
		if len(items) != 2 {
			t.Fatalf("expected items to contain 2 items, got %d", len(items))
		}
	})

	// Test BatchDelete
	t.Run("BatchDelete", func(t *testing.T) {
		// Setup test server
		server, _, entityService := setupTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request
			if r.Method != http.MethodDelete {
				t.Errorf("expected method to be DELETE, got %s", r.Method)
			}
			if r.URL.Path != "/TestEntities/batch" {
				t.Errorf("expected path to be /TestEntities/batch, got %s", r.URL.Path)
			}

			// Decode request body
			var requestBody []int64
			err := json.NewDecoder(r.Body).Decode(&requestBody)
			if err != nil {
				t.Fatalf("error decoding request body: %v", err)
			}

			// Verify request body
			if len(requestBody) != 2 {
				t.Fatalf("expected 2 ids, got %d", len(requestBody))
			}
			if requestBody[0] != 1 || requestBody[1] != 2 {
				t.Errorf("expected ids [1, 2], got %v", requestBody)
			}

			// Return response
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		// Call BatchDelete
		ctx := context.Background()
		ids := []int64{1, 2}
		err := entityService.BatchDelete(ctx, ids)

		// Verify result
		if err != nil {
			t.Fatalf("BatchDelete returned error: %v", err)
		}
	})
}

func TestEntityPagination(t *testing.T) {
	// Setup test server
	server, _, entityService := setupTestServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return response based on URL
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if r.URL.Path == "/pagination-test" {
			w.Write([]byte(`{"items": [{"id": 1, "name": "Item 1"}], "pageDetails": {"pageNumber": 1, "pageSize": 10, "count": 20, "nextPageUrl": "/next-page", "prevPageUrl": ""}}`))
		} else if r.URL.Path == "/next-page" {
			w.Write([]byte(`{"items": [{"id": 2, "name": "Item 2"}], "pageDetails": {"pageNumber": 2, "pageSize": 10, "count": 20, "nextPageUrl": "", "prevPageUrl": "/prev-page"}}`))
		} else if r.URL.Path == "/prev-page" {
			w.Write([]byte(`{"items": [{"id": 1, "name": "Item 1"}], "pageDetails": {"pageNumber": 1, "pageSize": 10, "count": 20, "nextPageUrl": "/next-page", "prevPageUrl": ""}}`))
		}
	}))
	defer server.Close()

	// Test Pagination
	t.Run("Pagination", func(t *testing.T) {
		ctx := context.Background()
		var result map[string]interface{}
		err := entityService.Pagination(ctx, "/pagination-test", &result)

		// Verify result
		if err != nil {
			t.Fatalf("Pagination returned error: %v", err)
		}

		// Verify result values
		items, ok := result["items"].([]interface{})
		if !ok {
			t.Fatalf("expected result to contain items array, got %v", result)
		}
		if len(items) != 1 {
			t.Fatalf("expected items to contain 1 item, got %d", len(items))
		}

		pageDetails, ok := result["pageDetails"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected result to contain pageDetails, got %v", result)
		}
		if pageDetails["pageNumber"].(float64) != 1 {
			t.Errorf("expected pageNumber to be 1, got %v", pageDetails["pageNumber"])
		}
	})

	// Test GetNextPage
	t.Run("GetNextPage", func(t *testing.T) {
		ctx := context.Background()
		pageDetails := autotask.PageDetails{
			NextPageUrl: "/next-page",
		}
		items, err := entityService.GetNextPage(ctx, pageDetails)

		// Verify result
		if err != nil {
			t.Fatalf("GetNextPage returned error: %v", err)
		}
		if len(items) != 1 {
			t.Fatalf("expected items to contain 1 item, got %d", len(items))
		}

		item, ok := items[0].(map[string]interface{})
		if !ok {
			t.Fatalf("expected item to be map[string]interface{}, got %T", items[0])
		}
		if item["id"].(float64) != 2 {
			t.Errorf("expected id to be 2, got %v", item["id"])
		}
	})

	// Test GetPreviousPage
	t.Run("GetPreviousPage", func(t *testing.T) {
		ctx := context.Background()
		pageDetails := autotask.PageDetails{
			PrevPageUrl: "/prev-page",
		}
		items, err := entityService.GetPreviousPage(ctx, pageDetails)

		// Verify result
		if err != nil {
			t.Fatalf("GetPreviousPage returned error: %v", err)
		}
		if len(items) != 1 {
			t.Fatalf("expected items to contain 1 item, got %d", len(items))
		}

		item, ok := items[0].(map[string]interface{})
		if !ok {
			t.Fatalf("expected item to be map[string]interface{}, got %T", items[0])
		}
		if item["id"].(float64) != 1 {
			t.Errorf("expected id to be 1, got %v", item["id"])
		}
	})

	// Test GetPreviousPage with no URL
	t.Run("GetPreviousPageNoUrl", func(t *testing.T) {
		// Create a direct client without using a server
		client := NewClient("test-user", "test-secret", "test-integration-code")

		// Create an entity service
		entityService := &EntityService{
			client:     client,
			entityName: "TestEntities",
		}

		// Create page details with no previous page URL
		pageDetails := autotask.PageDetails{
			PrevPageUrl: "",
		}

		// Call GetPreviousPage directly
		ctx := context.Background()
		items, err := entityService.GetPreviousPage(ctx, pageDetails)

		// Verify no error and no items
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if items != nil {
			t.Errorf("expected nil items, got %v", items)
		}
	})
}
