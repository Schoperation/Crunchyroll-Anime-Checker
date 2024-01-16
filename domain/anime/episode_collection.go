package anime

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/core"
)

// EpisodeCollection is a collection of episodes for an anime.
// It only includes the latest ones for all applicable locales.
// The latestSubs and latestDubs maps locales to the ID (Season Num-Episode Num) of the episode, in the episodes map.
type EpisodeCollection struct {
	animeId    AnimeId
	latestSubs map[core.Locale]string
	latestDubs map[core.Locale]string
	episodes   map[string]Episode
}

func NewEpisodeCollection(animeId AnimeId, latestEpisodes []LatestEpisodes, thumbnails map[string]Image) (EpisodeCollection, error) {
	episodeCollection := EpisodeCollection{
		animeId:    animeId,
		latestSubs: map[core.Locale]string{},
		latestDubs: map[core.Locale]string{},
		episodes:   map[string]Episode{},
	}

	for _, latestEpsForLocale := range latestEpisodes {
		thumbnail, exists := thumbnails[latestEpsForLocale.LatestSub().Key()]
		if !exists {
			return EpisodeCollection{}, fmt.Errorf("epcol: no thumbnail found for sub: %s", latestEpsForLocale.LatestSub().Key())
		}

		err := episodeCollection.AddSubForLocale(latestEpsForLocale.Locale(), latestEpsForLocale.LatestSub(), thumbnail)
		if err != nil {
			return EpisodeCollection{}, nil
		}

		thumbnail, exists = thumbnails[latestEpsForLocale.LatestDub().Key()]
		if !exists {
			return EpisodeCollection{}, fmt.Errorf("epcol: no thumbnail found for dub: %s", latestEpsForLocale.LatestDub().Key())
		}

		err = episodeCollection.AddDubForLocale(latestEpsForLocale.Locale(), latestEpsForLocale.LatestDub(), thumbnail)
		if err != nil {
			return EpisodeCollection{}, nil
		}
	}

	return episodeCollection, nil
}

func ReformEpisodeCollection(animeId AnimeId, latestEpisodes []LatestEpisodes, thumbnails map[string]Image) EpisodeCollection {
	episodeCollection := EpisodeCollection{
		animeId:    animeId,
		latestSubs: map[core.Locale]string{},
		latestDubs: map[core.Locale]string{},
		episodes:   map[string]Episode{},
	}

	for _, latestEpsForLocale := range latestEpisodes {
		thumbnail := thumbnails[latestEpsForLocale.LatestSub().Key()]
		_ = episodeCollection.AddSubForLocale(latestEpsForLocale.Locale(), latestEpsForLocale.LatestSub(), thumbnail)

		thumbnail = thumbnails[latestEpsForLocale.LatestDub().Key()]
		_ = episodeCollection.AddDubForLocale(latestEpsForLocale.Locale(), latestEpsForLocale.LatestDub(), thumbnail)
	}

	return episodeCollection
}

func (epcol *EpisodeCollection) GetLatestEpisodesForLocale(locale core.Locale) LatestEpisodes {
	latestSubSeason := 0
	latestSubEpisode := 0
	latestSubTitle := ""

	if latestSub, ok := epcol.latestSubs[locale]; ok {
		episode := epcol.episodes[latestSub]
		latestSubSeason = episode.Season()
		latestSubEpisode = episode.Number()
		latestSubTitle = episode.Titles().Title(locale)
	}

	latestDubSeason := 0
	latestDubEpisode := 0
	latestDubTitle := ""

	if latestDub, ok := epcol.latestDubs[locale]; ok {
		episode := epcol.episodes[latestDub]
		latestDubSeason = episode.Season()
		latestDubEpisode = episode.Number()
		latestDubTitle = episode.Titles().Title(locale)
	}

	return ReformLatestEpisodes(LatestEpisodesDto{
		AnimeId:          epcol.animeId.Int(),
		LocaleId:         locale.Id(),
		LatestSubSeason:  latestSubSeason,
		LatestSubEpisode: latestSubEpisode,
		LatestSubTitle:   latestSubTitle,
		LatestDubSeason:  latestDubSeason,
		LatestDubEpisode: latestDubEpisode,
		LatestDubTitle:   latestDubTitle,
	})
}

func (epcol *EpisodeCollection) AddSubForLocale(locale core.Locale, sub MinimalEpisode, thumbnail Image) error {
	if sub.IsBlank() {
		return nil
	}

	epcol.latestSubs[locale] = sub.Key()
	err := epcol.addEpisode(locale, sub, thumbnail)
	if err != nil {
		return err
	}

	return nil
}

func (epcol *EpisodeCollection) AddDubForLocale(locale core.Locale, dub MinimalEpisode, thumbnail Image) error {
	if dub.IsBlank() {
		return nil
	}

	epcol.latestDubs[locale] = dub.Key()
	err := epcol.addEpisode(locale, dub, thumbnail)
	if err != nil {
		return err
	}

	return nil
}

func (epcol *EpisodeCollection) addEpisode(locale core.Locale, episode MinimalEpisode, thumbnail Image) error {
	if ep, exists := epcol.episodes[episode.Key()]; exists {
		ep.AddTitle(TitleDto{
			LocaleId: locale.Id(),
			Title:    episode.Title(),
		})
		return nil
	}

	var err error
	epcol.episodes[episode.Key()], err = newEpisode(NewEpisodeArgs{
		Number:       episode.Number(),
		SeasonNumber: episode.Season(),
		Thumbnail:    thumbnail,
		Titles: []TitleDto{
			{
				LocaleId: locale.Id(),
				Title:    episode.Title(),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// cleanEpisodes Removes any unused episodes in the collection.
func (epcol *EpisodeCollection) CleanEpisodes() {
	usedKeys := make(map[string]bool, len(epcol.latestSubs)+len(epcol.latestDubs))
	for _, key := range epcol.latestSubs {
		usedKeys[key] = true
	}

	for _, key := range epcol.latestDubs {
		usedKeys[key] = true
	}

	for epKey := range epcol.episodes {
		if _, exists := usedKeys[epKey]; !exists {
			delete(epcol.episodes, epKey)
		}
	}
}
