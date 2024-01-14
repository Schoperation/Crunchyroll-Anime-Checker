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

type GetLatestEpisodesSubCommand struct {
	getAllSeasonsTranslator  getAllSeasonsTranslator
	getAllEpisodesTranslator getAllEpisodesTranslator
}

func NewGetLatestEpisodesSubCommand(
	getAllSeasonsTranslator getAllSeasonsTranslator,
	getAllEpisodesTranslator getAllEpisodesTranslator,
) GetLatestEpisodesSubCommand {
	return GetLatestEpisodesSubCommand{
		getAllSeasonsTranslator:  getAllSeasonsTranslator,
		getAllEpisodesTranslator: getAllEpisodesTranslator,
	}
}

func (subcmd GetLatestEpisodesSubCommand) Run(input GetLatestEpisodesSubCommandInput) (GetLatestEpisodesSubCommandOutput, error) {
	for _, crAnime := range input.UpdatedCrAnime {
		localAnime, exists := input.LocalAnime[crAnime.SeriesId()]
		if !exists {
			return GetLatestEpisodesSubCommandOutput{}, fmt.Errorf("no local anime with series ID %s", crAnime.SeriesId())
		}

		crSeasons, err := subcmd.getAllSeasonsTranslator.GetAllSeasonsBySeriesId(crAnime.SeriesId().String())
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
				fmt.Print(latestCrSub, latestSavedSub)
			}

		}
	}

	return GetLatestEpisodesSubCommandOutput{}, nil
}
