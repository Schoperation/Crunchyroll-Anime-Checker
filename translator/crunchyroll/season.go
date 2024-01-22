package crunchyroll

import (
	"schoperation/crunchyroll-anime-checker/domain/core"
	"schoperation/crunchyroll-anime-checker/domain/crunchyroll"
)

type crunchyrollAnimeSeasonClient interface {
	GetAllSeasonsBySeriesId(seriesId string) ([]crunchyroll.SeasonDto, error)
}

type SeasonTranslator struct {
	crunchyrollAnimeSeasonClient crunchyrollAnimeSeasonClient
}

func NewSeasonTranslator(crunchyrollAnimeSeasonClient crunchyrollAnimeSeasonClient) SeasonTranslator {
	return SeasonTranslator{
		crunchyrollAnimeSeasonClient: crunchyrollAnimeSeasonClient,
	}
}

func (translator SeasonTranslator) GetAllSeasonsBySeriesId(seriesId core.SeriesId) (crunchyroll.SeasonCollection, error) {
	dtos, err := translator.crunchyrollAnimeSeasonClient.GetAllSeasonsBySeriesId(seriesId.String())
	if err != nil {
		return crunchyroll.SeasonCollection{}, err
	}

	seasons := make([]crunchyroll.Season, len(dtos))
	for i, dto := range dtos {
		seasons[i] = crunchyroll.ReformSeason(dto)
	}

	return crunchyroll.NewSeasonCollection(seriesId, seasons)
}
