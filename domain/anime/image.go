package anime

import (
	"fmt"
	"net/url"
	"schoperation/crunchyrollanimestatus/domain/core"
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
	animeId       AnimeId
	imageType     core.ImageType
	seasonNumber  int
	episodeNumber int
	url           string
	encoded       string
	isDirty       bool
}

func NewImage(dto ImageDto) (Image, error) {
	imageType, err := core.NewImageTypeFromNumber(dto.ImageType)
	if err != nil {
		return Image{}, err
	}

	if imageType.IsThumbnail() {
		if dto.SeasonNumber <= 0 {
			return Image{}, fmt.Errorf("image season number must be greater than 0")
		}

		if dto.EpisodeNumber <= 0 {
			return Image{}, fmt.Errorf("image episode number must be greater than 0")
		}
	} else {
		dto.SeasonNumber = 0
		dto.EpisodeNumber = 0
	}

	if _, err := url.ParseRequestURI(dto.Url); err != nil {
		return Image{}, fmt.Errorf("image invalid URL for image: %v", err)
	}

	if strings.Trim(dto.Encoded, " ") == "" {
		return Image{}, fmt.Errorf("image encoded image must not be blank")
	}

	return Image{
		animeId:       0,
		imageType:     imageType,
		seasonNumber:  dto.SeasonNumber,
		episodeNumber: dto.EpisodeNumber,
		url:           dto.Url,
		encoded:       dto.Encoded,
		isDirty:       true,
	}, nil
}

func ReformImage(dto ImageDto) Image {
	return Image{
		animeId:       ReformAnimeId(dto.AnimeId),
		imageType:     core.ReformImageTypeFromNumber(dto.ImageType),
		seasonNumber:  dto.SeasonNumber,
		episodeNumber: dto.EpisodeNumber,
		url:           dto.Url,
		encoded:       dto.Encoded,
		isDirty:       false,
	}
}

func (image Image) AnimeId() AnimeId {
	return image.animeId
}

func (image Image) ImageType() core.ImageType {
	return image.imageType
}

func (image Image) Url() string {
	return image.url
}

func (image Image) Encoded() string {
	return image.encoded
}

func (image Image) IsDirty() bool {
	return image.isDirty
}

func (image *Image) AssignAnimeId(animeId AnimeId) {
	if !image.animeId.IsZero() {
		return
	}

	image.animeId = animeId
	image.isDirty = true
}

func (image *Image) UpdatePoster(newUrl, encoded string) error {
	if image.imageType != core.ImageTypePosterTall && image.imageType != core.ImageTypePosterWide {
		return fmt.Errorf("image must be a poster to only update url and encoded string")
	}

	if _, err := url.ParseRequestURI(newUrl); err != nil {
		return fmt.Errorf("image invalid URL for image: %v", err)
	}

	if strings.Trim(encoded, " ") == "" {
		return fmt.Errorf("image encoded image must not be blank")
	}

	image.url = newUrl
	image.encoded = encoded
	image.isDirty = true
	return nil
}
