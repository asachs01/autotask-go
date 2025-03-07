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

	// Example 1: Simple filter
	fmt.Println("Example 1: Simple filter - Active tickets with high priority")
	var ticketResponse1 struct {
		Items       []autotask.Ticket    `json:"items"`
		PageDetails autotask.PageDetails `json:"pageDetails"`
	}
	err := client.Tickets().Query(ctx, "Status!=5 AND Priority=1", &ticketResponse1)
	if err != nil {
		log.Fatalf("Error querying tickets: %v", err)
	}
	fmt.Printf("Found %d active high-priority tickets\n", len(ticketResponse1.Items))
	for i, t := range ticketResponse1.Items {
		if i < 3 { // Only print the first 3
			fmt.Printf("  - %s (ID: %d, Priority: %d)\n", t.Title, t.ID, t.Priority)
		}
	}

	// Example 2: Complex filter with OR condition
	fmt.Println("\nExample 2: Complex filter with OR condition - Tickets assigned to specific resources")
	var ticketResponse2 struct {
		Items       []autotask.Ticket    `json:"items"`
		PageDetails autotask.PageDetails `json:"pageDetails"`
	}
	// Replace with actual resource IDs from your environment
	err = client.Tickets().Query(ctx, "AssignedResourceID=123 OR AssignedResourceID=456", &ticketResponse2)
	if err != nil {
		log.Fatalf("Error querying tickets: %v", err)
	}
	fmt.Printf("Found %d tickets assigned to specific resources\n", len(ticketResponse2.Items))
	for i, t := range ticketResponse2.Items {
		if i < 3 { // Only print the first 3
			fmt.Printf("  - %s (ID: %d, Assigned To: %d)\n", t.Title, t.ID, t.AssignedResourceID)
		}
	}

	// Example 3: Complex filter with nested conditions
	fmt.Println("\nExample 3: Complex filter with nested conditions - Active tickets with specific criteria")
	var ticketResponse3 struct {
		Items       []autotask.Ticket    `json:"items"`
		PageDetails autotask.PageDetails `json:"pageDetails"`
	}
	// Find active tickets that are either high priority or assigned to a specific resource
	err = client.Tickets().Query(ctx, "Status!=5 AND (Priority=1 OR AssignedResourceID=123)", &ticketResponse3)
	if err != nil {
		log.Fatalf("Error querying tickets: %v", err)
	}
	fmt.Printf("Found %d tickets matching complex criteria\n", len(ticketResponse3.Items))
	for i, t := range ticketResponse3.Items {
		if i < 3 { // Only print the first 3
			fmt.Printf("  - %s (ID: %d, Priority: %d, Assigned To: %d)\n",
				t.Title, t.ID, t.Priority, t.AssignedResourceID)
		}
	}

	// Example 4: Using the 'contains' operator
	fmt.Println("\nExample 4: Using the 'contains' operator - Tickets with specific text in title")
	var ticketResponse4 struct {
		Items       []autotask.Ticket    `json:"items"`
		PageDetails autotask.PageDetails `json:"pageDetails"`
	}
	err = client.Tickets().Query(ctx, "Title contains 'error'", &ticketResponse4)
	if err != nil {
		log.Fatalf("Error querying tickets: %v", err)
	}
	fmt.Printf("Found %d tickets with 'error' in the title\n", len(ticketResponse4.Items))
	for i, t := range ticketResponse4.Items {
		if i < 3 { // Only print the first 3
			fmt.Printf("  - %s (ID: %d)\n", t.Title, t.ID)
		}
	}

	// Example 5: Date range filtering
	fmt.Println("\nExample 5: Date range filtering - Tickets created in the last 7 days")
	var ticketResponse5 struct {
		Items       []autotask.Ticket    `json:"items"`
		PageDetails autotask.PageDetails `json:"pageDetails"`
	}
	sevenDaysAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	err = client.Tickets().Query(ctx, "CreateDate>"+sevenDaysAgo, &ticketResponse5)
	if err != nil {
		log.Fatalf("Error querying tickets: %v", err)
	}
	fmt.Printf("Found %d tickets created in the last 7 days\n", len(ticketResponse5.Items))
	for i, t := range ticketResponse5.Items {
		if i < 3 { // Only print the first 3
			fmt.Printf("  - %s (ID: %d, Created: %s)\n", t.Title, t.ID, t.CreateDate)
		}
	}

	// Example 6: Projects with tasks that have time entries
	fmt.Println("\nExample 6: Projects with active tasks")
	var projectResponse struct {
		Items       []autotask.Project   `json:"items"`
		PageDetails autotask.PageDetails `json:"pageDetails"`
	}
	err = client.Projects().Query(ctx, "Status=1", &projectResponse)
	if err != nil {
		log.Fatalf("Error querying projects: %v", err)
	}

	if len(projectResponse.Items) > 0 {
		projectID := projectResponse.Items[0].ID
		fmt.Printf("Looking at tasks for project ID %d (%s)\n",
			projectID, projectResponse.Items[0].ProjectName)

		var taskResponse struct {
			Items       []autotask.Task      `json:"items"`
			PageDetails autotask.PageDetails `json:"pageDetails"`
		}
		// Find active tasks for this project
		err = client.Tasks().Query(ctx,
			fmt.Sprintf("ProjectID=%d AND Status!=5", projectID),
			&taskResponse)
		if err != nil {
			log.Fatalf("Error querying tasks: %v", err)
		}

		fmt.Printf("Found %d active tasks for this project\n", len(taskResponse.Items))
		for i, t := range taskResponse.Items {
			if i < 3 { // Only print the first 3
				fmt.Printf("  - %s (ID: %d, Status: %d)\n", t.Title, t.ID, t.Status)
			}
		}
	}
}
