package file

import (
	"encoding/csv"
	"fmt"
	"os"
)

type AnimeSenseiListWriter struct {
	path string
}

func NewAnimeSenseiListWriter(path string) AnimeSenseiListWriter {
	return AnimeSenseiListWriter{
		path: path,
	}
}

func (writer AnimeSenseiListWriter) WriteAll(seriesIds, slugTitles, titles []string) error {
	if len(seriesIds) != len(slugTitles) || len(seriesIds) != len(titles) {
		return fmt.Errorf("data slices must be equal lengths")
	}

	data := make([][]string, len(seriesIds))
	for i := range seriesIds {
		data[i] = []string{
			fileId(slugTitles[i]),
			seriesIds[i],
			slugTitles[i],
			titles[i],
		}
	}

	newAnimeSenseiList, err := os.Create(fmt.Sprintf("%s/anime_sensei_list.csv", writer.path))
	if err != nil {
		return err
	}

	newListAsCsv := csv.NewWriter(newAnimeSenseiList)
	newListAsCsv.Comma = '|'
	newListAsCsv.Write([]string{"file_id", "series_id", "slug_title", "title"})
	newListAsCsv.WriteAll(data)
	newAnimeSenseiList.Close()

	return nil
}
