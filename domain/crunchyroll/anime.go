package crunchyroll

import (
	"schoperation/crunchyrollanimestatus/domain/core"
	"slices"
	"strings"
	"time"
)

type AnimeDto struct {
	SeriesId    string
	SlugTitle   string
	Title       string
	IsSimulcast bool
	LastUpdated time.Time
	TallPosters []ImageDto
	WidePosters []ImageDto

	New          bool
	SeasonCount  int
	EpisodeCount int
}

type Anime struct {
	seriesId    core.SeriesId
	slugTitle   string
	title       string
	isSimulcast bool
	lastUpdated time.Time
	tallPoster  Image
	widePoster  Image
}

func ReformAnime(dto AnimeDto) Anime {
	tallPoster := ImageDto{}
	for _, poster := range dto.TallPosters {
		if poster.Width == 60 && poster.Height == 90 {
			tallPoster = poster
			break
		}
	}

	widePoster := ImageDto{}
	for _, poster := range dto.WidePosters {
		if poster.Width == 320 && poster.Height == 180 {
			widePoster = poster
			break
		}
	}

	return Anime{
		seriesId:    core.ReformSeriesId(dto.SeriesId),
		slugTitle:   dto.SlugTitle,
		title:       dto.Title,
		isSimulcast: dto.IsSimulcast,
		lastUpdated: dto.LastUpdated,
		tallPoster:  ReformImage(tallPoster),
		widePoster:  ReformImage(widePoster),
	}
}

func ReformAnimeCollection(dtos []AnimeDto) []Anime {
	animes := []Anime{}
	for _, dto := range dtos {
		if !shouldAddAnime(dto) {
			continue
		}

		animes = append(animes, ReformAnime(dto))
	}

	return animes
}

func shouldAddAnime(dto AnimeDto) bool {
	// Sometimes Crunchyroll marks a movie as a series. Lovely...
	// Usually they're one season with one episode.
	// Of course, this could also be a new show...
	if dto.SeasonCount == 1 && dto.EpisodeCount == 1 && !dto.New {
		return false
	}

	// Or... the slug ends in -movies.
	if strings.HasSuffix(dto.SlugTitle, "-movies") || strings.HasSuffix(dto.SlugTitle, "-movie") {
		return false
	}

	// Don't include OVAs since they're either one-time or aren't on a consistent schedule.
	// Plus they're a bit broken in this program; may implement later on if requested.
	if strings.HasSuffix(dto.SlugTitle, "-ova") {
		return false
	}

	// Specific exceptions that are either broken right now, or are pretty old and suck anyway.
	animeInQuotationMarks := []string{
		"G6EXH7VKM", // anifile (bruh)
		"GRG5HJN5W", // otalku (double bruh)
		"G6WE4W0N6", // chinese (acts weird with subtitles)
		"GRWEMGNER", // ""
		"GRP85E0MR", // ""
		"GRVND1G3Y", // ""
		"G1XHJV2PZ", // ""
		"G5PHNM77K", // ovas that slipped through
		"G5PHNM77K", // ""
		"G79H23XN2", // ""
		"G24H1NZXD", // ""
	}

	if slices.Contains(animeInQuotationMarks, dto.SeriesId) {
		return false
	}

	return true
}

func (anime Anime) SeriesId() core.SeriesId {
	return anime.seriesId
}

func (anime Anime) SlugTitle() string {
	return anime.slugTitle
}

func (anime Anime) Title() string {
	return anime.title
}

func (anime Anime) IsSimulcast() bool {
	return anime.isSimulcast
}

func (anime Anime) LastUpdated() time.Time {
	return anime.lastUpdated
}

func (anime Anime) TallPoster() Image {
	return anime.tallPoster
}

func (anime Anime) WidePoster() Image {
	return anime.widePoster
}
