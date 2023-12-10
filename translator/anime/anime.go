package anime

import (
	"schoperation/crunchyrollanimestatus/domain/anime"
)

type animeDao interface {
	GetAllMinimal() ([]anime.MinimalAnimeDto, error)
	GetAllByAnimeIds(animeIds []int) ([]anime.AnimeDto, error)
	InsertAll(dtos []anime.AnimeDto) error
	Update(dto anime.AnimeDto) error
}

type animeFactory interface {
	Reform(dto anime.AnimeDto) (anime.Anime, error)
}

type AnimeTranslator struct {
	animeDao     animeDao
	animeFactory animeFactory
}

func NewAnimeTranslator(animeDao animeDao, animeFactory animeFactory) AnimeTranslator {
	return AnimeTranslator{
		animeDao:     animeDao,
		animeFactory: animeFactory,
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

func (translator AnimeTranslator) GetAllByAnimeIds(animeIds []anime.AnimeId) ([]anime.Anime, error) {
	ids := make([]int, len(animeIds))
	for i, animeId := range animeIds {
		ids[i] = animeId.Int()
	}

	dtos, err := translator.animeDao.GetAllByAnimeIds(ids)
	if err != nil {
		return nil, err
	}

	anime := make([]anime.Anime, len(dtos))
	for i, dto := range dtos {
		anime[i], err = translator.animeFactory.Reform(dto)
		if err != nil {
			return nil, err
		}
	}

	return anime, nil
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
