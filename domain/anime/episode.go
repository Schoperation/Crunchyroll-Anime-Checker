package anime

import (
	"fmt"
	"time"
)

type NewEpisodeArgs struct {
	Number       int
	SeasonNumber int
	Thumbnail    Image
	Titles       []TitleDto
	LastUpdated  time.Time
}

// Episode represents a single episode that may be the latest sub/dub for multiple locales.
type Episode struct {
	number       int
	seasonNumber int
	thumbnail    Image
	titles       TitleCollection
	lastUpdated  time.Time
}

func newEpisode(args NewEpisodeArgs) (Episode, error) {
	if args.Number <= 0 {
		return Episode{}, fmt.Errorf("episode number must be greater than 0")
	}

	if args.SeasonNumber <= 0 {
		return Episode{}, fmt.Errorf("episode season number must be greater than 0")
	}

	titles, err := NewTitleCollection(args.Titles)
	if err != nil {
		return Episode{}, err
	}

	return Episode{
		number:       args.Number,
		seasonNumber: args.SeasonNumber,
		thumbnail:    args.Thumbnail,
		titles:       titles,
		lastUpdated:  time.Now().UTC(),
	}, nil
}

func ReformEpisode(args NewEpisodeArgs) Episode {
	return Episode{
		number:       args.Number,
		seasonNumber: args.SeasonNumber,
		thumbnail:    args.Thumbnail,
		titles:       ReformTitleCollection(args.Titles),
		lastUpdated:  args.LastUpdated,
	}
}

func (e Episode) Number() int {
	return e.number
}

func (e Episode) Season() int {
	return e.seasonNumber
}

func (e Episode) Thumbnail() Image {
	return e.thumbnail
}

func (e Episode) Titles() TitleCollection {
	return e.titles
}

func (e *Episode) AddTitle(dto TitleDto) error {
	return e.titles.Add(dto)
}
