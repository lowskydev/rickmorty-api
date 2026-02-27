//go:build ignore

package main

import (
	"fmt"
	"github.com/lowskydev/rickmorty-api/client"
)

func main() {
	chars, err := client.FetchCharactersByName("rick")
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	fmt.Printf("Total characters found: %d\n", len(chars))
	fmt.Printf("First: %s (%s)\n", chars[0].Name, chars[0].URL)
	fmt.Printf("Last:  %s (%s)\n", chars[len(chars)-1].Name, chars[len(chars)-1].URL)
}
