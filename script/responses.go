package script

import "time"

// File for all responses from Crunchyroll.

// Response for getting an alphabetical list of all anime series.
type allAnimeResponse struct {
	Total int      `json:"total"`
	Data  []series `json:"data"`
}

// Data entry in allAnimeResponse
// This has loads of fields, so only including the ones we need.
type series struct {
	Title      string    `json:"title"`
	SlugTitle  string    `json:"slug_title"`
	LastPublic time.Time `json:"last_public"`
	Id         string    `json:"id"` // series ID (G--------)
	Images     images    `json:"images"`
}

// Struct to store images for allAnimeResponse.
//
// For some reason it's a 2D array. But all entries use the first index of the first array.
type images struct {
	PosterTall [][]image `json:"poster_tall"`
	PosterWide [][]image `json:"poster_wide"`
}

// Struct to store an image.
type image struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Type   string `json:"type"`
	Source string `json:"source"`
}
