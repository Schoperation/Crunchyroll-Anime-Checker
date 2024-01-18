package subcommand

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"
	"schoperation/crunchyrollanimestatus/domain/crunchyroll"
	"slices"
)

type GetLatestEpisodesSubCommandInput struct {
	NewCrAnime     []crunchyroll.Anime
	UpdatedCrAnime []crunchyroll.Anime
	LocalAnime     map[core.SeriesId]anime.Anime
	Locales        []core.Locale
}

type GetLatestEpisodesSubCommandOutput struct {
	UpdatedLocalAnime     map[core.SeriesId]anime.Anime
	NewEpisodeCollections map[core.SeriesId]anime.EpisodeCollection
	Errors                map[core.SeriesId]error
}

type getAllSeasonsTranslator interface {
	GetAllSeasonsBySeriesId(seriesId core.SeriesId) (crunchyroll.SeasonCollection, error)
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
	errors := map[core.SeriesId]error{}
	for _, updatedCrAnime := range input.UpdatedCrAnime {
		localAnime, exists := input.LocalAnime[updatedCrAnime.SeriesId()]
		if !exists {
			return GetLatestEpisodesSubCommandOutput{}, fmt.Errorf("no local anime with series ID %s", updatedCrAnime.SeriesId())
		}

		crSeasons, err := subcmd.getAllSeasonsTranslator.GetAllSeasonsBySeriesId(updatedCrAnime.SeriesId())
		if err != nil {
			return GetLatestEpisodesSubCommandOutput{}, err
		}

		for _, locale := range input.Locales {
			localLatestEpisodes := localAnime.Episodes().GetLatestEpisodesForLocale(locale)

			latestSubSeason, subExists := crSeasons.LatestSub(locale)
			if subExists {
				crEpisodes, err := subcmd.getAllEpisodesTranslator.GetAllEpisodesBySeasonId(locale, latestSubSeason.Id())
				if err != nil {
					return GetLatestEpisodesSubCommandOutput{}, err
				}

				latestCrSub, exists := crEpisodes.LatestSub(locale)
				if !exists {
					return GetLatestEpisodesSubCommandOutput{}, fmt.Errorf("could not find latest episode in season that should have sub")
				}

				latestLocalSub := localLatestEpisodes.LatestSub()
				if !subcmd.areEpisodesTheSame(latestCrSub, latestLocalSub) {
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

			latestDubSeason, dubExists := crSeasons.LatestDub(locale)
			if dubExists {
				crEpisodes, err := subcmd.getAllEpisodesTranslator.GetAllEpisodesBySeasonId(locale, latestDubSeason.Id())
				if err != nil {
					return GetLatestEpisodesSubCommandOutput{}, err
				}

				latestCrDub, exists := crEpisodes.LatestDub(locale)
				if !exists {
					return GetLatestEpisodesSubCommandOutput{}, fmt.Errorf("could not find latest episode in season that should have dub")
				}

				latestLocalDub := localLatestEpisodes.LatestDub()
				if !subcmd.areEpisodesTheSame(latestCrDub, latestLocalDub) {
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

	newEpisodeCollections := make(map[core.SeriesId]anime.EpisodeCollection, len(input.NewCrAnime))
	for _, newCrAnime := range input.NewCrAnime {

		if !subcmd.validSeriesId(newCrAnime.SeriesId()) {
			continue
		}

		fmt.Printf("%s - %s\n", newCrAnime.SeriesId(), newCrAnime.SlugTitle())

		crSeasons, err := subcmd.getAllSeasonsTranslator.GetAllSeasonsBySeriesId(newCrAnime.SeriesId())
		if err != nil {
			errors[newCrAnime.SeriesId()] = err
			continue
		}

		newEpisodeCollection, err := anime.NewEpisodeCollection(anime.NewBlankAnimeId(), nil, nil)
		if err != nil {
			errors[newCrAnime.SeriesId()] = err
			continue
		}

		for _, locale := range input.Locales {
			latestSubSeason, subExists := crSeasons.LatestSub(locale)
			if subExists {
				crEpisodes, err := subcmd.getAllEpisodesTranslator.GetAllEpisodesBySeasonId(locale, latestSubSeason.Id())
				if err != nil {
					errors[newCrAnime.SeriesId()] = err
					break
				}

				latestCrSub, exists := crEpisodes.LatestSub(locale)
				if !exists {
					errors[newCrAnime.SeriesId()] = fmt.Errorf("could not find latest episode in season that should have sub")
					break
				}

				newSub, newSubThumbnail, err := subcmd.generateNewEpisode(anime.NewBlankAnimeId(), latestCrSub)
				if err != nil {
					errors[newCrAnime.SeriesId()] = err
					break
				}

				err = newEpisodeCollection.AddSubForLocale(locale, newSub, newSubThumbnail)
				if err != nil {
					errors[newCrAnime.SeriesId()] = err
					break
				}
			}

			latestDubSeason, dubExists := crSeasons.LatestDub(locale)
			if dubExists {
				crEpisodes, err := subcmd.getAllEpisodesTranslator.GetAllEpisodesBySeasonId(locale, latestDubSeason.Id())
				if err != nil {
					errors[newCrAnime.SeriesId()] = err
					break
				}

				latestCrDub, exists := crEpisodes.LatestDub(locale)
				if !exists {
					errors[newCrAnime.SeriesId()] = fmt.Errorf("could not find latest episode in season that should have dub")
					break
				}

				newDub, newDubThumbnail, err := subcmd.generateNewEpisode(anime.NewBlankAnimeId(), latestCrDub)
				if err != nil {
					errors[newCrAnime.SeriesId()] = err
					break
				}

				err = newEpisodeCollection.AddDubForLocale(locale, newDub, newDubThumbnail)
				if err != nil {
					errors[newCrAnime.SeriesId()] = err
					break
				}
			}

			// Temp code to debug blank english anime because CR's API is an inconsistent POS
			// TODO remove
			if !subExists && !dubExists {
				errors[newCrAnime.SeriesId()] = fmt.Errorf("could not find sub or dub")
			}

		}

		newEpisodeCollections[newCrAnime.SeriesId()] = newEpisodeCollection
	}

	return GetLatestEpisodesSubCommandOutput{
		UpdatedLocalAnime:     input.LocalAnime,
		NewEpisodeCollections: newEpisodeCollections,
		Errors:                errors,
	}, nil
}

func (subcmd GetLatestEpisodesSubCommand) validSeriesId(seriesId core.SeriesId) bool {
	validSeriesIds := []string{
		"GYEX5E1G6",
		"G62487M9Y",
		"GY2420WPR",
		"GRZX4X5MY",
		"G5PHNMWWN",
		"GDKHZE1GP",
		"G5PHNM77K",
		"GDKHZEPP3",
		"GDKHZE2Q7",
		"GRG5HJN5W",
		"G5PHNM7M1",
		"G79H23XN2",
		"GYGG9Z2WY",
		"GW4HM7GNK",
		"G9VHN91G5",
		"GY79XX9MR",
		"GY2P5MW0Y",
		"G24H1NZXD",
		"GR3K50PZR",
		"GRWEMGNER",
		"GXJHM3PVQ",
		"GYDKQX9Z6",
		"G24H1NW17",
		"GRP85E0MR",
		"G6EXH7VKM",
		"GR79V3VN6",
		"GVDHX8V05",
		"GY190V9ER",
		"GR2420096",
		"GXJHM3NEV",
		"G4PH0WEGW",
		"GYEXXWD36",
		"G1XHJV2PZ",
		"G67570P3R",
		"G6WE4W0N6",
		"G79H23WQ7",
		"GR8DN7N7R",
		"G79H23XVZ",
		"G0XHWM00Z",
		"GYGVNVEDY",
		"GRGGP17PR",
		"GR09QX10R",
		"GRPP77VJR",
		"GR24JZ886",
		"GR49G9VP6",
		"GRVND1G3Y",
		"G8DHV7478",
		"GXJHM397N",
		"G79H23ZD5",
	}

	return slices.Contains(validSeriesIds, seriesId.String())
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
