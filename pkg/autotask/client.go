package autotask

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	// BaseZoneInfoURL is the URL to get zone information
	BaseZoneInfoURL = "https://webservices.autotask.net/atservicesrest/v1.0/ZoneInformation"

	// DefaultUserAgent is the default user agent for the client
	DefaultUserAgent = "Autotask Go Client"

	// APIVersion is the version of the Autotask API
	APIVersion = "v1.0"
)

// client represents an Autotask API client
type client struct {
	// HTTP client used to communicate with the API
	httpClient *http.Client

	// Base URL for API requests
	baseURL *url.URL

	// User agent used when communicating with the Autotask API
	UserAgent string

	// API credentials
	username        string
	secret          string
	integrationCode string

	// Zone information
	zoneInfo  *ZoneInfo
	zoneMutex sync.Mutex

	// Rate limiter
	rateLimiter *RateLimiter

	// Logger
	logger *Logger

	// Entity clients
	companiesService          *companiesService
	ticketsService            *ticketsService
	contactsService           *contactsService
	webhooksService           *webhookService
	resourcesService          *resourcesService
	projectsService           *projectsService
	tasksService              *tasksService
	timeEntriesService        *timeEntriesService
	contractsService          *contractsService
	configurationItemsService *configurationItemsService
}

// NewClient returns a new Autotask API client
func NewClient(username, secret, integrationCode string) Client {
	httpClient := &http.Client{
		Timeout: time.Second * 60,
	}

	c := &client{
		httpClient:      httpClient,
		UserAgent:       DefaultUserAgent,
		username:        username,
		secret:          secret,
		integrationCode: integrationCode,
		rateLimiter:     NewRateLimiter(60),       // Default to 60 requests per minute
		logger:          New(LogLevelInfo, false), // Default to info level, debug off
	}

	// Initialize services
	c.companiesService = &companiesService{
		BaseEntityService: NewBaseEntityService(c, "Companies"),
	}
	c.ticketsService = &ticketsService{
		BaseEntityService: NewBaseEntityService(c, "Tickets"),
	}
	c.contactsService = &contactsService{
		BaseEntityService: NewBaseEntityService(c, "Contacts"),
	}
	c.webhooksService = &webhookService{
		BaseEntityService: NewBaseEntityService(c, "Webhooks"),
	}
	c.resourcesService = &resourcesService{
		BaseEntityService: NewBaseEntityService(c, "Resources"),
	}
	c.projectsService = &projectsService{
		BaseEntityService: NewBaseEntityService(c, "Projects"),
	}
	c.tasksService = &tasksService{
		BaseEntityService: NewBaseEntityService(c, "Tasks"),
	}
	c.timeEntriesService = &timeEntriesService{
		BaseEntityService: NewBaseEntityService(c, "TimeEntries"),
	}
	c.contractsService = &contractsService{
		BaseEntityService: NewBaseEntityService(c, "Contracts"),
	}
	c.configurationItemsService = &configurationItemsService{
		BaseEntityService: NewBaseEntityService(c, "ConfigurationItems"),
	}

	return c
}

// SetLogLevel sets the logging level
func (c *client) SetLogLevel(level LogLevel) {
	c.logger.SetLevel(level)
}

// SetDebugMode enables or disables debug logging
func (c *client) SetDebugMode(debug bool) {
	c.logger.SetDebugMode(debug)
}

// SetLogOutput sets the output writer for the logger
func (c *client) SetLogOutput(output *os.File) {
	c.logger.SetOutput(output)
}

// GetZoneInfo gets the zone information for the Autotask account
func (c *client) GetZoneInfo() (*ZoneInfo, error) {
	c.zoneMutex.Lock()
	defer c.zoneMutex.Unlock()

	if c.zoneInfo != nil {
		return c.zoneInfo, nil
	}

	// Build URL with user parameter
	zoneURL := fmt.Sprintf("%s?user=%s", BaseZoneInfoURL, url.QueryEscape(c.username))
	c.logger.Debug("Requesting zone info", map[string]interface{}{
		"url": zoneURL,
	})

	req, err := http.NewRequest("GET", zoneURL, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set both Basic auth and API headers
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.secret)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	req.Header.Set("UserName", c.username)
	req.Header.Set("Secret", c.secret)

	// Set API integration code
	req.Header.Set("ApiIntegrationCode", c.integrationCode)

	// Log request headers
	headers := make(map[string]string)
	for key, values := range req.Header {
		headers[key] = values[0]
	}
	c.logger.LogRequest(req.Method, req.URL.String(), headers)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Log response headers
	respHeaders := make(map[string]string)
	for key, values := range resp.Header {
		respHeaders[key] = values[0]
	}
	c.logger.LogResponse(resp.StatusCode, respHeaders)

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	// Read and print response body for debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c.logger.Debug("Zone info response", map[string]interface{}{
		"body": string(body),
	})

	// Create a new reader with the same body for json.Decoder
	bodyReader := bytes.NewReader(body)

	var zoneInfo ZoneInfo
	if err := json.NewDecoder(bodyReader).Decode(&zoneInfo); err != nil {
		return nil, err
	}

	c.logger.Debug("Parsed zone info", map[string]interface{}{
		"zone_info": zoneInfo,
	})

	c.zoneInfo = &zoneInfo
	// Add API version to base URL, ensuring lowercase
	baseURL := strings.Replace(zoneInfo.URL, "ATServicesRest", "atservicesrest", 1)
	baseURL = fmt.Sprintf("%sv1.0/", baseURL)
	c.logger.Debug("Using base URL", map[string]interface{}{
		"base_url": baseURL,
	})
	c.baseURL, err = url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return c.zoneInfo, nil
}

// NewRequest creates an API request with context
func (c *client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	// Get zone info if not already set
	if c.baseURL == nil {
		if _, err := c.GetZoneInfo(); err != nil {
			return nil, err
		}
	}

	// Resolve relative URL
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set both Basic auth and API headers
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.secret)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	req.Header.Set("UserName", c.username)
	req.Header.Set("Secret", c.secret)

	// Set API integration code
	req.Header.Set("ApiIntegrationCode", c.integrationCode)

	// Log request headers
	headers := make(map[string]string)
	for key, values := range req.Header {
		headers[key] = values[0]
	}
	c.logger.LogRequest(req.Method, req.URL.String(), headers)

	return req, nil
}

// Do sends an API request and returns the API response
func (c *client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	// Apply rate limiting
	c.rateLimiter.Wait()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Log response headers
	respHeaders := make(map[string]string)
	for key, values := range resp.Header {
		respHeaders[key] = values[0]
	}
	c.logger.LogResponse(resp.StatusCode, respHeaders)

	// Handle successful responses
	if resp.StatusCode == http.StatusNoContent {
		// For 204 No Content responses, there's no body to parse
		return resp, nil
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, c.handleErrorResponse(resp)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log the response body for debugging
	c.logger.Debug("Response body", map[string]interface{}{
		"body": string(body),
	})

	// If v is nil, we don't need to parse the response
	if v == nil {
		return resp, nil
	}

	// If v is a pointer to []byte, store the raw body
	if b, ok := v.(*[]byte); ok {
		*b = body
		return resp, nil
	}

	// Otherwise, unmarshal the JSON
	if err := json.Unmarshal(body, v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return resp, nil
}

// handleErrorResponse handles error responses from the API
func (c *client) handleErrorResponse(resp *http.Response) error {
	var errorResp ErrorResponse
	errorResp.Response = resp
	data, err := io.ReadAll(resp.Body)
	if err == nil && data != nil {
		if unmarshalErr := json.Unmarshal(data, &errorResp); unmarshalErr != nil {
			// If we can't unmarshal the error response, just log it and continue
			// We'll still return the error response with the status code
			c.logger.LogError(unmarshalErr)
		}
	}
	return &errorResp
}

// Companies returns the companies service
func (c *client) Companies() CompaniesService {
	return c.companiesService
}

// Tickets returns the tickets service
func (c *client) Tickets() TicketsService {
	return c.ticketsService
}

// Contacts returns the contacts service
func (c *client) Contacts() ContactsService {
	return c.contactsService
}

// Webhooks returns the webhooks service
func (c *client) Webhooks() WebhookService {
	return c.webhooksService
}

// Resources returns the resources service
func (c *client) Resources() ResourcesService {
	return c.resourcesService
}

// Projects returns the projects service
func (c *client) Projects() ProjectsService {
	return c.projectsService
}

// Tasks returns the tasks service
func (c *client) Tasks() TasksService {
	return c.tasksService
}

// TimeEntries returns the time entries service
func (c *client) TimeEntries() TimeEntriesService {
	return c.timeEntriesService
}

// Contracts returns the contracts service
func (c *client) Contracts() ContractsService {
	return c.contractsService
}

// ConfigurationItems returns the configuration items service
func (c *client) ConfigurationItems() ConfigurationItemsService {
	return c.configurationItemsService
}
