package command

import (
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

type crunchyrollAnimeTranslator interface {
	GetAllAnime(locale core.Locale) ([]crunchyroll.Anime, error)
}

type localAnimeTranslator interface {
	GetAllMinimal() (map[string]anime.MinimalAnime, error)
	SaveAll(newAnime []anime.Anime, updatedAnime []anime.Anime) error
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
	for _, anime := range crAnime {
		if _, exists := localMinimalAnime[anime.SeriesId()]; !exists {
			newCrAnimes = append(newCrAnimes, anime)
			continue
		}

		savedAnime := localMinimalAnime[anime.SeriesId()]
		if savedAnime.LastUpdated().Before(anime.LastUpdated()) {
			updatedCrAnimes = append(updatedCrAnimes, anime)
			continue
		}
	}

	if len(newCrAnimes) == 0 && len(updatedCrAnimes) == 0 {
		return RefreshAnimeCommandOutput{}, nil
	}

	newAnime := make([]anime.Anime, len(newCrAnimes))
	for i, newCrAnime := range newCrAnimes {
		newAnime[i], err = anime.NewAnime(anime.AnimeDto{
			AnimeId:          0,
			SeriesId:         newCrAnime.SeriesId(),
			SlugTitle:        newCrAnime.SlugTitle(),
			Title:            newCrAnime.Title(),
			LastUpdated:      time.Time{},
			SeasonIdentifier: "S#",
		})
		if err != nil {
			return RefreshAnimeCommandOutput{}, err
		}
	}

	err = cmd.localAnimeTranslator.SaveAll(newAnime, nil)
	if err != nil {
		return RefreshAnimeCommandOutput{}, err
	}

	return RefreshAnimeCommandOutput{
		NewAnimeCount:     len(newCrAnimes),
		UpdatedAnimeCount: len(updatedCrAnimes),
	}, nil
}
