package anime

import (
	"schoperation/crunchyrollanimestatus/domain/core"
	"time"
)

type LatestEpisodesDto struct {
	AnimeId          int
	LocaleId         int
	LastUpdated      time.Time
	LatestSubSeason  int
	LatestSubEpisode int
	LatestSubTitle   string
	LatestDubSeason  int
	LatestDubEpisode int
	LatestDubTitle   string
}

// LatestEpisodes holds the latest sub and dub episodes for a specified anime and locale.
type LatestEpisodes struct {
	animeId     AnimeId
	locale      core.Locale
	lastUpdated time.Time
	latestSub   MinimalEpisode
	latestDub   MinimalEpisode
}

func NewLatestEpisodes(dto LatestEpisodesDto) (LatestEpisodes, error) {
	animeId, err := NewAnimeId(dto.AnimeId)
	if err != nil {
		return LatestEpisodes{}, err
	}

	locale, err := core.NewLocale(dto.LocaleId)
	if err != nil {
		return LatestEpisodes{}, err
	}

	latestSub, err := NewMinimalEpisode(dto.LatestSubSeason, dto.LatestSubEpisode, dto.LatestSubTitle)
	if err != nil {
		return LatestEpisodes{}, err
	}

	latestDub, err := NewMinimalEpisode(dto.LatestDubSeason, dto.LatestDubEpisode, dto.LatestDubTitle)
	if err != nil {
		return LatestEpisodes{}, err
	}

	return LatestEpisodes{
		animeId:     animeId,
		locale:      locale,
		lastUpdated: time.Now().UTC(),
		latestSub:   latestSub,
		latestDub:   latestDub,
	}, nil
}

func ReformLatestEpisodes(dto LatestEpisodesDto) LatestEpisodes {
	return LatestEpisodes{
		animeId:     ReformAnimeId(dto.AnimeId),
		locale:      core.ReformLocale(dto.LocaleId),
		lastUpdated: dto.LastUpdated,
		latestSub:   ReformMinimalEpisode(dto.LatestSubSeason, dto.LatestSubEpisode, dto.LatestSubTitle),
		latestDub:   ReformMinimalEpisode(dto.LatestDubSeason, dto.LatestDubEpisode, dto.LatestDubTitle),
	}
}

func (le LatestEpisodes) Locale() core.Locale {
	return le.locale
}
