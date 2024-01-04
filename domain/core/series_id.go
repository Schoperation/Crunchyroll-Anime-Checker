package core

import (
	"fmt"
	"strings"
)

type SeriesId string

func NewSeriesId(id string) (SeriesId, error) {
	if strings.Trim(id, " ") == "" {
		return "", fmt.Errorf("series ID cannot be blank")
	}

	if !strings.HasPrefix(id, "G") {
		return "", fmt.Errorf("series ID must start with G")
	}

	return SeriesId(id), nil
}

func ReformSeriesId(id string) SeriesId {
	return SeriesId(id)
}

func (id SeriesId) String() string {
	return string(id)
}
