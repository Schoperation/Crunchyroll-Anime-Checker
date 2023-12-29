package core

import (
	"fmt"
	"strings"
)

type ImageType int

const (
	ImageTypeUnknown    ImageType = 0
	ImageTypePosterTall ImageType = 1
	ImageTypePosterWide ImageType = 2
	ImageTypeThumbnail  ImageType = 3
)

var imageTypeNames = map[ImageType]string{
	ImageTypeUnknown:    "unknown",
	ImageTypePosterTall: "poster_tall",
	ImageTypePosterWide: "poster_wide",
	ImageTypeThumbnail:  "thumbnail",
}

func NewImageTypeFromNumber(num int) (ImageType, error) {
	if num < 1 || num > 3 {
		return 0, fmt.Errorf("invalid image type int %d", num)
	}

	return ImageType(num), nil
}

func NewImageTypeFromString(imageType string) (ImageType, error) {
	for enum, name := range imageTypeNames {
		if strings.EqualFold(imageType, name) {
			return enum, nil
		}
	}

	return ImageTypeUnknown, fmt.Errorf("invalid image type %s", imageType)
}

func ReformImageTypeFromNumber(num int) ImageType {
	return ImageType(num)
}

func ReformImageTypeFromString(imageType string) ImageType {
	for enum, name := range imageTypeNames {
		if strings.EqualFold(imageType, name) {
			return enum
		}
	}

	return ImageTypeUnknown
}

func (imageType ImageType) Int() int {
	return int(imageType)
}

func (imageType ImageType) Name() string {
	return imageTypeNames[imageType]
}

func (imageType ImageType) IsThumbnail() bool {
	return imageType == ImageTypeThumbnail
}
