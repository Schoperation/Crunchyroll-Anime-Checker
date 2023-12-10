package anime

import "fmt"

// MinimalEpisode is a part of LatestEpisodes to hold a title, season, and ep number for either a sub or dub.
type MinimalEpisode struct {
	season int
	number int
	title  string
}

func NewMinimalEpisode(season, number int, title string) (MinimalEpisode, error) {
	if season < 0 {
		return MinimalEpisode{}, fmt.Errorf("season must be 0 or above")
	}

	if number < 0 {
		return MinimalEpisode{}, fmt.Errorf("episode snumber must be 0 or above")
	}

	return MinimalEpisode{
		season: season,
		number: number,
		title:  title,
	}, nil
}

func ReformMinimalEpisode(season, number int, title string) MinimalEpisode {
	return MinimalEpisode{
		season: season,
		number: number,
		title:  title,
	}
}

func (ep MinimalEpisode) Season() int {
	return ep.season
}

func (ep MinimalEpisode) Number() int {
	return ep.number
}

func (ep MinimalEpisode) Title() string {
	return ep.title
}

func (ep MinimalEpisode) IsZero() bool {
	return ep.season == 0 && ep.number == 0 && ep.title == ""
}
