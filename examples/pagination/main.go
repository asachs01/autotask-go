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

	// Example 1: Using the PaginationIterator to iterate through tickets
	fmt.Println("Example 1: Using the PaginationIterator to iterate through tickets")

	// Create an iterator with a small page size to demonstrate pagination
	iterator, err := autotask.NewPaginationIterator(ctx, client.Tickets(), "Status!=5", 5)
	if err != nil {
		log.Fatalf("Error creating pagination iterator: %v", err)
	}

	// Iterate through tickets
	fmt.Println("Iterating through tickets:")
	count := 0
	for iterator.Next() {
		item := iterator.Item()
		if ticketMap, ok := item.(map[string]interface{}); ok {
			fmt.Printf("  Page %d, Ticket %d: %v\n",
				iterator.CurrentPage()-1, // CurrentPage returns the next page to fetch
				count+1,
				ticketMap["title"])
		}
		count++
		if count >= 10 { // Show first 10 tickets across pages
			break
		}
	}
	fmt.Printf("Total tickets available: %d\n", iterator.TotalCount())

	// Example 2: Using FetchAllPagesWithCallback to process tickets in batches
	fmt.Println("\nExample 2: Using FetchAllPagesWithCallback to process tickets in batches")

	// Track tickets processed per page
	ticketsPerPage := make(map[int]int)
	maxPages := 3 // Limit to 3 pages

	err = autotask.FetchAllPagesWithCallback[autotask.Ticket](
		ctx,
		client.Tickets(),
		"Status!=5",
		func(items []autotask.Ticket, pageDetails autotask.PageDetails) error {
			ticketsPerPage[pageDetails.PageNumber] = len(items)
			fmt.Printf("  Processing page %d (%d items)\n", pageDetails.PageNumber, len(items))

			// Print first 3 tickets from this page
			for i, ticket := range items {
				if i < 3 {
					fmt.Printf("    Ticket %d: %s (ID: %d)\n", i+1, ticket.Title, ticket.ID)
				}
				if i == 2 && len(items) > 3 {
					fmt.Printf("    ... and %d more tickets\n", len(items)-3)
				}
			}

			// Stop after processing maxPages
			if pageDetails.PageNumber >= maxPages {
				return fmt.Errorf("reached max pages")
			}

			return nil
		},
	)
	if err != nil && err.Error() != "reached max pages" {
		log.Fatalf("Error processing tickets with callback: %v", err)
	}

	// Print summary
	var totalProcessed int
	for page, count := range ticketsPerPage {
		totalProcessed += count
		fmt.Printf("  Page %d: %d tickets\n", page, count)
	}
	fmt.Printf("Processed %d tickets from %d pages\n", totalProcessed, len(ticketsPerPage))

	// Example 3: Using FetchPage to get a specific page of results
	fmt.Println("\nExample 3: Using FetchPage to get a specific page of results")

	// Get page 2
	options := autotask.PaginationOptions{
		Page:     2,
		PageSize: 5,
	}
	page, err := autotask.FetchPage[autotask.Ticket](
		ctx,
		client.Tickets(),
		"Status!=5",
		options,
	)
	if err != nil {
		log.Fatalf("Error fetching page: %v", err)
	}

	fmt.Printf("Page %d - Found %d tickets\n", page.PageDetails.PageNumber, len(page.Items))
	for i, ticket := range page.Items {
		fmt.Printf("  Ticket %d: %s (ID: %d)\n", i+1, ticket.Title, ticket.ID)
	}

	// Show pagination information
	fmt.Println("\nPagination information:")
	fmt.Printf("  Current page: %d\n", page.PageDetails.PageNumber)
	fmt.Printf("  Page size: %d\n", page.PageDetails.PageSize)
	fmt.Printf("  Total items: %d\n", page.PageDetails.Count)
	fmt.Printf("  Has next page: %v\n", page.PageDetails.NextPageUrl != "")
	fmt.Printf("  Has previous page: %v\n", page.PageDetails.PrevPageUrl != "")
}
