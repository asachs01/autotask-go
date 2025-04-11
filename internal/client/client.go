package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/asachs01/autotask-go/pkg/autotask"
)

// Client represents an Autotask API client
type Client struct {
	// HTTP client used to communicate with the API
	client *http.Client

	// Base URL for API requests
	baseURL *url.URL

	// User agent used when communicating with the Autotask API
	UserAgent string

	// API credentials
	username        string
	secret          string
	integrationCode string

	// Rate limiter
	rateLimiter *autotask.RateLimiter
}

// NewClient returns a new Autotask API client
func NewClient(username, secret, integrationCode string) *Client {
	httpClient := &http.Client{
		Timeout: time.Second * 30,
	}

	return &Client{
		client:          httpClient,
		UserAgent:       "Autotask Go Client",
		username:        username,
		secret:          secret,
		integrationCode: integrationCode,
		rateLimiter:     autotask.NewRateLimiter(60), // Default to 60 requests per minute
	}
}

// newRequest creates a new HTTP request
func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	return req, nil
}

// do performs the HTTP request
func (c *Client) do(req *http.Request, v interface{}) error {
	c.rateLimiter.Wait()

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			// Log the error but don't return it since we're in a defer
			fmt.Fprintf(os.Stderr, "Error closing response body: %v\n", cerr)
		}
	}()

	if resp.StatusCode >= 400 {
		var errResp autotask.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("error decoding error response: %v", err)
		}
		errResp.Response = resp
		return &errResp
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}

	return nil
}
