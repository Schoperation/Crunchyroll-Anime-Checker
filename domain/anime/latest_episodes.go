package anime

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/core"
)

type LatestEpisodesDto struct {
	AnimeId          int
	LocaleId         int
	LatestSubSeason  int
	LatestSubEpisode int
	LatestSubTitle   string
	LatestDubSeason  int
	LatestDubEpisode int
	LatestDubTitle   string
}

// LatestEpisodes holds the latest sub and dub episodes for a specified anime and locale.
type LatestEpisodes struct {
	animeId   AnimeId
	locale    core.Locale
	latestSub MinimalEpisode
	latestDub MinimalEpisode
}

func NewLatestEpisodes(dto LatestEpisodesDto) (LatestEpisodes, error) {
	animeId, err := NewAnimeId(dto.AnimeId)
	if err != nil {
		return LatestEpisodes{}, err
	}

	locale, err := core.NewLocaleFromId(dto.LocaleId)
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

	if latestSub.IsBlank() && latestDub.IsBlank() {
		return LatestEpisodes{}, fmt.Errorf("latest episodes must have at least a sub or dub for locale %s, anime ID %d", locale.Name(), animeId.Int())
	}

	return LatestEpisodes{
		animeId:   animeId,
		locale:    locale,
		latestSub: latestSub,
		latestDub: latestDub,
	}, nil
}

func ReformLatestEpisodes(dto LatestEpisodesDto) LatestEpisodes {
	return LatestEpisodes{
		animeId:   ReformAnimeId(dto.AnimeId),
		locale:    core.ReformLocaleFromId(dto.LocaleId),
		latestSub: ReformMinimalEpisode(dto.LatestSubSeason, dto.LatestSubEpisode, dto.LatestSubTitle),
		latestDub: ReformMinimalEpisode(dto.LatestDubSeason, dto.LatestDubEpisode, dto.LatestDubTitle),
	}
}

func (le LatestEpisodes) AnimeId() AnimeId {
	return le.animeId
}

func (le LatestEpisodes) Locale() core.Locale {
	return le.locale
}

func (le LatestEpisodes) LatestSub() MinimalEpisode {
	return le.latestSub
}

func (le LatestEpisodes) LatestDub() MinimalEpisode {
	return le.latestDub
}
