package anime

import (
	"schoperation/crunchyrollanimestatus/domain/anime"
)

type latestEpisodesDao interface {
	GetAllByAnimeId(animeId int) ([]anime.LatestEpisodesDto, error)
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

func (translator LatestEpisodesTranslator) GetAllByAnimeId(animeId anime.AnimeId) ([]anime.LatestEpisodes, error) {
	dtos, err := translator.latestEpisodesDao.GetAllByAnimeId(animeId.Int())
	if err != nil {
		return nil, err
	}

	latestEpisodes := make([]anime.LatestEpisodes, len(dtos))
	for i, dto := range dtos {
		latestEpisodes[i] = anime.ReformLatestEpisodes(dto)
	}

	return latestEpisodes, nil
}

func (translator LatestEpisodesTranslator) SaveAll(newLatestEpisodes []anime.LatestEpisodes) error {
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
