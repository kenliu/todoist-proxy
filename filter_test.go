package main

import (
	"encoding/json"
	"testing"
)

func TestFilterSyncResponse_BasicProjectFilter(t *testing.T) {
	input := `{
		"sync_token": "abc123",
		"full_sync": true,
		"projects": [
			{"id": "1", "parent_id": null, "name": "Work"},
			{"id": "2", "parent_id": null, "name": "Personal"}
		],
		"items": [
			{"id": "10", "project_id": "1", "content": "Work task"},
			{"id": "11", "project_id": "2", "content": "Personal task"}
		]
	}`

	got, err := FilterSyncResponse([]byte(input), []string{"1"})
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(got, &result); err != nil {
		t.Fatal(err)
	}

	var projects []projectStub
	json.Unmarshal(result["projects"], &projects)
	if len(projects) != 1 || projects[0].ID != "1" {
		t.Errorf("expected 1 project (id=1), got %+v", projects)
	}

	var items []resourceStub
	json.Unmarshal(result["items"], &items)
	if len(items) != 1 || items[0].ProjectID != "1" {
		t.Errorf("expected 1 item (project_id=1), got %+v", items)
	}

	// sync_token should be preserved
	var token string
	json.Unmarshal(result["sync_token"], &token)
	if token != "abc123" {
		t.Errorf("expected sync_token=abc123, got %q", token)
	}
}

func TestFilterSyncResponse_HierarchyExpansion(t *testing.T) {
	// Parent "1" is allowed; children "1a" and "1b" (and grandchild "1a1") should be included.
	input := `{
		"sync_token": "tok",
		"full_sync": true,
		"projects": [
			{"id": "1",  "parent_id": null, "name": "Root"},
			{"id": "1a", "parent_id": "1",  "name": "Child A"},
			{"id": "1b", "parent_id": "1",  "name": "Child B"},
			{"id": "1a1","parent_id": "1a", "name": "Grandchild"},
			{"id": "2",  "parent_id": null, "name": "Other"}
		],
		"items": [
			{"id": "10", "project_id": "1",   "content": "Root task"},
			{"id": "11", "project_id": "1a",  "content": "Child A task"},
			{"id": "12", "project_id": "1a1", "content": "Grandchild task"},
			{"id": "13", "project_id": "2",   "content": "Other task"}
		]
	}`

	got, err := FilterSyncResponse([]byte(input), []string{"1"})
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]json.RawMessage
	json.Unmarshal(got, &result)

	var projects []projectStub
	json.Unmarshal(result["projects"], &projects)
	if len(projects) != 4 {
		t.Errorf("expected 4 projects (1, 1a, 1b, 1a1), got %d: %+v", len(projects), projects)
	}

	var items []resourceStub
	json.Unmarshal(result["items"], &items)
	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d: %+v", len(items), items)
	}
}

func TestFilterSyncResponse_PreservesUnknownFields(t *testing.T) {
	input := `{
		"sync_token": "tok",
		"full_sync": false,
		"user": {"id": "u1", "email": "test@example.com"},
		"projects": []
	}`

	got, err := FilterSyncResponse([]byte(input), []string{})
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]json.RawMessage
	json.Unmarshal(got, &result)

	if _, ok := result["user"]; !ok {
		t.Error("expected 'user' field to be preserved")
	}
}

func TestFilterSyncResponse_EmptyAllowedList(t *testing.T) {
	input := `{
		"sync_token": "tok",
		"full_sync": true,
		"projects": [{"id": "1", "parent_id": null, "name": "Work"}],
		"items": [{"id": "10", "project_id": "1", "content": "Task"}]
	}`

	got, err := FilterSyncResponse([]byte(input), []string{})
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]json.RawMessage
	json.Unmarshal(got, &result)

	var projects []projectStub
	json.Unmarshal(result["projects"], &projects)
	if len(projects) != 0 {
		t.Errorf("expected 0 projects, got %d", len(projects))
	}
}
