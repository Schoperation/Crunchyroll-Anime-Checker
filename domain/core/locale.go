package core

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
	id   int
	name string
}

var localeIds = map[string]int{
	"ja-jp": 1,
	"ko-kr": 2,
	"zh-cn": 3,
	"en-us": 4,
}

var localeNames = map[string]string{
	"ja-jp": "ja-JP",
	"ko-kr": "ko-KR",
	"zh-cn": "zh-CN",
	"en-us": "en-US",
}

func NewLocale(locale string) (Locale, error) {
	localeName, ok := localeNames[strings.ToLower(locale)]
	if !ok {
		return Locale{}, fmt.Errorf("could not parse locale %s", locale)
	}

	return Locale{
		id:   localeIds[strings.ToLower(locale)],
		name: localeName,
	}, nil
}

func ReformLocale(locale string) Locale {
	return Locale{
		id:   localeIds[strings.ToLower(locale)],
		name: localeNames[strings.ToLower(locale)],
	}
}

func (l Locale) Id() int {
	return l.id
}

func (l Locale) Name() string {
	return l.name
}

func NewEnglishLocale() Locale {
	return Locale{
		id:   localeIds["en-us"],
		name: localeNames["en-us"],
	}
}

func NewJapaneseLocale() Locale {
	return Locale{
		id:   localeIds["ja-jp"],
		name: localeNames["ja-jp"],
	}
}
