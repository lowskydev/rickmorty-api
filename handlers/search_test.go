package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
	// -1 means "no limit was specified"
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
