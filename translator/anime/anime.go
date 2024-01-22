package anime

import (
	"schoperation/crunchyroll-anime-checker/domain/anime"
	"schoperation/crunchyroll-anime-checker/domain/core"
)

type animeDao interface {
	GetAllMinimal() ([]anime.MinimalAnimeDto, error)
	GetAllByAnimeIds(animeIds []int) ([]anime.AnimeDto, error)
	InsertAll(dtos []anime.AnimeDto) ([]anime.MinimalAnimeDto, error)
	Update(dto anime.AnimeDto) error
}

type animeFactory interface {
	ReformAll(dtos []anime.AnimeDto) (map[core.SeriesId]anime.Anime, map[core.SeriesId]anime.Anime, error)
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

	minimalAnimes := map[core.SeriesId]anime.MinimalAnime{}
	for _, dto := range dtos {
		seriesId := core.ReformSeriesId(dto.SeriesId)
		minimalAnimes[seriesId] = anime.ReformMinimalAnime(dto)
	}

	return minimalAnimes, nil
}

func (translator AnimeTranslator) GetAllByAnimeIds(animeIds []anime.AnimeId) (map[core.SeriesId]anime.Anime, map[core.SeriesId]anime.Anime, error) {
	ids := make([]int, len(animeIds))
	for i, animeId := range animeIds {
		ids[i] = animeId.Int()
	}

	dtos, err := translator.animeDao.GetAllByAnimeIds(ids)
	if err != nil {
		return nil, nil, err
	}

	return translator.animeFactory.ReformAll(dtos)
}

func (translator AnimeTranslator) SaveAll(newAnime []anime.Anime, updatedAnime []anime.Anime) (map[core.SeriesId]anime.MinimalAnime, error) {
	newDtos := make([]anime.AnimeDto, len(newAnime))
	for i, series := range newAnime {
		newDtos[i] = series.Dto()
	}

	minimalAnimeDtos, err := translator.animeDao.InsertAll(newDtos)
	if err != nil {
		return nil, err
	}

	for _, series := range updatedAnime {
		err := translator.animeDao.Update(series.Dto())
		if err != nil {
			return nil, err
		}
	}

	minimalAnimes := make(map[core.SeriesId]anime.MinimalAnime, len(minimalAnimeDtos))
	for _, dto := range minimalAnimeDtos {
		seriesId := core.ReformSeriesId(dto.SeriesId)
		minimalAnimes[seriesId] = anime.ReformMinimalAnime(dto)
	}

	return minimalAnimes, nil
}
