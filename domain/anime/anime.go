package anime

import (
	"fmt"
	"strings"
	"time"
)

type AnimeDto struct {
	AnimeId          int
	SeriesId         string
	SlugTitle        string
	Title            string
	LastUpdated      time.Time
	SeasonIdentifier string
}

type Anime struct {
	animeId          int
	seriesId         string
	slugTitle        string
	title            string
	lastUpdated      time.Time
	seasonIdentifier string
	posterTall       Image
	posterWide       Image
	episodes         EpisodeCollection
}

func NewAnime(
	dto AnimeDto,
	posters []Image,
	episodes EpisodeCollection,
) (Anime, error) {
	if strings.Trim(dto.SeriesId, " ") == "" {
		return Anime{}, fmt.Errorf("series ID cannot be blank")
	}

	if strings.Trim(dto.SlugTitle, " ") == "" {
		return Anime{}, fmt.Errorf("slug title cannot be blank")
	}

	if strings.Trim(dto.Title, " ") == "" {
		return Anime{}, fmt.Errorf("title cannot be blank")
	}

	if strings.Trim(dto.SeasonIdentifier, " ") == "" {
		return Anime{}, fmt.Errorf("season identifier cannot be blank")
	}

	return Anime{
		animeId:          0,
		seriesId:         dto.SeriesId,
		slugTitle:        dto.SlugTitle,
		title:            dto.Title,
		lastUpdated:      time.Now().UTC(),
		seasonIdentifier: dto.SeasonIdentifier,
	}, nil
}

func ReformAnime(dto AnimeDto) Anime {
	return Anime{
		animeId:          dto.AnimeId,
		seriesId:         dto.SeriesId,
		slugTitle:        dto.SlugTitle,
		title:            dto.Title,
		lastUpdated:      dto.LastUpdated,
		seasonIdentifier: dto.SeasonIdentifier,
	}
}

func (anime Anime) Dto() AnimeDto {
	return AnimeDto{
		AnimeId:          anime.animeId,
		SeriesId:         anime.seriesId,
		SlugTitle:        anime.slugTitle,
		Title:            anime.title,
		LastUpdated:      anime.lastUpdated,
		SeasonIdentifier: anime.seasonIdentifier,
	}
}
