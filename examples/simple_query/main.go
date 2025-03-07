package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

	// Create a context with timeout for all operations
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Example: Basic query
	fmt.Println("Example: Basic query for active tickets")

	var ticketResponse struct {
		Items       []autotask.Ticket    `json:"items"`
		PageDetails autotask.PageDetails `json:"pageDetails"`
	}

	err := client.Tickets().Query(ctx, "Status!=5", &ticketResponse)
	if err != nil {
		log.Fatalf("Error querying tickets: %v", err)
	}

	fmt.Printf("Found %d tickets (total: %d)\n", len(ticketResponse.Items), ticketResponse.PageDetails.Count)
	for i, ticket := range ticketResponse.Items {
		if i < 5 { // Only show the first 5 for brevity
			fmt.Printf("  Ticket %d: %s (ID: %d)\n", i+1, ticket.Title, ticket.ID)
		}
	}
}
