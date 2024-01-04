package anime

import (
	"schoperation/crunchyrollanimestatus/domain/core"
	"time"
)

type MinimalAnimeDto struct {
	AnimeId     int
	SeriesId    string
	LastUpdated time.Time
}

type MinimalAnime struct {
	animeId     AnimeId
	seriesId    core.SeriesId
	lastUpdated time.Time
}

func ReformMinimalAnime(dto MinimalAnimeDto) MinimalAnime {
	return MinimalAnime{
		animeId:     ReformAnimeId(dto.AnimeId),
		seriesId:    core.ReformSeriesId(dto.SeriesId),
		lastUpdated: dto.LastUpdated,
	}
}

func (anime MinimalAnime) AnimeId() AnimeId {
	return anime.animeId
}

func (anime MinimalAnime) SeriesId() core.SeriesId {
	return anime.seriesId
}

func (anime MinimalAnime) LastUpdated() time.Time {
	return anime.lastUpdated
}
