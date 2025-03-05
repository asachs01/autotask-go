package main

import (
	"context"
	"encoding/json"
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

	// Example 1: Get the root company (ID 0)
	fmt.Println("Getting root company...")
	var company autotask.Company
	result, err := client.Companies().Get(ctx, 0)
	if err != nil {
		log.Fatalf("Error getting company: %v", err)
	}
	if resultMap, ok := result.(map[string]interface{}); ok {
		if err := mapToStruct(resultMap, &company); err != nil {
			log.Fatalf("Error converting map to struct: %v", err)
		}
	} else {
		log.Fatal("Expected map result from Get")
	}
	fmt.Printf("Root Company: %s (ID: %d)\n", company.CompanyName, company.ID)

	// Example 2: Query active companies
	fmt.Println("\nQuerying active companies...")
	var companyResponse struct {
		Items       []autotask.Company   `json:"items"`
		PageDetails autotask.PageDetails `json:"pageDetails"`
	}
	err = client.Companies().Query(ctx, "", &companyResponse)
	if err != nil {
		log.Fatalf("Error querying companies: %v", err)
	}
	fmt.Printf("Found %d active companies\n", len(companyResponse.Items))
	for i, c := range companyResponse.Items {
		if i < 5 { // Only print the first 5
			fmt.Printf("  - %s (ID: %d)\n", c.CompanyName, c.ID)
		}
	}

	// Example 3: Count active tickets
	fmt.Println("\nCounting active tickets...")
	count, err := client.Tickets().Count(ctx, "")
	if err != nil {
		log.Fatalf("Error counting tickets: %v", err)
	}
	fmt.Printf("Found %d active tickets\n", count)

	// Query tickets with pagination
	fmt.Println("\nQuerying tickets with pagination...")
	var ticketResponse struct {
		Items       []autotask.Ticket    `json:"items"`
		PageDetails autotask.PageDetails `json:"pageDetails"`
	}
	err = client.Tickets().Query(ctx, "", &ticketResponse)
	if err != nil {
		log.Fatalf("Error querying tickets: %v", err)
	}
	fmt.Printf("Page %d of tickets (showing %d of %d total)\n",
		ticketResponse.PageDetails.PageNumber,
		len(ticketResponse.Items),
		ticketResponse.PageDetails.Count)

	// Example 5: Get the next page of tickets if available
	if ticketResponse.PageDetails.NextPageUrl != "" {
		fmt.Println("\nGetting next page of tickets...")
		nextPageTickets, err := client.Tickets().GetNextPage(ctx, ticketResponse.PageDetails)
		if err != nil {
			log.Fatalf("Error getting next page: %v", err)
		}
		fmt.Printf("Next page contains %d tickets\n", len(nextPageTickets))
	}

	// Example 6: Query resources
	fmt.Println("\nQuerying resources...")
	var resourceResponse struct {
		Items       []autotask.Resource  `json:"items"`
		PageDetails autotask.PageDetails `json:"pageDetails"`
	}
	err = client.Resources().Query(ctx, "", &resourceResponse)
	if err != nil {
		log.Fatalf("Error querying resources: %v", err)
	}
	fmt.Printf("Found %d active resources\n", len(resourceResponse.Items))
	for i, r := range resourceResponse.Items {
		if i < 5 { // Only print the first 5
			fmt.Printf("  - %s %s (ID: %d)\n", r.FirstName, r.LastName, r.ID)
		}
	}
}

// Helper function to convert map to struct
func mapToStruct(m map[string]interface{}, v interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
