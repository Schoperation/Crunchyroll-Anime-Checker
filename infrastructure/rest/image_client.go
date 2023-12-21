package rest

import (
	"fmt"
	"io"
	"net/http"
)

type ImageClient struct{}

func NewImageClient() ImageClient {
	return ImageClient{}
}

func (client ImageClient) GetImageByURL(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotModified {
		return nil, fmt.Errorf("got bad response when retrieving image: %d", response.StatusCode)
	}

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
