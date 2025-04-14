package autotask

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertNotEqual asserts that two values are not equal
func AssertNotEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if expected == actual {
		t.Errorf("%s: expected %v to not equal %v", message, expected, actual)
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient("test-user", "test-secret", "test-integration-code")

	// Verify client is not nil
	AssertNotNil(t, client, "client should not be nil")

	// Test service accessors to verify initialization
	AssertNotNil(t, client.Companies(), "Companies service should not be nil")
	AssertNotNil(t, client.Tickets(), "Tickets service should not be nil")
	AssertNotNil(t, client.Contacts(), "Contacts service should not be nil")
	AssertNotNil(t, client.Webhooks(), "Webhooks service should not be nil")
	AssertNotNil(t, client.Resources(), "Resources service should not be nil")
	AssertNotNil(t, client.Projects(), "Projects service should not be nil")
	AssertNotNil(t, client.Tasks(), "Tasks service should not be nil")
	AssertNotNil(t, client.TimeEntries(), "TimeEntries service should not be nil")
	AssertNotNil(t, client.Contracts(), "Contracts service should not be nil")
	AssertNotNil(t, client.ConfigurationItems(), "ConfigurationItems service should not be nil")
}

func TestGetZoneInfo(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Create a client
	client := server.NewTestClient()

	// Get zone info
	zoneInfo, err := client.GetZoneInfo()

	// Verify no error
	AssertNil(t, err, "error should be nil")

	// Verify zone info
	AssertNotNil(t, zoneInfo, "zone info should not be nil")
	AssertEqual(t, "MockZone", zoneInfo.ZoneName, "zone name should match")
	AssertContains(t, zoneInfo.URL, server.Server.URL, "URL should contain server URL")
	AssertContains(t, zoneInfo.WebURL, server.Server.URL, "WebURL should contain server URL")
	AssertEqual(t, 1, zoneInfo.CI, "CI should match")
}

func TestNewRequestAndDo(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Add a handler for the test endpoint
	server.AddHandler("/test", func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		AssertEqual(t, http.MethodGet, r.Method, "method should match")

		// Send a response
		server.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "success",
		})
	})

	// Create a client
	client := server.NewTestClient()

	// Create a context
	ctx := context.Background()

	// Create a request using the client's NewRequest method
	req, err := client.NewRequest(ctx, http.MethodGet, "/test", nil)
	AssertNil(t, err, "error should be nil")

	// Verify headers
	AssertEqual(t, DefaultUserAgent, req.Header.Get("User-Agent"), "User-Agent header should match")
	AssertEqual(t, "application/json", req.Header.Get("Content-Type"), "Content-Type header should match")
	AssertEqual(t, "application/json", req.Header.Get("Accept"), "Accept header should match")
	AssertNotEqual(t, "", req.Header.Get("Authorization"), "Authorization header should not be empty")
	AssertEqual(t, "test-user", req.Header.Get("UserName"), "UserName header should match")
	AssertEqual(t, "test-secret", req.Header.Get("Secret"), "Secret header should match")
	AssertEqual(t, "test-integration-code", req.Header.Get("ApiIntegrationCode"), "ApiIntegrationCode header should match")

	// Send the request
	var response map[string]string
	resp, err := client.Do(req, &response)

	// Verify no error
	AssertNil(t, err, "error should be nil")

	// Verify response
	AssertNotNil(t, resp, "response should not be nil")
	AssertEqual(t, http.StatusOK, resp.StatusCode, "status code should match")
	AssertEqual(t, "success", response["message"], "response message should match")
}

func TestDoWithError(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Add a handler for the test endpoint
	server.AddHandler("/test-error", func(w http.ResponseWriter, r *http.Request) {
		// Send an error response
		server.RespondWithError(w, http.StatusBadRequest, "Bad request", []string{"Error 1", "Error 2"})
	})

	// Create a client
	client := server.NewTestClient()

	// Create a context
	ctx := context.Background()

	// Create a request
	req, err := client.NewRequest(ctx, http.MethodGet, "/test-error", nil)
	AssertNil(t, err, "error should be nil")

	// Send the request
	var response map[string]string
	_, err = client.Do(req, &response)

	// Verify error
	AssertNotNil(t, err, "error should not be nil")

	// Verify error type
	errorResp, ok := err.(*ErrorResponse)
	AssertTrue(t, ok, "error should be of type *ErrorResponse")

	// Verify error fields
	AssertEqual(t, "Bad request", errorResp.Message, "error message should match")
	AssertLen(t, errorResp.Errors, 2, "error should have 2 errors")
	AssertEqual(t, "Error 1", errorResp.Errors[0], "first error should match")
	AssertEqual(t, "Error 2", errorResp.Errors[1], "second error should match")
}

func TestServiceAccessors(t *testing.T) {
	// Create a client
	client := NewClient("test-user", "test-secret", "test-integration-code")

	// Test Companies service
	companies := client.Companies()
	AssertNotNil(t, companies, "Companies service should not be nil")
	AssertEqual(t, "Companies", companies.GetEntityName(), "entity name should match")

	// Test Tickets service
	tickets := client.Tickets()
	AssertNotNil(t, tickets, "Tickets service should not be nil")
	AssertEqual(t, "Tickets", tickets.GetEntityName(), "entity name should match")

	// Test Contacts service
	contacts := client.Contacts()
	AssertNotNil(t, contacts, "Contacts service should not be nil")
	AssertEqual(t, "Contacts", contacts.GetEntityName(), "entity name should match")

	// Test Webhooks service
	webhooks := client.Webhooks()
	AssertNotNil(t, webhooks, "Webhooks service should not be nil")
	AssertEqual(t, "Webhooks", webhooks.GetEntityName(), "entity name should match")

	// Test Resources service
	resources := client.Resources()
	AssertNotNil(t, resources, "Resources service should not be nil")
	AssertEqual(t, "Resources", resources.GetEntityName(), "entity name should match")

	// Test Projects service
	projects := client.Projects()
	AssertNotNil(t, projects, "Projects service should not be nil")
	AssertEqual(t, "Projects", projects.GetEntityName(), "entity name should match")

	// Test Tasks service
	tasks := client.Tasks()
	AssertNotNil(t, tasks, "Tasks service should not be nil")
	AssertEqual(t, "Tasks", tasks.GetEntityName(), "entity name should match")

	// Test TimeEntries service
	timeEntries := client.TimeEntries()
	AssertNotNil(t, timeEntries, "TimeEntries service should not be nil")
	AssertEqual(t, "TimeEntries", timeEntries.GetEntityName(), "entity name should match")

	// Test Contracts service
	contracts := client.Contracts()
	AssertNotNil(t, contracts, "Contracts service should not be nil")
	AssertEqual(t, "Contracts", contracts.GetEntityName(), "entity name should match")

	// Test ConfigurationItems service
	configItems := client.ConfigurationItems()
	AssertNotNil(t, configItems, "ConfigurationItems service should not be nil")
	AssertEqual(t, "ConfigurationItems", configItems.GetEntityName(), "entity name should match")
}

func TestSetLogLevel(t *testing.T) {
	// Create a client
	client := NewClient("test-user", "test-secret", "test-integration-code")

	// Set log level
	client.SetLogLevel(LogLevelDebug)

	// We can't directly verify the log level since the client is unexported,
	// but we can verify that the method doesn't panic
}

func TestSetDebugMode(t *testing.T) {
	// Create a client
	client := NewClient("test-user", "test-secret", "test-integration-code")

	// Set debug mode
	client.SetDebugMode(true)

	// We can't directly verify the debug mode since the client is unexported,
	// but we can verify that the method doesn't panic
}

func TestErrorResponse(t *testing.T) {
	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "https://example.com/test", nil)

	// Create an error response with a valid Response field
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "400 Bad Request",
		Request:    req,
	}

	errorResp := &ErrorResponse{
		Response: resp,
		Message:  "Bad request",
		Errors:   []string{"Error 1", "Error 2"},
	}

	// Get the error message
	errorMsg := errorResp.Error()

	// Print the actual error message for debugging
	t.Logf("Actual error message: %s", errorMsg)

	// Verify error message contains expected parts
	AssertContains(t, errorMsg, "GET", "error message should contain the request method")
	AssertContains(t, errorMsg, "example.com", "error message should contain the request URL")
	AssertContains(t, errorMsg, "400", "error message should contain the status code")
	AssertContains(t, errorMsg, "Error 1", "error message should contain the first error")
}

func TestQueryEndpoint(t *testing.T) {
	tests := []struct {
		name       string
		setupTest  func() (*httptest.Server, Client)
		testQuery  func(Client) error
		wantFilter map[string]interface{}
	}{
		{
			name: "empty filter query",
			setupTest: func() (*httptest.Server, Client) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Verify the request method
					assert.Equal(t, http.MethodPost, r.Method)

					// Verify the endpoint
					assert.Equal(t, "/ATServicesRest/V1.0/TestEntity/query", r.URL.Path)

					// Parse and verify the request body
					var requestBody map[string]interface{}
					err := json.NewDecoder(r.Body).Decode(&requestBody)
					require.NoError(t, err)

					// Verify the filter structure
					filters, ok := requestBody["filter"].([]interface{})
					require.True(t, ok, "filter should be an array")
					require.Len(t, filters, 1, "filter should have one item")

					filterObj, ok := filters[0].(map[string]interface{})
					require.True(t, ok, "filter item should be an object")
					assert.Equal(t, "and", filterObj["op"])

					items, ok := filterObj["items"].([]interface{})
					require.True(t, ok, "items should be an array")
					require.Len(t, items, 1, "items should have one item")

					item, ok := items[0].(map[string]interface{})
					require.True(t, ok, "item should be an object")
					assert.Equal(t, "exist", item["op"])
					assert.Equal(t, "id", item["field"])

					// Return a mock response
					response := map[string]interface{}{
						"items": []map[string]interface{}{
							{"id": 1, "name": "Test Item 1"},
							{"id": 2, "name": "Test Item 2"},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					err = json.NewEncoder(w).Encode(response)
					require.NoError(t, err)
				}))

				// Create a test client
				client := NewClient(
					"test-user",
					"test-secret",
					"test-integration-code",
				)

				return server, client
			},
			testQuery: func(client Client) error {
				ctx := context.Background()
				var result struct {
					Items []map[string]interface{} `json:"items"`
				}
				req, err := client.NewRequest(ctx, http.MethodPost, "/ATServicesRest/V1.0/TestEntity/query", map[string]interface{}{
					"filter": []map[string]interface{}{
						{
							"op": "and",
							"items": []map[string]interface{}{
								{
									"op":    "exist",
									"field": "id",
								},
							},
						},
					},
				})
				if err != nil {
					return err
				}

				_, err = client.Do(req, &result)
				if err != nil {
					return err
				}

				require.Len(t, result.Items, 2)
				assert.Equal(t, float64(1), result.Items[0]["id"])
				assert.Equal(t, "Test Item 1", result.Items[0]["name"])
				assert.Equal(t, float64(2), result.Items[1]["id"])
				assert.Equal(t, "Test Item 2", result.Items[1]["name"])

				return nil
			},
		},
		{
			name: "date filter query",
			setupTest: func() (*httptest.Server, Client) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Verify the request method
					assert.Equal(t, http.MethodPost, r.Method)

					// Verify the endpoint
					assert.Equal(t, "/ATServicesRest/V1.0/TestEntity/query", r.URL.Path)

					// Parse and verify the request body
					var requestBody map[string]interface{}
					err := json.NewDecoder(r.Body).Decode(&requestBody)
					require.NoError(t, err)

					// Verify the filter structure
					filters, ok := requestBody["filter"].([]interface{})
					require.True(t, ok, "filter should be an array")
					require.Len(t, filters, 1, "filter should have one item")

					filterObj, ok := filters[0].(map[string]interface{})
					require.True(t, ok, "filter item should be an object")
					assert.Equal(t, "and", filterObj["op"])

					items, ok := filterObj["items"].([]interface{})
					require.True(t, ok, "items should be an array")
					require.Len(t, items, 1, "items should have one item")

					item, ok := items[0].(map[string]interface{})
					require.True(t, ok, "item should be an object")
					assert.Equal(t, "gte", item["op"])
					assert.Equal(t, "testDate", item["field"])

					// Verify the date format
					dateStr, ok := item["value"].(string)
					require.True(t, ok, "value should be a string")
					_, err = time.Parse(time.RFC3339, dateStr)
					require.NoError(t, err)

					// Return a mock response
					response := map[string]interface{}{
						"items": []map[string]interface{}{
							{"id": 1, "testDate": "2024-01-01T00:00:00Z"},
							{"id": 2, "testDate": "2024-01-02T00:00:00Z"},
						},
					}
					w.Header().Set("Content-Type", "application/json")
					err = json.NewEncoder(w).Encode(response)
					require.NoError(t, err)
				}))

				// Create a test client
				client := NewClient(
					"test-user",
					"test-secret",
					"test-integration-code",
				)

				return server, client
			},
			testQuery: func(client Client) error {
				ctx := context.Background()
				since := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				var result struct {
					Items []map[string]interface{} `json:"items"`
				}
				req, err := client.NewRequest(ctx, http.MethodPost, "/ATServicesRest/V1.0/TestEntity/query", map[string]interface{}{
					"filter": []map[string]interface{}{
						{
							"op": "and",
							"items": []map[string]interface{}{
								{
									"op":    "gte",
									"field": "testDate",
									"value": since.Format(time.RFC3339),
								},
							},
						},
					},
				})
				if err != nil {
					return err
				}

				_, err = client.Do(req, &result)
				if err != nil {
					return err
				}

				require.Len(t, result.Items, 2)
				assert.Equal(t, float64(1), result.Items[0]["id"])
				assert.Equal(t, "2024-01-01T00:00:00Z", result.Items[0]["testDate"])
				assert.Equal(t, float64(2), result.Items[1]["id"])
				assert.Equal(t, "2024-01-02T00:00:00Z", result.Items[1]["testDate"])

				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, client := tt.setupTest()
			defer server.Close()

			err := tt.testQuery(client)
			require.NoError(t, err)
		})
	}
}
