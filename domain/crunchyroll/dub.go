package crunchyroll

import "schoperation/crunchyrollanimestatus/domain/core"

// Representation of a "version" from Crunchyroll, primarily used to store dub info.
type DubDto struct {
	AudioLocale string
	GUID        string
	Original    bool // Usually identifies the Japanese version.
}

type Dub struct {
	audioLocale core.Locale
	seasonId    string
	original    bool
}

func ReformDub(dto DubDto) Dub {
	locale, err := core.NewLocaleByString(dto.AudioLocale)
	if err != nil {
		return Dub{}
	}

	return Dub{
		audioLocale: locale,
		seasonId:    dto.GUID,
		original:    dto.Original,
	}
}

func (dub Dub) Locale() core.Locale {
	return dub.audioLocale
}

func (dub Dub) SeasonId() string {
	return dub.seasonId
}

func (dub Dub) IsOriginal() bool {
	return dub.original
}
