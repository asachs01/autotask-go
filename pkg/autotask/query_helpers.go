package autotask

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Query executes a query against the Autotask API
func (c *client) Query(ctx context.Context, entityName string, params interface{}, response interface{}) error {
	url := entityName + "/query"

	reqBytes, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := c.NewRequest(ctx, "POST", url, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = c.Do(req, response)
	if err != nil {
		return err
	}

	return nil
}

// QueryWithDateFilter executes a query with a date-based filter
func (c *client) QueryWithDateFilter(ctx context.Context, entityName, fieldName string, since time.Time) ([]map[string]interface{}, error) {
	reqBody := map[string]interface{}{
		"filter": []map[string]interface{}{
			{
				"op":    "gte",
				"field": fieldName,
				"value": since.Format(time.RFC3339),
			},
		},
	}

	var response struct {
		Items []map[string]interface{} `json:"items"`
	}

	if err := c.Query(ctx, entityName, reqBody, &response); err != nil {
		return nil, fmt.Errorf("failed to execute date filter query: %w", err)
	}

	return response.Items, nil
}

// QueryWithEmptyFilter executes a query without any filter
func (c *client) QueryWithEmptyFilter(ctx context.Context, entityName string) ([]map[string]interface{}, error) {
	reqBody := map[string]interface{}{
		"filter": []map[string]interface{}{
			{
				"op":    "exist",
				"field": "id",
			},
		},
	}

	var response struct {
		Items []map[string]interface{} `json:"items"`
	}

	if err := c.Query(ctx, entityName, reqBody, &response); err != nil {
		return nil, fmt.Errorf("failed to execute empty filter query: %w", err)
	}

	return response.Items, nil
}

// BatchGetEntities retrieves entities in batches using their IDs
func (c *client) BatchGetEntities(ctx context.Context, entityName string, ids []int64, batchSize int) ([]map[string]interface{}, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	if batchSize <= 0 {
		batchSize = 50 // Default batch size
	}

	var allEntities []map[string]interface{}
	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}
		batch := ids[i:end]

		reqBody := map[string]interface{}{
			"filter": []map[string]interface{}{
				{
					"op":    "in",
					"field": "id",
					"value": batch,
				},
			},
		}

		var response struct {
			Items []map[string]interface{} `json:"items"`
		}

		if err := c.Query(ctx, entityName, reqBody, &response); err != nil {
			return nil, fmt.Errorf("failed to execute batch query: %w", err)
		}

		allEntities = append(allEntities, response.Items...)
	}

	return allEntities, nil
}

// QueryIDRange executes a query to retrieve entities within a specific ID range
func (c *client) QueryIDRange(ctx context.Context, entityName string, startID, endID int64) ([]map[string]interface{}, error) {
	reqBody := map[string]interface{}{
		"filter": []map[string]interface{}{
			{
				"op":    "gte",
				"field": "id",
				"value": startID,
			},
			{
				"op":    "lte",
				"field": "id",
				"value": endID,
			},
		},
	}

	var response struct {
		Items []map[string]interface{} `json:"items"`
	}

	if err := c.Query(ctx, entityName, reqBody, &response); err != nil {
		return nil, fmt.Errorf("failed to execute ID range query: %w", err)
	}

	return response.Items, nil
}
