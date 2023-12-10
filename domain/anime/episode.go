package anime

import (
	"fmt"
)

type NewEpisodeArgs struct {
	AnimeId      AnimeId
	Number       int
	SeasonNumber int
	Thumbnail    Image
	Titles       []TitleDto
}

// Episode represents a single episode that may be the latest sub/dub for multiple locales.
type Episode struct {
	animeId      AnimeId
	number       int
	seasonNumber int
	thumbnail    Image
	titles       TitleCollection
}

func NewEpisode(args NewEpisodeArgs) (Episode, error) {
	if args.Number <= 0 {
		return Episode{}, fmt.Errorf("episode number must be greater than 0")
	}

	if args.SeasonNumber <= 0 {
		return Episode{}, fmt.Errorf("season number must be greater than 0")
	}

	titles, err := NewTitleCollection(args.Titles)
	if err != nil {
		return Episode{}, err
	}

	return Episode{
		animeId:      args.AnimeId,
		number:       args.Number,
		seasonNumber: args.SeasonNumber,
		thumbnail:    args.Thumbnail,
		titles:       titles,
	}, nil
}

func ReformEpisode(args NewEpisodeArgs) Episode {
	return Episode{
		animeId:      args.AnimeId,
		number:       args.Number,
		seasonNumber: args.SeasonNumber,
		thumbnail:    args.Thumbnail,
		titles:       ReformTitleCollection(args.Titles),
	}
}

func (e Episode) AnimeId() AnimeId {
	return e.animeId
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
