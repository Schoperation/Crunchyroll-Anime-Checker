package anime

import (
	"fmt"
	"schoperation/crunchyroll-anime-checker/domain/anime"
)

type thumbnailDao interface {
	GetAllByAnimeIds(animeIds []int) ([]anime.ImageDto, error)
	InsertAll(dtos []anime.ImageDto) error
	Delete(animeId, seasonNumber, episodeNumber int) error
}

type thumbnailFileWriter interface {
	WriteAll(dtos []anime.ImageDto) error
}

type ThumbnailTranslator struct {
	thumbnailDao        thumbnailDao
	thumbnailFileWriter thumbnailFileWriter
}

func NewThumbnailTranslator(thumbnailDao thumbnailDao, thumbnailFileWriter thumbnailFileWriter) ThumbnailTranslator {
	return ThumbnailTranslator{
		thumbnailDao:        thumbnailDao,
		thumbnailFileWriter: thumbnailFileWriter,
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

func (translator ThumbnailTranslator) SaveAll(newThumbnails []anime.Image, deletedThumbnails []anime.Image) error {
	newDtos := make([]anime.ImageDto, len(newThumbnails))
	for i, thumbnail := range newThumbnails {
		newDtos[i] = thumbnail.Dto()
	}

	err := translator.thumbnailDao.InsertAll(newDtos)
	if err != nil {
		return err
	}

	for _, thumbnail := range deletedThumbnails {
		err := translator.thumbnailDao.Delete(thumbnail.AnimeId().Int(), thumbnail.SeasonNumber(), thumbnail.EpisodeNumber())
		if err != nil {
			return err
		}
	}

	return nil
}

func (translator ThumbnailTranslator) CreateThumbnailFiles(thumbnails []anime.Image, slugTitles map[anime.AnimeId]string) error {
	dtos := make([]anime.ImageDto, len(thumbnails))
	for i, thumbnail := range thumbnails {
		slugTitle, exists := slugTitles[thumbnail.AnimeId()]
		if !exists {
			return fmt.Errorf("missing slug title for anime ID %d", thumbnail.AnimeId())
		}

		dto := thumbnail.Dto()
		dto.SlugTitle = slugTitle
		dtos[i] = dto
	}

	err := translator.thumbnailFileWriter.WriteAll(dtos)
	if err != nil {
		return err
	}

	return nil
}
