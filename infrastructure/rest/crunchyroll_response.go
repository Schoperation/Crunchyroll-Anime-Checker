package rest

import "time"

// Separate file for responses from Crunchyroll since the client file was hurting my eyes...
// And to help with caching

type allAnimeResponse struct {
	Total int      `json:"total"`
	Data  []series `json:"data"`
}

type series struct {
	Id             string         `json:"id"` // series ID (G--------)
	SlugTitle      string         `json:"slug_title"`
	Title          string         `json:"title"`
	LastPublic     time.Time      `json:"last_public"`
	Images         images         `json:"images"`
	SeriesMetaData seriesMetaData `json:"series_metadata"`
}

type seriesMetaData struct {
	SeasonCount  int `json:"season_count"`
	EpisodeCount int `json:"episode_count"`
}

// For some reason it's a 2D array. But all entries use the first index of the first array...
type images struct {
	PosterTall [][]image `json:"poster_tall"`
	PosterWide [][]image `json:"poster_wide"`
}

type image struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Type   string `json:"type"`
	Source string `json:"source"`
}

type seasonsResponse struct {
	Total int      `json:"total"`
	Data  []season `json:"data"`
}

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
