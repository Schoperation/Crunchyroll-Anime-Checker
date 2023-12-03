package crunchyroll

type ImageDto struct {
	Width     int
	Height    int
	ImageType string
	Source    string
}

type Image struct {
	width     int
	height    int
	imageType string
	source    string
}

func ReformImage(dto ImageDto) Image {
	return Image{
		width:     dto.Width,
		height:    dto.Height,
		imageType: dto.ImageType,
		source:    dto.Source,
	}
}
