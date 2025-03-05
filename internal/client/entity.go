package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/asachs01/autotask-go/pkg/autotask"
)

// EntityService implements the base entity service interface
type EntityService struct {
	client     *Client
	entityName string
}

// Get gets an entity by ID
func (s *EntityService) Get(ctx context.Context, id int64) (interface{}, error) {
	path := fmt.Sprintf("/%s/%d", s.entityName, id)
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.client.do(req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Query queries entities with a filter
func (s *EntityService) Query(ctx context.Context, filter string, result interface{}) error {
	path := "/" + s.entityName + "/query"
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	return s.client.do(req, result)
}

// Create creates a new entity
func (s *EntityService) Create(ctx context.Context, entity interface{}) (interface{}, error) {
	path := "/" + s.entityName
	req, err := s.client.newRequest(ctx, http.MethodPost, path, entity)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.client.do(req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Update updates an existing entity
func (s *EntityService) Update(ctx context.Context, id int64, entity interface{}) (interface{}, error) {
	path := fmt.Sprintf("/%s/%d", s.entityName, id)
	req, err := s.client.newRequest(ctx, http.MethodPatch, path, entity)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := s.client.do(req, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Delete deletes an entity by ID
func (s *EntityService) Delete(ctx context.Context, id int64) error {
	path := fmt.Sprintf("/%s/%d", s.entityName, id)
	req, err := s.client.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}

// Count counts entities matching a filter
func (s *EntityService) Count(ctx context.Context, filter string) (int, error) {
	path := "/" + s.entityName + "/count"
	req, err := s.client.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return 0, err
	}

	var count struct {
		Count int `json:"count"`
	}
	if err := s.client.do(req, &count); err != nil {
		return 0, err
	}

	return count.Count, nil
}

// Pagination handles paginated results
func (s *EntityService) Pagination(ctx context.Context, url string, result interface{}) error {
	req, err := s.client.newRequest(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	return s.client.do(req, result)
}

// BatchCreate creates multiple entities in a single request
func (s *EntityService) BatchCreate(ctx context.Context, entities []interface{}, result interface{}) error {
	path := "/" + s.entityName + "/batch"
	req, err := s.client.newRequest(ctx, http.MethodPost, path, entities)
	if err != nil {
		return err
	}

	return s.client.do(req, result)
}

// BatchUpdate updates multiple entities in a single request
func (s *EntityService) BatchUpdate(ctx context.Context, entities []interface{}, result interface{}) error {
	path := "/" + s.entityName + "/batch"
	req, err := s.client.newRequest(ctx, http.MethodPatch, path, entities)
	if err != nil {
		return err
	}

	return s.client.do(req, result)
}

// BatchDelete deletes multiple entities in a single request
func (s *EntityService) BatchDelete(ctx context.Context, ids []int64) error {
	path := "/" + s.entityName + "/batch"
	req, err := s.client.newRequest(ctx, http.MethodDelete, path, ids)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}

// GetNextPage gets the next page of results
func (s *EntityService) GetNextPage(ctx context.Context, pageDetails autotask.PageDetails) ([]interface{}, error) {
	if pageDetails.NextPageUrl == "" {
		return nil, nil
	}

	var result autotask.ListResponse
	err := s.Pagination(ctx, pageDetails.NextPageUrl, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// GetPreviousPage gets the previous page of results
func (s *EntityService) GetPreviousPage(ctx context.Context, pageDetails autotask.PageDetails) ([]interface{}, error) {
	if pageDetails.PrevPageUrl == "" {
		return nil, nil
	}

	var result autotask.ListResponse
	err := s.Pagination(ctx, pageDetails.PrevPageUrl, &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}
