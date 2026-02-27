package models

type Character struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	URL      string   `json:"url"`
	Episodes []string `json:"episode"`
}
