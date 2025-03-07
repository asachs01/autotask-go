package autotask

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Response represents a generic response from the Autotask API.
type Response struct {
	Item interface{} `json:"item"`
}

// ListResponse represents a generic list response from the Autotask API.
type ListResponse struct {
	Items       []interface{} `json:"items"`
	PageDetails PageDetails   `json:"pageDetails,omitempty"`
}

// BaseEntityService is a base service that provides common functionality for all entity services.
type BaseEntityService struct {
	Client     Client
	EntityName string
}

// NewBaseEntityService creates a new base entity service
func NewBaseEntityService(client Client, entityName string) BaseEntityService {
	return BaseEntityService{
		Client:     client,
		EntityName: entityName,
	}
}

// Get gets an entity by ID.
func (s *BaseEntityService) Get(ctx context.Context, id int64) (interface{}, error) {
	url := fmt.Sprintf("%s/%d", s.EntityName, id)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var result Response
	_, err = s.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}

	return result.Item, nil
}

// Query queries entities with a filter.
func (s *BaseEntityService) Query(ctx context.Context, filter string, result interface{}) error {
	// Parse the filter string into a query filter or filter group
	var queryFilter interface{}

	if filter != "" {
		// Use the new ParseFilterString function to handle complex filters
		queryFilter = ParseFilterString(filter)
	} else {
		// Default filter for active items
		switch s.EntityName {
		case "Companies":
			queryFilter = NewQueryFilter("IsActive", OperatorEquals, true)
		case "Tickets":
			queryFilter = NewQueryFilter("Status", OperatorNotEquals, 5) // 5 is typically "Completed"
		case "Resources":
			queryFilter = NewQueryFilter("IsActive", OperatorEquals, true)
		case "Projects":
			queryFilter = NewQueryFilter("Status", OperatorNotEquals, 5) // 5 is typically "Completed"
		case "Tasks":
			queryFilter = NewQueryFilter("Status", OperatorNotEquals, 5) // 5 is typically "Completed"
		case "TimeEntries":
			// Default to entries from the last 30 days
			thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
			queryFilter = NewQueryFilter("DateWorked", OperatorGreaterOrEqual, thirtyDaysAgo)
		case "Contracts":
			queryFilter = NewQueryFilter("Status", OperatorEquals, 1) // 1 is typically "Active"
		case "ConfigurationItems":
			queryFilter = NewQueryFilter("Active", OperatorEquals, true)
		}
	}

	// Create query parameters using the enhanced type
	params := NewEntityQueryParams(queryFilter).WithMaxRecords(500)

	// Use the correct endpoint structure according to the API docs
	url := s.EntityName + "/query"

	// Convert params to JSON string for URL parameter
	searchJSON, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal search params: %w", err)
	}

	// Add search parameter to URL
	url = fmt.Sprintf("%s?search=%s", url, string(searchJSON))

	req, err := s.Client.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	// Get the raw response body
	var respBody []byte
	_, err = s.Client.Do(req, &respBody)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	// Unmarshal the response into the result
	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// Create creates a new entity.
func (s *BaseEntityService) Create(ctx context.Context, entity interface{}) (interface{}, error) {
	url := s.EntityName
	req, err := s.Client.NewRequest(ctx, http.MethodPost, url, entity)
	if err != nil {
		return nil, err
	}

	var result Response
	_, err = s.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}

	return result.Item, nil
}

// Update updates an existing entity.
func (s *BaseEntityService) Update(ctx context.Context, id int64, entity interface{}) (interface{}, error) {
	url := fmt.Sprintf("%s/%d", s.EntityName, id)
	req, err := s.Client.NewRequest(ctx, http.MethodPatch, url, entity)
	if err != nil {
		return nil, err
	}

	var result Response
	_, err = s.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}

	return result.Item, nil
}

// Delete deletes an entity by ID.
func (s *BaseEntityService) Delete(ctx context.Context, id int64) error {
	url := fmt.Sprintf("%s/%d", s.EntityName, id)
	req, err := s.Client.NewRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	_, err = s.Client.Do(req, nil)
	return err
}

// Count counts entities matching a filter.
func (s *BaseEntityService) Count(ctx context.Context, filter string) (int, error) {
	// Convert the filter string to the required JSON format
	var fieldName string
	switch s.EntityName {
	case "Companies":
		fieldName = "IsActive"
	case "Tickets":
		fieldName = "Status"
	case "Resources":
		fieldName = "IsActive"
	default:
		fieldName = filter
	}

	filterObj := NewQueryFilter(fieldName, OperatorEquals, true)

	// Create query parameters using the existing type
	params := NewEntityQueryParams(filterObj)

	// Use the correct endpoint structure according to the API docs
	url := s.EntityName + "/query/count"

	// Convert params to JSON string for URL parameter
	searchJSON, err := json.Marshal(params)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal search params: %w", err)
	}

	// Add search parameter to URL
	url = fmt.Sprintf("%s?search=%s", url, string(searchJSON))

	req, err := s.Client.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	var count struct {
		Count int `json:"count"`
	}
	_, err = s.Client.Do(req, &count)
	if err != nil {
		return 0, err
	}

	return count.Count, nil
}

// Pagination handles paginated results.
func (s *BaseEntityService) Pagination(ctx context.Context, url string, result interface{}) error {
	req, err := s.Client.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	_, err = s.Client.Do(req, result)
	return err
}

// BatchCreate creates multiple entities in a single request.
func (s *BaseEntityService) BatchCreate(ctx context.Context, entities []interface{}, result interface{}) error {
	url := s.EntityName + "/batch"
	req, err := s.Client.NewRequest(ctx, http.MethodPost, url, entities)
	if err != nil {
		return err
	}

	_, err = s.Client.Do(req, result)
	return err
}

// BatchUpdate updates multiple entities in a single request.
func (s *BaseEntityService) BatchUpdate(ctx context.Context, entities []interface{}, result interface{}) error {
	url := s.EntityName + "/batch"
	req, err := s.Client.NewRequest(ctx, http.MethodPatch, url, entities)
	if err != nil {
		return err
	}

	_, err = s.Client.Do(req, result)
	return err
}

// BatchDelete deletes multiple entities in a single request.
func (s *BaseEntityService) BatchDelete(ctx context.Context, ids []int64) error {
	url := s.EntityName + "/batch"
	req, err := s.Client.NewRequest(ctx, http.MethodDelete, url, ids)
	if err != nil {
		return err
	}

	_, err = s.Client.Do(req, nil)
	return err
}

// GetNextPage gets the next page of results
func (s *BaseEntityService) GetNextPage(ctx context.Context, pageDetails PageDetails) ([]interface{}, error) {
	var result ListResponse
	err := s.Pagination(ctx, pageDetails.NextPageUrl, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// GetPreviousPage gets the previous page of results
func (s *BaseEntityService) GetPreviousPage(ctx context.Context, pageDetails PageDetails) ([]interface{}, error) {
	var result ListResponse
	err := s.Pagination(ctx, pageDetails.PrevPageUrl, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}
