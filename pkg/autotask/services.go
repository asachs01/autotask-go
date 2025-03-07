package autotask

import (
	"context"
	"net/http"
)

// companiesService implements the CompaniesService interface
type companiesService struct {
	BaseEntityService
}

// Get gets a company by ID.
func (s *companiesService) Get(ctx context.Context, id int64) (interface{}, error) {
	return s.BaseEntityService.Get(ctx, id)
}

// Query queries companies with a filter.
func (s *companiesService) Query(ctx context.Context, filter string, result interface{}) error {
	return s.BaseEntityService.Query(ctx, filter, result)
}

// Create creates a new company.
func (s *companiesService) Create(ctx context.Context, entity interface{}) (interface{}, error) {
	return s.BaseEntityService.Create(ctx, entity)
}

// Update updates an existing company.
func (s *companiesService) Update(ctx context.Context, id int64, entity interface{}) (interface{}, error) {
	return s.BaseEntityService.Update(ctx, id, entity)
}

// Delete deletes a company by ID.
func (s *companiesService) Delete(ctx context.Context, id int64) error {
	return s.BaseEntityService.Delete(ctx, id)
}

// Count counts companies matching a filter.
func (s *companiesService) Count(ctx context.Context, filter string) (int, error) {
	return s.BaseEntityService.Count(ctx, filter)
}

// GetNextPage gets the next page of results.
func (s *companiesService) GetNextPage(ctx context.Context, pageDetails PageDetails) ([]interface{}, error) {
	return s.BaseEntityService.GetNextPage(ctx, pageDetails)
}

// GetPreviousPage gets the previous page of results.
func (s *companiesService) GetPreviousPage(ctx context.Context, pageDetails PageDetails) ([]interface{}, error) {
	return s.BaseEntityService.GetPreviousPage(ctx, pageDetails)
}

// ticketsService implements the TicketsService interface
type ticketsService struct {
	BaseEntityService
}

// contactsService implements the ContactsService interface
type contactsService struct {
	BaseEntityService
}

// webhookService implements the WebhookService interface
type webhookService struct {
	BaseEntityService
}

// RegisterHandler registers a webhook handler
func (s *webhookService) RegisterHandler(eventType string, handler WebhookHandler) {
	// TODO: Implement webhook handler registration
}

// HandleWebhook handles incoming webhook requests
func (s *webhookService) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement webhook request handling
}

// CreateWebhook creates a new webhook
func (s *webhookService) CreateWebhook(ctx context.Context, url string, events []string) error {
	webhook := struct {
		URL    string   `json:"url"`
		Events []string `json:"events"`
	}{
		URL:    url,
		Events: events,
	}

	_, err := s.Create(ctx, webhook)
	return err
}

// DeleteWebhook deletes a webhook
func (s *webhookService) DeleteWebhook(ctx context.Context, id int64) error {
	return s.Delete(ctx, id)
}

// ListWebhooks lists all webhooks
func (s *webhookService) ListWebhooks(ctx context.Context) ([]interface{}, error) {
	var result ListResponse
	err := s.Query(ctx, "", &result)
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}

// resourcesService handles communication with the resources related methods of the Autotask API.
type resourcesService struct {
	BaseEntityService
}

// projectsService handles communication with the projects related methods of the Autotask API.
type projectsService struct {
	BaseEntityService
}

// Get gets a project by ID.
func (s *projectsService) Get(ctx context.Context, id int64) (interface{}, error) {
	return s.BaseEntityService.Get(ctx, id)
}

// Query queries projects with a filter.
func (s *projectsService) Query(ctx context.Context, filter string, result interface{}) error {
	return s.BaseEntityService.Query(ctx, filter, result)
}

// Create creates a new project.
func (s *projectsService) Create(ctx context.Context, entity interface{}) (interface{}, error) {
	return s.BaseEntityService.Create(ctx, entity)
}

// Update updates an existing project.
func (s *projectsService) Update(ctx context.Context, id int64, entity interface{}) (interface{}, error) {
	return s.BaseEntityService.Update(ctx, id, entity)
}

// Delete deletes a project by ID.
func (s *projectsService) Delete(ctx context.Context, id int64) error {
	return s.BaseEntityService.Delete(ctx, id)
}

// tasksService handles communication with the tasks related methods of the Autotask API.
type tasksService struct {
	BaseEntityService
}

// Get gets a task by ID.
func (s *tasksService) Get(ctx context.Context, id int64) (interface{}, error) {
	return s.BaseEntityService.Get(ctx, id)
}

// Query queries tasks with a filter.
func (s *tasksService) Query(ctx context.Context, filter string, result interface{}) error {
	return s.BaseEntityService.Query(ctx, filter, result)
}

// Create creates a new task.
func (s *tasksService) Create(ctx context.Context, entity interface{}) (interface{}, error) {
	return s.BaseEntityService.Create(ctx, entity)
}

// Update updates an existing task.
func (s *tasksService) Update(ctx context.Context, id int64, entity interface{}) (interface{}, error) {
	return s.BaseEntityService.Update(ctx, id, entity)
}

// Delete deletes a task by ID.
func (s *tasksService) Delete(ctx context.Context, id int64) error {
	return s.BaseEntityService.Delete(ctx, id)
}

// timeEntriesService handles communication with the time entries related methods of the Autotask API.
type timeEntriesService struct {
	BaseEntityService
}

// Get gets a time entry by ID.
func (s *timeEntriesService) Get(ctx context.Context, id int64) (interface{}, error) {
	return s.BaseEntityService.Get(ctx, id)
}

// Query queries time entries with a filter.
func (s *timeEntriesService) Query(ctx context.Context, filter string, result interface{}) error {
	return s.BaseEntityService.Query(ctx, filter, result)
}

// Create creates a new time entry.
func (s *timeEntriesService) Create(ctx context.Context, entity interface{}) (interface{}, error) {
	return s.BaseEntityService.Create(ctx, entity)
}

// Update updates an existing time entry.
func (s *timeEntriesService) Update(ctx context.Context, id int64, entity interface{}) (interface{}, error) {
	return s.BaseEntityService.Update(ctx, id, entity)
}

// Delete deletes a time entry by ID.
func (s *timeEntriesService) Delete(ctx context.Context, id int64) error {
	return s.BaseEntityService.Delete(ctx, id)
}

// contractsService handles communication with the contracts related methods of the Autotask API.
type contractsService struct {
	BaseEntityService
}

// Get gets a contract by ID.
func (s *contractsService) Get(ctx context.Context, id int64) (interface{}, error) {
	return s.BaseEntityService.Get(ctx, id)
}

// Query queries contracts with a filter.
func (s *contractsService) Query(ctx context.Context, filter string, result interface{}) error {
	return s.BaseEntityService.Query(ctx, filter, result)
}

// Create creates a new contract.
func (s *contractsService) Create(ctx context.Context, entity interface{}) (interface{}, error) {
	return s.BaseEntityService.Create(ctx, entity)
}

// Update updates an existing contract.
func (s *contractsService) Update(ctx context.Context, id int64, entity interface{}) (interface{}, error) {
	return s.BaseEntityService.Update(ctx, id, entity)
}

// Delete deletes a contract by ID.
func (s *contractsService) Delete(ctx context.Context, id int64) error {
	return s.BaseEntityService.Delete(ctx, id)
}

// configurationItemsService handles communication with the configuration items related methods of the Autotask API.
type configurationItemsService struct {
	BaseEntityService
}

// Get gets a configuration item by ID.
func (s *configurationItemsService) Get(ctx context.Context, id int64) (interface{}, error) {
	return s.BaseEntityService.Get(ctx, id)
}

// Query queries configuration items with a filter.
func (s *configurationItemsService) Query(ctx context.Context, filter string, result interface{}) error {
	return s.BaseEntityService.Query(ctx, filter, result)
}

// Create creates a new configuration item.
func (s *configurationItemsService) Create(ctx context.Context, entity interface{}) (interface{}, error) {
	return s.BaseEntityService.Create(ctx, entity)
}

// Update updates an existing configuration item.
func (s *configurationItemsService) Update(ctx context.Context, id int64, entity interface{}) (interface{}, error) {
	return s.BaseEntityService.Update(ctx, id, entity)
}

// Delete deletes a configuration item by ID.
func (s *configurationItemsService) Delete(ctx context.Context, id int64) error {
	return s.BaseEntityService.Delete(ctx, id)
}
