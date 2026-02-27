package client

import (
	"encoding/json"
	"fmt"
	"github.com/lowskydev/rickmorty-api/models"
	"net/http"
)

const baseURL = "https://rickandmortyapi.com/api"

// sends a GET request to url and decodes the JSON response into target
func getJSON(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: unexpected status %d", url, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode response from %s: %w", url, err)
	}

	return nil
}

// returns all characters whose name contains the search term
func FetchCharactersByName(name string) ([]models.Character, error) {
	var all []models.Character

	// Start at page 1
	// After each page API tells the URL for the next one
	url := fmt.Sprintf("%s/character?name=%s", baseURL, name)

	for url != "" {
		var page models.CharacterPage

		if err := getJSON(url, &page); err != nil {
			return nil, err
		}

		if len(page.Results) == 0 {
			break
		}

		all = append(all, page.Results...)

		// page.Info.Next is either URL "https://...?page=2"
		// or "" if this was the last page
		url = page.Info.Next
	}

	return all, nil
}
