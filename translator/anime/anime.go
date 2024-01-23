package anime

import (
	"schoperation/crunchyroll-anime-checker/domain/anime"
	"schoperation/crunchyroll-anime-checker/domain/core"
	"slices"
)

type animeDao interface {
	GetAllMinimal() ([]anime.MinimalAnimeDto, error)
	GetAll() ([]anime.AnimeDto, error)
	GetAllByAnimeIds(animeIds []int) ([]anime.AnimeDto, error)
	InsertAll(dtos []anime.AnimeDto) ([]anime.MinimalAnimeDto, error)
	Update(dto anime.AnimeDto) error
}

type senseiListWriter interface {
	WriteAll(seriesIds, slugTitles, titles []string) error
}

type animeFactory interface {
	ReformAll(dtos []anime.AnimeDto) (map[core.SeriesId]anime.Anime, map[core.SeriesId]anime.Anime, error)
}

type AnimeTranslator struct {
	animeDao         animeDao
	senseiListWriter senseiListWriter
	animeFactory     animeFactory
}

func NewAnimeTranslator(
	animeDao animeDao,
	senseiListWriter senseiListWriter,
	animeFactory animeFactory,
) AnimeTranslator {
	return AnimeTranslator{
		animeDao:         animeDao,
		senseiListWriter: senseiListWriter,
		animeFactory:     animeFactory,
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

func (translator AnimeTranslator) GetAll() ([]anime.Anime, error) {
	dtos, err := translator.animeDao.GetAll()
	if err != nil {
		return nil, err
	}

	animes, _, err := translator.animeFactory.ReformAll(dtos)
	if err != nil {
		return nil, err
	}

	animeSlice := make([]anime.Anime, len(animes))
	i := 0
	for _, localAnime := range animes {
		animeSlice[i] = localAnime
		i++
	}

	slices.SortFunc(animeSlice, func(a, b anime.Anime) int {
		if a.SlugTitle() < b.SlugTitle() {
			return -1
		}

		if a.SlugTitle() > b.SlugTitle() {
			return 1
		}

		return 0
	})

	return animeSlice, nil
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

func (translator AnimeTranslator) CreateSenseiList(animes []anime.Anime) error {
	seriesIds := make([]string, len(animes))
	slugTitles := make([]string, len(animes))
	titles := make([]string, len(animes))

	for i, localAnime := range animes {
		seriesIds[i] = localAnime.SeriesId().String()
		slugTitles[i] = localAnime.SlugTitle()
		titles[i] = localAnime.Title()
	}

	err := translator.senseiListWriter.WriteAll(seriesIds, slugTitles, titles)
	if err != nil {
		return err
	}

	return nil
}
