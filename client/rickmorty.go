package client

import (
	"encoding/json"
	"fmt"
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
