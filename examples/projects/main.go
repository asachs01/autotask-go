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

	// Example 1: Query active projects
	fmt.Println("Querying active projects...")
	var projectResponse struct {
		Items       []autotask.Project   `json:"items"`
		PageDetails autotask.PageDetails `json:"pageDetails"`
	}
	err := client.Projects().Query(ctx, "Status!=5", &projectResponse) // Status 5 is typically "Completed"
	if err != nil {
		log.Fatalf("Error querying projects: %v", err)
	}
	fmt.Printf("Found %d active projects\n", len(projectResponse.Items))
	for i, p := range projectResponse.Items {
		if i < 5 { // Only print the first 5
			fmt.Printf("  - %s (ID: %d)\n", p.ProjectName, p.ID)
		}
	}

	// Example 2: Get project details
	if len(projectResponse.Items) > 0 {
		projectID := projectResponse.Items[0].ID
		fmt.Printf("\nGetting details for project ID %d...\n", projectID)
		var project autotask.Project
		result, err := client.Projects().Get(ctx, projectID)
		if err != nil {
			log.Fatalf("Error getting project: %v", err)
		}
		if resultMap, ok := result.(map[string]interface{}); ok {
			if err := mapToStruct(resultMap, &project); err != nil {
				log.Fatalf("Error converting map to struct: %v", err)
			}
		} else {
			log.Fatal("Expected map result from Get")
		}
		fmt.Printf("Project: %s\n", project.ProjectName)
		fmt.Printf("  Description: %s\n", project.Description)
		fmt.Printf("  Status: %d\n", project.Status)
		fmt.Printf("  Start Date: %s\n", project.StartDate)
		fmt.Printf("  End Date: %s\n", project.EndDate)
		fmt.Printf("  Estimated Hours: %.2f\n", project.EstimatedHours)
		fmt.Printf("  Completed Percentage: %.2f%%\n", project.CompletedPercentage)

		// Example 3: Query tasks for this project
		fmt.Printf("\nQuerying tasks for project ID %d...\n", projectID)
		var taskResponse struct {
			Items       []autotask.Task      `json:"items"`
			PageDetails autotask.PageDetails `json:"pageDetails"`
		}
		err = client.Tasks().Query(ctx, fmt.Sprintf("ProjectID=%d", projectID), &taskResponse)
		if err != nil {
			log.Fatalf("Error querying tasks: %v", err)
		}
		fmt.Printf("Found %d tasks for this project\n", len(taskResponse.Items))
		for i, t := range taskResponse.Items {
			if i < 5 { // Only print the first 5
				fmt.Printf("  - %s (ID: %d, Status: %d)\n", t.Title, t.ID, t.Status)
			}
		}

		// Example 4: Query time entries for this project's tasks
		if len(taskResponse.Items) > 0 {
			taskID := taskResponse.Items[0].ID
			fmt.Printf("\nQuerying time entries for task ID %d...\n", taskID)
			var timeEntryResponse struct {
				Items       []autotask.TimeEntry `json:"items"`
				PageDetails autotask.PageDetails `json:"pageDetails"`
			}
			err = client.TimeEntries().Query(ctx, fmt.Sprintf("TaskID=%d", taskID), &timeEntryResponse)
			if err != nil {
				log.Fatalf("Error querying time entries: %v", err)
			}
			fmt.Printf("Found %d time entries for this task\n", len(timeEntryResponse.Items))
			var totalHours float64
			for i, te := range timeEntryResponse.Items {
				totalHours += te.HoursWorked
				if i < 5 { // Only print the first 5
					fmt.Printf("  - Date: %s, Hours: %.2f, Resource ID: %d\n", te.DateWorked, te.HoursWorked, te.ResourceID)
				}
			}
			fmt.Printf("Total hours worked on this task: %.2f\n", totalHours)
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
