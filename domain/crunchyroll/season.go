package crunchyroll

import "schoperation/crunchyrollanimestatus/domain/core"

type SeasonDto struct {
	Id              string
	Number          int
	SequenceNumber  int
	Keywords        []string
	Identifier      string
	SubtitleLocales []string
	Dubs            []DubDto
}

type Season struct {
	id              string
	number          int
	sequenceNumber  int
	keywords        []string
	identifier      string
	subtitleLocales map[core.Locale]bool
	dubs            map[core.Locale]Dub
}

func ReformSeason(dto SeasonDto) Season {
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

	return Season{
		id:              dto.Id,
		number:          dto.Number,
		sequenceNumber:  dto.SequenceNumber,
		keywords:        dto.Keywords,
		identifier:      dto.Identifier,
		subtitleLocales: subtitleLocales,
		dubs:            dubs,
	}
}

func (season Season) SequenceNumber() int {
	return season.sequenceNumber
}

func (season Season) Keywords() []string {
	return season.keywords
}

func (season Season) HasSubForLocale(locale core.Locale) bool {
	_, ok := season.subtitleLocales[locale]
	return ok
}

func (season Season) HasDubForLocale(locale core.Locale) bool {
	_, ok := season.dubs[locale]
	return ok
}
