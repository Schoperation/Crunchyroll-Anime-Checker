package anime

import "time"

type MinimalAnimeDto struct {
	SeriesId    string
	SlugTitle   string
	LastUpdated time.Time
}
