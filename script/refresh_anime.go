package script

import (
	"encoding/csv"
	"fmt"
	"os"
)

type RefreshAnimeCmd struct {
	name string
}

func NewRefreshAnimeCmd(name string) Command {
	return RefreshAnimeCmd{
		name: name,
	}
}

func (cmd RefreshAnimeCmd) Name() string {
	return cmd.name
}

func (cmd RefreshAnimeCmd) Run(client CrunchyrollClient) error {
	var allAnime allAnimeResponse
	err := client.GetWithQueryParams("content/v2/discover/browse", &allAnime, map[string]string{
		"start":   "0",
		"n":       "2000",
		"type":    "series",
		"sort_by": "alphabetical",
		"ratings": "true",
	})
	if err != nil {
		return err
	}

	newList := make([][]string, allAnime.Total)
	for i, series := range allAnime.Data {
		// series_id, slug_title, title
		newList[i] = []string{
			series.Id,
			series.SlugTitle,
			series.Title,
		}
	}

	newAnimeSenseiList, err := os.Create(fmt.Sprintf("%s/anime_sensei_list_new.csv", client.listsPath))
	if err != nil {
		return err
	}

	newListAsCsv := csv.NewWriter(newAnimeSenseiList)
	newListAsCsv.Comma = '|'
	newListAsCsv.Write([]string{"series_id", "slug_title", "title"})
	newListAsCsv.WriteAll(newList)
	newAnimeSenseiList.Close()

	err = os.Rename(fmt.Sprintf("%s/anime_sensei_list_new.csv", client.listsPath), fmt.Sprintf("%s/anime_sensei_list.csv", client.listsPath))
	if err != nil {
		return err
	}

	return nil
}
