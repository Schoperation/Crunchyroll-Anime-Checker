package factory

import (
	"fmt"
	"schoperation/crunchyroll-anime-checker/domain/anime"
	"schoperation/crunchyroll-anime-checker/domain/core"
)

type allPostersFetcher interface {
	GetAllByAnimeIds(animeIds []anime.AnimeId) (map[anime.AnimeId][]anime.Image, error)
}

type allLatestEpisodesFetcher interface {
	GetAllByAnimeIds(animeIds []anime.AnimeId) (map[anime.AnimeId][]anime.LatestEpisodes, error)
}

type allThumbnailsFetcher interface {
	GetAllByAnimeIds(animeIds []anime.AnimeId) (map[anime.AnimeId]map[string]anime.Image, error)
}

type AnimeFactory struct {
	allPostersFetcher        allPostersFetcher
	allLatestEpisodesFetcher allLatestEpisodesFetcher
	allThumbnailsFetcher     allThumbnailsFetcher
}

func NewAnimeFactory(
	allPostersFetcher allPostersFetcher,
	allLatestEpisodesFetcher allLatestEpisodesFetcher,
	allThumbnailsFetcher allThumbnailsFetcher,
) AnimeFactory {
	return AnimeFactory{
		allPostersFetcher:        allPostersFetcher,
		allLatestEpisodesFetcher: allLatestEpisodesFetcher,
		allThumbnailsFetcher:     allThumbnailsFetcher,
	}
}

func (factory AnimeFactory) ReformAll(dtos []anime.AnimeDto) (map[core.SeriesId]anime.Anime, map[core.SeriesId]anime.Anime, error) {
	animeIds := make([]anime.AnimeId, len(dtos))
	for i, dto := range dtos {
		animeIds[i] = anime.ReformAnimeId(dto.AnimeId)
	}

	allPosters, err := factory.allPostersFetcher.GetAllByAnimeIds(animeIds)
	if err != nil {
		return nil, nil, err
	}

	allLatestEpisodes, err := factory.allLatestEpisodesFetcher.GetAllByAnimeIds(animeIds)
	if err != nil {
		return nil, nil, err
	}

	allThumbnails, err := factory.allThumbnailsFetcher.GetAllByAnimeIds(animeIds)
	if err != nil {
		return nil, nil, err
	}

	// Create and return a copy of the map here so we don't have to do any funky pointer logic...
	animes := make(map[core.SeriesId]anime.Anime, len(dtos))
	originalAnime := make(map[core.SeriesId]anime.Anime, len(dtos))

	for j, dto := range dtos {
		animeId := animeIds[j]
		latestEpisodes, exists := allLatestEpisodes[animeId]
		if !exists {
			return nil, nil, fmt.Errorf("could not find reformed latest episodes for anime ID %d", animeId.Int())
		}

		thumbnails, exists := allThumbnails[animeId]
		if !exists {
			return nil, nil, fmt.Errorf("could not find reformed thumbnails for anime ID %d", animeId.Int())
		}

		posters, exists := allPosters[animeId]
		if !exists {
			return nil, nil, fmt.Errorf("could not find reformed posters for anime ID %d", animeId.Int())
		}

		episodes := anime.ReformEpisodeCollection(animeId, latestEpisodes, thumbnails)
		seriesId := core.ReformSeriesId(dto.SeriesId)
		animes[seriesId] = anime.ReformAnime(dto, posters, episodes)

		origEpisodes := anime.ReformEpisodeCollection(animeId, latestEpisodes, thumbnails)
		originalAnime[seriesId] = anime.ReformAnime(dto, posters, origEpisodes)
	}

	return animes, originalAnime, nil
}
