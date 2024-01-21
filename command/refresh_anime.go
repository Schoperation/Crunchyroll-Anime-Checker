package command

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/command/subcommand"
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"
	"schoperation/crunchyrollanimestatus/domain/crunchyroll"
	"time"
)

type RefreshAnimeCommandInput struct {
}

type RefreshAnimeCommandOutput struct {
	NewAnimeCount     int
	UpdatedAnimeCount int
}

type crunchyrollAnimeFetcher interface {
	GetAllAnime(locale core.Locale) ([]crunchyroll.Anime, error)
}

type localAnimeFetcher interface {
	GetAllMinimal() (map[core.SeriesId]anime.MinimalAnime, error)
	GetAllByAnimeIds(animeIds []anime.AnimeId) (map[core.SeriesId]anime.Anime, error)
}

type refreshPostersSubCommand interface {
	Run(input subcommand.RefreshPostersSubCommandInput) (subcommand.RefreshPostersSubCommandOutput, error)
}

type getLatestEpisodesSubCommand interface {
	Run(input subcommand.GetLatestEpisodesSubCommandInput) (subcommand.GetLatestEpisodesSubCommandOutput, error)
}

type animeSaver interface {
	SaveAll(locales []core.Locale, newAnimes []anime.Anime, updatedAnimes []anime.Anime) error
}

type RefreshAnimeCommand struct {
	crunchyrollAnimeFetcher     crunchyrollAnimeFetcher
	localAnimeFetcher           localAnimeFetcher
	refreshPostersSubCommand    refreshPostersSubCommand
	getLatestEpisodesSubCommand getLatestEpisodesSubCommand
	animeSaver                  animeSaver
}

func NewRefreshAnimeCommand(
	crunchyrollAnimeFetcher crunchyrollAnimeFetcher,
	localAnimeFetcher localAnimeFetcher,
	refreshPostersSubCommand refreshPostersSubCommand,
	getLatestEpisodesSubCommand getLatestEpisodesSubCommand,
	animeSaver animeSaver,
) RefreshAnimeCommand {
	return RefreshAnimeCommand{
		crunchyrollAnimeFetcher:     crunchyrollAnimeFetcher,
		localAnimeFetcher:           localAnimeFetcher,
		refreshPostersSubCommand:    refreshPostersSubCommand,
		getLatestEpisodesSubCommand: getLatestEpisodesSubCommand,
		animeSaver:                  animeSaver,
	}
}

func (cmd RefreshAnimeCommand) Run(input RefreshAnimeCommandInput) (RefreshAnimeCommandOutput, error) {
	locales := []core.Locale{
		core.NewEnglishLocale(),
		core.NewSpanishLocale(),
	}

	fmt.Printf("Retrieving anime...\n")
	startTime := time.Now().UTC()
	crAnimes, err := cmd.crunchyrollAnimeFetcher.GetAllAnime(core.NewEnglishLocale())
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}

	localAnimes, err := cmd.localAnimeFetcher.GetAllMinimal()
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

	localAnimeToBeUpdated, err := cmd.localAnimeFetcher.GetAllByAnimeIds(animeIds)
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}
	fmt.Printf("Finished anime retrieval in %v\n", time.Since(startTime))

	fmt.Printf("Refreshing posters...\n")
	startTime = time.Now().UTC()
	posterCmdOutput, err := cmd.refreshPostersSubCommand.Run(subcommand.RefreshPostersSubCommandInput{
		NewCrAnime:     newCrAnimes,
		UpdatedCrAnime: updatedCrAnimes,
		LocalAnime:     localAnimeToBeUpdated,
	})
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}
	fmt.Printf("Finished posters in %v\n", time.Since(startTime))

	fmt.Printf("Refreshing latest episodes...\n")
	startTime = time.Now().UTC()
	latestEpisodeCmdOutput, err := cmd.getLatestEpisodesSubCommand.Run(subcommand.GetLatestEpisodesSubCommandInput{
		NewCrAnime:     newCrAnimes,
		UpdatedCrAnime: updatedCrAnimes,
		LocalAnime:     localAnimeToBeUpdated,
		Locales:        locales,
	})
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}
	fmt.Printf("Finished latest episodes in %v\n", time.Since(startTime))

	fmt.Printf("Saving anime...\n")
	startTime = time.Now().UTC()
	newLocalAnimes := make([]anime.Anime, 1)
	for i, newCrAnime := range newCrAnimes {

		// TODO temp testing
		if newCrAnime.SeriesId().String() != "G1XHJV0KV" {
			continue
		}

		posters := posterCmdOutput.NewPosters[newCrAnime.SeriesId()]
		episodes := latestEpisodeCmdOutput.NewEpisodeCollections[newCrAnime.SeriesId()]

		newAnime, err := anime.NewAnime(anime.AnimeDto{
			AnimeId:     0,
			SeriesId:    newCrAnime.SeriesId().String(),
			SlugTitle:   newCrAnime.SlugTitle(),
			Title:       newCrAnime.Title(),
			IsSimulcast: newCrAnime.IsSimulcast(),
			LastUpdated: newCrAnime.LastUpdated(),
		},
			posters,
			episodes)
		if err != nil {
			return RefreshAnimeCommandOutput{}, err
		}

		newLocalAnimes[i] = newAnime
	}

	err = cmd.animeSaver.SaveAll(locales, newLocalAnimes, nil)
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}

	fmt.Printf("Finished saving anime in %v\n", time.Since(startTime))

	return RefreshAnimeCommandOutput{
		NewAnimeCount:     len(newCrAnimes),
		UpdatedAnimeCount: len(updatedCrAnimes),
	}, nil
}
