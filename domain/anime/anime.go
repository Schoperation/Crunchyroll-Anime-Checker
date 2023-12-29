package anime

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/core"
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

// Anime is what we're all here for!
type Anime struct {
	animeId          AnimeId
	seriesId         string
	slugTitle        string
	title            string
	lastUpdated      time.Time
	seasonIdentifier string
	posters          []Image
	episodes         EpisodeCollection
}

func NewAnime(
	dto AnimeDto,
	posters []Image,
	episodes EpisodeCollection,
) (Anime, error) {
	if strings.Trim(dto.SeriesId, " ") == "" {
		return Anime{}, fmt.Errorf("anime series ID cannot be blank")
	}

	if strings.Trim(dto.SlugTitle, " ") == "" {
		return Anime{}, fmt.Errorf("anime slug title cannot be blank")
	}

	if strings.Trim(dto.Title, " ") == "" {
		return Anime{}, fmt.Errorf("anime title cannot be blank")
	}

	if strings.Trim(dto.SeasonIdentifier, " ") == "" {
		return Anime{}, fmt.Errorf("anime season identifier cannot be blank")
	}

	if len(posters) != 2 {
		return Anime{}, fmt.Errorf("anime must have 2 posters")
	}

	hasPosterTall := false
	hasPosterWide := false
	for _, poster := range posters {
		switch poster.ImageType() {
		case core.ImageTypePosterTall:
			hasPosterTall = true
		case core.ImageTypePosterWide:
			hasPosterWide = true
		}
	}

	if !hasPosterTall {
		return Anime{}, fmt.Errorf("anime must have a tall poster")
	}

	if !hasPosterWide {
		return Anime{}, fmt.Errorf("anime must have a wide poster")
	}

	return Anime{
		animeId:          0,
		seriesId:         dto.SeriesId,
		slugTitle:        dto.SlugTitle,
		title:            dto.Title,
		lastUpdated:      time.Now().UTC(),
		seasonIdentifier: dto.SeasonIdentifier,
		posters:          posters,
		episodes:         episodes,
	}, nil
}

func ReformAnime(
	dto AnimeDto,
	posters []Image,
	episodes EpisodeCollection,
) Anime {
	return Anime{
		animeId:          ReformAnimeId(dto.AnimeId),
		seriesId:         dto.SeriesId,
		slugTitle:        dto.SlugTitle,
		title:            dto.Title,
		lastUpdated:      dto.LastUpdated,
		seasonIdentifier: dto.SeasonIdentifier,
		posters:          posters,
		episodes:         episodes,
	}
}

func (anime Anime) SeriesId() string {
	return anime.seriesId
}

func (anime Anime) Posters() []Image {
	return anime.posters
}

func (anime Anime) Dto() AnimeDto {
	return AnimeDto{
		AnimeId:          anime.animeId.Int(),
		SeriesId:         anime.seriesId,
		SlugTitle:        anime.slugTitle,
		Title:            anime.title,
		LastUpdated:      anime.lastUpdated,
		SeasonIdentifier: anime.seasonIdentifier,
	}
}
