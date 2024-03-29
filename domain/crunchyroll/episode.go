package crunchyroll

import (
	"schoperation/crunchyroll-anime-checker/domain/core"
	"strings"
)

type EpisodeDto struct {
	Number          int
	Season          int
	Title           string
	SeasonId        string
	IsSubbed        bool
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
	// Sometimes CR leaves out subtitle locales (e.g. A Certain Magical Index, Arawaka Under The Bridge) even those it's subbed in English.
	// If IsSubbed is true then we can assume there are English subtitles; most likely embedded into the episodes themselves.
	if dto.IsSubbed && len(dto.SubtitleLocales) == 0 {
		dto.SubtitleLocales = []string{core.NewEnglishLocale().Name()}
	}

	// Sometimes CR leaves out the episode number if it's some bizarre decima (e.g. BURN THE WITCH)...
	if dto.Number <= 0 {
		dto.Number = 1
	}

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

	if strings.Trim(smallestThumbnail.Source, " ") == "" {
		smallestThumbnail.Source = core.DefaultPosterUrl
		smallestThumbnail.ImageType = core.ImageTypeThumbnail.Name()
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
