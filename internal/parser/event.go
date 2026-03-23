// This if for parse event.json file in runner directory
// Used for parse repo name, runner group, and workflow name
package parser

import (
	"encoding/json"
	"os"
	"strings"
)

// EventInfo holds parsed GitHub event data
type EventInfo struct {
	WorkflowName string `json:"workflow"`

	Repository struct {
		RepoName     string `json:"name"`
		RepoFullName string `json:"full_name"`
		PushedAt     string `json:"pushed_at"` // RFC3339 format
	} `json:"repository"`

	Organization *struct {
		OrgName string `json:"login"`
	} `json:"organization,omitempty"`

	Enterprise *struct {
		Slug string `json:"slug"`
	} `json:"enterprise,omitempty"`
}

func ReadEventJSON(path string) (*EventInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var event EventInfo
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

func EventRepository(event *EventInfo) string {
	if event == nil {
		return "unknown"
	}
	if strings.TrimSpace(event.Repository.RepoFullName) != "" {
		return strings.TrimSpace(event.Repository.RepoFullName)
	}
	if strings.TrimSpace(event.Repository.RepoName) != "" {
		return strings.TrimSpace(event.Repository.RepoName)
	}
	return "unknown"
}

func EventRepositoryOwner(event *EventInfo) string {
	if event == nil {
		return "unknown"
	}
	if event.Organization != nil && strings.TrimSpace(event.Organization.OrgName) != "" {
		return strings.TrimSpace(event.Organization.OrgName)
	}
	fullName := EventRepository(event)
	parts := strings.SplitN(fullName, "/", 2)
	if len(parts) == 2 && strings.TrimSpace(parts[0]) != "" {
		return strings.TrimSpace(parts[0])
	}
	return "unknown"
}
