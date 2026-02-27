package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lowskydev/rickmorty-api/client"
	"github.com/lowskydev/rickmorty-api/models"
)

func Search(w http.ResponseWriter, r *http.Request) {
	term, limit, err := parseSearchParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	type fetchResult struct {
		results []models.SearchResult
		err     error
	}

	charCh := make(chan fetchResult, 1)
	locCh := make(chan fetchResult, 1)
	epCh := make(chan fetchResult, 1)

	// Launch three goroutines. Each one calls a different API endpoint
	go func() {
		chars, err := client.FetchCharactersByName(term)
		if err != nil {
			charCh <- fetchResult{err: err}
			return
		}
		var results []models.SearchResult
		for _, c := range chars {
			results = append(results, models.SearchResult{
				Name: c.Name,
				Type: "character",
				URL:  c.URL,
			})
		}
		charCh <- fetchResult{results: results}
	}()

	go func() {
		locs, err := client.FetchLocationsByName(term)
		if err != nil {
			locCh <- fetchResult{err: err}
			return
		}
		var results []models.SearchResult
		for _, l := range locs {
			results = append(results, models.SearchResult{
				Name: l.Name,
				Type: "location",
				URL:  l.URL,
			})
		}
		locCh <- fetchResult{results: results}
	}()

	go func() {
		eps, err := client.FetchEpisodesByName(term)
		if err != nil {
			epCh <- fetchResult{err: err}
			return
		}
		var results []models.SearchResult
		for _, e := range eps {
			results = append(results, models.SearchResult{
				Name: e.Name,
				Type: "episode",
				URL:  e.URL,
			})
		}
		epCh <- fetchResult{results: results}
	}()

	charRes := <-charCh
	locRes := <-locCh
	epRes := <-epCh

	if charRes.err != nil || locRes.err != nil || epRes.err != nil {
		http.Error(w, "error fetching data from Rick and Morty API", http.StatusInternalServerError)
		return
	}

	var all []models.SearchResult
	all = append(all, charRes.results...)
	all = append(all, locRes.results...)
	all = append(all, epRes.results...)

	// Apply limit - keep only the first `limit` items
	if limit >= 0 && len(all) > limit {
		all = all[:limit]
	}

	if all == nil {
		all = []models.SearchResult{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(all)
}

func parseSearchParams(r *http.Request) (term string, limit int, err error) {
	term = r.URL.Query().Get("term")
	if term == "" {
		return "", 0, fmt.Errorf("missing required query parameter: term")
	}

	limit = -1
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			return "", 0, fmt.Errorf("limit must be a positive integer")
		}
	}

	return term, limit, nil
}
