package crunchyroll

import (
	"strings"
	"time"
)

type AnimeDto struct {
	SeriesId     string
	SlugTitle    string
	Title        string
	LastUpdated  time.Time
	SeasonCount  int
	EpisodeCount int
	TallPosters  []ImageDto
	WidePosters  []ImageDto
}

type Anime struct {
	seriesId    string
	slugTitle   string
	title       string
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
		seriesId:    dto.SeriesId,
		slugTitle:   dto.SlugTitle,
		title:       dto.Title,
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
	if dto.SeasonCount == 1 && dto.EpisodeCount == 1 {
		return false
	}

	// Or... the slug ends in -movies
	if strings.HasSuffix(dto.SlugTitle, "-movies") || strings.HasSuffix(dto.SlugTitle, "-movie") {
		return false
	}

	// Try not to include OVAs; since they're basically one-time
	if strings.HasSuffix(dto.SlugTitle, "-ova") {
		return false
	}

	return true
}
