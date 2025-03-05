package autotask

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	// Create a filter object based on the filter string
	var filterObj QueryFilter
	if filter != "" {
		// Parse the filter string to extract field and value
		parts := strings.Split(filter, "=")
		if len(parts) == 2 {
			field := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if value == "true" {
				filterObj = NewQueryFilter(field, OperatorEquals, true)
			} else if value == "false" {
				filterObj = NewQueryFilter(field, OperatorEquals, false)
			} else if strings.HasPrefix(value, "<") {
				// Handle less than
				num, _ := strconv.Atoi(strings.TrimPrefix(value, "<"))
				filterObj = NewQueryFilter(field, OperatorLessThan, num)
			} else if strings.HasPrefix(value, ">") {
				// Handle greater than
				num, _ := strconv.Atoi(strings.TrimPrefix(value, ">"))
				filterObj = NewQueryFilter(field, OperatorGreaterThan, num)
			} else {
				// Default to equals
				filterObj = NewQueryFilter(field, OperatorEquals, value)
			}
		}
	} else {
		// Default filter for active items
		switch s.EntityName {
		case "Companies":
			filterObj = NewQueryFilter("IsActive", OperatorEquals, true)
		case "Tickets":
			filterObj = NewQueryFilter("Status", OperatorNotEquals, 5) // 5 is typically "Completed"
		case "Resources":
			filterObj = NewQueryFilter("IsActive", OperatorEquals, true)
		}
	}

	// Create query parameters using the existing type
	params := NewEntityQueryParams(filterObj).WithMaxRecords(500)

	// Use the correct endpoint structure according to the API docs
	url := fmt.Sprintf("%s/query", s.EntityName)

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
	url := fmt.Sprintf("%s", s.EntityName)
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
	url := fmt.Sprintf("%s/query/count", s.EntityName)

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
	url := fmt.Sprintf("%s/batch", s.EntityName)
	req, err := s.Client.NewRequest(ctx, http.MethodPost, url, entities)
	if err != nil {
		return err
	}

	_, err = s.Client.Do(req, result)
	return err
}

// BatchUpdate updates multiple entities in a single request.
func (s *BaseEntityService) BatchUpdate(ctx context.Context, entities []interface{}, result interface{}) error {
	url := fmt.Sprintf("%s/batch", s.EntityName)
	req, err := s.Client.NewRequest(ctx, http.MethodPatch, url, entities)
	if err != nil {
		return err
	}

	_, err = s.Client.Do(req, result)
	return err
}

// BatchDelete deletes multiple entities in a single request.
func (s *BaseEntityService) BatchDelete(ctx context.Context, ids []int64) error {
	url := fmt.Sprintf("%s/batch", s.EntityName)
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
