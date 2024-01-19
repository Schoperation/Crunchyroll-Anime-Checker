package anime

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/core"
	"strings"
	"time"
)

const NumPostersPerAnime = 2

type AnimeDto struct {
	AnimeId     int
	SeriesId    string
	SlugTitle   string
	Title       string
	IsSimulcast bool
	LastUpdated time.Time

	// Used in testing. Ignore otherwise.
	IsDirty bool
}

// Anime is what we're all here for!
type Anime struct {
	animeId     AnimeId
	seriesId    core.SeriesId
	slugTitle   string
	title       string
	isSimulcast bool
	lastUpdated time.Time
	posters     []Image
	episodes    EpisodeCollection
	isDirty     bool
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

	err = validateAnimePosters(posters)
	if err != nil {
		return Anime{}, err
	}

	return Anime{
		animeId:     0,
		seriesId:    seriesId,
		slugTitle:   dto.SlugTitle,
		title:       dto.Title,
		isSimulcast: dto.IsSimulcast,
		lastUpdated: now(),
		posters:     posters,
		episodes:    episodes,
		isDirty:     true,
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
		isSimulcast: dto.IsSimulcast,
		lastUpdated: dto.LastUpdated,
		posters:     posters,
		episodes:    episodes,
		isDirty:     dto.IsDirty,
	}
}

func (anime *Anime) AnimeId() AnimeId {
	return anime.animeId
}

func (anime *Anime) SeriesId() core.SeriesId {
	return anime.seriesId
}

func (anime *Anime) Posters() []Image {
	return anime.posters
}

func (anime *Anime) Episodes() *EpisodeCollection {
	return &anime.episodes
}

func (anime *Anime) IsDirty() bool {
	return anime.isDirty
}

func (anime *Anime) Dto() AnimeDto {
	return AnimeDto{
		AnimeId:     anime.animeId.Int(),
		SeriesId:    anime.seriesId.String(),
		SlugTitle:   anime.slugTitle,
		Title:       anime.title,
		LastUpdated: anime.lastUpdated,
	}
}

func (anime *Anime) SetDirty() {
	anime.isDirty = true
	anime.lastUpdated = now()
}

func (anime *Anime) AssignAnimeId(animeId AnimeId) {
	if !anime.animeId.IsZero() {
		return
	}

	if animeId.IsZero() {
		return
	}

	anime.animeId = animeId
}

func (anime *Anime) UpdatePosters(newPosters []Image) error {
	err := validateAnimePosters(newPosters)
	if err != nil {
		return err
	}

	anime.posters = newPosters
	anime.SetDirty()
	return nil
}

func validateAnimePosters(posters []Image) error {
	if len(posters) != NumPostersPerAnime {
		return fmt.Errorf("anime must have %d posters", NumPostersPerAnime)
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
		return fmt.Errorf("anime must have a tall poster")
	}

	if !hasPosterWide {
		return fmt.Errorf("anime must have a wide poster")
	}

	return nil
}
