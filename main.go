package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	// TODOIST_PROXY_ALLOW is a comma-separated list of Todoist project IDs whose
	// data (and the data of their descendant projects) will be included in sync
	// responses. All other project data is stripped before returning to clients.
	allowedIDs := parseAllowedProjects(os.Getenv("TODOIST_PROXY_ALLOW"))
	if len(allowedIDs) == 0 {
		log.Println("warning: TODOIST_PROXY_ALLOW is not set; all projects will be filtered out")
	} else {
		log.Printf("allowing %d project(s): %v", len(allowedIDs), allowedIDs)
	}

	proxy := newReverseProxy()

	mux := http.NewServeMux()
	// The sync endpoint is intercepted so the response can be filtered.
	// Every other path is forwarded to api.todoist.com without modification.
	mux.Handle("/api/v1/sync", newSyncHandler(allowedIDs))
	mux.Handle("/", proxy)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

// parseAllowedProjects splits a comma-separated env var value into a slice of
// trimmed, non-empty project ID strings.
func parseAllowedProjects(env string) []string {
	if env == "" {
		return nil
	}
	parts := strings.Split(env, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if id := strings.TrimSpace(p); id != "" {
			result = append(result, id)
		}
	}
	return result
}
