package anime

import (
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"
)

type animeDao interface {
	GetAllMinimal() ([]anime.MinimalAnimeDto, error)
	GetAllByAnimeIds(animeIds []int) ([]anime.AnimeDto, error)
	InsertAll(dtos []anime.AnimeDto) error
	Update(dto anime.AnimeDto) error
}

type animeFactory interface {
	Reform(dto anime.AnimeDto) (anime.Anime, error)
	ReformAll(dtos []anime.AnimeDto) (map[core.SeriesId]anime.Anime, error)
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

func (translator AnimeTranslator) GetAllMinimal() (map[core.SeriesId]anime.MinimalAnime, error) {
	dtos, err := translator.animeDao.GetAllMinimal()
	if err != nil {
		return nil, err
	}

	minimalAnime := map[core.SeriesId]anime.MinimalAnime{}
	for _, dto := range dtos {
		reformedMinimalAnime := anime.ReformMinimalAnime(dto)
		minimalAnime[reformedMinimalAnime.SeriesId()] = reformedMinimalAnime
	}

	return minimalAnime, nil
}

func (translator AnimeTranslator) GetAllByAnimeIds(animeIds []anime.AnimeId) (map[core.SeriesId]anime.Anime, error) {
	ids := make([]int, len(animeIds))
	for i, animeId := range animeIds {
		ids[i] = animeId.Int()
	}

	dtos, err := translator.animeDao.GetAllByAnimeIds(ids)
	if err != nil {
		return nil, err
	}

	return translator.animeFactory.ReformAll(dtos)
}

func (translator AnimeTranslator) SaveAll(newAnime []anime.Anime, updatedAnime []anime.Anime) error {
	newDtos := make([]anime.AnimeDto, len(newAnime))
	for i, series := range newAnime {
		newDtos[i] = series.Dto()
	}

	err := translator.animeDao.InsertAll(newDtos)
	if err != nil {
		return err
	}

	return nil
}
