package client

import "testing"

func TestFetchCharactersByName_Found(t *testing.T) {
	chars, err := FetchCharactersByName("rick")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chars) <= 20 {
		t.Errorf("expected more than 20 results (pagination check), got %d", len(chars))
	}
	for i, c := range chars {
		if c.Name == "" {
			t.Errorf("chars[%d].Name is empty", i)
		}
		if c.URL == "" {
			t.Errorf("chars[%d].URL is empty", i)
		}
	}
}

func TestFetchCharactersByName_NotFound(t *testing.T) {
	chars, err := FetchCharactersByName("zzzzzzzzzzzzz")
	if err != nil {
		t.Fatalf("expected no error for 404 response, got: %v", err)
	}
	if len(chars) != 0 {
		t.Errorf("expected 0 results, got %d", len(chars))
	}
}

func TestFetchLocationsByName_Found(t *testing.T) {
	locs, err := FetchLocationsByName("earth")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(locs) == 0 {
		t.Error("expected at least one location, got none")
	}
}

func TestFetchEpisodesByName_Found(t *testing.T) {
	eps, err := FetchEpisodesByName("rick")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(eps) == 0 {
		t.Error("expected at least one episode, got none")
	}
}
