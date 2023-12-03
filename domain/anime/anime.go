package anime

import "time"

type AnimeDto struct {
	Title       string
	SlugTitle   string
	SeriesId    string
	LastUpdated time.Time
}

type Anime struct {
	title       string
	slugTitle   string
	seriesId    string
	lastUpdated time.Time
	latestSub   Episode
	latestDub   Episode
}
