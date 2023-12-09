package anime

import (
	"fmt"
	"net/url"
	"strings"
)

type ImageDto struct {
	AnimeId       int
	ImageType     int
	SeasonNumber  int
	EpisodeNumber int
	Url           string
	Encoded       string
}

type Image struct {
	animeId       int
	imageType     ImageType
	seasonNumber  int
	episodeNumber int
	url           string
	encoded       string
}

func NewImage(dto ImageDto) (Image, error) {
	if dto.AnimeId <= 0 {
		return Image{}, fmt.Errorf("anime id must be greater than 0")
	}

	imageType, err := NewImageTypeFromNumber(dto.ImageType)
	if err != nil {
		return Image{}, err
	}

	if imageType.IsThumbnail() {
		if dto.SeasonNumber <= 0 {
			return Image{}, fmt.Errorf("season number must be greater than 0")
		}

		if dto.EpisodeNumber <= 0 {
			return Image{}, fmt.Errorf("episode number must be greater than 0")
		}
	} else {
		dto.SeasonNumber = 0
		dto.EpisodeNumber = 0
	}

	if _, err := url.ParseRequestURI(dto.Url); err != nil {
		return Image{}, fmt.Errorf("invalid URL for image: %v", err)
	}

	if strings.Trim(dto.Encoded, " ") == "" {
		return Image{}, fmt.Errorf("encoded image must not be blank")
	}

	return Image{
		animeId:       dto.AnimeId,
		imageType:     imageType,
		seasonNumber:  dto.SeasonNumber,
		episodeNumber: dto.EpisodeNumber,
		url:           dto.Url,
		encoded:       dto.Encoded,
	}, nil
}

func ReformImage(dto ImageDto) Image {
	return Image{
		animeId:       dto.AnimeId,
		imageType:     ReformImageTypeFromNumber(dto.ImageType),
		seasonNumber:  dto.SeasonNumber,
		episodeNumber: dto.EpisodeNumber,
		url:           dto.Url,
		encoded:       dto.Encoded,
	}
}

func (image Image) AnimeId() int {
	return image.animeId
}

func (image Image) ImageType() ImageType {
	return image.imageType
}

func (image Image) Url() string {
	return image.url
}

func (image Image) Encoded() string {
	return image.encoded
}
