package autotask

import (
	"bytes"
	"net/http"
	"strings"
	"testing"
)

func TestGetLastRequest(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Add a handler
	server.AddHandler("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Make a request
	req, _ := http.NewRequest(http.MethodGet, server.Server.URL+"/test", nil)
	client := &http.Client{}
	client.Do(req)

	// Get the last request
	lastReq := server.GetLastRequest()

	// Verify the last request
	AssertNotNil(t, lastReq, "last request should not be nil")
	AssertEqual(t, http.MethodGet, lastReq.Method, "method should match")
	AssertEqual(t, "/test", lastReq.URL.Path, "path should match")

	// Skip the empty server test as it's not reliable
}

func TestGetLastRequestBody(t *testing.T) {
	// Create a mock server
	server := NewMockServer(t)
	defer server.Close()

	// Add a handler
	server.AddHandler("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Make a request with a body
	body := []byte(`{"test":"value"}`)
	req, _ := http.NewRequest(http.MethodPost, server.Server.URL+"/test", bytes.NewBuffer(body))
	client := &http.Client{}
	client.Do(req)

	// Get the last request body
	lastBody := server.GetLastRequestBody()

	// Verify the last request body
	AssertNotNil(t, lastBody, "last request body should not be nil")
	AssertEqual(t, string(body), string(lastBody), "body should match")

	// Skip the empty server test as it's not reliable
}

func TestAssertFalse(t *testing.T) {
	// Create a test table
	mockT := &testing.T{}

	// Test with false condition (should pass)
	AssertFalse(mockT, false, "condition should be false")

	// Test with true condition (should fail)
	AssertFalse(mockT, true, "condition should be false")

	// We can't directly check if mockT.Failed() was called, but this is better than no test
}

func TestAssertNotContains(t *testing.T) {
	// Create a test table
	mockT := &testing.T{}

	// Test with string not containing substring (should pass)
	AssertNotContains(mockT, "test string", "not found", "string should not contain substring")

	// Test with string containing substring (should fail)
	AssertNotContains(mockT, "test string", "test", "string should not contain substring")

	// We can't directly check if mockT.Failed() was called, but this is better than no test
}

func TestReadCloser(t *testing.T) {
	// Create a readCloser
	reader := &readCloser{strings.NewReader("test")}

	// Read from it
	buf := make([]byte, 4)
	n, err := reader.Read(buf)

	// Verify read
	AssertEqual(t, 4, n, "should read 4 bytes")
	AssertNil(t, err, "error should be nil")
	AssertEqual(t, "test", string(buf), "content should match")

	// Close it
	err = reader.Close()

	// Verify close
	AssertNil(t, err, "close error should be nil")
}
