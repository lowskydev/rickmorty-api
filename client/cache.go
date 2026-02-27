package client

import (
	"sync"
	"time"

	"github.com/lowskydev/rickmorty-api/models"
)

// stores the results of FetchAllCharacters and FetchAllEpisodes
// so we only hit the API once instead of on every request
type cache struct {
	mu         sync.Mutex // prevents two requests from fetching simultaneously
	characters []models.Character
	episodes   []models.Episode
	fetchedAt  time.Time
	ttl        time.Duration // how long before we refresh
}

var globalCache = &cache{
	ttl: 1 * time.Hour,
}

// returns cached characters, fetching from the API if the
// cache is empty or expired
func GetAllCharacters() ([]models.Character, error) {
	globalCache.mu.Lock()
	defer globalCache.mu.Unlock()

	if globalCache.characters != nil && time.Since(globalCache.fetchedAt) < globalCache.ttl {
		return globalCache.characters, nil
	}

	chars, err := FetchAllCharacters()
	if err != nil {
		return nil, err
	}

	globalCache.characters = chars
	globalCache.fetchedAt = time.Now()
	return chars, nil
}

// returns cached episodes, fetching from the API if the
// cache is empty or expired
func GetAllEpisodes() ([]models.Episode, error) {
	globalCache.mu.Lock()
	defer globalCache.mu.Unlock()

	if globalCache.episodes != nil && time.Since(globalCache.fetchedAt) < globalCache.ttl {
		return globalCache.episodes, nil
	}

	eps, err := FetchAllEpisodes()
	if err != nil {
		return nil, err
	}

	globalCache.episodes = eps
	globalCache.fetchedAt = time.Now()
	return eps, nil
}
