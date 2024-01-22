package crunchyroll

import (
	"schoperation/crunchyroll-anime-checker/domain/core"
)

type ImageDto struct {
	Width     int
	Height    int
	ImageType string
	Source    string
}

type Image struct {
	width     int
	height    int
	imageType core.ImageType
	source    string
}

func ReformImage(dto ImageDto) Image {
	return Image{
		width:     dto.Width,
		height:    dto.Height,
		imageType: core.ReformImageTypeFromString(dto.ImageType),
		source:    dto.Source,
	}
}

func (image Image) ImageType() core.ImageType {
	return image.imageType
}

func (image Image) Source() string {
	return image.source
}
