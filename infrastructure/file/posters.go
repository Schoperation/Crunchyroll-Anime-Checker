package file

import (
	"encoding/json"
	"fmt"
	"os"
	"schoperation/crunchyroll-anime-checker/domain/anime"
	"schoperation/crunchyroll-anime-checker/domain/core"
)

type PosterWriter struct {
	path string
}

func NewPosterWriter(path string) PosterWriter {
	return PosterWriter{
		path: path,
	}
}

type postersFileModel struct {
	TotalCount           int                     `json:"total_count"`
	DefaultPosterUrl     string                  `json:"default_poster_url"`
	DefaultPosterEncoded string                  `json:"default_poster_encoded"`
	Posters              map[string]postersModel `json:"posters"`
}

type postersModel struct {
	PosterTallUrl     string `json:"poster_tall_url"`
	PosterTallEncoded string `json:"poster_tall_encoded"`
	PosterWideUrl     string `json:"poster_wide_url"`
	PosterWideEncoded string `json:"poster_wide_encoded"`
}

func (writer PosterWriter) WriteAllWithIdentifier(identifier string, dtos []anime.PostersDto) error {
	posters := make(map[string]postersModel, len(dtos))

	for _, dto := range dtos {
		posters[dto.SlugTitle] = postersModel{
			PosterTallUrl:     dto.PosterTallUrl,
			PosterTallEncoded: dto.PosterTallEncoded,
			PosterWideUrl:     dto.PosterWideUrl,
			PosterWideEncoded: dto.PosterWideEncoded,
		}
	}

	fileModel := postersFileModel{
		TotalCount:           len(dtos),
		DefaultPosterUrl:     core.DefaultPosterUrl,
		DefaultPosterEncoded: core.DefaultPosterEncoded,
		Posters:              posters,
	}

	err := os.Mkdir(fmt.Sprintf("%s/posters", writer.path), 0770)
	if err != nil && err == os.ErrExist {
		return err
	}

	newFile, err := os.Create(fmt.Sprintf("%s/posters/%s_new.json", writer.path, identifier))
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(fileModel, "", "    ")
	if err != nil {
		return err
	}

	_, err = newFile.Write(bytes)
	if err != nil {
		return err
	}

	err = os.Rename(fmt.Sprintf("%s/posters/%s_new.json", writer.path, identifier), fmt.Sprintf("%s/posters/%s.json", writer.path, identifier))
	if err != nil {
		return err
	}

	return nil
}
