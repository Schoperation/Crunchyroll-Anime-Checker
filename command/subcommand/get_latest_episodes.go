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
		localeAnime, exists := input.LocalAnime[crAnime.SeriesId()]
		if !exists {
			return GetLatestEpisodesSubCommandOutput{}, fmt.Errorf("no local anime with series ID %s", crAnime.SeriesId())
		}

		for _, locale := range input.Locales {
			_ = localeAnime.Episodes().GetLatestEpisodesForLocale(locale)
		}
	}

	return GetLatestEpisodesSubCommandOutput{}, nil
}
