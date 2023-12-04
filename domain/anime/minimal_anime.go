package anime

import "time"

type MinimalAnimeDto struct {
	AnimeId     int
	SeriesId    string
	LastUpdated time.Time
}

type MinimalAnime struct {
	animeId     int
	seriesId    string
	lastUpdated time.Time
}

func ReformMinimalAnime(dto MinimalAnimeDto) MinimalAnime {
	return MinimalAnime{
		animeId:     dto.AnimeId,
		seriesId:    dto.SeriesId,
		lastUpdated: dto.LastUpdated,
	}
}

func (anime MinimalAnime) AnimeId() int {
	return anime.animeId
}

func (anime MinimalAnime) SeriesId() string {
	return anime.seriesId
}

func (anime MinimalAnime) LastUpdated() time.Time {
	return anime.lastUpdated
}
