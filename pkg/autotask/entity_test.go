package autotask

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestNewBaseEntityService(t *testing.T) {
	// Create a mock client
	mockClient := NewClient("test-user", "test-secret", "test-integration-code")

	// Create a base entity service
	service := NewBaseEntityService(mockClient, "TestEntity")

	// Verify service is not nil
	AssertNotNil(t, service, "service should not be nil")

	// Verify entity name
	AssertEqual(t, "TestEntity", service.GetEntityName(), "entity name should match")

	// Verify client
	AssertEqual(t, mockClient, service.GetClient(), "client should match")
}

func TestEntityGet(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Add a handler for the entity endpoint
	server.AddHandler("/TestEntities/123", func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		AssertEqual(t, http.MethodGet, r.Method, "method should match")

		// Create a mock entity
		entity := map[string]interface{}{
			"id":   float64(123),
			"name": "Mock TestEntity 123",
		}

		// Create a response with the entity as the item
		response := Response{
			Item: entity,
		}

		// Send a response
		server.RespondWithJSON(w, http.StatusOK, response)
	})

	// Create a client
	client := server.NewTestClient()

	// Create a base entity service
	service := NewBaseEntityService(client, "TestEntities")

	// Get an entity
	ctx := context.Background()
	entity, err := service.Get(ctx, 123)

	// Verify no error
	AssertNil(t, err, "error should be nil")

	// Verify entity
	AssertNotNil(t, entity, "entity should not be nil")

	// Convert entity to map
	entityMap, ok := entity.(map[string]interface{})
	AssertTrue(t, ok, "entity should be of type map[string]interface{}")

	// Verify entity fields
	AssertEqual(t, float64(123), entityMap["id"], "entity ID should match")
	AssertEqual(t, "Mock TestEntity 123", entityMap["name"], "entity name should match")
}

func TestEntityQuery(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Add a handler for the entity endpoint
	server.AddHandler("/TestEntities/query", func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		AssertEqual(t, http.MethodGet, r.Method, "method should match")

		// Verify content type
		AssertEqual(t, "application/json", r.Header.Get("Content-Type"), "Content-Type header should match")

		// Parse query parameters
		query := r.URL.Query()
		searchParam := query.Get("search")

		// Verify search parameter is not empty
		AssertNotNil(t, searchParam, "search parameter should not be nil")

		// Decode the search parameter
		var searchObj map[string]interface{}
		err := json.Unmarshal([]byte(searchParam), &searchObj)
		AssertNil(t, err, "error decoding search parameter should be nil")

		// Extract filter from search object
		filter, ok := searchObj["filter"].(string)
		if !ok {
			// If filter is not a string, it might be a complex filter object
			// For this test, we'll just check that filter exists
			_, ok = searchObj["filter"]
			AssertTrue(t, ok, "filter should exist in search object")
		} else {
			// If it is a string, verify it
			AssertEqual(t, "name='Test'", filter, "filter should match")
		}

		// Send a response
		entities := CreateMockEntities(2, "TestEntity", 123)
		response := CreateMockListResponse(entities, 1, 10, 2)
		server.RespondWithJSON(w, http.StatusOK, response)
	})

	// Create a client
	client := server.NewTestClient()

	// Create a base entity service
	service := NewBaseEntityService(client, "TestEntities")

	// Query entities
	ctx := context.Background()
	var result ListResponse
	err := service.Query(ctx, "name='Test'", &result)

	// Verify no error
	AssertNil(t, err, "error should be nil")

	// Verify result
	AssertNotNil(t, result, "result should not be nil")
	AssertEqual(t, 2, len(result.Items), "result should have 2 items")

	// Verify page details
	AssertEqual(t, 1, result.PageDetails.PageNumber, "page number should match")
	AssertEqual(t, 10, result.PageDetails.PageSize, "page size should match")
	AssertEqual(t, 2, result.PageDetails.Count, "count should match")
}

func TestEntityCreate(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Add a handler for the entity endpoint
	server.AddHandler("/TestEntities", func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		AssertEqual(t, http.MethodPost, r.Method, "method should match")

		// Verify content type
		AssertEqual(t, "application/json", r.Header.Get("Content-Type"), "Content-Type header should match")

		// Read request body
		var entity map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&entity)
		AssertNil(t, err, "error decoding request body should be nil")

		// Verify entity fields
		AssertEqual(t, "Test Entity", entity["name"], "entity name should match")

		// Create response entity with ID
		responseEntity := map[string]interface{}{
			"id":   float64(123),
			"name": "Test Entity",
		}

		// Create a response with the entity as the item
		response := Response{
			Item: responseEntity,
		}

		// Send a response
		server.RespondWithJSON(w, http.StatusOK, response)
	})

	// Create a client
	client := server.NewTestClient()

	// Create a base entity service
	service := NewBaseEntityService(client, "TestEntities")

	// Create an entity
	ctx := context.Background()
	entity := map[string]interface{}{
		"name": "Test Entity",
	}
	result, err := service.Create(ctx, entity)

	// Verify no error
	AssertNil(t, err, "error should be nil")

	// Verify result
	AssertNotNil(t, result, "result should not be nil")

	// Convert result to map
	resultMap, ok := result.(map[string]interface{})
	AssertTrue(t, ok, "result should be of type map[string]interface{}")

	// Verify result fields
	AssertEqual(t, float64(123), resultMap["id"], "result ID should match")
	AssertEqual(t, "Test Entity", resultMap["name"], "result name should match")
}

func TestEntityUpdate(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Add a handler for the entity endpoint
	server.AddHandler("/TestEntities/123", func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		AssertEqual(t, http.MethodPatch, r.Method, "method should match")

		// Verify content type
		AssertEqual(t, "application/json", r.Header.Get("Content-Type"), "Content-Type header should match")

		// Read request body
		var entity map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&entity)
		AssertNil(t, err, "error decoding request body should be nil")

		// Verify entity fields
		AssertEqual(t, "Updated Entity", entity["name"], "entity name should match")

		// Create response entity with ID
		responseEntity := map[string]interface{}{
			"id":   float64(123),
			"name": "Updated Entity",
		}

		// Create a response with the entity as the item
		response := Response{
			Item: responseEntity,
		}

		// Send a response
		server.RespondWithJSON(w, http.StatusOK, response)
	})

	// Create a client
	client := server.NewTestClient()

	// Create a base entity service
	service := NewBaseEntityService(client, "TestEntities")

	// Update an entity
	ctx := context.Background()
	entity := map[string]interface{}{
		"name": "Updated Entity",
	}
	result, err := service.Update(ctx, 123, entity)

	// Verify no error
	AssertNil(t, err, "error should be nil")

	// Verify result
	AssertNotNil(t, result, "result should not be nil")

	// Convert result to map
	resultMap, ok := result.(map[string]interface{})
	AssertTrue(t, ok, "result should be of type map[string]interface{}")

	// Verify result fields
	AssertEqual(t, float64(123), resultMap["id"], "result ID should match")
	AssertEqual(t, "Updated Entity", resultMap["name"], "result name should match")
}

func TestEntityDelete(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Add a handler for the entity endpoint
	server.AddHandler("/TestEntities/123", func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		AssertEqual(t, http.MethodDelete, r.Method, "method should match")

		// Send a response with no content
		w.WriteHeader(http.StatusNoContent)
	})

	// Create a client
	client := server.NewTestClient()

	// Create a base entity service
	service := NewBaseEntityService(client, "TestEntities")

	// Delete an entity
	ctx := context.Background()
	err := service.Delete(ctx, 123)

	// Verify no error
	AssertNil(t, err, "error should be nil")
}

func TestEntityCount(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Add a handler for the entity endpoint
	server.AddHandler("/TestEntities/query/count", func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		AssertEqual(t, http.MethodGet, r.Method, "method should match")

		// Verify content type
		AssertEqual(t, "application/json", r.Header.Get("Content-Type"), "Content-Type header should match")

		// Parse query parameters
		query := r.URL.Query()
		searchParam := query.Get("search")

		// Verify search parameter is not empty
		AssertNotNil(t, searchParam, "search parameter should not be nil")

		// Decode the search parameter
		var searchObj map[string]interface{}
		err := json.Unmarshal([]byte(searchParam), &searchObj)
		AssertNil(t, err, "error decoding search parameter should be nil")

		// Extract filter from search object
		// For the Count method, the filter might be a complex object
		_, ok := searchObj["filter"]
		AssertTrue(t, ok, "filter should exist in search object")

		// Send a response
		server.RespondWithJSON(w, http.StatusOK, map[string]int{
			"count": 42,
		})
	})

	// Create a client
	client := server.NewTestClient()

	// Create a base entity service
	service := NewBaseEntityService(client, "TestEntities")

	// Count entities
	ctx := context.Background()
	count, err := service.Count(ctx, "name='Test'")

	// Verify no error
	AssertNil(t, err, "error should be nil")

	// Verify count
	AssertEqual(t, 42, count, "count should match")
}
