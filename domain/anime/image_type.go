package anime

import "fmt"

type ImageType int

const (
	ImageTypePosterTall ImageType = 1
	ImageTypePosterWide ImageType = 2
	ImageTypeThumbnail  ImageType = 3
)

var imageTypeNames = map[ImageType]string{
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

func ReformImageTypeFromNumber(num int) ImageType {
	return ImageType(num)
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
