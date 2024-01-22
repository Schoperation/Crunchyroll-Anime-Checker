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
type Locale int

const (
	LocaleJaJP = 1
	LocaleKoKR = 2
	LocaleZhCN = 3
	LocaleEnUS = 4
)

var localeNames = map[int]string{
	LocaleJaJP: "ja-JP",
	LocaleKoKR: "ko-KR",
	LocaleZhCN: "zh-CN",
	LocaleEnUS: "en-US",
}

func NewLocaleFromId(localeId int) (Locale, error) {
	_, ok := localeNames[localeId]
	if !ok {
		return 0, fmt.Errorf("could not parse locale id %d", localeId)
	}

	return Locale(localeId), nil
}

func NewLocaleFromString(localeString string) (Locale, error) {
	for locale, name := range localeNames {
		if strings.EqualFold(name, localeString) {
			return Locale(locale), nil
		}
	}

	return 0, fmt.Errorf("could not parse locale %s", localeString)
}

func ReformLocaleFromId(localeId int) Locale {
	return Locale(localeId)
}

func (l Locale) Id() int {
	return int(l)
}

func (l Locale) Name() string {
	return localeNames[int(l)]
}

func NewJapaneseLocale() Locale {
	return ReformLocaleFromId(LocaleJaJP)
}

func NewEnglishLocale() Locale {
	return ReformLocaleFromId(LocaleEnUS)
}

func SupportedLocales() []Locale {
	locales := make([]Locale, len(localeNames)-3)
	i := 0

	for id := range localeNames {
		if id < LocaleEnUS {
			continue
		}

		locales[i] = ReformLocaleFromId(id)
		i++
	}

	return locales
}
