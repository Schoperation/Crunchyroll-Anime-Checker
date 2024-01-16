package command

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/command/subcommand"
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"
	"schoperation/crunchyrollanimestatus/domain/crunchyroll"
)

type RefreshAnimeCommandInput struct {
}

type RefreshAnimeCommandOutput struct {
	NewAnimeCount     int
	UpdatedAnimeCount int
}

type crunchyrollAnimeTranslator interface {
	GetAllAnime(locale core.Locale) ([]crunchyroll.Anime, error)
}

type localAnimeTranslator interface {
	GetAllMinimal() (map[core.SeriesId]anime.MinimalAnime, error)
	GetAllByAnimeIds(animeIds []anime.AnimeId) (map[core.SeriesId]anime.Anime, error)
}

type refreshPostersSubCommand interface {
	Run(input subcommand.RefreshPostersSubCommandInput) (subcommand.RefreshPostersSubCommandOutput, error)
}

type getLatestEpisodesSubCommand interface {
	Run(input subcommand.GetLatestEpisodesSubCommandInput) (subcommand.GetLatestEpisodesSubCommandOutput, error)
}

type RefreshAnimeCommand struct {
	crunchyrollAnimeTranslator  crunchyrollAnimeTranslator
	localAnimeTranslator        localAnimeTranslator
	refreshPostersSubCommand    refreshPostersSubCommand
	getLatestEpisodesSubCommand getLatestEpisodesSubCommand
}

func NewRefreshAnimeCommand(
	crunchyrollAnimeTranslator crunchyrollAnimeTranslator,
	localAnimeTranslator localAnimeTranslator,
	refreshPostersSubCommand refreshPostersSubCommand,
	getLatestEpisodesSubCommand getLatestEpisodesSubCommand,
) RefreshAnimeCommand {
	return RefreshAnimeCommand{
		crunchyrollAnimeTranslator:  crunchyrollAnimeTranslator,
		localAnimeTranslator:        localAnimeTranslator,
		refreshPostersSubCommand:    refreshPostersSubCommand,
		getLatestEpisodesSubCommand: getLatestEpisodesSubCommand,
	}
}

func (cmd RefreshAnimeCommand) Run(input RefreshAnimeCommandInput) (RefreshAnimeCommandOutput, error) {
	crAnimes, err := cmd.crunchyrollAnimeTranslator.GetAllAnime(core.NewEnglishLocale())
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}

	localAnimes, err := cmd.localAnimeTranslator.GetAllMinimal()
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}

	var newCrAnimes []crunchyroll.Anime
	var updatedCrAnimes []crunchyroll.Anime
	var animeIds []anime.AnimeId
	for _, crAnime := range crAnimes {
		localAnime, exists := localAnimes[crAnime.SeriesId()]
		if !exists {
			newCrAnimes = append(newCrAnimes, crAnime)
			continue
		}

		if localAnime.LastUpdated().Before(crAnime.LastUpdated()) {
			updatedCrAnimes = append(updatedCrAnimes, crAnime)
			animeIds = append(animeIds, localAnime.AnimeId())
		}
	}

	if len(newCrAnimes) == 0 && len(updatedCrAnimes) == 0 {
		return RefreshAnimeCommandOutput{}, nil
	}

	fmt.Printf("%d new anime to add, %d anime to update...\n", len(newCrAnimes), len(updatedCrAnimes))

	localAnimeToBeUpdated, err := cmd.localAnimeTranslator.GetAllByAnimeIds(animeIds)
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}

	/*
		fmt.Printf("Refreshing posters...")
		posterCmdOutput, err := cmd.refreshPostersSubCommand.Run(subcommand.RefreshPostersSubCommandInput{
			NewCrAnime:     newCrAnimes,
			UpdatedCrAnime: updatedCrAnimes,
			LocalAnime:     localAnimeToBeUpdated,
		})
		if err != nil {
			return RefreshAnimeCommandOutput{}, err
		}
	*/

	fmt.Println("Refreshing latest episodes...")
	latestEpisodeCmdOutput, err := cmd.getLatestEpisodesSubCommand.Run(subcommand.GetLatestEpisodesSubCommandInput{
		NewCrAnime:     newCrAnimes,
		UpdatedCrAnime: updatedCrAnimes,
		LocalAnime:     localAnimeToBeUpdated,
		Locales:        []core.Locale{core.NewEnglishLocale()},
	})
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}

	fmt.Println(len(latestEpisodeCmdOutput.NewEpisodeCollections))

	return RefreshAnimeCommandOutput{
		NewAnimeCount:     len(newCrAnimes),
		UpdatedAnimeCount: len(updatedCrAnimes),
	}, nil
}
