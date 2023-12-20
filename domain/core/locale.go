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

var locales = map[int]string{
	1: "ja-JP",
	2: "ko-KR",
	3: "zh-CN",
	4: "en-US",
}

func NewLocale(localeId int) (Locale, error) {
	localeName, ok := locales[localeId]
	if !ok {
		return Locale{}, fmt.Errorf("could not parse locale id %d", localeId)
	}

	return Locale{
		id:   localeId,
		name: localeName,
	}, nil
}

func NewLocaleByString(localeString string) (Locale, error) {
	for i, name := range locales {
		if strings.EqualFold(name, localeString) {
			return Locale{
				id:   i,
				name: name,
			}, nil
		}
	}

	return Locale{}, fmt.Errorf("could not parse locale %s", localeString)
}

func ReformLocale(localeId int) Locale {
	return Locale{
		id:   localeId,
		name: locales[localeId],
	}
}

func (l Locale) Id() int {
	return l.id
}

func (l Locale) Name() string {
	return l.name
}

func NewEnglishLocale() Locale {
	return ReformLocale(4)
}

func NewJapaneseLocale() Locale {
	return ReformLocale(1)
}
