package autotask

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// MockResponse represents a mock API response
type MockResponse struct {
	StatusCode int
	Body       interface{}
	Headers    map[string]string
}

// MockServer represents a mock Autotask API server for testing
type MockServer struct {
	Server           *httptest.Server
	Mux              *http.ServeMux
	ResponseHandlers map[string]func(w http.ResponseWriter, r *http.Request)
	Requests         []*http.Request
	RequestBodies    [][]byte
	t                *testing.T
}

// NewMockServer creates a new mock Autotask API server
func NewMockServer(t *testing.T) *MockServer {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mock := &MockServer{
		Server:           server,
		Mux:              mux,
		ResponseHandlers: make(map[string]func(w http.ResponseWriter, r *http.Request)),
		Requests:         make([]*http.Request, 0),
		RequestBodies:    make([][]byte, 0),
		t:                t,
	}

	// Setup default zone info handler
	mock.AddHandler("/ZoneInformation", func(w http.ResponseWriter, r *http.Request) {
		zoneInfo := ZoneInfo{
			ZoneName: "MockZone",
			URL:      fmt.Sprintf("%s/atservicesrest/", server.URL),
			WebURL:   fmt.Sprintf("%s/web/", server.URL),
			CI:       1,
		}
		mock.RespondWithJSON(w, http.StatusOK, zoneInfo)
	})

	return mock
}

// AddHandler adds a response handler for a specific path
func (m *MockServer) AddHandler(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	m.ResponseHandlers[path] = handler
	m.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		// Record the request for later inspection
		m.Requests = append(m.Requests, r)

		// Read and record the request body
		if r.Body != nil {
			body := make([]byte, r.ContentLength)
			if _, err := r.Body.Read(body); err != nil {
				m.t.Errorf("Failed to read request body: %v", err)
				return
			}
			m.RequestBodies = append(m.RequestBodies, body)

			// Reset the body for the handler
			r.Body = &readCloser{strings.NewReader(string(body))}
		}

		// Call the handler
		handler(w, r)
	})

	// Also handle the path with the API version prefix if it doesn't already have it
	if !strings.HasPrefix(path, "/atservicesrest/v1.0") && !strings.HasPrefix(path, "/ZoneInformation") {
		apiPath := "/atservicesrest/v1.0" + path
		m.Mux.HandleFunc(apiPath, func(w http.ResponseWriter, r *http.Request) {
			// Record the request for later inspection
			m.Requests = append(m.Requests, r)

			// Read and record the request body
			if r.Body != nil {
				body := make([]byte, r.ContentLength)
				if _, err := r.Body.Read(body); err != nil {
					m.t.Errorf("Failed to read request body: %v", err)
					return
				}
				m.RequestBodies = append(m.RequestBodies, body)

				// Reset the body for the handler
				r.Body = &readCloser{strings.NewReader(string(body))}
			}

			// Call the handler
			handler(w, r)
		})
	}
}

// RespondWithJSON sends a JSON response with the given status code and body
func (m *MockServer) RespondWithJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if body != nil {
		if err := json.NewEncoder(w).Encode(body); err != nil {
			m.t.Fatalf("failed to encode response body: %v", err)
		}
	}
}

// RespondWithError sends an error response with the given status code and message
func (m *MockServer) RespondWithError(w http.ResponseWriter, statusCode int, message string, errors []string) {
	errorResponse := ErrorResponse{
		Message: message,
		Errors:  errors,
	}
	m.RespondWithJSON(w, statusCode, errorResponse)
}

// Close closes the mock server
func (m *MockServer) Close() {
	m.Server.Close()
}

// NewTestClient creates a new client that uses the mock server
func (m *MockServer) NewTestClient() Client {
	c := NewClient("test-user", "test-secret", "test-integration-code").(*client)

	// Override the base zone info URL to use the mock server
	baseURL, _ := url.Parse(m.Server.URL)
	c.baseURL = baseURL

	// Set the zone info directly to avoid making an HTTP request
	c.zoneInfo = &ZoneInfo{
		ZoneName: "MockZone",
		URL:      fmt.Sprintf("%s/atservicesrest/", m.Server.URL),
		WebURL:   fmt.Sprintf("%s/web/", m.Server.URL),
		CI:       1,
	}

	return c
}

// GetLastRequest returns the last request made to the mock server
func (m *MockServer) GetLastRequest() *http.Request {
	if len(m.Requests) == 0 {
		return nil
	}
	return m.Requests[len(m.Requests)-1]
}

// GetLastRequestBody returns the body of the last request made to the mock server
func (m *MockServer) GetLastRequestBody() []byte {
	if len(m.RequestBodies) == 0 {
		return nil
	}
	return m.RequestBodies[len(m.RequestBodies)-1]
}

// readCloser is a helper type that implements io.ReadCloser
type readCloser struct {
	*strings.Reader
}

// Close implements io.Closer
func (r *readCloser) Close() error {
	return nil
}

// CreateMockListResponse creates a mock list response with the given items
func CreateMockListResponse(items []interface{}, pageNumber, pageSize, totalCount int) map[string]interface{} {
	nextPageUrl := ""
	prevPageUrl := ""

	if pageNumber*pageSize < totalCount {
		nextPageUrl = fmt.Sprintf("?page=%d&pageSize=%d", pageNumber+1, pageSize)
	}

	if pageNumber > 1 {
		prevPageUrl = fmt.Sprintf("?page=%d&pageSize=%d", pageNumber-1, pageSize)
	}

	pageDetails := PageDetails{
		PageNumber:  pageNumber,
		PageSize:    pageSize,
		Count:       totalCount,
		NextPageUrl: nextPageUrl,
		PrevPageUrl: prevPageUrl,
	}

	return map[string]interface{}{
		"items":       items,
		"pageDetails": pageDetails,
	}
}

// CreateMockEntity creates a mock entity with the given ID
func CreateMockEntity(id int64, entityType string) map[string]interface{} {
	return map[string]interface{}{
		"id":   id,
		"name": fmt.Sprintf("Mock %s %d", entityType, id),
	}
}

// CreateMockEntities creates a list of mock entities
func CreateMockEntities(count int, entityType string, startID int64) []interface{} {
	entities := make([]interface{}, count)
	for i := 0; i < count; i++ {
		entities[i] = CreateMockEntity(startID+int64(i), entityType)
	}
	return entities
}

// AssertEqual asserts that two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if expected != actual {
		t.Errorf("%s: expected %v, got %v", message, expected, actual)
	}
}

// AssertNotNil asserts that a value is not nil
func AssertNotNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value == nil {
		t.Errorf("%s: expected non-nil value", message)
	}
}

// AssertNil asserts that a value is nil
func AssertNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value != nil {
		t.Errorf("%s: expected nil value, got %v", message, value)
	}
}

// AssertTrue asserts that a condition is true
func AssertTrue(t *testing.T, condition bool, message string) {
	t.Helper()
	if !condition {
		t.Errorf("%s: expected true", message)
	}
}

// AssertFalse asserts that a condition is false
func AssertFalse(t *testing.T, condition bool, message string) {
	t.Helper()
	if condition {
		t.Errorf("%s: expected false", message)
	}
}

// AssertContains asserts that a string contains a substring
func AssertContains(t *testing.T, haystack, needle string, message string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Errorf("%s: expected %q to contain %q", message, haystack, needle)
	}
}

// AssertNotContains asserts that a string does not contain a substring
func AssertNotContains(t *testing.T, haystack, needle string, message string) {
	t.Helper()
	if strings.Contains(haystack, needle) {
		t.Errorf("%s: expected %q to not contain %q", message, haystack, needle)
	}
}

// AssertLen asserts that a slice or map has the expected length
func AssertLen(t *testing.T, value interface{}, expected int, message string) {
	t.Helper()

	var actual int

	switch v := value.(type) {
	case []interface{}:
		actual = len(v)
	case map[string]interface{}:
		actual = len(v)
	case string:
		actual = len(v)
	case []string:
		actual = len(v)
	case []int:
		actual = len(v)
	case []int64:
		actual = len(v)
	default:
		t.Fatalf("%s: unsupported type for length assertion: %T", message, value)
	}

	if actual != expected {
		t.Errorf("%s: expected length %d, got %d", message, expected, actual)
	}
}
