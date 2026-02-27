//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lowskydev/rickmorty-api/models"
)

func main() {
	resp, err := http.Get("https://rickandmortyapi.com/api/character/1")
	if err != nil {
		fmt.Println("request failed:", err)
		return
	}
	defer resp.Body.Close()

	var c models.Character
	json.NewDecoder(resp.Body).Decode(&c)

	fmt.Printf("ID:       %d\n", c.ID)
	fmt.Printf("Name:     %s\n", c.Name)
	fmt.Printf("URL:      %s\n", c.URL)
	fmt.Printf("Episodes: %d episodes\n", len(c.Episodes))
}
