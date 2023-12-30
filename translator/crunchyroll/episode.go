package crunchyroll

import (
	"schoperation/crunchyrollanimestatus/domain/core"
	"schoperation/crunchyrollanimestatus/domain/crunchyroll"
)

type crunchyrollAnimeEpisodesClient interface {
	GetAllEpisodesBySeasonId(locale, seasonId string) ([]crunchyroll.EpisodeDto, error)
}

type EpisodeTranslator struct {
	crunchyrollAnimeEpisodesClient crunchyrollAnimeEpisodesClient
}

func NewEpisodeTranslator(crunchyrollAnimeEpisodesClient crunchyrollAnimeEpisodesClient) EpisodeTranslator {
	return EpisodeTranslator{
		crunchyrollAnimeEpisodesClient: crunchyrollAnimeEpisodesClient,
	}
}

func (translator EpisodeTranslator) GetAllEpisodesBySeasonId(locale core.Locale, seasonId string) (crunchyroll.EpisodeCollection, error) {
	dtos, err := translator.crunchyrollAnimeEpisodesClient.GetAllEpisodesBySeasonId(locale.Name(), seasonId)
	if err != nil {
		return crunchyroll.EpisodeCollection{}, err
	}

	episodes := make([]crunchyroll.Episode, len(dtos))
	for i, dto := range dtos {
		episodes[i] = crunchyroll.ReformEpisode(dto)
	}

	return crunchyroll.NewEpisodeCollection(seasonId, episodes)
}
