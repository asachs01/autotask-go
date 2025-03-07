package autotask

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

func TestNewPaginationIterator(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Add a handler for the query endpoint
	server.AddHandler("/TestEntities/query", func(w http.ResponseWriter, r *http.Request) {
		// Create mock entities
		entities := CreateMockEntities(2, "TestEntity", 123)

		// Create mock response with pagination details
		response := CreateMockListResponse(entities, 1, 10, 25)

		// Send response
		server.RespondWithJSON(w, http.StatusOK, response)
	})

	// Create a client
	client := server.NewTestClient()

	// Create a base entity service
	baseService := NewBaseEntityService(client, "TestEntities")
	service := &baseService

	// Create a pagination iterator
	ctx := context.Background()
	iterator, err := NewPaginationIterator(ctx, service, "name='Test'", 10)

	// Verify no error
	AssertNil(t, err, "error should be nil")

	// Verify iterator
	AssertNotNil(t, iterator, "iterator should not be nil")
	AssertEqual(t, 10, iterator.pageSize, "page size should match")
	AssertEqual(t, "name='Test'", iterator.filter, "filter should match")
	AssertEqual(t, -1, iterator.currentIndex, "current index should be -1")
}

func TestPaginationIteratorNext(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Add a handler for the query endpoint
	server.AddHandler("/TestEntities/query", func(w http.ResponseWriter, r *http.Request) {
		// Read request body to determine if it's the first or second page
		var requestBody map[string]interface{}
		if r.Body != nil {
			err := json.NewDecoder(r.Body).Decode(&requestBody)
			if err != nil {
				// If there's an error decoding, it might be a GET request with query params
				// Just use default values
				requestBody = make(map[string]interface{})
			}
		}

		// Get page parameter if it exists
		var pageNum int
		if page, ok := requestBody["page"].(float64); ok {
			pageNum = int(page)
		} else {
			pageNum = 1 // Default to first page
		}

		// Create mock entities with IDs based on the page
		startID := int64((pageNum-1)*2 + 123)
		entities := CreateMockEntities(2, "TestEntity", startID)

		// Create mock response with pagination details
		response := CreateMockListResponse(entities, pageNum, 2, 6)

		// Add next/prev page URLs to the response
		pageDetails := response["pageDetails"].(PageDetails)
		if pageNum < 3 {
			pageDetails.NextPageUrl = "/TestEntities/query?page=" + strconv.Itoa(pageNum+1)
		}
		if pageNum > 1 {
			pageDetails.PrevPageUrl = "/TestEntities/query?page=" + strconv.Itoa(pageNum-1)
		}
		response["pageDetails"] = pageDetails

		// Send response
		server.RespondWithJSON(w, http.StatusOK, response)
	})

	// Create a client
	client := server.NewTestClient()

	// Create a base entity service
	baseService := NewBaseEntityService(client, "TestEntities")
	service := &baseService

	// Create a pagination iterator
	ctx := context.Background()
	iterator, err := NewPaginationIterator(ctx, service, "name='Test'", 2)

	// Verify no error
	AssertNil(t, err, "error should be nil")

	// Verify first page
	hasNext := iterator.Next()
	AssertTrue(t, hasNext, "should have next item")
	AssertNil(t, iterator.Error(), "error should be nil")

	item := iterator.Item()
	AssertNotNil(t, item, "item should not be nil")

	// Verify second item
	hasNext = iterator.Next()
	AssertTrue(t, hasNext, "should have next item")
	AssertNil(t, iterator.Error(), "error should be nil")

	item = iterator.Item()
	AssertNotNil(t, item, "item should not be nil")

	// Verify loading next page
	hasNext = iterator.Next()
	AssertTrue(t, hasNext, "should have next item")
	AssertNil(t, iterator.Error(), "error should be nil")

	item = iterator.Item()
	AssertNotNil(t, item, "item should not be nil")
}

func TestFetchPage(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Add a handler for the query endpoint
	server.AddHandler("/TestEntities/query", func(w http.ResponseWriter, r *http.Request) {
		// Default values
		pageNum := 1
		pageSize := 10

		// Check if this is a GET request with query parameters
		if r.Method == http.MethodGet {
			query := r.URL.Query()

			// Parse search parameter if it exists
			if searchParam := query.Get("search"); searchParam != "" {
				var searchObj map[string]interface{}
				err := json.Unmarshal([]byte(searchParam), &searchObj)
				if err == nil {
					// Extract page and pageSize from search object
					if page, ok := searchObj["page"].(float64); ok {
						pageNum = int(page)
					}
					if size, ok := searchObj["maxRecords"].(float64); ok {
						pageSize = int(size)
					}
				}
			}
		} else {
			// For POST requests, read from body
			var requestBody map[string]interface{}
			if r.Body != nil {
				err := json.NewDecoder(r.Body).Decode(&requestBody)
				if err == nil {
					// Extract page and pageSize
					if page, ok := requestBody["page"].(float64); ok {
						pageNum = int(page)
					}
					if size, ok := requestBody["pageSize"].(float64); ok {
						pageSize = int(size)
					}
				}
			}
		}

		// Create mock entities
		startID := int64((pageNum-1)*pageSize + 123)
		entities := CreateMockEntities(pageSize, "TestEntity", startID)

		// Create mock response with pagination details
		response := CreateMockListResponse(entities, pageNum, pageSize, 50)

		// Add next/prev page URLs to the response
		pageDetails := response["pageDetails"].(PageDetails)
		if pageNum < 5 { // Assuming we have 5 pages total
			pageDetails.NextPageUrl = "/TestEntities/query?page=" + strconv.Itoa(pageNum+1)
		}
		if pageNum > 1 {
			pageDetails.PrevPageUrl = "/TestEntities/query?page=" + strconv.Itoa(pageNum-1)
		}
		response["pageDetails"] = pageDetails

		// Send response
		server.RespondWithJSON(w, http.StatusOK, response)
	})

	// Create a client
	client := server.NewTestClient()

	// Create a base entity service
	baseService := NewBaseEntityService(client, "TestEntities")
	service := &baseService

	// Fetch a page
	ctx := context.Background()
	options := DefaultPaginationOptions()
	options.Page = 2
	options.PageSize = 10 // Match the pageSize in the handler

	page, err := FetchPage[map[string]interface{}](ctx, service, "name='Test'", options)

	// Verify no error
	AssertNil(t, err, "error should be nil")

	// Verify page
	AssertNotNil(t, page, "page should not be nil")
	AssertEqual(t, 10, len(page.Items), "page should have 10 items")
	AssertEqual(t, 2, page.PageDetails.PageNumber, "page number should be 2")
	AssertEqual(t, 10, page.PageDetails.PageSize, "page size should be 10")
	AssertEqual(t, 50, page.PageDetails.Count, "total count should be 50")

	// Verify first item
	firstItem := page.Items[0]
	AssertEqual(t, float64(123), firstItem["id"], "first item ID should be 123")
}

func TestFetchAllPages(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Add a handler for the query endpoint
	server.AddHandler("/TestEntities/query", func(w http.ResponseWriter, r *http.Request) {
		// Default values
		pageNum := 1
		pageSize := 10

		// Check if this is a GET request with query parameters
		if r.Method == http.MethodGet {
			query := r.URL.Query()

			// Parse search parameter if it exists
			if searchParam := query.Get("search"); searchParam != "" {
				var searchObj map[string]interface{}
				err := json.Unmarshal([]byte(searchParam), &searchObj)
				if err == nil {
					// Extract page and pageSize from search object
					if page, ok := searchObj["page"].(float64); ok {
						pageNum = int(page)
					}
					if size, ok := searchObj["maxRecords"].(float64); ok {
						pageSize = int(size)
					}
				}
			}
		} else {
			// For POST requests, read from body
			var requestBody map[string]interface{}
			if r.Body != nil {
				err := json.NewDecoder(r.Body).Decode(&requestBody)
				if err == nil {
					// Extract page and pageSize
					if page, ok := requestBody["page"].(float64); ok {
						pageNum = int(page)
					}
					if size, ok := requestBody["pageSize"].(float64); ok {
						pageSize = int(size)
					}
				}
			}
		}

		// Create mock entities
		startID := int64((pageNum-1)*pageSize + 123)
		entities := CreateMockEntities(pageSize, "TestEntity", startID)

		// Create mock response with pagination details
		totalCount := 25 // Total of 3 pages
		response := CreateMockListResponse(entities, pageNum, pageSize, totalCount)

		// Add next/prev page URLs to the response
		pageDetails := response["pageDetails"].(PageDetails)
		if pageNum*pageSize < totalCount {
			pageDetails.NextPageUrl = "/TestEntities/query?page=" + strconv.Itoa(pageNum+1)
		}
		if pageNum > 1 {
			pageDetails.PrevPageUrl = "/TestEntities/query?page=" + strconv.Itoa(pageNum-1)
		}
		response["pageDetails"] = pageDetails

		// Send response
		server.RespondWithJSON(w, http.StatusOK, response)
	})

	// Create a client
	client := server.NewTestClient()

	// Create a base entity service
	baseService := NewBaseEntityService(client, "TestEntities")
	service := &baseService

	// Fetch all pages
	ctx := context.Background()

	allItems, err := FetchAllPages[map[string]interface{}](ctx, service, "name='Test'")

	// Verify no error
	AssertNil(t, err, "error should be nil")

	// Verify all items
	AssertNotNil(t, allItems, "all items should not be nil")
	AssertEqual(t, 100, len(allItems), "should have 100 items")

	// Verify first and last items
	firstItem := allItems[0]
	AssertEqual(t, float64(123), firstItem["id"], "first item ID should be 123")

	lastItem := allItems[len(allItems)-1]
	AssertEqual(t, float64(222), lastItem["id"], "last item ID should be 222")
}

func TestFetchAllPagesWithCallback(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Add a handler for the query endpoint
	server.AddHandler("/TestEntities/query", func(w http.ResponseWriter, r *http.Request) {
		// Default values
		pageNum := 1
		pageSize := 10

		// Check if this is a GET request with query parameters
		if r.Method == http.MethodGet {
			query := r.URL.Query()

			// Parse search parameter if it exists
			if searchParam := query.Get("search"); searchParam != "" {
				var searchObj map[string]interface{}
				err := json.Unmarshal([]byte(searchParam), &searchObj)
				if err == nil {
					// Extract page and pageSize from search object
					if page, ok := searchObj["page"].(float64); ok {
						pageNum = int(page)
					}
					if size, ok := searchObj["maxRecords"].(float64); ok {
						pageSize = int(size)
					}
				}
			}
		} else {
			// For POST requests, read from body
			var requestBody map[string]interface{}
			if r.Body != nil {
				err := json.NewDecoder(r.Body).Decode(&requestBody)
				if err == nil {
					// Extract page and pageSize
					if page, ok := requestBody["page"].(float64); ok {
						pageNum = int(page)
					}
					if size, ok := requestBody["pageSize"].(float64); ok {
						pageSize = int(size)
					}
				}
			}
		}

		// Create mock entities
		startID := int64((pageNum-1)*pageSize + 123)
		entities := CreateMockEntities(pageSize, "TestEntity", startID)

		// Create mock response with pagination details
		totalCount := 25 // Total of 3 pages
		response := CreateMockListResponse(entities, pageNum, pageSize, totalCount)

		// Add next/prev page URLs to the response
		pageDetails := response["pageDetails"].(PageDetails)
		if pageNum*pageSize < totalCount {
			pageDetails.NextPageUrl = "/TestEntities/query?page=" + strconv.Itoa(pageNum+1)
		}
		if pageNum > 1 {
			pageDetails.PrevPageUrl = "/TestEntities/query?page=" + strconv.Itoa(pageNum-1)
		}
		response["pageDetails"] = pageDetails

		// Send response
		server.RespondWithJSON(w, http.StatusOK, response)
	})

	// Create a client
	client := server.NewTestClient()

	// Create a base entity service
	baseService := NewBaseEntityService(client, "TestEntities")
	service := &baseService

	// Fetch all pages with callback
	ctx := context.Background()

	// Track callback invocations
	callbackCount := 0
	totalItems := 0

	err := FetchAllPagesWithCallback[map[string]interface{}](ctx, service, "name='Test'", func(items []map[string]interface{}, pageDetails PageDetails) error {
		callbackCount++
		totalItems += len(items)

		// Verify page details
		AssertEqual(t, callbackCount, pageDetails.PageNumber, "page number should match callback count")
		AssertEqual(t, 100, pageDetails.PageSize, "page size should be 100")
		AssertEqual(t, 25, pageDetails.Count, "total count should be 25")

		return nil
	})

	// Verify no error
	AssertNil(t, err, "error should be nil")

	// Verify callback was called for each page
	AssertEqual(t, 1, callbackCount, "callback should be called 1 time")
	AssertEqual(t, 100, totalItems, "total items should be 100")
}

func TestDefaultPaginationOptions(t *testing.T) {
	// Test default options
	options := DefaultPaginationOptions()

	// Verify default values
	AssertEqual(t, 1, options.Page, "default page should be 1")
	AssertEqual(t, 100, options.PageSize, "default page size should be 100")
}
