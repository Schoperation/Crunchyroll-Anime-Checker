package subcommand

import (
	"fmt"
	"schoperation/crunchyroll-anime-checker/domain/anime"
	"schoperation/crunchyroll-anime-checker/domain/core"
	"schoperation/crunchyroll-anime-checker/domain/crunchyroll"
)

type RefreshLatestEpisodesSubCommandInput struct {
	NewCrAnime     []crunchyroll.Anime
	UpdatedCrAnime []crunchyroll.Anime
	LocalAnime     map[core.SeriesId]anime.Anime
	Locales        []core.Locale
}

type RefreshLatestEpisodesSubCommandOutput struct {
	UpdatedLocalAnime     map[core.SeriesId]anime.Anime
	NewEpisodeCollections map[core.SeriesId]anime.EpisodeCollection
}

type seriesSeasonsFetcher interface {
	GetAllSeasonsBySeriesId(seriesId core.SeriesId) (crunchyroll.SeasonCollection, error)
}

type seasonEpisodesFetcher interface {
	GetAllEpisodesBySeasonId(locale core.Locale, seasonId string) (crunchyroll.EpisodeCollection, error)
}

type encodedThumbnailFetcher interface {
	GetEncodedImageByURL(url string) (string, error)
}

type RefreshLatestEpisodesSubCommand struct {
	seriesSeasonsFetcher    seriesSeasonsFetcher
	seasonEpisodesFetcher   seasonEpisodesFetcher
	encodedThumbnailFetcher encodedThumbnailFetcher
}

func NewRefreshLatestEpisodesSubCommand(
	seriesSeasonsFetcher seriesSeasonsFetcher,
	seasonEpisodesFetcher seasonEpisodesFetcher,
	encodedThumbnailFetcher encodedThumbnailFetcher,
) RefreshLatestEpisodesSubCommand {
	return RefreshLatestEpisodesSubCommand{
		seriesSeasonsFetcher:    seriesSeasonsFetcher,
		seasonEpisodesFetcher:   seasonEpisodesFetcher,
		encodedThumbnailFetcher: encodedThumbnailFetcher,
	}
}

func (subcmd RefreshLatestEpisodesSubCommand) Run(input RefreshLatestEpisodesSubCommandInput) (RefreshLatestEpisodesSubCommandOutput, map[core.SeriesId]error) {
	errors := map[core.SeriesId]error{}

	for _, updatedCrAnime := range input.UpdatedCrAnime {
		localAnime, exists := input.LocalAnime[updatedCrAnime.SeriesId()]
		if !exists {
			errors[updatedCrAnime.SeriesId()] = fmt.Errorf("no local anime found")
			continue
		}

		crSeasons, err := subcmd.seriesSeasonsFetcher.GetAllSeasonsBySeriesId(updatedCrAnime.SeriesId())
		if err != nil {
			errors[updatedCrAnime.SeriesId()] = err
			continue
		}

		for _, locale := range input.Locales {
			localLatestEpisodes, err := localAnime.Episodes().GetLatestEpisodesForLocale(locale)
			if err != nil {
				localLatestEpisodes = anime.ReformLatestEpisodes(anime.LatestEpisodesDto{})
			}

			latestSubSeason, subExists := crSeasons.LatestSub(locale)
			if subExists {
				crEpisodes, err := subcmd.seasonEpisodesFetcher.GetAllEpisodesBySeasonId(locale, latestSubSeason.Id())
				if err != nil {
					errors[updatedCrAnime.SeriesId()] = err
					break
				}

				latestCrSub, exists := crEpisodes.LatestSub(locale)
				if !exists {
					errors[updatedCrAnime.SeriesId()] = fmt.Errorf("could not find latest episode in season that should have sub")
					break
				}

				latestLocalSub := localLatestEpisodes.LatestSub()
				if !subcmd.areEpisodesTheSame(latestCrSub, latestLocalSub) {
					newSub, newSubThumbnail, err := subcmd.generateNewEpisode(localAnime.AnimeId(), latestCrSub)
					if err != nil {
						errors[updatedCrAnime.SeriesId()] = err
						break
					}

					err = localAnime.Episodes().AddSubForLocale(locale, newSub, newSubThumbnail)
					if err != nil {
						errors[updatedCrAnime.SeriesId()] = err
						break
					}
				}
			}

			latestDubSeason, dubExists := crSeasons.LatestDub(locale)
			if dubExists {
				crEpisodes, err := subcmd.seasonEpisodesFetcher.GetAllEpisodesBySeasonId(locale, latestDubSeason.Id())
				if err != nil {
					errors[updatedCrAnime.SeriesId()] = err
					break
				}

				latestCrDub, exists := crEpisodes.LatestDub(locale)
				if !exists {
					errors[updatedCrAnime.SeriesId()] = fmt.Errorf("could not find latest episode in season that should have dub")
					break
				}

				latestLocalDub := localLatestEpisodes.LatestDub()
				if !subcmd.areEpisodesTheSame(latestCrDub, latestLocalDub) {
					newDub, newDubThumbnail, err := subcmd.generateNewEpisode(localAnime.AnimeId(), latestCrDub)
					if err != nil {
						errors[updatedCrAnime.SeriesId()] = err
						break
					}

					err = localAnime.Episodes().AddDubForLocale(locale, newDub, newDubThumbnail)
					if err != nil {
						errors[updatedCrAnime.SeriesId()] = err
						break
					}
				}
			}
		}

		if _, errored := errors[updatedCrAnime.SeriesId()]; errored {
			continue
		}

		localAnime.SetDirty()
		input.LocalAnime[updatedCrAnime.SeriesId()] = localAnime
	}

	newEpisodeCollections := make(map[core.SeriesId]anime.EpisodeCollection, len(input.NewCrAnime))
	for _, newCrAnime := range input.NewCrAnime {
		crSeasons, err := subcmd.seriesSeasonsFetcher.GetAllSeasonsBySeriesId(newCrAnime.SeriesId())
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
				crEpisodes, err := subcmd.seasonEpisodesFetcher.GetAllEpisodesBySeasonId(locale, latestSubSeason.Id())
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
				crEpisodes, err := subcmd.seasonEpisodesFetcher.GetAllEpisodesBySeasonId(locale, latestDubSeason.Id())
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
			// TODO temp testing
			if !subExists && !dubExists {
				errors[newCrAnime.SeriesId()] = fmt.Errorf("could not find sub or dub")
			}

		}

		if _, errored := errors[newCrAnime.SeriesId()]; errored {
			continue
		}

		newEpisodeCollections[newCrAnime.SeriesId()] = newEpisodeCollection
	}

	return RefreshLatestEpisodesSubCommandOutput{
		UpdatedLocalAnime:     input.LocalAnime,
		NewEpisodeCollections: newEpisodeCollections,
	}, errors
}

func (subcmd RefreshLatestEpisodesSubCommand) areEpisodesTheSame(crEp crunchyroll.Episode, localEp anime.MinimalEpisode) bool {
	return crEp.Number() == localEp.Number() &&
		crEp.Season() == localEp.Season() &&
		crEp.Title() == localEp.Title()
}

func (subcmd RefreshLatestEpisodesSubCommand) generateNewEpisode(animeId anime.AnimeId, crEpisode crunchyroll.Episode) (anime.MinimalEpisode, anime.Image, error) {
	newEpisode, err := anime.NewMinimalEpisode(crEpisode.Season(), crEpisode.Number(), crEpisode.Title())
	if err != nil {
		return anime.MinimalEpisode{}, anime.Image{}, err
	}

	encodedThumbnail, err := subcmd.encodedThumbnailFetcher.GetEncodedImageByURL(crEpisode.Thumbnail().Source())
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
