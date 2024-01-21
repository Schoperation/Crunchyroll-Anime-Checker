package factory

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"
)

type posterTranslator interface {
	GetAllByAnimeIds(animeIds []anime.AnimeId) (map[anime.AnimeId][]anime.Image, error)
}

type latestEpisodesTranslator interface {
	GetAllByAnimeIds(animeIds []anime.AnimeId) (map[anime.AnimeId][]anime.LatestEpisodes, error)
}

type thumbnailTranslator interface {
	GetAllByAnimeIds(animeIds []anime.AnimeId) (map[anime.AnimeId]map[string]anime.Image, error)
}

type AnimeFactory struct {
	posterTranslator         posterTranslator
	latestEpisodesTranslator latestEpisodesTranslator
	thumbnailTranslator      thumbnailTranslator
}

func NewAnimeFactory(
	posterTranslator posterTranslator,
	latestEpisodesTranslator latestEpisodesTranslator,
	thumbnailTranslator thumbnailTranslator,
) AnimeFactory {
	return AnimeFactory{
		posterTranslator:         posterTranslator,
		latestEpisodesTranslator: latestEpisodesTranslator,
		thumbnailTranslator:      thumbnailTranslator,
	}
}

func (factory AnimeFactory) ReformAll(dtos []anime.AnimeDto) (map[core.SeriesId]anime.Anime, error) {
	animeIds := make([]anime.AnimeId, len(dtos))
	for i, dto := range dtos {
		animeIds[i] = anime.ReformAnimeId(dto.AnimeId)
	}

	allPosters, err := factory.posterTranslator.GetAllByAnimeIds(animeIds)
	if err != nil {
		return nil, err
	}

	allLatestEpisodes, err := factory.latestEpisodesTranslator.GetAllByAnimeIds(animeIds)
	if err != nil {
		return nil, err
	}

	allThumbnails, err := factory.thumbnailTranslator.GetAllByAnimeIds(animeIds)
	if err != nil {
		return nil, err
	}

	animes := make(map[core.SeriesId]anime.Anime, len(dtos))
	for j, dto := range dtos {
		animeId := animeIds[j]
		latestEpisodes, exists := allLatestEpisodes[animeId]
		if !exists {
			return nil, fmt.Errorf("could not find reformed latest episodes for anime ID %d", animeId.Int())
		}

		thumbnails, exists := allThumbnails[animeId]
		if !exists {
			return nil, fmt.Errorf("could not find reformed thumbnails for anime ID %d", animeId.Int())
		}

		posters, exists := allPosters[animeId]
		if !exists {
			return nil, fmt.Errorf("could not find reformed posters for anime ID %d", animeId.Int())
		}

		episodes := anime.ReformEpisodeCollection(animeId, latestEpisodes, thumbnails)
		seriesId := core.ReformSeriesId(dto.SeriesId)
		animes[seriesId] = anime.ReformAnime(dto, posters, episodes)
	}

	return animes, nil
}
