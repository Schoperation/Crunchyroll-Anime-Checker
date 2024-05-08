package file

import (
	"encoding/json"
	"fmt"
	"os"
	"schoperation/crunchyroll-anime-checker/domain/anime"
	"schoperation/crunchyroll-anime-checker/domain/core"
)

type ThumbnailWriter struct {
	path string
}

func NewThumbnailWriter(path string) ThumbnailWriter {
	return ThumbnailWriter{
		path: path,
	}
}

type thumbnailsFileModel struct {
	TotalCount              int                                  `json:"total_count"`
	DefaultThumbnailUrl     string                               `json:"default_thumbnail_url"`
	DefaultThumbnailEncoded string                               `json:"default_thumbnail_encoded"`
	Thumbnails              map[string]map[string]thumbnailModel `json:"thumbnails"`
}

type thumbnailModel struct {
	Url     string `json:"url"`
	Encoded string `json:"encoded"`
}

func (writer ThumbnailWriter) WriteAll(dtos []anime.ImageDto) error {
	fileIds := getFileIds()

	// Map-ception
	thumbnailMaps := make(map[string]map[string]map[string]thumbnailModel, len(fileIds))
	for _, id := range fileIds {
		thumbnailMaps[id] = make(map[string]map[string]thumbnailModel)
	}

	for _, dto := range dtos {
		id := fileId(dto.SlugTitle)

		animes := thumbnailMaps[id]
		thumbnails := map[string]thumbnailModel{}

		if savedThumbnails, exist := animes[dto.SlugTitle]; exist {
			thumbnails = savedThumbnails
		}

		thumbnails[fmt.Sprintf("%d-%d", dto.SeasonNumber, dto.EpisodeNumber)] = thumbnailModel{
			Url:     dto.Url,
			Encoded: dto.Encoded,
		}

		animes[dto.SlugTitle] = thumbnails
		thumbnailMaps[id] = animes
	}

	err := os.Mkdir(fmt.Sprintf("%s/thumbnails", writer.path), 0770)
	if err != nil && err == os.ErrExist {
		return err
	}

	for _, id := range fileIds {
		fileModel := thumbnailsFileModel{
			TotalCount:              len(thumbnailMaps[id]),
			DefaultThumbnailUrl:     core.DefaultPosterUrl,
			DefaultThumbnailEncoded: core.DefaultPosterEncoded,
			Thumbnails:              thumbnailMaps[id],
		}

		bytes, err := json.MarshalIndent(fileModel, "", "    ")
		if err != nil {
			return err
		}

		newFile, err := os.Create(fmt.Sprintf("%s/thumbnails/%s.json", writer.path, id))
		if err != nil {
			return err
		}

		_, err = newFile.Write(bytes)
		if err != nil {
			return err
		}
	}

	return nil
}
