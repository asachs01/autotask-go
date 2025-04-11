package autotask

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	neturl "net/url"
)

// PaginationIterator provides an iterator pattern for paginated results
type PaginationIterator struct {
	service      EntityService
	filter       string
	currentPage  int
	pageSize     int
	totalCount   int
	hasMore      bool
	nextPageUrl  string
	prevPageUrl  string
	items        []interface{}
	currentIndex int
	ctx          context.Context
}

// NewPaginationIterator creates a new pagination iterator
func NewPaginationIterator(ctx context.Context, service EntityService, filter string, pageSize int) (*PaginationIterator, error) {
	iterator := &PaginationIterator{
		service:      service,
		filter:       filter,
		currentPage:  1,
		pageSize:     pageSize,
		currentIndex: -1,
		ctx:          ctx,
	}

	// Load the first page
	err := iterator.loadNextPage()
	if err != nil {
		return nil, err
	}

	return iterator, nil
}

// Next advances the iterator to the next item
// Returns false when there are no more items
func (p *PaginationIterator) Next() bool {
	p.currentIndex++

	// If we've reached the end of the current page and there are more pages
	if p.currentIndex >= len(p.items) && p.hasMore {
		err := p.loadNextPage()
		if err != nil {
			return false
		}
		p.currentIndex = 0
	}

	return p.currentIndex < len(p.items)
}

// Item returns the current item
func (p *PaginationIterator) Item() interface{} {
	if p.currentIndex < 0 || p.currentIndex >= len(p.items) {
		return nil
	}
	return p.items[p.currentIndex]
}

// Error returns any error that occurred during iteration
func (p *PaginationIterator) Error() error {
	return nil
}

// TotalCount returns the total number of items across all pages
func (p *PaginationIterator) TotalCount() int {
	return p.totalCount
}

// CurrentPage returns the current page number
func (p *PaginationIterator) CurrentPage() int {
	return p.currentPage
}

// loadNextPage loads the next page of results
func (p *PaginationIterator) loadNextPage() error {
	var urlPath string
	var req *http.Request
	var err error

	// Create a response structure
	var response struct {
		Items       []interface{} `json:"items"`
		PageDetails PageDetails   `json:"pageDetails"`
	}

	// If this is the first page, use the regular query endpoint
	if p.currentPage == 1 {
		// Parse the filter string
		var queryFilter interface{}
		if p.filter != "" {
			queryFilter = ParseFilterString(p.filter)
		}

		// Create query parameters
		params := NewEntityQueryParams(queryFilter).
			WithMaxRecords(p.pageSize)

		// Use the correct endpoint structure according to the API docs
		urlPath = p.service.GetEntityName() + "/query"

		// Convert params to JSON string for URL parameter
		searchJSON, err := json.Marshal(params)
		if err != nil {
			return fmt.Errorf("failed to marshal search params: %w", err)
		}

		// Add search parameter to URL
		searchJSONStr := string(searchJSON)
		escapedJSON := neturl.QueryEscape(searchJSONStr)
		urlPath = fmt.Sprintf("%s?search=%s", urlPath, escapedJSON)

		req, err = p.service.GetClient().NewRequest(p.ctx, http.MethodGet, urlPath, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
	} else if p.nextPageUrl != "" {
		// For subsequent pages, use the nextPageUrl from the previous response
		// The nextPageUrl already contains all necessary parameters
		req, err = p.service.GetClient().NewRequest(p.ctx, http.MethodGet, p.nextPageUrl, nil)
		if err != nil {
			return fmt.Errorf("failed to create request for next page: %w", err)
		}
	} else {
		// No more pages
		return nil
	}

	_, err = p.service.GetClient().Do(req, &response)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	// Ensure the page number is correctly set
	// The API might not always return the correct page number
	response.PageDetails.PageNumber = p.currentPage

	// Update the iterator state
	p.items = response.Items
	p.totalCount = response.PageDetails.Count
	p.nextPageUrl = response.PageDetails.NextPageUrl
	p.prevPageUrl = response.PageDetails.PrevPageUrl
	p.hasMore = response.PageDetails.NextPageUrl != "" && len(response.Items) > 0
	p.currentPage++

	return nil
}

// PaginatedResults is a generic structure for paginated results
type PaginatedResults[T any] struct {
	Items       []T         `json:"items"`
	PageDetails PageDetails `json:"pageDetails"`
}

// PaginationOptions contains options for pagination
type PaginationOptions struct {
	PageSize int
	Page     int
}

// DefaultPaginationOptions returns default pagination options
func DefaultPaginationOptions() PaginationOptions {
	return PaginationOptions{
		PageSize: 100,
		Page:     1,
	}
}

// FetchAllPages is a convenience method to fetch all pages of results
// This is useful when you need all results and don't want to manually handle pagination
func FetchAllPages[T any](ctx context.Context, service EntityService, filter string) ([]T, error) {
	var allItems []T
	var nextPageUrl string
	pageSize := 100 // Default page size

	// First page
	var response struct {
		Items       []T         `json:"items"`
		PageDetails PageDetails `json:"pageDetails"`
	}

	// Parse the filter string
	var queryFilter interface{}
	if filter != "" {
		queryFilter = ParseFilterString(filter)
	}

	// Create query parameters
	params := NewEntityQueryParams(queryFilter).
		WithMaxRecords(pageSize)

	// Use the correct endpoint structure according to the API docs
	urlPath := service.GetEntityName() + "/query"

	// Convert params to JSON string for URL parameter
	searchJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search params: %w", err)
	}

	// Add search parameter to URL
	searchJSONStr := string(searchJSON)
	escapedJSON := neturl.QueryEscape(searchJSONStr)
	urlPath = fmt.Sprintf("%s?search=%s", urlPath, escapedJSON)

	req, err := service.GetClient().NewRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	_, err = service.GetClient().Do(req, &response)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	// Add items to the result
	allItems = append(allItems, response.Items...)
	nextPageUrl = response.PageDetails.NextPageUrl

	// Fetch subsequent pages using the nextPageUrl
	for nextPageUrl != "" && len(response.Items) > 0 {
		req, err := service.GetClient().NewRequest(ctx, http.MethodGet, nextPageUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request for next page: %w", err)
		}

		_, err = service.GetClient().Do(req, &response)
		if err != nil {
			return nil, fmt.Errorf("query failed: %w", err)
		}

		// Add items to the result
		allItems = append(allItems, response.Items...)
		nextPageUrl = response.PageDetails.NextPageUrl
	}

	return allItems, nil
}

// FetchAllPagesWithCallback is a convenience method to fetch all pages of results
// and process each page with a callback function
// This is useful for processing large result sets without keeping all items in memory
func FetchAllPagesWithCallback[T any](
	ctx context.Context,
	service EntityService,
	filter string,
	callback func(items []T, pageDetails PageDetails) error,
) error {
	var nextPageUrl string
	pageSize := 100  // Default page size
	currentPage := 1 // Start with first page

	// First page
	var response struct {
		Items       []T         `json:"items"`
		PageDetails PageDetails `json:"pageDetails"`
	}

	// Parse the filter string
	var queryFilter interface{}
	if filter != "" {
		queryFilter = ParseFilterString(filter)
	}

	// Create query parameters
	params := NewEntityQueryParams(queryFilter).
		WithMaxRecords(pageSize)

	// Use the correct endpoint structure according to the API docs
	urlPath := service.GetEntityName() + "/query"

	// Convert params to JSON string for URL parameter
	searchJSON, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal search params: %w", err)
	}

	// Add search parameter to URL
	searchJSONStr := string(searchJSON)
	escapedJSON := neturl.QueryEscape(searchJSONStr)
	urlPath = fmt.Sprintf("%s?search=%s", urlPath, escapedJSON)

	req, err := service.GetClient().NewRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	_, err = service.GetClient().Do(req, &response)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	// Ensure the page number is correctly set
	response.PageDetails.PageNumber = currentPage

	// Process this page with the callback
	err = callback(response.Items, response.PageDetails)
	if err != nil {
		return err
	}

	nextPageUrl = response.PageDetails.NextPageUrl
	currentPage++

	// Fetch subsequent pages using the nextPageUrl
	for nextPageUrl != "" && len(response.Items) > 0 {
		req, err := service.GetClient().NewRequest(ctx, http.MethodGet, nextPageUrl, nil)
		if err != nil {
			return fmt.Errorf("failed to create request for next page: %w", err)
		}

		_, err = service.GetClient().Do(req, &response)
		if err != nil {
			return fmt.Errorf("query failed: %w", err)
		}

		// Ensure the page number is correctly set
		response.PageDetails.PageNumber = currentPage

		// Process this page with the callback
		err = callback(response.Items, response.PageDetails)
		if err != nil {
			return err
		}

		nextPageUrl = response.PageDetails.NextPageUrl
		currentPage++
	}

	return nil
}

// FetchPage is a convenience method to fetch a specific page of results
// Note: According to Autotask API documentation, you should use the pagination URLs
// provided in the response rather than manually constructing page requests.
// This method is provided for convenience but may not work as expected for all cases.
func FetchPage[T any](
	ctx context.Context,
	service EntityService,
	filter string,
	options PaginationOptions,
) (*PaginatedResults[T], error) {
	// Create a response structure
	var response PaginatedResults[T]

	// Parse the filter string
	var queryFilter interface{}
	if filter != "" {
		queryFilter = ParseFilterString(filter)
	}

	// Create query parameters
	params := NewEntityQueryParams(queryFilter).
		WithMaxRecords(options.PageSize).
		WithPage(options.Page)

	// Use the correct endpoint structure according to the API docs
	urlPath := service.GetEntityName() + "/query"

	// Convert params to JSON string for URL parameter
	searchJSON, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search params: %w", err)
	}

	// Add search parameter to URL
	searchJSONStr := string(searchJSON)
	escapedJSON := neturl.QueryEscape(searchJSONStr)
	urlPath = fmt.Sprintf("%s?search=%s", urlPath, escapedJSON)

	req, err := service.GetClient().NewRequest(ctx, http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	_, err = service.GetClient().Do(req, &response)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	// Ensure the page number is correctly set
	response.PageDetails.PageNumber = options.Page

	// If options.Page > 1, we need to follow the nextPageUrl to get to the requested page
	nextPageUrl := response.PageDetails.NextPageUrl
	currentPage := 1

	for currentPage < options.Page && nextPageUrl != "" {
		req, err := service.GetClient().NewRequest(ctx, http.MethodGet, nextPageUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request for page %d: %w", currentPage+1, err)
		}

		_, err = service.GetClient().Do(req, &response)
		if err != nil {
			return nil, fmt.Errorf("query failed for page %d: %w", currentPage+1, err)
		}

		currentPage++
		nextPageUrl = response.PageDetails.NextPageUrl

		// If we've reached the requested page, break
		if currentPage >= options.Page {
			break
		}

		// If there are no more pages, break
		if nextPageUrl == "" {
			break
		}
	}

	// Ensure the page number is correctly set
	response.PageDetails.PageNumber = currentPage

	return &response, nil
}
