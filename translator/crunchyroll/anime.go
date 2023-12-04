package crunchyroll

import (
	"schoperation/crunchyrollanimestatus/domain/core"
	"schoperation/crunchyrollanimestatus/domain/crunchyroll"
)

type crunchyrollClient interface {
	GetAllAnime(locale string) ([]crunchyroll.AnimeDto, error)
}

type AnimeTranslator struct {
	crunchyrollClient crunchyrollClient
}

func NewAnimeTranslator(crunchyrollClient crunchyrollClient) AnimeTranslator {
	return AnimeTranslator{
		crunchyrollClient: crunchyrollClient,
	}
}

func (translator AnimeTranslator) GetAllAnime(locale core.Locale) ([]crunchyroll.Anime, error) {
	dtos, err := translator.crunchyrollClient.GetAllAnime(locale.Name())
	if err != nil {
		return nil, err
	}

	return crunchyroll.ReformAnimeCollection(dtos), nil
}
