package anime

import (
	"fmt"
)

type EpisodeDto struct {
	Number           int
	SeasonNumber     int
	ThumbnailUrl     string
	ThumbnailEncoded string
}

type Episode struct {
	number       int
	seasonNumber int
	thumbnail    Image
	titles       TitleCollection
}

func NewEpisode(dto EpisodeDto, titleDtos []TitleDto) (Episode, error) {
	if dto.Number <= 0 {
		return Episode{}, fmt.Errorf("episode number must be greater than 0")
	}

	if dto.SeasonNumber <= 0 {
		return Episode{}, fmt.Errorf("season number must be greater than 0")
	}

	image, err := NewImage(ImageDto{
		ImageType: ImageTypeThumbnail.Int(),
		Url:       dto.ThumbnailUrl,
		Encoded:   dto.ThumbnailEncoded,
	})
	if err != nil {
		return Episode{}, err
	}

	titles, err := NewTitleCollection(titleDtos)
	if err != nil {
		return Episode{}, err
	}

	return Episode{
		number:       dto.Number,
		seasonNumber: dto.SeasonNumber,
		thumbnail:    image,
		titles:       titles,
	}, nil
}

func ReformEpisode(dto EpisodeDto, titleDtos []TitleDto) Episode {
	return Episode{
		number:       dto.Number,
		seasonNumber: dto.SeasonNumber,
		thumbnail: ReformImage(ImageDto{
			ImageType: ImageTypeThumbnail.Int(),
			Url:       dto.ThumbnailUrl,
			Encoded:   dto.ThumbnailEncoded,
		}),
		titles: ReformTitleCollection(titleDtos),
	}
}

func (e Episode) Number() int {
	return e.number
}

func (e Episode) Season() int {
	return e.seasonNumber
}

func (e Episode) Thumbnail() Image {
	return e.thumbnail
}

func (e Episode) Titles() TitleCollection {
	return e.titles
}
