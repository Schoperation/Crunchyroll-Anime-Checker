package crunchyroll

import "schoperation/crunchyrollanimestatus/domain/core"

type EpisodeDto struct {
	Number          int
	Season          int
	Title           string
	SeasonId        string
	SubtitleLocales []string
	Dubs            []DubDto
	Thumbnails      []ImageDto
}

type Episode struct {
	number          int
	season          int
	title           string
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

	smallestThumbnail := ImageDto{}
	for _, thumbnail := range dto.Thumbnails {
		if thumbnail.Width == 320 && thumbnail.Height == 180 {
			smallestThumbnail = thumbnail
			break
		}
	}

	return Episode{
		number:          dto.Number,
		title:           dto.Title,
		season:          dto.Season,
		seasonId:        dto.SeasonId,
		subtitleLocales: subtitleLocales,
		dubs:            dubs,
		thumbnail:       ReformImage(smallestThumbnail),
	}
}

func (episode Episode) Number() int {
	return episode.number
}

func (episode Episode) Season() int {
	return episode.season
}

func (episode Episode) Title() string {
	return episode.title
}

func (episode Episode) Thumbnail() Image {
	return episode.thumbnail
}

func (episode Episode) hasSubForLocale(locale core.Locale) bool {
	_, ok := episode.subtitleLocales[locale]
	return ok
}

func (episode Episode) hasDubForLocale(locale core.Locale) bool {
	_, ok := episode.dubs[locale]
	return ok
}
