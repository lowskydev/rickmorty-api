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
		data, err := client.GetAllCharacters()
		charsCh <- charsResult{data, err}
	}()
	go func() {
		data, err := client.GetAllEpisodes()
		epsCh <- epsResult{data, err}
	}()

	charsRes := <-charsCh
	epsRes := <-epsCh

	if charsRes.err != nil || epsRes.err != nil {
		if charsRes.err != nil {
			fmt.Printf("characters error: %v\n", charsRes.err)
		}
		if epsRes.err != nil {
			fmt.Printf("episodes error: %v\n", epsRes.err)
		}
		http.Error(w, "error fetching data from Rick and Morty API", http.StatusInternalServerError)
		return
	}

	charByURL := make(map[string]models.Character)
	for _, c := range charsRes.data {
		charByURL[c.URL] = c
	}

	pairCount := countPairs(epsRes.data)

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

func countPairs(episodes []models.Episode) map[string]int {
	pairCount := make(map[string]int)
	for _, ep := range episodes {
		chars := ep.Characters
		for i := 0; i < len(chars); i++ {
			for j := i + 1; j < len(chars); j++ {
				key := pairKey(chars[i], chars[j])
				pairCount[key]++
			}
		}
	}
	return pairCount
}

func pairKey(urlA, urlB string) string {
	if urlA > urlB {
		urlA, urlB = urlB, urlA
	}
	return fmt.Sprintf("%s|%s", urlA, urlB)
}

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
