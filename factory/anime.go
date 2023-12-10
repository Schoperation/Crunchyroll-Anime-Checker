package factory

import (
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"
)

type posterTranslator interface {
	GetAllByAnimeId(animeId anime.AnimeId) ([]anime.Image, error)
}

type latestEpisodesTranslator interface {
	GetAllByAnimeId(animeId anime.AnimeId) (map[core.Locale]anime.LatestEpisodes, error)
}

type thumbnailTranslator interface {
	GetAllByAnimeId(animeId anime.AnimeId) (map[string]anime.Image, error)
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

}
