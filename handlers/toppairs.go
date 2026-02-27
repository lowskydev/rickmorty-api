package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/lowskydev/rickmorty-api/client"
	"github.com/lowskydev/rickmorty-api/models"
)

func TopPairs(w http.ResponseWriter, r *http.Request) {
	minEp, maxEp, limit, err := parseTopPairsParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch all characters and all episodes concurrently
	type charsResult struct {
		data []models.Character
		err  error
	}
	type epsResult struct {
		data []models.Episode
		err  error
	}

	charsCh := make(chan charsResult, 1)
	epsCh := make(chan epsResult, 1)

	go func() {
		data, err := client.FetchAllCharacters()
		charsCh <- charsResult{data, err}
	}()
	go func() {
		data, err := client.FetchAllEpisodes()
		epsCh <- epsResult{data, err}
	}()

	charsRes := <-charsCh
	epsRes := <-epsCh

	if charsRes.err != nil || epsRes.err != nil {
		http.Error(w, "error fetching data from Rick and Morty API", http.StatusInternalServerError)
		return
	}

	// Build a map of character URL -> Character for fast lookups
	charByURL := make(map[string]models.Character)
	for _, c := range charsRes.data {
		charByURL[c.URL] = c
	}

	// Count how many episodes each pair of characters shares
	pairCount := make(map[string]int)
	for _, ep := range epsRes.data {
		chars := ep.Characters
		// Generate every unique combination of 2 characters in this episode
		for i := 0; i < len(chars); i++ {
			for j := i + 1; j < len(chars); j++ {
				key := pairKey(chars[i], chars[j])
				pairCount[key]++
			}
		}
	}

	// Convert the map into a slice of PairResult, applying min/max filters
	var pairs []models.PairResult
	for key, count := range pairCount {
		if minEp >= 0 && count < minEp {
			continue
		}
		if maxEp >= 0 && count > maxEp {
			continue
		}

		urls := strings.SplitN(key, "|", 2)
		c1, ok1 := charByURL[urls[0]]
		c2, ok2 := charByURL[urls[1]]
		if !ok1 || !ok2 {
			continue
		}

		pairs = append(pairs, models.PairResult{
			Character1: models.CharacterRef{Name: c1.Name, URL: c1.URL},
			Character2: models.CharacterRef{Name: c2.Name, URL: c2.URL},
			Episodes:   count,
		})
	}

	// Sort descending by episode count
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Episodes > pairs[j].Episodes
	})

	if limit >= 0 && len(pairs) > limit {
		pairs = pairs[:limit]
	}

	if pairs == nil {
		pairs = []models.PairResult{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pairs)
}

// pairKey builds a stable unique key for a pair of character URLs
func pairKey(urlA, urlB string) string {
	if urlA > urlB {
		urlA, urlB = urlB, urlA
	}
	return fmt.Sprintf("%s|%s", urlA, urlB)
}

// parseTopPairsParams reads min, max, limit from query params
func parseTopPairsParams(r *http.Request) (minEp, maxEp, limit int, err error) {
	minEp, maxEp, limit = -1, -1, 20

	if s := r.URL.Query().Get("min"); s != "" {
		minEp, err = strconv.Atoi(s)
		if err != nil || minEp < 0 {
			return 0, 0, 0, fmt.Errorf("min must be a non-negative integer")
		}
	}
	if s := r.URL.Query().Get("max"); s != "" {
		maxEp, err = strconv.Atoi(s)
		if err != nil || maxEp < 0 {
			return 0, 0, 0, fmt.Errorf("max must be a non-negative integer")
		}
	}
	if s := r.URL.Query().Get("limit"); s != "" {
		limit, err = strconv.Atoi(s)
		if err != nil || limit < 0 {
			return 0, 0, 0, fmt.Errorf("limit must be a non-negative integer")
		}
	}
	return
}
