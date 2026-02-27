package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/lowskydev/rickmorty-api/models"
)

const baseURL = "https://rickandmortyapi.com/api"

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

func FetchCharactersByName(name string) ([]models.Character, error) {
	var all []models.Character
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
		url = page.Info.Next
	}
	return all, nil
}

func FetchLocationsByName(name string) ([]models.Location, error) {
	var all []models.Location
	url := fmt.Sprintf("%s/location?name=%s", baseURL, name)

	for url != "" {
		var page models.LocationPage
		if err := getJSON(url, &page); err != nil {
			return nil, err
		}
		if len(page.Results) == 0 {
			break
		}
		all = append(all, page.Results...)
		url = page.Info.Next
	}
	return all, nil
}

func FetchEpisodesByName(name string) ([]models.Episode, error) {
	var all []models.Episode
	url := fmt.Sprintf("%s/episode?name=%s", baseURL, name)

	for url != "" {
		var page models.EpisodePage
		if err := getJSON(url, &page); err != nil {
			return nil, err
		}
		if len(page.Results) == 0 {
			break
		}
		all = append(all, page.Results...)
		url = page.Info.Next
	}
	return all, nil
}

func FetchAllCharacters() ([]models.Character, error) {
	var all []models.Character
	url := fmt.Sprintf("%s/character", baseURL)

	for url != "" {
		var page models.CharacterPage
		if err := getJSON(url, &page); err != nil {
			return nil, err
		}
		if len(page.Results) == 0 {
			break
		}
		all = append(all, page.Results...)
		url = page.Info.Next
		if url != "" {
			time.Sleep(250 * time.Millisecond)
		}
	}
	return all, nil
}

func FetchAllEpisodes() ([]models.Episode, error) {
	var all []models.Episode
	url := fmt.Sprintf("%s/episode", baseURL)

	for url != "" {
		var page models.EpisodePage
		if err := getJSON(url, &page); err != nil {
			return nil, err
		}
		if len(page.Results) == 0 {
			break
		}
		all = append(all, page.Results...)
		url = page.Info.Next
		if url != "" {
			time.Sleep(250 * time.Millisecond)
		}
	}
	return all, nil
}
