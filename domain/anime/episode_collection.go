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

func NewEpisodeCollection(latestEpisodes []LatestEpisodes, thumbnails map[string]Image) (EpisodeCollection, error) {
	episodeCollection := EpisodeCollection{
		latestSubs: map[core.Locale]string{},
		latestDubs: map[core.Locale]string{},
		episodes:   map[string]Episode{},
	}

	for _, latestEpsForLocale := range latestEpisodes {
		episodeCollection.animeId = latestEpsForLocale.animeId

		thumbnail, exists := thumbnails[fmt.Sprintf("%d-%d", latestEpsForLocale.latestSub.Season(), latestEpsForLocale.latestSub.Number())]
		if !exists {
			return EpisodeCollection{}, fmt.Errorf("epcol: no thumbnail found for sub: S%dE%d", latestEpsForLocale.latestSub.Season(), latestEpsForLocale.latestSub.Number())
		}

		err := episodeCollection.AddSubForLocale(latestEpsForLocale.Locale(), latestEpsForLocale.latestSub, thumbnail)
		if err != nil {
			return EpisodeCollection{}, nil
		}

		thumbnail, exists = thumbnails[fmt.Sprintf("%d-%d", latestEpsForLocale.latestDub.Season(), latestEpsForLocale.latestDub.Number())]
		if !exists {
			return EpisodeCollection{}, fmt.Errorf("epcol: no thumbnail found for dub: S%dE%d", latestEpsForLocale.latestDub.Season(), latestEpsForLocale.latestDub.Number())
		}

		err = episodeCollection.AddDubForLocale(latestEpsForLocale.Locale(), latestEpsForLocale.latestDub, thumbnail)
		if err != nil {
			return EpisodeCollection{}, nil
		}
	}

	return episodeCollection, nil
}

func ReformEpisodeCollection(latestEpisodes []LatestEpisodes, thumbnails map[string]Image) EpisodeCollection {
	episodeCollection := EpisodeCollection{
		latestSubs: map[core.Locale]string{},
		latestDubs: map[core.Locale]string{},
		episodes:   map[string]Episode{},
	}

	for _, latestEpsForLocale := range latestEpisodes {
		episodeCollection.animeId = latestEpsForLocale.animeId

		thumbnail := thumbnails[fmt.Sprintf("%d-%d", latestEpsForLocale.latestSub.Season(), latestEpsForLocale.latestSub.Number())]
		_ = episodeCollection.AddSubForLocale(latestEpsForLocale.Locale(), latestEpsForLocale.latestSub, thumbnail)

		thumbnail = thumbnails[fmt.Sprintf("%d-%d", latestEpsForLocale.latestDub.Season(), latestEpsForLocale.latestDub.Number())]
		_ = episodeCollection.AddDubForLocale(latestEpsForLocale.Locale(), latestEpsForLocale.latestDub, thumbnail)
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

	epcol.latestSubs[locale] = fmt.Sprintf("%d-%d", sub.Season(), sub.Number())
	err := epcol.addEpisode(locale, sub, thumbnail)
	if err != nil {
		return err
	}

	// TODO delete old latest dub if applicable

	return nil
}

func (epcol *EpisodeCollection) AddDubForLocale(locale core.Locale, dub MinimalEpisode, thumbnail Image) error {
	if dub.IsBlank() {
		return nil
	}

	epcol.latestDubs[locale] = fmt.Sprintf("%d-%d", dub.Season(), dub.Number())
	err := epcol.addEpisode(locale, dub, thumbnail)
	if err != nil {
		return err
	}

	// TODO delete old latest dub if applicable

	return nil
}

func (epcol *EpisodeCollection) addEpisode(locale core.Locale, episode MinimalEpisode, thumbnail Image) error {
	key := fmt.Sprintf("%d-%d", episode.Season(), episode.Number())
	if ep, exists := epcol.episodes[key]; exists {
		ep.AddTitle(TitleDto{
			LocaleId: locale.Id(),
			Title:    episode.Title(),
		})
		return nil
	}

	var err error
	epcol.episodes[key], err = newEpisode(NewEpisodeArgs{
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
