# Autotask Go Client

A Go client library for the Autotask PSA REST API.

## Installation

```bash
go get github.com/asachs01/autotask-go
```

## Usage

```go
import "github.com/asachs01/autotask-go/pkg/autotask"

// Create a new client
client := autotask.NewClient(
	os.Getenv("AUTOTASK_USERNAME"),
	os.Getenv("AUTOTASK_SECRET"),
	os.Getenv("AUTOTASK_INTEGRATION_CODE"),
)

// Create a context
ctx := context.Background()

// Get a company by ID
company, err := client.Companies().Get(ctx, 0)
if err != nil {
	log.Fatal(err)
}

// Query active companies
var companyResponse struct {
	Items       []autotask.Company   `json:"items"`
	PageDetails autotask.PageDetails `json:"pageDetails"`
}
err = client.Companies().Query(ctx, "", &companyResponse)
if err != nil {
	log.Fatal(err)
}

// Query tickets with a filter
var ticketResponse struct {
	Items       []autotask.Ticket    `json:"items"`
	PageDetails autotask.PageDetails `json:"pageDetails"`
}
err = client.Tickets().Query(ctx, "Status!=5", &ticketResponse)
if err != nil {
	log.Fatal(err)
}

// Query with complex filters
var highPriorityTickets struct {
	Items       []autotask.Ticket    `json:"items"`
	PageDetails autotask.PageDetails `json:"pageDetails"`
}
// Find active tickets that are either high priority or assigned to a specific resource
err = client.Tickets().Query(ctx, "Status!=5 AND (Priority=1 OR AssignedResourceID=123)", &highPriorityTickets)
if err != nil {
	log.Fatal(err)
}

// Using pagination helpers
// Option 1: Iterator pattern
iterator, err := autotask.NewPaginationIterator(ctx, client.Tickets(), "Status!=5", 10)
if err != nil {
	log.Fatal(err)
}
for iterator.Next() {
	item := iterator.Item()
	// Process each item
}

// Option 2: Fetch a specific page
options := autotask.PaginationOptions{
	Page:     2,
	PageSize: 10,
}
page, err := autotask.FetchPage[autotask.Ticket](
	ctx,
	client.Tickets(),
	"Status!=5",
	options,
)
if err != nil {
	log.Fatal(err)
}
// Process page.Items

// Option 3: Process pages with a callback
err = autotask.FetchAllPagesWithCallback[autotask.Ticket](
	ctx,
	client.Tickets(),
	"Status!=5",
	func(items []autotask.Ticket, pageDetails autotask.PageDetails) error {
		// Process each page of items
		return nil
	},
)
```

## Features

- Full support for Autotask PSA REST API v1.0
- Automatic zone detection and routing
- Rate limiting support
- Pagination support
- Configurable logging
- Context support for timeouts and cancellation
- Advanced query filtering with logical operators (AND, OR) and nested conditions
- Enhanced pagination helpers with iterator patterns and convenience methods
- Proper support for Autotask's pagination mechanism using nextPageUrl/prevPageUrl

## Supported Entities

- Companies
- Tickets
- Contacts
- Resources
- Webhooks
- Projects
- Tasks
- Time Entries
- Contracts
- Configuration Items

## Authentication

The client requires three pieces of information for authentication:

1. Username
2. Secret (API Key)
3. Integration Code

These can be obtained from your Autotask administrator or the Autotask API admin interface.

## Examples

See the [examples](examples) directory for complete examples of using the client.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.