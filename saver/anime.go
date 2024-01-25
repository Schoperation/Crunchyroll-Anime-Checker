package saver

import (
	"schoperation/crunchyroll-anime-checker/domain/anime"
	"schoperation/crunchyroll-anime-checker/domain/core"
)

type animeSaver interface {
	SaveAll(newAnime []anime.Anime, updatedAnime []anime.Anime) (map[core.SeriesId]anime.MinimalAnime, error)
}

type posterSaver interface {
	SaveAll(newPosters []anime.Image, updatedPosters []anime.Image) error
}

type latestEpisodesSaver interface {
	SaveAll(newLatestEpisodes []anime.LatestEpisodes, updatedLatestEpisodes []anime.LatestEpisodes) error
}

type thumbnailSaver interface {
	SaveAll(newThumbnails []anime.Image, deletedThumbnails []anime.Image) error
}

type AnimeSaver struct {
	animeSaver          animeSaver
	posterSaver         posterSaver
	latestEpisodesSaver latestEpisodesSaver
	thumbnailSaver      thumbnailSaver
}

func NewAnimeSaver(
	animeSaver animeSaver,
	posterSaver posterSaver,
	latestEpisodesSaver latestEpisodesSaver,
	thumbnailSaver thumbnailSaver,
) AnimeSaver {
	return AnimeSaver{
		animeSaver:          animeSaver,
		posterSaver:         posterSaver,
		latestEpisodesSaver: latestEpisodesSaver,
		thumbnailSaver:      thumbnailSaver,
	}
}

func (saver AnimeSaver) SaveAll(
	locales []core.Locale,
	newAnimes []anime.Anime,
	updatedAnimes []anime.Anime,
	originalLocalAnimes map[core.SeriesId]anime.Anime,
) error {
	newMinimalAnime, err := saver.animeSaver.SaveAll(newAnimes, updatedAnimes)
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

		for _, locale := range locales {
			le, err := newAnime.Episodes().GetLatestEpisodesForLocale(locale)
			if err != nil {
				return err
			}

			newLatestEpisodes = append(newLatestEpisodes, le)
		}

		newThumbnails = append(newThumbnails, newAnime.Episodes().Thumbnails()...)
	}

	var updatedPosters []anime.Image
	var updatedLatestEpisodes []anime.LatestEpisodes
	var deletedThumbnails []anime.Image

	for _, updatedAnime := range updatedAnimes {
		originalAnime := originalLocalAnimes[updatedAnime.SeriesId()]

		updatedPosters = append(updatedPosters, updatedAnime.Posters()...)

		// TODO allow removed locales for latest episodes?
		for _, locale := range updatedAnime.Episodes().Locales() {
			_, err := originalAnime.Episodes().GetLatestEpisodesForLocale(locale)
			if err != nil {
				newLe, err := updatedAnime.Episodes().GetLatestEpisodesForLocale(locale)
				if err != nil {
					return err
				}

				newLatestEpisodes = append(newLatestEpisodes, newLe)
				continue
			}

			updatedLe, err := updatedAnime.Episodes().GetLatestEpisodesForLocale(locale)
			if err != nil {
				return err
			}

			updatedLatestEpisodes = append(updatedLatestEpisodes, updatedLe)
		}

		existingThumbnails := make(map[string]anime.Image)
		for _, origThumbnail := range originalAnime.Episodes().Thumbnails() {
			existingThumbnails[origThumbnail.Key()] = origThumbnail
		}

		for _, newThumbnail := range updatedAnime.Episodes().Thumbnails() {
			if _, existing := existingThumbnails[newThumbnail.Key()]; existing {
				continue
			}

			newThumbnails = append(newThumbnails, newThumbnail)
		}

		deletedThumbnails = append(deletedThumbnails, updatedAnime.Episodes().CleanEpisodes()...)
	}

	err = saver.posterSaver.SaveAll(newPosters, updatedPosters)
	if err != nil {
		return err
	}

	err = saver.latestEpisodesSaver.SaveAll(newLatestEpisodes, updatedLatestEpisodes)
	if err != nil {
		return err
	}

	err = saver.thumbnailSaver.SaveAll(newThumbnails, deletedThumbnails)
	if err != nil {
		return err
	}

	return nil
}
