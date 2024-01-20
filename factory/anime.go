package factory

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"
)

type posterTranslator interface {
	GetAllByAnimeId(animeId anime.AnimeId) ([]anime.Image, error)
	GetAllByAnimeIds(animeIds []anime.AnimeId) (map[anime.AnimeId][]anime.Image, error)
}

type latestEpisodesTranslator interface {
	GetAllByAnimeId(animeId anime.AnimeId) ([]anime.LatestEpisodes, error)
	GetAllByAnimeIds(animeIds []anime.AnimeId) (map[anime.AnimeId][]anime.LatestEpisodes, error)
}

type thumbnailTranslator interface {
	GetAllByAnimeId(animeId anime.AnimeId) (map[string]anime.Image, error)
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

func (factory AnimeFactory) Reform(dto anime.AnimeDto) (anime.Anime, error) {
	animeId := anime.ReformAnimeId(dto.AnimeId)

	posters, err := factory.posterTranslator.GetAllByAnimeId(animeId)
	if err != nil {
		return anime.Anime{}, err
	}

	latestEpisodes, err := factory.latestEpisodesTranslator.GetAllByAnimeId(animeId)
	if err != nil {
		return anime.Anime{}, err
	}

	thumbnails, err := factory.thumbnailTranslator.GetAllByAnimeId(animeId)
	if err != nil {
		return anime.Anime{}, err
	}

	episodes := anime.ReformEpisodeCollection(animeId, latestEpisodes, thumbnails)

	return anime.ReformAnime(dto, posters, episodes), nil
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
