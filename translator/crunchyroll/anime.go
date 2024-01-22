package crunchyroll

import (
	"schoperation/crunchyroll-anime-checker/domain/core"
	"schoperation/crunchyroll-anime-checker/domain/crunchyroll"
)

type crunchyrollAnimeClient interface {
	GetAllAnime(locale string) ([]crunchyroll.AnimeDto, error)
}

type AnimeTranslator struct {
	crunchyrollAnimeClient crunchyrollAnimeClient
}

func NewAnimeTranslator(crunchyrollAnimeClient crunchyrollAnimeClient) AnimeTranslator {
	return AnimeTranslator{
		crunchyrollAnimeClient: crunchyrollAnimeClient,
	}
}

func (translator AnimeTranslator) GetAllAnime(locale core.Locale) ([]crunchyroll.Anime, error) {
	dtos, err := translator.crunchyrollAnimeClient.GetAllAnime(locale.Name())
	if err != nil {
		return nil, err
	}

	return crunchyroll.ReformAnimeCollection(dtos), nil
}
