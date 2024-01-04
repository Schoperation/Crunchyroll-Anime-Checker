package subcommand

import (
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"
	"schoperation/crunchyrollanimestatus/domain/crunchyroll"
)

type GetLatestEpisodesSubCommandInput struct {
	NewCrAnime     []crunchyroll.Anime
	UpdatedCrAnime []crunchyroll.Anime
	LocalAnime     map[core.SeriesId]anime.Anime
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
	return GetLatestEpisodesSubCommandOutput{}, nil
}
