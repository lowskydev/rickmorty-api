package models

type Character struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	URL      string   `json:"url"`
	Episodes []string `json:"episode"`
}

type Location struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Episode struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	URL        string   `json:"url"`
	Characters []string `json:"characters"`
}

type PageInfo struct {
	Count int    `json:"count"`
	Pages int    `json:"pages"`
	Next  string `json:"next"`
}

type CharacterPage struct {
	Info    PageInfo    `json:"info"`
	Results []Character `json:"results"`
}

type LocationPage struct {
	Info    PageInfo   `json:"info"`
	Results []Location `json:"results"`
}

type EpisodePage struct {
	Info    PageInfo  `json:"info"`
	Results []Episode `json:"results"`
}

// --- System response types ---

// for /search endpoint
type SearchResult struct {
	Name string `json:"name"`
	Type string `json:"type"` // "character", "location", "episode"
	URL  string `json:"url"`
}

// used inside PairResult
type CharacterRef struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// for /top-pairs endpoint
type PairResult struct {
	Character1 CharacterRef `json:"character1"`
	Character2 CharacterRef `json:"character2"`
	Episodes   int          `json:"episodes"`
}
