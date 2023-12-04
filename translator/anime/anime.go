package anime

import (
	"schoperation/crunchyrollanimestatus/domain/anime"
)

type animeDao interface {
	GetAllMinimal() ([]anime.MinimalAnimeDto, error)
	InsertAll(dtos []anime.AnimeDto) error
}

type AnimeTranslator struct {
	animeDao animeDao
}

func NewAnimeTranslator(animeDao animeDao) AnimeTranslator {
	return AnimeTranslator{
		animeDao: animeDao,
	}
}

func (translator AnimeTranslator) GetAllMinimal() (map[string]anime.MinimalAnime, error) {
	dtos, err := translator.animeDao.GetAllMinimal()
	if err != nil {
		return nil, err
	}

	minimalAnime := map[string]anime.MinimalAnime{}
	for _, dto := range dtos {
		minimalAnime[dto.SeriesId] = anime.ReformMinimalAnime(dto)
	}

	return minimalAnime, nil
}

func (translator AnimeTranslator) SaveAll(newAnime []anime.Anime, updatedAnime []anime.Anime) error {
	dtos := make([]anime.AnimeDto, len(newAnime))
	for i, series := range newAnime {
		dtos[i] = series.Dto()
	}

	err := translator.animeDao.InsertAll(dtos)
	if err != nil {
		return err
	}

	return nil
}
