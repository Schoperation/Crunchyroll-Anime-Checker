package command

import (
	"schoperation/crunchyroll-anime-checker/domain/anime"
	"schoperation/crunchyroll-anime-checker/domain/core"
)

type GenerateAnimeFilesCommandInput struct {
}

type GenerateAnimeFilesCommandOutput struct {
}

type allAnimeFetcher interface {
	GetAll() ([]anime.Anime, error)
}

type senseiListCreator interface {
	CreateSenseiList(animes []anime.Anime) error
}

type latestEpisodesFileCreator interface {
	CreateFileForLocale(locale core.Locale, latestEpisodes []anime.LatestEpisodes, slugTitles map[anime.AnimeId]string) error
}

type GenerateAnimeFilesCommand struct {
	allAnimeFetcher           allAnimeFetcher
	senseiListCreator         senseiListCreator
	latestEpisodesFileCreator latestEpisodesFileCreator
}

func NewGenerateAnimeFilesCommand(
	allAnimeFetcher allAnimeFetcher,
	senseiListCreator senseiListCreator,
	latestEpisodesFileCreator latestEpisodesFileCreator,
) GenerateAnimeFilesCommand {
	return GenerateAnimeFilesCommand{
		allAnimeFetcher:           allAnimeFetcher,
		senseiListCreator:         senseiListCreator,
		latestEpisodesFileCreator: latestEpisodesFileCreator,
	}
}

func (cmd GenerateAnimeFilesCommand) Run(input GenerateAnimeFilesCommandInput) (GenerateAnimeFilesCommandOutput, error) {
	locales := core.SupportedLocales()

	animes, err := cmd.allAnimeFetcher.GetAll()
	if err != nil {
		return GenerateAnimeFilesCommandOutput{}, err
	}

	err = cmd.senseiListCreator.CreateSenseiList(animes)
	if err != nil {
		return GenerateAnimeFilesCommandOutput{}, err
	}

	slugTitles := make(map[anime.AnimeId]string)

	for _, locale := range locales {
		var leSlice []anime.LatestEpisodes
		for _, localAnime := range animes {
			le, err := localAnime.Episodes().GetLatestEpisodesForLocale(locale)
			if err != nil {
				continue
			}

			leSlice = append(leSlice, le)

			if _, exists := slugTitles[localAnime.AnimeId()]; !exists {
				slugTitles[localAnime.AnimeId()] = localAnime.SlugTitle()
			}
		}

		err := cmd.latestEpisodesFileCreator.CreateFileForLocale(locale, leSlice, slugTitles)
		if err != nil {
			return GenerateAnimeFilesCommandOutput{}, err
		}
	}

	return GenerateAnimeFilesCommandOutput{}, nil
}
