package anime

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/core"
	"strings"
	"time"
)

type AnimeDto struct {
	AnimeId     int
	SeriesId    string
	SlugTitle   string
	Title       string
	LastUpdated time.Time
}

// Anime is what we're all here for!
type Anime struct {
	animeId     AnimeId
	seriesId    core.SeriesId
	slugTitle   string
	title       string
	lastUpdated time.Time
	posters     []Image
	episodes    EpisodeCollection
	isDirty     bool
	isNew       bool
}

func NewAnime(
	dto AnimeDto,
	posters []Image,
	episodes EpisodeCollection,
) (Anime, error) {
	seriesId, err := core.NewSeriesId(dto.SeriesId)
	if err != nil {
		return Anime{}, err
	}

	if strings.Trim(dto.SlugTitle, " ") == "" {
		return Anime{}, fmt.Errorf("anime slug title cannot be blank")
	}

	if strings.Trim(dto.Title, " ") == "" {
		return Anime{}, fmt.Errorf("anime title cannot be blank")
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
		animeId:     0,
		seriesId:    seriesId,
		slugTitle:   dto.SlugTitle,
		title:       dto.Title,
		lastUpdated: time.Now().UTC(),
		posters:     posters,
		episodes:    episodes,
		isDirty:     true,
		isNew:       true,
	}, nil
}

func ReformAnime(
	dto AnimeDto,
	posters []Image,
	episodes EpisodeCollection,
) Anime {
	return Anime{
		animeId:     ReformAnimeId(dto.AnimeId),
		seriesId:    core.ReformSeriesId(dto.SeriesId),
		slugTitle:   dto.SlugTitle,
		title:       dto.Title,
		lastUpdated: dto.LastUpdated,
		posters:     posters,
		episodes:    episodes,
		isDirty:     false,
		isNew:       false,
	}
}

func (anime Anime) SeriesId() core.SeriesId {
	return anime.seriesId
}

func (anime Anime) Posters() []Image {
	return anime.posters
}

func (anime Anime) Dto() AnimeDto {
	return AnimeDto{
		AnimeId:     anime.animeId.Int(),
		SeriesId:    anime.seriesId.String(),
		SlugTitle:   anime.slugTitle,
		Title:       anime.title,
		LastUpdated: anime.lastUpdated,
	}
}
