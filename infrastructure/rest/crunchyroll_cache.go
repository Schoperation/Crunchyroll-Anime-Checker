package rest

import "fmt"

const episodeResponsesCacheLimit = 3

// Basic in-memory cache to store responses we may need multiples of.
type crunchyrollCache struct {
	episodeResponses      map[string]episodesResponse
	episodeResponsesOrder map[int]string
	oldestOrder           int
	newestOrder           int
}

func newCrunchyrollCache() crunchyrollCache {
	return crunchyrollCache{
		episodeResponses:      make(map[string]episodesResponse, episodeResponsesCacheLimit),
		episodeResponsesOrder: make(map[int]string, episodeResponsesCacheLimit),
		oldestOrder:           1,
		newestOrder:           0,
	}
}

func (cache *crunchyrollCache) GetEpisodesResponse(locale, seasonId string) (episodesResponse, bool) {
	response, exists := cache.episodeResponses[cache.key(locale, seasonId)]
	if exists {
		return response, true
	}

	return episodesResponse{}, false
}

func (cache *crunchyrollCache) SaveEpisodesResponse(locale, seasonId string, resp episodesResponse) {
	if len(cache.episodeResponses) >= episodeResponsesCacheLimit {
		oldestSeasonId := cache.episodeResponsesOrder[cache.oldestOrder]
		delete(cache.episodeResponses, oldestSeasonId)
		delete(cache.episodeResponsesOrder, cache.oldestOrder)
		cache.oldestOrder++
	}

	_, exists := cache.episodeResponses[cache.key(locale, seasonId)]
	if exists {
		return
	}

	cache.episodeResponses[cache.key(locale, seasonId)] = resp
	cache.episodeResponsesOrder[cache.newestOrder+1] = seasonId
	cache.newestOrder++
}

func (cache *crunchyrollCache) key(locale, seasonId string) string {
	return fmt.Sprintf("%s_%s", locale, seasonId)
}
