package crunchyroll

import (
	"cmp"
	"fmt"
	"schoperation/crunchyroll-anime-checker/domain/core"
	"slices"
	"strings"
)

// EpisodeCollection is a collection of episodes, primarily used to retrieve the latest sub or dub for a specified locale.
// It is automatically sorted in descending order for easy iteration.
type EpisodeCollection struct {
	seasonId string
	episodes []Episode
}

func NewEpisodeCollection(seasonId string, episodes []Episode) (EpisodeCollection, error) {
	if strings.Trim(seasonId, " ") == "" {
		return EpisodeCollection{}, fmt.Errorf("season id must not be blank")
	}

	if len(episodes) == 0 {
		return EpisodeCollection{}, fmt.Errorf("season must have episodes; season ID %s", seasonId)
	}

	slices.SortFunc(episodes, func(x, y Episode) int {
		return cmp.Compare(x.Number(), y.Number())
	})

	return EpisodeCollection{
		seasonId: seasonId,
		episodes: episodes,
	}, nil
}

func (col EpisodeCollection) LatestSub(locale core.Locale) (Episode, bool) {
	for i := len(col.episodes) - 1; i >= 0; i-- {
		if col.episodes[i].hasSubForLocale(locale) {
			return col.episodes[i], true
		}
	}

	return Episode{}, false
}

func (col EpisodeCollection) LatestDub(locale core.Locale) (Episode, bool) {
	for i := len(col.episodes) - 1; i >= 0; i-- {
		if col.episodes[i].hasDubForLocale(locale) {
			return col.episodes[i], true
		}
	}

	return Episode{}, false
}
