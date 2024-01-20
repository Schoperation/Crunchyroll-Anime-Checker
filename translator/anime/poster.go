package anime

import "schoperation/crunchyrollanimestatus/domain/anime"

type posterDao interface {
	GetAllByAnimeId(animeId int) ([]anime.ImageDto, error)
	GetAllByAnimeIds(animeIds []int) ([]anime.ImageDto, error)
	InsertAll(dtos []anime.ImageDto) error
}

type PosterTranslator struct {
	posterDao posterDao
}

func NewPosterTranslator(posterDao posterDao) PosterTranslator {
	return PosterTranslator{
		posterDao: posterDao,
	}
}

func (translator PosterTranslator) GetAllByAnimeId(animeId anime.AnimeId) ([]anime.Image, error) {
	dtos, err := translator.posterDao.GetAllByAnimeId(animeId.Int())
	if err != nil {
		return nil, err
	}

	images := make([]anime.Image, len(dtos))
	for i, dto := range dtos {
		images[i] = anime.ReformImage(dto)
	}

	return images, nil
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

func (translator PosterTranslator) SaveAll(newPosters []anime.Image) error {
	newDtos := make([]anime.ImageDto, len(newPosters))
	for i, poster := range newPosters {
		newDtos[i] = poster.Dto()
	}

	err := translator.posterDao.InsertAll(newDtos)
	if err != nil {
		return err
	}

	return nil
}
