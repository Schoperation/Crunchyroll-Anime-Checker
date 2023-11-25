package script

import "time"

// This file defines the structures of the JSON lists used to store anime and their assets.

// anime_atlas_blank.json
//
// Main file to store latest episodes of anime.
//
// The map key for Anime is a slug title (e.g. spy-x-family).
type AnimeAtlas struct {
	TotalCount int              `json:"total_count"`
	Anime      map[string]Anime `json:"anime"`
}

type Anime struct {
	Name        string    `json:"name"`
	LastUpdated time.Time `json:"last_updated"`
	Sub         Episode   `json:"sub"`
	Dub         Episode   `json:"dub"`
}

type Episode struct {
	Season  int    `json:"season"`
	Episode int    `json:"episode"`
	Title   string `json:"title"`
}

// anime_posters.json
//
// File to store posters; the images on the landing page of an anime series.
//
// The map key for Posters is a slug title (e.g. spy-x-family).
type AnimePosters struct {
	TotalCount int               `json:"total_count"`
	Posters    map[string]Poster `json:"posters"`
}

type Poster struct {
	PosterTallURL     string `json:"poster_tall_url"`
	PosterTallEncoded string `json:"poster_tall_encoded"`
	PosterWideURL     string `json:"poster_wide_url"`
	PosterWideEncoded string `json:"poster_wide_encoded"`
}

// anime_episode_thumbnails.json
//
// File to store episode thumbnails.
//
// The first map key is a slug title (spy-x-family).
// The second key is a combo of season and episode (e.g. 1-2 is season 1, episode 2).
type AnimeEpisodeThumbnails struct {
	TotalCount        int                                    `json:"total_count"`
	EpisodeThumbnails map[string]map[string]EpisodeThumbnail `json:"episode_thumbnails"`
}

type EpisodeThumbnail struct {
	URL     string `json:"url"`
	Encoded string `json:"encoded"`
}
