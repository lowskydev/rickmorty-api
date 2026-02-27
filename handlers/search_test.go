package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lowskydev/rickmorty-api/models"
)

func TestParseSearchParams_MissingTerm(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/search", nil)
	_, _, err := parseSearchParams(r)
	if err == nil {
		t.Error("expected an error for missing term, got nil")
	}
}

func TestParseSearchParams_BadLimit(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/search?term=rick&limit=abc", nil)
	_, _, err := parseSearchParams(r)
	if err == nil {
		t.Error("expected an error for bad limit, got nil")
	}
}

func TestParseSearchParams_NegativeLimit(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/search?term=rick&limit=-5", nil)
	_, _, err := parseSearchParams(r)
	if err == nil {
		t.Error("expected an error for negative limit, got nil")
	}
}

func TestParseSearchParams_NoLimit(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/search?term=rick", nil)
	term, limit, err := parseSearchParams(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if term != "rick" {
		t.Errorf("expected term 'rick', got '%s'", term)
	}
	if limit != -1 {
		t.Errorf("expected limit -1, got %d", limit)
	}
}

func TestParseSearchParams_WithLimit(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/search?term=rick&limit=10", nil)
	_, limit, err := parseSearchParams(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if limit != 10 {
		t.Errorf("expected limit 10, got %d", limit)
	}
}

func TestSearch_ReturnsResults(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/search?term=rick&limit=5", nil)
	w := httptest.NewRecorder()

	Search(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}

	var results []models.SearchResult
	if err := json.NewDecoder(w.Body).Decode(&results); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if len(results) > 5 {
		t.Errorf("limit=5 but got %d results", len(results))
	}
	for i, res := range results {
		if res.Name == "" || res.Type == "" || res.URL == "" {
			t.Errorf("results[%d] has empty fields: %+v", i, res)
		}
	}
}

func TestSearch_NoResults(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/search?term=zzzzzzzz", nil)
	w := httptest.NewRecorder()

	Search(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var results []models.SearchResult
	json.NewDecoder(w.Body).Decode(&results)

	if results == nil {
		t.Error("expected empty array [], got null")
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
