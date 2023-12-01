package crunchyroll

import (
	"fmt"
	"strings"
)

/*
Locale represents a language supported by Crunchyroll and this project.
Used in various calls to ensure we get the correct data for a particular language.
TODO Need to support more than English.
*/
type Locale struct {
	name string
}

var locales = map[string]string{
	"ja-jp": "ja-JP",
	"ko-kr": "ko-KR",
	"zh-cn": "zh-CN",
	"en-us": "en-US",
}

func NewLocale(locale string) (Locale, error) {
	localeName, ok := locales[strings.ToLower(locale)]
	if !ok {
		return Locale{}, fmt.Errorf("could not parse locale %s", locale)
	}

	return Locale{
		name: localeName,
	}, nil
}

func (l Locale) Name() string {
	return l.name
}

// Specific constructors for original languages of anime (or media, since now we're going outside of Japan)

func NewJapaneseLocale() Locale {
	return Locale{
		name: locales["ja-jp"],
	}
}

func NewKoreanLocale() Locale {
	return Locale{
		name: locales["ko-kr"],
	}
}

func NewChineseLocale() Locale {
	return Locale{
		name: locales["zh-cn"],
	}
}
