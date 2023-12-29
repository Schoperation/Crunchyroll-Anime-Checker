package crunchyroll

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/core"
	"strings"
)

// SeasonCollection is a collection of seasons for a Crunchyroll anime.
// This is primarily used to determine the latest season for a particular locale, and
// to help filter out scrupulous seasons (e.g. OVAs)
type SeasonCollection struct {
	seriesId string
	seasons  []Season
}

func NewSeasonCollection(seriesId string, seasons []Season) (SeasonCollection, error) {
	if strings.Trim(seriesId, " ") == "" {
		return SeasonCollection{}, fmt.Errorf("series id must not be blank")
	}

	if len(seasons) == 0 {
		return SeasonCollection{}, fmt.Errorf("anime must have seasons; series ID %s", seriesId)
	}

	return SeasonCollection{
		seriesId: seriesId,
		seasons:  seasons,
	}, nil
}

func (col SeasonCollection) LatestSub(locale core.Locale) Season {
	return Season{}
}

func (col SeasonCollection) LatestDub(locale core.Locale) Season {
	return Season{}
}

// Determines if a season has tangible episodes rather than a movie, OVA, interview, etc.
func (col SeasonCollection) isValidSeason(season Season) bool {
	return true
}
