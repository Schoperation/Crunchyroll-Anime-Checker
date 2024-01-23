package anime

import (
	"fmt"
	"schoperation/crunchyroll-anime-checker/domain/anime"
	"schoperation/crunchyroll-anime-checker/domain/core"
)

type latestEpisodesDao interface {
	GetAllByAnimeIds(animeIds []int) ([]anime.LatestEpisodesDto, error)
	InsertAll(dtos []anime.LatestEpisodesDto) error
	Update(dto anime.LatestEpisodesDto) error
}

type latestEpisodesFileWriter interface {
	WriteAllByLocale(localeName string, dtos []anime.LatestEpisodesDto) error
}

type LatestEpisodesTranslator struct {
	latestEpisodesDao        latestEpisodesDao
	latestEpisodesFileWriter latestEpisodesFileWriter
}

func NewLatestEpisodesTranslator(latestEpisodesDao latestEpisodesDao, latestEpisodesFileWriter latestEpisodesFileWriter) LatestEpisodesTranslator {
	return LatestEpisodesTranslator{
		latestEpisodesDao:        latestEpisodesDao,
		latestEpisodesFileWriter: latestEpisodesFileWriter,
	}
}

func (translator LatestEpisodesTranslator) GetAllByAnimeIds(animeIds []anime.AnimeId) (map[anime.AnimeId][]anime.LatestEpisodes, error) {
	ids := make([]int, len(animeIds))
	for i, animeId := range animeIds {
		ids[i] = animeId.Int()
	}

	dtos, err := translator.latestEpisodesDao.GetAllByAnimeIds(ids)
	if err != nil {
		return nil, err
	}

	latestEpisodesMap := make(map[anime.AnimeId][]anime.LatestEpisodes)
	for _, dto := range dtos {
		animeId := anime.ReformAnimeId(dto.AnimeId)

		if latestEpisodes, exists := latestEpisodesMap[animeId]; exists {
			latestEpisodesMap[animeId] = append(latestEpisodes, anime.ReformLatestEpisodes(dto))
			continue
		}

		latestEpisodesMap[animeId] = []anime.LatestEpisodes{anime.ReformLatestEpisodes(dto)}
	}

	return latestEpisodesMap, nil
}

func (translator LatestEpisodesTranslator) SaveAll(newLatestEpisodes []anime.LatestEpisodes, updatedLatestEpisodes []anime.LatestEpisodes) error {
	newDtos := make([]anime.LatestEpisodesDto, len(newLatestEpisodes))
	for i, le := range newLatestEpisodes {
		newDtos[i] = le.Dto()
	}

	err := translator.latestEpisodesDao.InsertAll(newDtos)
	if err != nil {
		return err
	}

	for _, le := range updatedLatestEpisodes {
		err := translator.latestEpisodesDao.Update(le.Dto())
		if err != nil {
			return err
		}
	}

	return nil
}

func (translator LatestEpisodesTranslator) CreateFileForLocale(locale core.Locale, latestEpisodes []anime.LatestEpisodes, slugTitles map[anime.AnimeId]string) error {
	dtos := make([]anime.LatestEpisodesDto, len(latestEpisodes))
	for i, le := range latestEpisodes {
		slugTitle, exists := slugTitles[le.AnimeId()]
		if !exists {
			return fmt.Errorf("missing slug title for anime ID %d", le.AnimeId())
		}

		dto := le.Dto()
		dto.SlugTitle = slugTitle
		dtos[i] = dto
	}

	err := translator.latestEpisodesFileWriter.WriteAllByLocale(locale.Name(), dtos)
	if err != nil {
		return err
	}

	return nil
}
