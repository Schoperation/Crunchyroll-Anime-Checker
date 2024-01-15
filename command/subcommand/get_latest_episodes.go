package subcommand

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"
	"schoperation/crunchyrollanimestatus/domain/crunchyroll"
)

type GetLatestEpisodesSubCommandInput struct {
	NewCrAnime     []crunchyroll.Anime
	UpdatedCrAnime []crunchyroll.Anime
	LocalAnime     map[core.SeriesId]anime.Anime
	Locales        []core.Locale
}

type GetLatestEpisodesSubCommandOutput struct {
}

type getAllSeasonsTranslator interface {
	GetAllSeasonsBySeriesId(seriesId string) (crunchyroll.SeasonCollection, error)
}

type getAllEpisodesTranslator interface {
	GetAllEpisodesBySeasonId(locale core.Locale, seasonId string) (crunchyroll.EpisodeCollection, error)
}

type getEncodedThumbnailTranslator interface {
	GetEncodedImageByURL(url string) (string, error)
}

type GetLatestEpisodesSubCommand struct {
	getAllSeasonsTranslator       getAllSeasonsTranslator
	getAllEpisodesTranslator      getAllEpisodesTranslator
	getEncodedThumbnailTranslator getEncodedThumbnailTranslator
}

func NewGetLatestEpisodesSubCommand(
	getAllSeasonsTranslator getAllSeasonsTranslator,
	getAllEpisodesTranslator getAllEpisodesTranslator,
	getEncodedThumbnailTranslator getEncodedThumbnailTranslator,
) GetLatestEpisodesSubCommand {
	return GetLatestEpisodesSubCommand{
		getAllSeasonsTranslator:       getAllSeasonsTranslator,
		getAllEpisodesTranslator:      getAllEpisodesTranslator,
		getEncodedThumbnailTranslator: getEncodedThumbnailTranslator,
	}
}

func (subcmd GetLatestEpisodesSubCommand) Run(input GetLatestEpisodesSubCommandInput) (GetLatestEpisodesSubCommandOutput, error) {
	for _, updatedCrAnime := range input.UpdatedCrAnime {
		localAnime, exists := input.LocalAnime[updatedCrAnime.SeriesId()]
		if !exists {
			return GetLatestEpisodesSubCommandOutput{}, fmt.Errorf("no local anime with series ID %s", updatedCrAnime.SeriesId())
		}

		crSeasons, err := subcmd.getAllSeasonsTranslator.GetAllSeasonsBySeriesId(updatedCrAnime.SeriesId().String())
		if err != nil {
			return GetLatestEpisodesSubCommandOutput{}, err
		}

		for _, locale := range input.Locales {
			savedLatestEpisodes := localAnime.Episodes().GetLatestEpisodesForLocale(locale)

			latestSubSeason, exists := crSeasons.LatestSub(locale)
			if exists {
				crEpisodes, err := subcmd.getAllEpisodesTranslator.GetAllEpisodesBySeasonId(locale, latestSubSeason.Id())
				if err != nil {
					return GetLatestEpisodesSubCommandOutput{}, err
				}

				latestCrSub, exists := crEpisodes.LatestSub(locale)
				if !exists {
					return GetLatestEpisodesSubCommandOutput{}, fmt.Errorf("could not find latest episode in season that should have sub")
				}

				latestSavedSub := savedLatestEpisodes.LatestSub()
				if !subcmd.areEpisodesTheSame(latestCrSub, latestSavedSub) {
					newSub, newSubThumbnail, err := subcmd.generateNewEpisode(localAnime.AnimeId(), latestCrSub)
					if err != nil {
						return GetLatestEpisodesSubCommandOutput{}, err
					}

					err = localAnime.Episodes().AddSubForLocale(locale, newSub, newSubThumbnail)
					if err != nil {
						return GetLatestEpisodesSubCommandOutput{}, err
					}
				}
			}

			latestDubSeason, exists := crSeasons.LatestDub(locale)
			if exists {
				crEpisodes, err := subcmd.getAllEpisodesTranslator.GetAllEpisodesBySeasonId(locale, latestDubSeason.Id())
				if err != nil {
					return GetLatestEpisodesSubCommandOutput{}, err
				}

				latestCrDub, exists := crEpisodes.LatestDub(locale)
				if !exists {
					return GetLatestEpisodesSubCommandOutput{}, fmt.Errorf("could not find latest episode in season that should have dub")
				}

				latestSavedDub := savedLatestEpisodes.LatestDub()
				if !subcmd.areEpisodesTheSame(latestCrDub, latestSavedDub) {
					newDub, newDubThumbnail, err := subcmd.generateNewEpisode(localAnime.AnimeId(), latestCrDub)
					if err != nil {
						return GetLatestEpisodesSubCommandOutput{}, err
					}

					err = localAnime.Episodes().AddDubForLocale(locale, newDub, newDubThumbnail)
					if err != nil {
						return GetLatestEpisodesSubCommandOutput{}, err
					}
				}
			}
		}

		localAnime.Episodes().CleanEpisodes()
		input.LocalAnime[updatedCrAnime.SeriesId()] = localAnime
	}

	for _, newCrAnime := range input.NewCrAnime {
		crSeasons, err := subcmd.getAllSeasonsTranslator.GetAllSeasonsBySeriesId(newCrAnime.SeriesId().String())
		if err != nil {
			return GetLatestEpisodesSubCommandOutput{}, err
		}
	}

	return GetLatestEpisodesSubCommandOutput{}, nil
}

func (subcmd GetLatestEpisodesSubCommand) areEpisodesTheSame(crEp crunchyroll.Episode, savedEp anime.MinimalEpisode) bool {
	return crEp.Number() == savedEp.Number() && crEp.Season() == savedEp.Season()
}

func (subcmd GetLatestEpisodesSubCommand) generateNewEpisode(animeId anime.AnimeId, crEpisode crunchyroll.Episode) (anime.MinimalEpisode, anime.Image, error) {
	newEpisode, err := anime.NewMinimalEpisode(crEpisode.Season(), crEpisode.Number(), crEpisode.Title())
	if err != nil {
		return anime.MinimalEpisode{}, anime.Image{}, err
	}

	encodedThumbnail, err := subcmd.getEncodedThumbnailTranslator.GetEncodedImageByURL(crEpisode.Thumbnail().Source())
	if err != nil {
		return anime.MinimalEpisode{}, anime.Image{}, err
	}

	newThumbnail, err := anime.NewImage(anime.ImageDto{
		AnimeId:       animeId.Int(),
		ImageType:     crEpisode.Thumbnail().ImageType().Int(),
		SeasonNumber:  newEpisode.Season(),
		EpisodeNumber: newEpisode.Number(),
		Url:           crEpisode.Thumbnail().Source(),
		Encoded:       encodedThumbnail,
	})
	if err != nil {
		return anime.MinimalEpisode{}, anime.Image{}, err
	}

	return newEpisode, newThumbnail, err
}
