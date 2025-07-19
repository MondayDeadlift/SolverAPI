package model

import "time"

type CodewarsKata struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Slug      string   `json:"slug"`
	URL       string   `json:"url"`
	Tags      []string `json:"tags"`
	Languages []string `json:"languages"`
}

type Kata struct {
	CodewarsKata
	AddedAt time.Time `json:"added_at"`
}
