package saver

import (
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"
)

type animeTranslator interface {
	SaveAll(newAnime []anime.Anime, updatedAnime []anime.Anime) (map[core.SeriesId]anime.MinimalAnime, error)
}

type posterTranslator interface {
	SaveAll(newPosters []anime.Image, updatedPosters []anime.Image) error
}

type latestEpisodesTranslator interface {
	SaveAll(newLatestEpisodes []anime.LatestEpisodes, updatedLatestEpisodes []anime.LatestEpisodes) error
}

type thumbnailTranslator interface {
	SaveAll(newThumbnails []anime.Image, updatedThumbnails []anime.Image) error
}

type AnimeSaver struct {
	animeTranslator          animeTranslator
	posterTranslator         posterTranslator
	latestEpisodesTranslator latestEpisodesTranslator
	thumbnailTranslator      thumbnailTranslator
}

func NewAnimeSaver(
	animeTranslator animeTranslator,
	posterTranslator posterTranslator,
	latestEpisodesTranslator latestEpisodesTranslator,
	thumbnailTranslator thumbnailTranslator,
) AnimeSaver {
	return AnimeSaver{
		animeTranslator:          animeTranslator,
		posterTranslator:         posterTranslator,
		latestEpisodesTranslator: latestEpisodesTranslator,
		thumbnailTranslator:      thumbnailTranslator,
	}
}

// TODO take in updated anime
func (saver AnimeSaver) SaveAll(locales []core.Locale, newAnimes []anime.Anime, updatedAnimes []anime.Anime) error {
	newMinimalAnime, err := saver.animeTranslator.SaveAll(newAnimes, updatedAnimes)
	if err != nil {
		return err
	}

	var newPosters []anime.Image
	var newLatestEpisodes []anime.LatestEpisodes
	var newThumbnails []anime.Image

	for _, newAnime := range newAnimes {
		minimalAnime := newMinimalAnime[newAnime.SeriesId()]

		err := newAnime.AssignAnimeId(minimalAnime.AnimeId())
		if err != nil {
			return err
		}

		newPosters = append(newPosters, newAnime.Posters()...)

		var leToAdd []anime.LatestEpisodes
		for _, locale := range locales {
			leToAdd = append(leToAdd, newAnime.Episodes().GetLatestEpisodesForLocale(locale))
		}

		newLatestEpisodes = append(newLatestEpisodes, leToAdd...)

		newThumbnails = append(newThumbnails, newAnime.Episodes().Thumbnails()...)
	}

	err = saver.posterTranslator.SaveAll(newPosters, nil)
	if err != nil {
		return err
	}

	err = saver.latestEpisodesTranslator.SaveAll(newLatestEpisodes, nil)
	if err != nil {
		return err
	}

	err = saver.thumbnailTranslator.SaveAll(newThumbnails, nil)
	if err != nil {
		return err
	}

	return nil
}
