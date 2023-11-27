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
	Title          string         `json:"title"`
	SlugTitle      string         `json:"slug_title"`
	LastPublic     time.Time      `json:"last_public"`
	Id             string         `json:"id"` // series ID (G--------)
	Images         images         `json:"images"`
	SeriesMetaData seriesMetaData `json:"series_metadata"`
}

// Used in a series struct for allAnimeResponse.
type seriesMetaData struct {
	SeasonCount     int      `json:"season_count"`
	EpisodeCount    int      `json:"episode_count"`
	AudioLocales    []string `json:"audio_locales"`
	SubtitleLocales []string `json:"subtitle_locales"`
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

// ==========================================================================================
// Response for getting an anime's list of seasons.
type seasonsResponse struct {
	Total int      `json:"total"`
	Data  []season `json:"data"`
}

// Struct for storing a season when retrieving an anime's list of seasons.
type season struct {
	Id              string    `json:"id"`
	Identifier      string    `json:"identifier"`
	SeasonNumber    int       `json:"season_number"`
	AudioLocales    []string  `json:"audio_locales"`
	SubtitleLocales []string  `json:"subtitle_locales"`
	Versions        []version `json:"versions"`
}

// An array of "versions", or just different seasons in different locales.
// Returned in the response for getting an anime's list of seasons.
type version struct {
	AudioLocale string `json:"audio_locale"`
	GUID        string `json:"guid"`     // season ID (G--------)
	Original    bool   `json:"original"` // Usually identifies the Japanese version
	Variant     string `json:"variant"`
}

// ==========================================================================================
// Response for getting an anime's list of episodes in a season.
type seasonEpisodesResponse struct {
	Total int             `json:"total"`
	Data  []seasonEpisode `json:"data"`
}

type seasonEpisode struct {
	Number          int       `json:"episode_number"`
	Title           string    `json:"title"`
	SubtitleLocales []string  `json:"subtitle_locales"`
	Versions        []version `json:"versions"`
}
