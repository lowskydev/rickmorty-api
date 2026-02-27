package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func Search(w http.ResponseWriter, r *http.Request) {
	term, limit, err := parseSearchParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"term":  term,
		"limit": limit,
	})
}

// reads and validates ?term= and ?limit= from the request URL
func parseSearchParams(r *http.Request) (term string, limit int, err error) {
	term = r.URL.Query().Get("term")
	if term == "" {
		return "", 0, fmt.Errorf("missing required query parameter: term")
	}

	limit = -1 // default: no limit
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			return "", 0, fmt.Errorf("limit must be a positive integer")
		}
	}

	return term, limit, nil
}
