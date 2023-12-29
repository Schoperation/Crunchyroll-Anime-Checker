package anime

import (
	"fmt"
)

type AnimeId int

func NewAnimeId(id int) (AnimeId, error) {
	if id <= 0 {
		return AnimeId(0), fmt.Errorf("anime id must be greater than 0")
	}

	return AnimeId(id), nil
}

func ReformAnimeId(id int) AnimeId {
	return AnimeId(id)
}

func (id AnimeId) Int() int {
	return int(id)
}

func (id AnimeId) IsZero() bool {
	return id == 0
}
