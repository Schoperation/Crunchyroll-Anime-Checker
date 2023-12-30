package crunchyroll

import "schoperation/crunchyrollanimestatus/domain/core"

type EpisodeDto struct {
	Number          int
	SeasonId        string
	SubtitleLocales []string
	Dubs            []DubDto
	Thumbnails      []ImageDto
}

type Episode struct {
	number          int
	seasonId        string
	subtitleLocales map[core.Locale]bool
	dubs            map[core.Locale]Dub
	thumbnail       Image
}

func ReformEpisode(dto EpisodeDto) Episode {
	subtitleLocales := map[core.Locale]bool{}
	for _, sub := range dto.SubtitleLocales {
		locale, err := core.NewLocaleFromString(sub)
		if err != nil {
			continue
		}

		subtitleLocales[locale] = true
	}

	dubs := map[core.Locale]Dub{}
	for _, dubDto := range dto.Dubs {
		dub := ReformDub(dubDto)
		if dub.SeasonId() == "" {
			continue
		}

		dubs[dub.Locale()] = dub
	}

	smallestThumbanil := ImageDto{}
	for _, thumbnail := range dto.Thumbnails {
		if thumbnail.Width == 320 && thumbnail.Height == 180 {
			smallestThumbanil = thumbnail
			break
		}
	}

	return Episode{
		number:          dto.Number,
		seasonId:        dto.SeasonId,
		subtitleLocales: subtitleLocales,
		dubs:            dubs,
		thumbnail:       ReformImage(smallestThumbanil),
	}
}

func (episode Episode) Number() int {
	return episode.number
}

func (episode Episode) HasSubForLocale(locale core.Locale) bool {
	_, ok := episode.subtitleLocales[locale]
	return ok
}

func (episode Episode) HasDubForLocale(locale core.Locale) bool {
	_, ok := episode.dubs[locale]
	return ok
}
