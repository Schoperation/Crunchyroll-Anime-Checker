package subcommand

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"
	"schoperation/crunchyrollanimestatus/domain/crunchyroll"
)

type RefreshBasicInfoInput struct {
	NewCrAnime     []crunchyroll.Anime
	UpdatedCrAnime []crunchyroll.Anime
	LocalAnime     map[core.SeriesId]anime.Anime
}

type RefreshBasicInfoOutput struct {
	UpdatedLocalAnime map[core.SeriesId]anime.Anime
	NewAnimeDtos      map[core.SeriesId]anime.AnimeDto
}

type RefreshBasicInfoSubCommand struct{}

func NewRefreshBasicInfoSubCommand() RefreshBasicInfoSubCommand {
	return RefreshBasicInfoSubCommand{}
}

func (subcmd RefreshBasicInfoSubCommand) Run(input RefreshBasicInfoInput) (RefreshBasicInfoOutput, map[core.SeriesId]error) {
	errors := map[core.SeriesId]error{}

	for _, updatedCrAnime := range input.UpdatedCrAnime {
		fmt.Printf("\tUpdating %s - %s\n", updatedCrAnime.SeriesId(), updatedCrAnime.SlugTitle())

		localAnime, exists := input.LocalAnime[updatedCrAnime.SeriesId()]
		if !exists {
			errors[updatedCrAnime.SeriesId()] = fmt.Errorf("no local anime found")
			continue
		}

		err := localAnime.UpdateBasicInfo(anime.AnimeDto{
			AnimeId:     localAnime.AnimeId().Int(),
			SeriesId:    updatedCrAnime.SeriesId().String(),
			SlugTitle:   updatedCrAnime.SlugTitle(),
			Title:       updatedCrAnime.Title(),
			IsSimulcast: updatedCrAnime.IsSimulcast(),
		})
		if err != nil {
			errors[updatedCrAnime.SeriesId()] = err
			continue
		}

		input.LocalAnime[updatedCrAnime.SeriesId()] = localAnime
	}

	newAnimeDtos := make(map[core.SeriesId]anime.AnimeDto, len(input.NewCrAnime))
	for _, newCrAnime := range input.NewCrAnime {
		fmt.Printf("\tCreating %s - %s\n", newCrAnime.SeriesId(), newCrAnime.SlugTitle())

		newAnimeDtos[newCrAnime.SeriesId()] = anime.AnimeDto{
			AnimeId:     0,
			SeriesId:    newCrAnime.SeriesId().String(),
			SlugTitle:   newCrAnime.SlugTitle(),
			Title:       newCrAnime.Title(),
			IsSimulcast: newCrAnime.IsSimulcast(),
		}
	}

	return RefreshBasicInfoOutput{
		UpdatedLocalAnime: input.LocalAnime,
		NewAnimeDtos:      newAnimeDtos,
	}, errors
}
