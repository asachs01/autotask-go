package autotask

import (
	"context"
	"net/http"
	"testing"
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
