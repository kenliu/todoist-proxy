package main

import (
	"encoding/json"
)

// projectStub holds just the fields needed for hierarchy expansion.
type projectStub struct {
	ID       string  `json:"id"`
	ParentID *string `json:"parent_id"`
}

// resourceStub holds just the project_id field for filtering resource arrays.
type resourceStub struct {
	ProjectID string `json:"project_id"`
}

// FilterSyncResponse filters a raw Todoist Sync API JSON response, keeping only
// projects in the allowed set (and their descendants) and resources belonging to
// those projects. Unknown top-level fields are preserved unchanged.
func FilterSyncResponse(body []byte, seedIDs []string) ([]byte, error) {
	// Parse the top-level response into a map so unknown fields are preserved.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	// Build the full allowed set by expanding seed IDs through the project hierarchy.
	allowedSet, err := expandAllowedProjects(raw["projects"], seedIDs)
	if err != nil {
		return nil, err
	}

	// Filter each project-scoped resource array.
	resourceKeys := []string{"projects", "items", "sections", "labels", "filters", "reminders"}
	for _, key := range resourceKeys {
		data, ok := raw[key]
		if !ok {
			continue
		}
		filtered, err := filterByProject(data, key == "projects", allowedSet)
		if err != nil {
			return nil, err
		}
		raw[key] = filtered
	}

	return json.Marshal(raw)
}

// expandAllowedProjects takes the raw projects JSON array and a seed list of
// project IDs, and returns the full set of allowed IDs (seeds + all descendants).
func expandAllowedProjects(projectsRaw json.RawMessage, seedIDs []string) (map[string]bool, error) {
	allowed := make(map[string]bool, len(seedIDs))
	for _, id := range seedIDs {
		allowed[id] = true
	}

	if projectsRaw == nil {
		return allowed, nil
	}

	var projects []projectStub
	if err := json.Unmarshal(projectsRaw, &projects); err != nil {
		return nil, err
	}

	// Build parent→children map.
	children := make(map[string][]string)
	for _, p := range projects {
		if p.ParentID != nil {
			children[*p.ParentID] = append(children[*p.ParentID], p.ID)
		}
	}

	// BFS to expand allowed set to include all descendants.
	queue := make([]string, 0, len(seedIDs))
	for _, id := range seedIDs {
		queue = append(queue, id)
	}
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		for _, child := range children[id] {
			if !allowed[child] {
				allowed[child] = true
				queue = append(queue, child)
			}
		}
	}

	return allowed, nil
}

// filterByProject filters a JSON array, keeping elements whose project_id is in
// allowedSet. For the projects array itself, it matches on "id" instead.
func filterByProject(data json.RawMessage, isProjectsArray bool, allowedSet map[string]bool) (json.RawMessage, error) {
	var rawItems []json.RawMessage
	if err := json.Unmarshal(data, &rawItems); err != nil {
		return data, nil // not an array; return as-is
	}

	result := make([]json.RawMessage, 0, len(rawItems))
	for _, item := range rawItems {
		var keep bool
		if isProjectsArray {
			var p projectStub
			if err := json.Unmarshal(item, &p); err == nil {
				keep = allowedSet[p.ID]
			}
		} else {
			var r resourceStub
			if err := json.Unmarshal(item, &r); err == nil {
				keep = r.ProjectID == "" || allowedSet[r.ProjectID]
			}
		}
		if keep {
			result = append(result, item)
		}
	}

	return json.Marshal(result)
}
