package anime

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/anime"
)

type thumbnailDao interface {
	GetAllByAnimeId(animeId int) ([]anime.ImageDto, error)
}

type ThumbnailTranslator struct {
	thumbnailDao thumbnailDao
}

func NewThumbnailTranslator(thumbnailDao thumbnailDao) ThumbnailTranslator {
	return ThumbnailTranslator{
		thumbnailDao: thumbnailDao,
	}
}

func (translator ThumbnailTranslator) GetAllByAnimeId(animeId anime.AnimeId) (map[string]anime.Image, error) {
	dtos, err := translator.thumbnailDao.GetAllByAnimeId(animeId.Int())
	if err != nil {
		return nil, err
	}

	images := make(map[string]anime.Image, len(dtos))
	for _, dto := range dtos {
		images[fmt.Sprintf("%d-%d", dto.SeasonNumber, dto.EpisodeNumber)] = anime.ReformImage(dto)

	}

	return images, nil
}
