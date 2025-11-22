package notion

import (
	"os"
	"testing"
)

func TestGetProjectsIntegration(t *testing.T) {

	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		t.Skip("NOTION_TOKEN not set; skipping integration test")
	}

	// Create the real Notion client
	cli, err := NewNotionClient(token)
	if err != nil {
		t.Fatalf("failed to create Notion client: %v", err)
	}

	projects, err := cli.GetProjects()
	if err != nil {
		t.Fatalf("failed to get projects: %v", err)
	}

	if len(projects) == 0 {
		t.Fatal("expected at least one project")
	}

	// Optional: print project names
	for _, p := range projects {
		t.Logf("project: %v, %v", p.Id, p.Title)
		//t.Logf("Project: %s (ID: %s)", p.Name, p.ID)
	}
}

func TestGetTasksIntegration(t *testing.T) {
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		t.Skip("NOTION_TOKEN not set; skipping integration test")
	}

	taskId := "351350bb-84d5-46e3-b740-1962f24e5b2d"
	// Create the real Notion client
	cli, err := NewNotionClient(token)
	if err != nil {
		t.Fatalf("failed to create Notion client: %v", err)
	}
	tasks, err := cli.GetConcerts(taskId)
	if err != nil {
		t.Fatalf("failed to get tasks: %v", err)
	}
	if len(tasks) == 0 {
		t.Fatal("expected at least one task")
	}
	for _, task := range tasks {
		t.Logf("task: %v, %v, date: %v", task.Title, task.StartDateAndTime, task.LocationId)
	}
}

func TestGetLocationsIntegration(t *testing.T) {
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		t.Skip("NOTION_TOKEN not set; skipping integration test")
	}

	// Create the real Notion client
	cli, err := NewNotionClient(token)
	if err != nil {
		t.Fatalf("failed to create Notion client: %v", err)
	}
	locations, err := cli.GetLocations()
	if err != nil {
		t.Fatalf("failed to get locations: %v", err)
	}
	if len(locations) == 0 {
		t.Fatal("expected at least one task")
	}
	for _, loc := range locations {
		t.Logf("location: %v, %v, %v", loc.Id, loc.Location, loc.Address)
	}
}
