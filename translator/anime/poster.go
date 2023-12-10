package anime

import "schoperation/crunchyrollanimestatus/domain/anime"

type posterDao interface {
	GetAllByAnimeId(animeId int) ([]anime.ImageDto, error)
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
