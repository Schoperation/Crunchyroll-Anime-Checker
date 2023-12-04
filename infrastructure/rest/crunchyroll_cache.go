package rest

// Basic in-memory cache to store responses we may need multiples of.
type crunchyrollCache struct {
	cachedAllAnimeResponse allAnimeResponse
}

func newCrunchyrollCache() crunchyrollCache {
	return crunchyrollCache{
		cachedAllAnimeResponse: allAnimeResponse{Total: 0, Data: nil},
	}
}

func (cache *crunchyrollCache) GetAllAnimeResponse() (allAnimeResponse, bool) {
	return cache.cachedAllAnimeResponse, cache.cachedAllAnimeResponse.Total != 0
}

func (cache *crunchyrollCache) SaveAllAnimeResponse(resp allAnimeResponse) {
	cache.cachedAllAnimeResponse = resp
}
