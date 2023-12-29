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
	GetAllMinimal() (map[string]anime.MinimalAnime, error)
	GetAllByAnimeIds(animeIds []anime.AnimeId) ([]anime.Anime, error)
}

type refreshPostersSubCommand interface {
	Run(input subcommand.RefreshPostersSubCommandInput) (subcommand.RefreshPostersSubCommandOutput, error)
}

type RefreshAnimeCommand struct {
	crunchyrollAnimeTranslator crunchyrollAnimeTranslator
	localAnimeTranslator       localAnimeTranslator
	refreshPostersSubCommand   refreshPostersSubCommand
}

func NewRefreshAnimeCommand(
	crunchyrollAnimeTranslator crunchyrollAnimeTranslator,
	localAnimeTranslator localAnimeTranslator,
	refreshPostersSubCommand refreshPostersSubCommand,
) RefreshAnimeCommand {
	return RefreshAnimeCommand{
		crunchyrollAnimeTranslator: crunchyrollAnimeTranslator,
		localAnimeTranslator:       localAnimeTranslator,
		refreshPostersSubCommand:   refreshPostersSubCommand,
	}
}

func (cmd RefreshAnimeCommand) Run(input RefreshAnimeCommandInput) (RefreshAnimeCommandOutput, error) {
	crAnime, err := cmd.crunchyrollAnimeTranslator.GetAllAnime(core.NewEnglishLocale())
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}

	localMinimalAnime, err := cmd.localAnimeTranslator.GetAllMinimal()
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}

	var newCrAnimes []crunchyroll.Anime
	updatedCrAnimes := make(map[string]crunchyroll.Anime)
	var animeIds []anime.AnimeId
	for _, anime := range crAnime {
		savedAnime, exists := localMinimalAnime[anime.SeriesId()]
		if !exists {
			newCrAnimes = append(newCrAnimes, anime)
			continue
		}

		if savedAnime.LastUpdated().Before(anime.LastUpdated()) {
			updatedCrAnimes[anime.SeriesId()] = anime
			animeIds = append(animeIds, savedAnime.AnimeId())
		}
	}

	if len(newCrAnimes) == 0 && len(updatedCrAnimes) == 0 {
		return RefreshAnimeCommandOutput{}, nil
	}

	animeToBeUpdated, err := cmd.localAnimeTranslator.GetAllByAnimeIds(animeIds)
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}

	posterCmdOutput, err := cmd.refreshPostersSubCommand.Run(subcommand.RefreshPostersSubCommandInput{
		NewCrAnime:     newCrAnimes,
		UpdatedCrAnime: updatedCrAnimes,
		SavedAnime:     animeToBeUpdated,
	})
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}

	fmt.Println(len(posterCmdOutput.NewPosters))

	return RefreshAnimeCommandOutput{
		NewAnimeCount:     len(newCrAnimes),
		UpdatedAnimeCount: len(updatedCrAnimes),
	}, nil
}
