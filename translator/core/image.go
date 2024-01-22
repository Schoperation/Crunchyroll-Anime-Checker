package core

import "schoperation/crunchyroll-anime-checker/domain/core"

type getEncodedImageClient interface {
	GetImageByURL(url string) ([]byte, error)
}

type ImageTranslator struct {
	getEncodedImageClient getEncodedImageClient
}

func NewImageTranslator(getEncodedImageClient getEncodedImageClient) ImageTranslator {
	return ImageTranslator{
		getEncodedImageClient: getEncodedImageClient,
	}
}

func (translator ImageTranslator) GetEncodedImageByURL(url string) (string, error) {
	imageBytes, err := translator.getEncodedImageClient.GetImageByURL(url)
	if err != nil {
		return "", err
	}

	imageEncoder := core.NewImageEncoder()
	return imageEncoder.Encode(imageBytes), nil
}
