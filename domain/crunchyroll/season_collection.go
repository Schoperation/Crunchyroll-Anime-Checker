package crunchyroll

import (
	"cmp"
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/core"
	"slices"
)

// SeasonCollection is a collection of seasons for a Crunchyroll anime.
// This is primarily used to determine the latest season for a particular locale, and
// to help filter out scrupulous seasons (e.g. OVA seasons with one episode)
// It is sorted in descending order automatically for easier iteration.
type SeasonCollection struct {
	seriesId core.SeriesId
	seasons  []Season
}

func NewSeasonCollection(seriesId core.SeriesId, seasons []Season) (SeasonCollection, error) {
	if len(seasons) == 0 {
		return SeasonCollection{}, fmt.Errorf("anime must have seasons; series ID %s", seriesId)
	}

	slices.SortFunc(seasons, func(x, y Season) int {
		return cmp.Compare(x.SequenceNumber(), y.SequenceNumber())
	})

	return SeasonCollection{
		seriesId: seriesId,
		seasons:  seasons,
	}, nil
}

func (col SeasonCollection) LatestSub(locale core.Locale) (Season, bool) {
	for i := len(col.seasons) - 1; i >= 0; i-- {
		if !isValidSeason(col.seasons[i]) {
			continue
		}

		if col.seasons[i].hasSubForLocale(locale) {
			return col.seasons[i], true
		}
	}

	return Season{}, false
}

func (col SeasonCollection) LatestDub(locale core.Locale) (Season, bool) {
	for i := len(col.seasons) - 1; i >= 0; i-- {
		if !isValidSeason(col.seasons[i]) {
			continue
		}

		if col.seasons[i].hasDubForLocale(locale) {
			return col.seasons[i], true
		}
	}

	return Season{}, false
}

// Determines if a season has tangible episodes rather than a movie, OVA, interview, etc.
func isValidSeason(season Season) bool {
	return !slices.Contains(season.Keywords(), "movie")
}
