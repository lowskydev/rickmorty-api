package handlers

import (
	"encoding/json"
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

// integration test against the real API
func TestTopPairs_ReturnsResults(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/top-pairs?limit=5", nil)
	w := httptest.NewRecorder()

	TopPairs(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var pairs []models.PairResult
	if err := json.NewDecoder(w.Body).Decode(&pairs); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(pairs) == 0 {
		t.Fatal("expected at least one pair, got none")
	}
	if len(pairs) > 5 {
		t.Errorf("limit=5 but got %d pairs", len(pairs))
	}

	// Verify results are sorted descending
	for i := 1; i < len(pairs); i++ {
		if pairs[i].Episodes > pairs[i-1].Episodes {
			t.Errorf("pairs not sorted descending: pairs[%d].Episodes=%d > pairs[%d].Episodes=%d",
				i, pairs[i].Episodes, i-1, pairs[i-1].Episodes)
		}
	}

	// Verify each pair has all fields populated
	for i, p := range pairs {
		if p.Character1.Name == "" || p.Character1.URL == "" {
			t.Errorf("pairs[%d].Character1 has empty fields", i)
		}
		if p.Character2.Name == "" || p.Character2.URL == "" {
			t.Errorf("pairs[%d].Character2 has empty fields", i)
		}
		if p.Episodes == 0 {
			t.Errorf("pairs[%d].Episodes is 0", i)
		}
	}
}
