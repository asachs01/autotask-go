package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asachs01/autotask-go/pkg/autotask"
)

func main() {
	// Create a new Autotask client
	client := autotask.NewClient(
		os.Getenv("AUTOTASK_USERNAME"),
		os.Getenv("AUTOTASK_SECRET"),
		os.Getenv("AUTOTASK_INTEGRATION_CODE"),
	)

	// Enable debug mode
	client.SetDebugMode(true)

	// Set webhook secret for verification
	client.Webhooks().SetWebhookSecret(os.Getenv("WEBHOOK_SECRET"))

	// Register webhook handlers for different event types
	client.Webhooks().RegisterHandler("ticket.created", handleTicketCreated)
	client.Webhooks().RegisterHandler("ticket.updated", handleTicketUpdated)
	client.Webhooks().RegisterHandler("ticket.deleted", handleTicketDeleted)

	// Create a webhook (if needed)
	// This registers your webhook URL with Autotask
	// Note: You only need to do this once, not on every startup
	if os.Getenv("REGISTER_WEBHOOK") == "true" {
		ctx := context.Background()
		webhookURL := os.Getenv("WEBHOOK_URL")
		events := []string{"ticket.created", "ticket.updated", "ticket.deleted"}

		fmt.Println("Registering webhook at:", webhookURL)
		err := client.Webhooks().CreateWebhook(ctx, webhookURL, events)
		if err != nil {
			log.Fatalf("Failed to create webhook: %v", err)
		}
		fmt.Println("Webhook registered successfully")
	}

	// Set up HTTP server to handle webhook callbacks
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		client.Webhooks().HandleWebhook(w, r)
	})

	// Start the HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	// Start the server in a goroutine
	go func() {
		fmt.Println("Starting webhook server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	<-stop
	fmt.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}

// Handler for ticket.created events
func handleTicketCreated(event *autotask.WebhookEvent) error {
	fmt.Printf("Ticket created: ID=%d\n", event.EntityID)

	// You can unmarshal the event.Data into a specific struct if needed
	var ticketData map[string]interface{}
	if err := json.Unmarshal(event.Data, &ticketData); err != nil {
		return fmt.Errorf("failed to parse ticket data: %w", err)
	}

	fmt.Printf("Ticket details: %+v\n", ticketData)
	return nil
}

// Handler for ticket.updated events
func handleTicketUpdated(event *autotask.WebhookEvent) error {
	fmt.Printf("Ticket updated: ID=%d\n", event.EntityID)
	return nil
}

// Handler for ticket.deleted events
func handleTicketDeleted(event *autotask.WebhookEvent) error {
	fmt.Printf("Ticket deleted: ID=%d\n", event.EntityID)
	return nil
}
