package command

import (
	"fmt"
	"schoperation/crunchyroll-anime-checker/command/subcommand"
	"schoperation/crunchyroll-anime-checker/domain/anime"
	"schoperation/crunchyroll-anime-checker/domain/core"
	"schoperation/crunchyroll-anime-checker/domain/crunchyroll"
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
	GetAllByAnimeIds(animeIds []anime.AnimeId) (map[core.SeriesId]anime.Anime, map[core.SeriesId]anime.Anime, error)
}

type refreshBasicInfoSubCommand interface {
	Run(input subcommand.RefreshBasicInfoInput) (subcommand.RefreshBasicInfoOutput, map[core.SeriesId]error)
}

type refreshPostersSubCommand interface {
	Run(input subcommand.RefreshPostersSubCommandInput) (subcommand.RefreshPostersSubCommandOutput, map[core.SeriesId]error)
}

type refreshLatestEpisodesSubCommand interface {
	Run(input subcommand.RefreshLatestEpisodesSubCommandInput) (subcommand.RefreshLatestEpisodesSubCommandOutput, map[core.SeriesId]error)
}

type animeSaver interface {
	SaveAll(newAnimes []anime.Anime, updatedAnimes []anime.Anime, originalLocalAnimes map[core.SeriesId]anime.Anime) error
}

type RefreshAnimeCommand struct {
	crunchyrollAnimeFetcher         crunchyrollAnimeFetcher
	localAnimeFetcher               localAnimeFetcher
	refreshBasicInfoSubCommand      refreshBasicInfoSubCommand
	refreshPostersSubCommand        refreshPostersSubCommand
	refreshLatestEpisodesSubCommand refreshLatestEpisodesSubCommand
	animeSaver                      animeSaver
}

func NewRefreshAnimeCommand(
	crunchyrollAnimeFetcher crunchyrollAnimeFetcher,
	localAnimeFetcher localAnimeFetcher,
	refreshBasicInfoSubCommand refreshBasicInfoSubCommand,
	refreshPostersSubCommand refreshPostersSubCommand,
	refreshLatestEpisodesSubCommand refreshLatestEpisodesSubCommand,
	animeSaver animeSaver,
) RefreshAnimeCommand {
	return RefreshAnimeCommand{
		crunchyrollAnimeFetcher:         crunchyrollAnimeFetcher,
		localAnimeFetcher:               localAnimeFetcher,
		refreshBasicInfoSubCommand:      refreshBasicInfoSubCommand,
		refreshPostersSubCommand:        refreshPostersSubCommand,
		refreshLatestEpisodesSubCommand: refreshLatestEpisodesSubCommand,
		animeSaver:                      animeSaver,
	}
}

func (cmd RefreshAnimeCommand) Run(input RefreshAnimeCommandInput) (RefreshAnimeCommandOutput, error) {
	locales := core.SupportedLocales()

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

	fmt.Printf("%d new anime to add, %d anime to update.\n", len(newCrAnimes), len(updatedCrAnimes))

	localAnimeToBeUpdated, originalLocalAnimes, err := cmd.localAnimeFetcher.GetAllByAnimeIds(animeIds)
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}

	fmt.Printf("Finished anime retrieval in %v.\n\n", time.Since(startTime))

	fmt.Printf("Refreshing basic info...\n")
	startTime = time.Now().UTC()

	basicInfoCmdOutput, errs := cmd.refreshBasicInfoSubCommand.Run(subcommand.RefreshBasicInfoInput{
		NewCrAnime:     newCrAnimes,
		UpdatedCrAnime: updatedCrAnimes,
		LocalAnime:     localAnimeToBeUpdated,
	})
	if errs != nil {
		cmd.printErrors(errs)
	}
	fmt.Printf("Finished basic info refreshing in %v.\n\n", time.Since(startTime))

	fmt.Printf("Refreshing posters...\n")
	startTime = time.Now().UTC()

	posterCmdOutput, errs := cmd.refreshPostersSubCommand.Run(subcommand.RefreshPostersSubCommandInput{
		NewCrAnime:     newCrAnimes,
		UpdatedCrAnime: updatedCrAnimes,
		LocalAnime:     basicInfoCmdOutput.UpdatedLocalAnime,
	})
	if errs != nil {
		cmd.printErrors(errs)
	}
	fmt.Printf("Finished posters in %v.\n\n", time.Since(startTime))

	fmt.Printf("Refreshing latest episodes...\n")
	startTime = time.Now().UTC()

	latestEpisodeCmdOutput, errs := cmd.refreshLatestEpisodesSubCommand.Run(subcommand.RefreshLatestEpisodesSubCommandInput{
		NewCrAnime:     newCrAnimes,
		UpdatedCrAnime: updatedCrAnimes,
		LocalAnime:     posterCmdOutput.UpdatedLocalAnime,
		Locales:        locales,
	})
	if errs != nil {
		cmd.printErrors(errs)
	}
	fmt.Printf("Finished latest episodes in %v.\n\n", time.Since(startTime))

	fmt.Printf("Saving anime...\n")
	startTime = time.Now().UTC()

	var newLocalAnimes []anime.Anime
	for seriesId, dto := range basicInfoCmdOutput.NewAnimeDtos {
		posters, exists := posterCmdOutput.NewPosters[seriesId]
		if !exists {
			fmt.Printf("\tcould not find posters for anime %s, skipping\n", seriesId)
			continue
		}

		episodes, exists := latestEpisodeCmdOutput.NewEpisodeCollections[seriesId]
		if !exists {
			fmt.Printf("\tcould not find episodes for anime %s, skipping\n", seriesId)
			continue
		}

		newAnime, err := anime.NewAnime(dto, posters, episodes)
		if err != nil {
			fmt.Printf("\tcould not create anime %s, skipping (%v)\n", seriesId, err)
			continue
		}

		newLocalAnimes = append(newLocalAnimes, newAnime)
	}

	updatedLocalAnimes := make([]anime.Anime, len(latestEpisodeCmdOutput.UpdatedLocalAnime))
	j := 0
	for _, updatedLocalAnime := range latestEpisodeCmdOutput.UpdatedLocalAnime {
		updatedLocalAnimes[j] = updatedLocalAnime
		j++
	}

	err = cmd.animeSaver.SaveAll(newLocalAnimes, updatedLocalAnimes, originalLocalAnimes)
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}

	fmt.Printf("Finished saving anime in %v.\n\n", time.Since(startTime))

	return RefreshAnimeCommandOutput{
		NewAnimeCount:     len(newCrAnimes),
		UpdatedAnimeCount: len(updatedCrAnimes),
	}, nil
}

func (cmd RefreshAnimeCommand) printErrors(errs map[core.SeriesId]error) {
	for seriesId, err := range errs {
		fmt.Printf("\tError w/ series %s: %v\n", seriesId, err)
	}
}
