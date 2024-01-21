package anime

import (
	"schoperation/crunchyrollanimestatus/domain/anime"
)

type thumbnailDao interface {
	GetAllByAnimeIds(animeIds []int) ([]anime.ImageDto, error)
	InsertAll(dtos []anime.ImageDto) error
}

type ThumbnailTranslator struct {
	thumbnailDao thumbnailDao
}

func NewThumbnailTranslator(thumbnailDao thumbnailDao) ThumbnailTranslator {
	return ThumbnailTranslator{
		thumbnailDao: thumbnailDao,
	}
}

func (translator ThumbnailTranslator) GetAllByAnimeIds(animeIds []anime.AnimeId) (map[anime.AnimeId]map[string]anime.Image, error) {
	ids := make([]int, len(animeIds))
	for i, animeId := range animeIds {
		ids[i] = animeId.Int()
	}

	dtos, err := translator.thumbnailDao.GetAllByAnimeIds(ids)
	if err != nil {
		return nil, err
	}

	thumbnailMap := make(map[anime.AnimeId]map[string]anime.Image)
	for _, dto := range dtos {
		animeId := anime.ReformAnimeId(dto.AnimeId)
		newThumbnail := anime.ReformImage(dto)

		if thumbnails, exists := thumbnailMap[animeId]; exists {
			thumbnails[newThumbnail.Key()] = newThumbnail
			thumbnailMap[animeId] = thumbnails
			continue
		}

		thumbnailMap[animeId] = map[string]anime.Image{newThumbnail.Key(): newThumbnail}
	}

	return thumbnailMap, nil
}

func (translator ThumbnailTranslator) SaveAll(newThumbnails []anime.Image, updatedThumbnails []anime.Image) error {
	newDtos := make([]anime.ImageDto, len(newThumbnails))
	for i, thumbnail := range newThumbnails {
		newDtos[i] = thumbnail.Dto()
	}

	err := translator.thumbnailDao.InsertAll(newDtos)
	if err != nil {
		return err
	}

	return nil
}
