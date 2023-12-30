package anime

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/core"
	"strings"
)

type TitleDto struct {
	LocaleId int
	Title    string
}

// TitleCollection is a collection of titles mapped to respective locales.
// It is used mostly for episodes that may be reused for multiple locales, but with different titles (because of different languages, duh).
// This way we won't need multiple Episodes just for different titles. Especially with all of those encoded thumbnails...
type TitleCollection struct {
	col map[core.Locale]string
}

func NewTitleCollection(dtos []TitleDto) (TitleCollection, error) {
	newCollection := TitleCollection{
		col: map[core.Locale]string{},
	}

	for _, dto := range dtos {
		err := newCollection.Add(dto)
		if err != nil {
			return TitleCollection{}, err
		}
	}

	return newCollection, nil
}

func ReformTitleCollection(dtos []TitleDto) TitleCollection {
	newCollection := TitleCollection{
		col: map[core.Locale]string{},
	}

	for _, dto := range dtos {
		newCollection.col[core.ReformLocaleFromId(dto.LocaleId)] = dto.Title
	}

	return newCollection
}

func (collection *TitleCollection) Add(dto TitleDto) error {
	locale, err := core.NewLocaleFromId(dto.LocaleId)
	if err != nil {
		return err
	}

	if strings.Trim(dto.Title, " ") == "" {
		return fmt.Errorf("title must not be blank in title collection")
	}

	collection.col[locale] = dto.Title
	return nil
}

func (collection TitleCollection) Title(locale core.Locale) string {
	title, ok := collection.col[locale]
	if !ok {
		return ""
	}

	return title
}
