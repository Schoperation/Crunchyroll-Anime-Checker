package anime

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/anime"
)

type thumbnailDao interface {
	GetAllByAnimeId(animeId int) ([]anime.ImageDto, error)
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

func (translator ThumbnailTranslator) SaveAll(newThumbnails []anime.Image) error {
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
