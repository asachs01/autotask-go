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

	// Example: Basic pagination with the FetchPage helper
	fmt.Println("Example: Basic pagination with the FetchPage helper")

	// First page - using FetchPage helper
	page1, err := autotask.FetchPage[autotask.Ticket](
		ctx,
		client.Tickets(),
		"Status!=5",
		autotask.PaginationOptions{
			Page:     1,
			PageSize: 5,
		},
	)
	if err != nil {
		log.Fatalf("Error fetching page 1: %v", err)
	}

	fmt.Printf("Page 1 - Found %d tickets (total: %d)\n", len(page1.Items), page1.PageDetails.Count)
	for i, ticket := range page1.Items {
		fmt.Printf("  Ticket %d: %s (ID: %d)\n", i+1, ticket.Title, ticket.ID)
	}

	// Check if there's a next page
	if page1.PageDetails.NextPageUrl != "" {
		fmt.Println("\nFetching page 2...")

		// Second page - using FetchPage helper
		page2, err := autotask.FetchPage[autotask.Ticket](
			ctx,
			client.Tickets(),
			"Status!=5",
			autotask.PaginationOptions{
				Page:     2,
				PageSize: 5,
			},
		)
		if err != nil {
			log.Fatalf("Error fetching page 2: %v", err)
		}

		fmt.Printf("Page 2 - Found %d tickets\n", len(page2.Items))
		for i, ticket := range page2.Items {
			fmt.Printf("  Ticket %d: %s (ID: %d)\n", i+1, ticket.Title, ticket.ID)
		}
	} else {
		fmt.Println("\nNo more pages available.")
	}
}
