package file

import (
	"encoding/json"
	"fmt"
	"os"
	"schoperation/crunchyroll-anime-checker/domain/anime"
)

type LatestEpisodesWriter struct {
	path string
}

func NewLatestEpisodesWriter(path string) LatestEpisodesWriter {
	return LatestEpisodesWriter{
		path: path,
	}
}

type latestEpisodesFileModel struct {
	TotalCount     int                            `json:"total_count"`
	LatestEpisodes map[string]latestEpisodesModel `json:"latest_episodes"`
}

type latestEpisodesModel struct {
	Sub episodeModel `json:"sub"`
	Dub episodeModel `json:"dub"`
}

type episodeModel struct {
	Season  int    `json:"season"`
	Episode int    `json:"episode"`
	Title   string `json:"title"`
}

func (writer LatestEpisodesWriter) WriteAllByLocale(localeName string, dtos []anime.LatestEpisodesDto) error {
	latestEpisodes := make(map[string]latestEpisodesModel)

	for _, dto := range dtos {
		latestEpisodes[dto.SlugTitle] = latestEpisodesModel{
			Sub: episodeModel{
				Season:  dto.LatestSubSeason,
				Episode: dto.LatestSubEpisode,
				Title:   dto.LatestSubTitle,
			},
			Dub: episodeModel{
				Season:  dto.LatestDubSeason,
				Episode: dto.LatestDubEpisode,
				Title:   dto.LatestDubTitle,
			},
		}
	}

	fileModel := latestEpisodesFileModel{
		TotalCount:     len(dtos),
		LatestEpisodes: latestEpisodes,
	}

	err := os.Mkdir(fmt.Sprintf("%s/latest_episodes", writer.path), 0770)
	if err != nil && err == os.ErrExist {
		return err
	}

	newFile, err := os.Create(fmt.Sprintf("%s/latest_episodes/%s_new.json", writer.path, localeName))
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

	err = os.Rename(fmt.Sprintf("%s/latest_episodes/%s_new.json", writer.path, localeName), fmt.Sprintf("%s/latest_episodes/%s.json", writer.path, localeName))
	if err != nil {
		return err
	}

	return nil
}
