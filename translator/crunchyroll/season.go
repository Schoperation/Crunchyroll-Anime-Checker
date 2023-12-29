package crunchyroll

import (
	"schoperation/crunchyrollanimestatus/domain/crunchyroll"
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

func (translator SeasonTranslator) GetAllSeasonsBySeriesId(seriesId string) ([]crunchyroll.Season, error) {
	dtos, err := translator.crunchyrollAnimeSeasonClient.GetAllSeasonsBySeriesId(seriesId)
	if err != nil {
		return nil, err
	}

	seasons := make([]crunchyroll.Season, len(dtos))
	for i, dto := range dtos {
		seasons[i] = crunchyroll.ReformSeason(dto)
	}

	return seasons, nil
}
