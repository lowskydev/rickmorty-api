package client

import (
	"sync"
	"time"

	"github.com/lowskydev/rickmorty-api/models"
)

type cache struct {
	mu sync.Mutex

	characters    []models.Character
	charFetchedAt time.Time

	episodes    []models.Episode
	epFetchedAt time.Time

	ttl time.Duration
}

var globalCache = &cache{
	ttl: 1 * time.Hour,
}

func GetAllCharacters() ([]models.Character, error) {
	globalCache.mu.Lock()
	defer globalCache.mu.Unlock()

	if globalCache.characters != nil && time.Since(globalCache.charFetchedAt) < globalCache.ttl {
		return globalCache.characters, nil
	}

	chars, err := FetchAllCharacters()
	if err != nil {
		return nil, err
	}

	globalCache.characters = chars
	globalCache.charFetchedAt = time.Now()
	return chars, nil
}

func GetAllEpisodes() ([]models.Episode, error) {
	globalCache.mu.Lock()
	defer globalCache.mu.Unlock()

	if globalCache.episodes != nil && time.Since(globalCache.epFetchedAt) < globalCache.ttl {
		return globalCache.episodes, nil
	}

	eps, err := FetchAllEpisodes()
	if err != nil {
		return nil, err
	}

	globalCache.episodes = eps
	globalCache.epFetchedAt = time.Now()
	return eps, nil
}
