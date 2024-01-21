package anime

import (
	"schoperation/crunchyrollanimestatus/domain/anime"
)

type latestEpisodesDao interface {
	GetAllByAnimeIds(animeIds []int) ([]anime.LatestEpisodesDto, error)
	InsertAll(dtos []anime.LatestEpisodesDto) error
}

type LatestEpisodesTranslator struct {
	latestEpisodesDao latestEpisodesDao
}

func NewLatestEpisodesTranslator(latestEpisodesDao latestEpisodesDao) LatestEpisodesTranslator {
	return LatestEpisodesTranslator{
		latestEpisodesDao: latestEpisodesDao,
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

	return nil
}
