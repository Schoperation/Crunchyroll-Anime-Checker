package crunchyroll

import "schoperation/crunchyrollanimestatus/domain/core"

type SeasonDto struct {
	Id              string
	Number          int
	Identifier      string
	SubtitleLocales []string
	Dubs            []DubDto
}

type Season struct {
	id              string
	number          int
	identifier      string
	subtitleLocales map[core.Locale]bool
	dubs            map[core.Locale]Dub
}

func ReformSeason(dto SeasonDto) Season {
	subtitleLocales := map[core.Locale]bool{}
	for _, sub := range dto.SubtitleLocales {
		locale, err := core.NewLocaleByString(sub)
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

	return Season{
		id:              dto.Id,
		number:          dto.Number,
		identifier:      dto.Identifier,
		subtitleLocales: subtitleLocales,
		dubs:            dubs,
	}
}
