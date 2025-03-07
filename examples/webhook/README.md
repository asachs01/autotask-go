# Autotask Webhook Example

This example demonstrates how to use the Autotask webhook functionality to receive and process events from Autotask.

## Features

- Webhook registration with Autotask
- Webhook signature verification for security
- Event handling for different event types
- Graceful server shutdown

## Prerequisites

- Go 1.16 or higher
- Autotask API credentials
- A publicly accessible URL for receiving webhooks (you can use tools like ngrok for development)

## Environment Variables

Set the following environment variables before running the example:

```
AUTOTASK_USERNAME=your_username
AUTOTASK_SECRET=your_secret
AUTOTASK_INTEGRATION_CODE=your_integration_code
WEBHOOK_SECRET=your_webhook_secret
WEBHOOK_URL=https://your-public-url.com/webhook
REGISTER_WEBHOOK=true  # Set to "true" to register the webhook with Autotask
```

## Running the Example

1. Set the environment variables
2. Run the example:

```bash
go run main.go
```

The server will start on port 8080 and listen for webhook events at the `/webhook` endpoint.

## Webhook Registration

When you set `REGISTER_WEBHOOK=true`, the example will register your webhook URL with Autotask. This only needs to be done once, so you can set it to `false` after the initial registration.

## Webhook Verification

The example uses HMAC-SHA256 signature verification to ensure that webhook requests are coming from Autotask. Set the `WEBHOOK_SECRET` environment variable to a secure random string and configure the same secret in your Autotask webhook settings.

## Event Handling

The example registers handlers for the following event types:

- `ticket.created`
- `ticket.updated`
- `ticket.deleted`

You can add more handlers for other event types as needed.

## Customizing Event Handlers

To customize the event handlers, modify the `handleTicketCreated`, `handleTicketUpdated`, and `handleTicketDeleted` functions in `main.go`.

## Production Considerations

For production use, consider the following:

- Use HTTPS for your webhook endpoint
- Store your webhook secret securely
- Implement retry logic for failed webhook processing
- Add monitoring and logging
- Consider using a more robust HTTP server configuration 