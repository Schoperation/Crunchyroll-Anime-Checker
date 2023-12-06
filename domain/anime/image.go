package anime

import (
	"fmt"
	"net/url"
	"strings"
)

type ImageDto struct {
	ImageType int
	Url       string
	Encoded   string
}

type Image struct {
	imageType ImageType
	url       string
	encoded   string
}

func NewImage(dto ImageDto) (Image, error) {
	imageType, err := NewImageTypeFromNumber(dto.ImageType)
	if err != nil {
		return Image{}, err
	}

	if _, err := url.ParseRequestURI(dto.Url); err != nil {
		return Image{}, fmt.Errorf("invalid URL for image: %v", err)
	}

	if strings.Trim(dto.Encoded, " ") == "" {
		return Image{}, fmt.Errorf("encoded image must not be blank")
	}

	return Image{
		imageType: imageType,
		url:       dto.Url,
		encoded:   dto.Encoded,
	}, nil
}

func ReformImage(dto ImageDto) Image {
	return Image{
		imageType: ReformImageTypeFromNumber(dto.ImageType),
		url:       dto.Url,
		encoded:   dto.Encoded,
	}
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
