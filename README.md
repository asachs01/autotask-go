# Autotask Go Client

A Go client library for the Autotask PSA REST API.

## Requirements

- Go 1.22 or higher

## Installation

```bash
go get github.com/asachs01/autotask-go@latest
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

// Optional: Configure logging
client.SetLogLevel(autotask.LogLevelDebug)
client.SetDebugMode(true)

// Create a context
ctx := context.Background()

// Query companies
var companyResponse struct {
	Items []map[string]interface{} `json:"items"`
}
err := client.Companies().Query(ctx, "Status!=5", &companyResponse)
if err != nil {
	log.Fatal(err)
}

// Query with date filter
since := time.Now().AddDate(0, -1, 0) // 1 month ago
tickets, err := client.QueryWithDateFilter(ctx, "Tickets", "LastActivityDate", since)
if err != nil {
	log.Fatal(err)
}

// Batch get entities
ids := []int64{1, 2, 3, 4, 5}
entities, err := client.BatchGetEntities(ctx, "Companies", ids, 50)
if err != nil {
	log.Fatal(err)
}

// Query with ID range
startID := int64(1000)
endID := int64(2000)
rangeEntities, err := client.QueryIDRange(ctx, "Tickets", startID, endID)
if err != nil {
	log.Fatal(err)
}
```

## Features

- Full support for Autotask PSA REST API v1.0
- Automatic zone detection and routing
- Rate limiting (default: 60 requests per minute)
- Configurable logging with debug mode
- Context support for timeouts and cancellation
- Query helpers for common operations:
  - Date-based filtering
  - Empty filters
  - Batch entity retrieval
  - ID range queries
- Support for all major Autotask entities

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

## Rate Limiting

The client includes built-in rate limiting to prevent API throttling. By default, it's set to 60 requests per minute, which can be adjusted if needed.

## Logging

The client provides configurable logging with different log levels:
- Info (default)
- Debug
- Error

You can enable debug mode for more detailed logging and set custom log outputs.

## Examples

See the [examples](examples) directory for complete examples of using the client.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.