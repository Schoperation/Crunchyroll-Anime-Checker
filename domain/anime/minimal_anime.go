package anime

import "time"

type MinimalAnimeDto struct {
	AnimeId     int
	SeriesId    string
	SlugTitle   string
	LastUpdated time.Time
}
