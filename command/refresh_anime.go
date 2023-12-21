package command

import (
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

type RefreshAnimeCommand struct {
	crunchyrollAnimeTranslator crunchyrollAnimeTranslator
	localAnimeTranslator       localAnimeTranslator
}

func NewRefreshAnimeCommand(
	crunchyrollAnimeTranslator crunchyrollAnimeTranslator,
	localAnimeTranslator localAnimeTranslator,
) RefreshAnimeCommand {
	return RefreshAnimeCommand{
		crunchyrollAnimeTranslator: crunchyrollAnimeTranslator,
		localAnimeTranslator:       localAnimeTranslator,
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
	var updatedCrAnimes []crunchyroll.Anime
	var animeIds []anime.AnimeId
	for _, anime := range crAnime {
		savedAnime, exists := localMinimalAnime[anime.SeriesId()]
		if !exists {
			newCrAnimes = append(newCrAnimes, anime)
			continue
		}

		if savedAnime.LastUpdated().Before(anime.LastUpdated()) {
			updatedCrAnimes = append(updatedCrAnimes, anime)
			animeIds = append(animeIds, savedAnime.AnimeId())
		}
	}

	if len(newCrAnimes) == 0 && len(updatedCrAnimes) == 0 {
		return RefreshAnimeCommandOutput{}, nil
	}

	animeToBeUpdated, err := 

	return RefreshAnimeCommandOutput{
		NewAnimeCount:     len(newCrAnimes),
		UpdatedAnimeCount: len(updatedCrAnimes),
	}, nil
}
