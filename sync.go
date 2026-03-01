package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

const todoistSyncURL = "https://api.todoist.com/api/v1/sync"

// newSyncHandler returns an http.HandlerFunc that proxies requests to the
// Todoist Sync API and then filters the response to only include data for
// the allowed projects (and their descendants). The client's credentials are
// forwarded as-is; this proxy does not store or inspect them.
func newSyncHandler(allowedIDs []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read the incoming request body.
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}

		// Forward the request to Todoist.
		upstream, err := http.NewRequestWithContext(r.Context(), http.MethodPost, todoistSyncURL, bytes.NewReader(body))
		if err != nil {
			http.Error(w, "failed to build upstream request", http.StatusInternalServerError)
			return
		}
		// Copy relevant headers (Authorization, Content-Type, User-Agent).
		for _, h := range []string{"Authorization", "Content-Type", "User-Agent", "X-Request-Id"} {
			if v := r.Header.Get(h); v != "" {
				upstream.Header.Set(h, v)
			}
		}

		resp, err := http.DefaultClient.Do(upstream)
		if err != nil {
			http.Error(w, "upstream request failed", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "failed to read upstream response", http.StatusBadGateway)
			return
		}

		// Only filter successful JSON responses.
		if resp.StatusCode == http.StatusOK {
			filtered, err := FilterSyncResponse(respBody, allowedIDs)
			if err != nil {
				log.Printf("filter error: %v", err)
				// Fall back to unfiltered response on parse error.
				filtered = respBody
			}
			respBody = filtered
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
	}
}
