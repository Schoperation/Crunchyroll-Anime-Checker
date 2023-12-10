package anime

import (
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"
)

type latestEpisodesDao interface {
	GetAllByAnimeId(animeId int) ([]anime.LatestEpisodesDto, error)
}

type LatestEpisodesTranslator struct {
	latestEpisodesDao latestEpisodesDao
}

func NewLatestEpisodesTranslator(latestEpisodesDao latestEpisodesDao) LatestEpisodesTranslator {
	return LatestEpisodesTranslator{
		latestEpisodesDao: latestEpisodesDao,
	}
}

func (translator LatestEpisodesTranslator) GetAllByAnimeId(animeId anime.AnimeId) (map[core.Locale]anime.LatestEpisodes, error) {
	dtos, err := translator.latestEpisodesDao.GetAllByAnimeId(animeId.Int())
	if err != nil {
		return nil, err
	}

	latestEpisodesMap := make(map[core.Locale]anime.LatestEpisodes, len(dtos))
	for _, dto := range dtos {
		latestEpisodes := anime.ReformLatestEpisodes(dto)
		latestEpisodesMap[latestEpisodes.Locale()] = latestEpisodes
	}

	return latestEpisodesMap, nil
}
