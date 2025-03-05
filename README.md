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
```

## Features

- Full support for Autotask PSA REST API v1.0
- Automatic zone detection and routing
- Rate limiting support
- Pagination support
- Configurable logging
- Context support for timeouts and cancellation

## Supported Entities

- Companies
- Tickets
- Contacts
- Resources
- Webhooks

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