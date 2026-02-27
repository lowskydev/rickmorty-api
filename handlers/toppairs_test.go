package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lowskydev/rickmorty-api/models"
)

func TestParseTopPairsParams_Defaults(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/top-pairs", nil)
	minEp, maxEp, limit, err := parseTopPairsParams(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if minEp != -1 {
		t.Errorf("expected minEp -1, got %d", minEp)
	}
	if maxEp != -1 {
		t.Errorf("expected maxEp -1, got %d", maxEp)
	}
	if limit != 20 {
		t.Errorf("expected default limit 20, got %d", limit)
	}
}

func TestParseTopPairsParams_BadMin(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/top-pairs?min=abc", nil)
	_, _, _, err := parseTopPairsParams(r)
	if err == nil {
		t.Error("expected error for bad min, got nil")
	}
}

func TestParseTopPairsParams_BadMax(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/top-pairs?max=abc", nil)
	_, _, _, err := parseTopPairsParams(r)
	if err == nil {
		t.Error("expected error for bad max, got nil")
	}
}

func TestParseTopPairsParams_BadLimit(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/top-pairs?limit=abc", nil)
	_, _, _, err := parseTopPairsParams(r)
	if err == nil {
		t.Error("expected error for bad limit, got nil")
	}
}

func TestPairKey_Symmetric(t *testing.T) {
	k1 := pairKey("https://rickandmortyapi.com/api/character/1", "https://rickandmortyapi.com/api/character/2")
	k2 := pairKey("https://rickandmortyapi.com/api/character/2", "https://rickandmortyapi.com/api/character/1")
	if k1 != k2 {
		t.Errorf("pairKey should be symmetric, got %q and %q", k1, k2)
	}
}

func TestPairKey_DifferentPairsDifferentKeys(t *testing.T) {
	k1 := pairKey("https://rickandmortyapi.com/api/character/1", "https://rickandmortyapi.com/api/character/2")
	k2 := pairKey("https://rickandmortyapi.com/api/character/1", "https://rickandmortyapi.com/api/character/3")
	if k1 == k2 {
		t.Error("different pairs should produce different keys")
	}
}

func TestCountPairs_BasicCount(t *testing.T) {
	episodes := []models.Episode{
		{Characters: []string{"url/1", "url/2", "url/3"}},
		{Characters: []string{"url/1", "url/2"}},
	}

	counts := countPairs(episodes)

	k := pairKey("url/1", "url/2")
	if counts[k] != 2 {
		t.Errorf("expected pair (1,2) count 2, got %d", counts[k])
	}
}

func TestCountPairs_NoDuplicates(t *testing.T) {
	episodes := []models.Episode{
		{Characters: []string{"url/1", "url/2", "url/3"}},
	}

	counts := countPairs(episodes)

	if len(counts) != 3 {
		t.Errorf("expected 3 pairs from 3 characters, got %d", len(counts))
	}
}

func TestCountPairs_Empty(t *testing.T) {
	counts := countPairs([]models.Episode{})
	if len(counts) != 0 {
		t.Errorf("expected 0 pairs for empty episodes, got %d", len(counts))
	}
}
