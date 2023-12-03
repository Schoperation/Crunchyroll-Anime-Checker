package crunchyroll

import "schoperation/crunchyrollanimestatus/domain/crunchyroll"

type getAllAnimeClient interface {
	GetAllAnime(locale string) ([]crunchyroll.AnimeDto, error)
}

type AnimeTranslator struct {
	getAllAnimeClient getAllAnimeClient
}

func NewAnimeTranslator(getAllAnimeClient getAllAnimeClient) AnimeTranslator {
	return AnimeTranslator{
		getAllAnimeClient: getAllAnimeClient,
	}
}

func (translator AnimeTranslator) GetAllAnime(locale crunchyroll.Locale) ([]crunchyroll.Anime, error) {
	dtos, err := translator.getAllAnimeClient.GetAllAnime(locale.Name())
	if err != nil {
		return nil, err
	}

	return crunchyroll.ReformAnimeCollection(dtos), nil
}
