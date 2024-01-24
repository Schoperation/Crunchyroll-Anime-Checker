package anime

import (
	"fmt"
	"schoperation/crunchyroll-anime-checker/domain/anime"
	"schoperation/crunchyroll-anime-checker/domain/core"
)

type posterDao interface {
	GetAllByAnimeIds(animeIds []int) ([]anime.ImageDto, error)
	InsertAll(dtos []anime.ImageDto) error
	Update(dto anime.ImageDto) error
}

type posterFileWriter interface {
	WriteAll(dtos []anime.PostersDto) error
}

type PosterTranslator struct {
	posterDao        posterDao
	posterFileWriter posterFileWriter
}

func NewPosterTranslator(posterDao posterDao, posterFileWriter posterFileWriter) PosterTranslator {
	return PosterTranslator{
		posterDao:        posterDao,
		posterFileWriter: posterFileWriter,
	}
}

func (translator PosterTranslator) GetAllByAnimeIds(animeIds []anime.AnimeId) (map[anime.AnimeId][]anime.Image, error) {
	ids := make([]int, len(animeIds))
	for i, animeId := range animeIds {
		ids[i] = animeId.Int()
	}

	dtos, err := translator.posterDao.GetAllByAnimeIds(ids)
	if err != nil {
		return nil, err
	}

	posterMap := make(map[anime.AnimeId][]anime.Image)
	for _, dto := range dtos {
		animeId := anime.ReformAnimeId(dto.AnimeId)

		if posters, exists := posterMap[animeId]; exists {
			posterMap[animeId] = append(posters, anime.ReformImage(dto))
			continue
		}

		posterMap[animeId] = []anime.Image{anime.ReformImage(dto)}
	}

	return posterMap, nil
}

func (translator PosterTranslator) SaveAll(newPosters []anime.Image, updatedPosters []anime.Image) error {
	newDtos := make([]anime.ImageDto, len(newPosters))
	for i, poster := range newPosters {
		newDtos[i] = poster.Dto()
	}

	err := translator.posterDao.InsertAll(newDtos)
	if err != nil {
		return err
	}

	for _, poster := range updatedPosters {
		err := translator.posterDao.Update(poster.Dto())
		if err != nil {
			return err
		}
	}

	return nil
}

func (translator PosterTranslator) CreatePosterFiles(posters []anime.Image, slugTitles map[anime.AnimeId]string) error {
	posterDtoMap := make(map[anime.AnimeId]anime.PostersDto, len(slugTitles))
	for _, poster := range posters {
		slugTitle, exists := slugTitles[poster.AnimeId()]
		if !exists {
			return fmt.Errorf("missing slug title for anime ID %d", poster.AnimeId())
		}

		dto := anime.PostersDto{
			SlugTitle: slugTitle,
		}

		savedDto, exists := posterDtoMap[poster.AnimeId()]
		if exists {
			dto = savedDto
		}

		switch poster.ImageType() {
		case core.ImageTypePosterTall:
			dto.PosterTallUrl = poster.Url()
			dto.PosterTallEncoded = poster.Encoded()
		case core.ImageTypePosterWide:
			dto.PosterWideUrl = poster.Url()
			dto.PosterWideEncoded = poster.Encoded()
		}

		posterDtoMap[poster.AnimeId()] = dto
	}

	posterDtos := make([]anime.PostersDto, len(posterDtoMap))
	i := 0
	for _, dto := range posterDtoMap {
		posterDtos[i] = dto
		i++
	}

	err := translator.posterFileWriter.WriteAll(posterDtos)
	if err != nil {
		return err
	}

	return nil
}
