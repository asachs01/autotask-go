package autotask

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// LogLevel represents the logging level
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// PageDetails represents pagination information
type PageDetails struct {
	PageNumber  int    `json:"pageNumber"`
	PageSize    int    `json:"pageSize"`
	Count       int    `json:"count"`
	NextPageUrl string `json:"nextPageUrl"`
	PrevPageUrl string `json:"prevPageUrl"`
}

// EntityService represents the base interface for all entity services
type EntityService interface {
	// Get retrieves an entity by ID
	Get(ctx context.Context, id int64) (interface{}, error)

	// Query retrieves entities matching the filter
	Query(ctx context.Context, filter string, result interface{}) error

	// Create creates a new entity
	Create(ctx context.Context, entity interface{}) (interface{}, error)

	// Update updates an existing entity
	Update(ctx context.Context, id int64, entity interface{}) (interface{}, error)

	// Delete deletes an entity
	Delete(ctx context.Context, id int64) error

	// Count returns the number of entities matching the filter
	Count(ctx context.Context, filter string) (int, error)

	// Pagination handles paginated results
	Pagination(ctx context.Context, url string, result interface{}) error

	// BatchCreate creates multiple entities in a single request
	BatchCreate(ctx context.Context, entities []interface{}, result interface{}) error

	// BatchUpdate updates multiple entities in a single request
	BatchUpdate(ctx context.Context, entities []interface{}, result interface{}) error

	// BatchDelete deletes multiple entities in a single request
	BatchDelete(ctx context.Context, ids []int64) error

	// GetNextPage gets the next page of results
	GetNextPage(ctx context.Context, pageDetails PageDetails) ([]interface{}, error)

	// GetPreviousPage gets the previous page of results
	GetPreviousPage(ctx context.Context, pageDetails PageDetails) ([]interface{}, error)
}

// CompaniesService represents the companies service interface
type CompaniesService interface {
	EntityService
}

// TicketsService represents the tickets service interface
type TicketsService interface {
	EntityService
}

// ContactsService represents the contacts service interface
type ContactsService interface {
	EntityService
}

// WebhookHandler is a function type that handles webhook events
type WebhookHandler func(event *WebhookEvent) error

// WebhookEvent represents a webhook event from Autotask
type WebhookEvent struct {
	EventType string          `json:"eventType"`
	Entity    string          `json:"entity"`
	EntityID  int64           `json:"entityId"`
	Timestamp string          `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

// WebhookService represents the webhook service interface
type WebhookService interface {
	EntityService
	RegisterHandler(eventType string, handler WebhookHandler)
	HandleWebhook(w http.ResponseWriter, r *http.Request)
	CreateWebhook(ctx context.Context, url string, events []string) error
	DeleteWebhook(ctx context.Context, id int64) error
	ListWebhooks(ctx context.Context) ([]interface{}, error)
}

// ResourcesService represents the resources service interface
type ResourcesService interface {
	EntityService
}

// Client represents the main interface for the Autotask API client
type Client interface {
	// Companies returns the companies service
	Companies() CompaniesService

	// Tickets returns the tickets service
	Tickets() TicketsService

	// Contacts returns the contacts service
	Contacts() ContactsService

	// Resources returns the resources service
	Resources() ResourcesService

	// Webhooks returns the webhooks service
	Webhooks() WebhookService

	// SetLogLevel sets the logging level
	SetLogLevel(level LogLevel)

	// SetDebugMode enables or disables debug logging
	SetDebugMode(debug bool)

	// SetLogOutput sets the output writer for the logger
	SetLogOutput(output *os.File)

	// GetZoneInfo gets the zone information for the Autotask account
	GetZoneInfo() (*ZoneInfo, error)

	// NewRequest creates a new HTTP request
	NewRequest(ctx context.Context, method, url string, body interface{}) (*http.Request, error)

	// Do sends an HTTP request and returns the response
	Do(req *http.Request, v interface{}) (*http.Response, error)
}

// ZoneInfo represents the zone information for an Autotask account
type ZoneInfo struct {
	ZoneName string `json:"zoneName"`
	URL      string `json:"url"`
	WebURL   string `json:"webUrl"`
	CI       int    `json:"ci"`
}

// ErrorResponse represents an error response from the Autotask API
type ErrorResponse struct {
	Response *http.Response
	Message  string   `json:"Message"`
	Errors   []string `json:"errors"`
}

// Error implements the error interface
func (r *ErrorResponse) Error() string {
	if len(r.Errors) > 0 {
		return fmt.Sprintf("%v %v: %d %v",
			r.Response.Request.Method, r.Response.Request.URL,
			r.Response.StatusCode, r.Errors[0])
	}
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
}
