package subcommand

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"
	"schoperation/crunchyrollanimestatus/domain/crunchyroll"
)

type RefreshPostersSubCommandInput struct {
	NewCrAnime     []crunchyroll.Anime
	UpdatedCrAnime []crunchyroll.Anime
	LocalAnime     map[core.SeriesId]anime.Anime
}

type RefreshPostersSubCommandOutput struct {
	UpdatedLocalAnime map[core.SeriesId]anime.Anime
	NewPosters        map[core.SeriesId][]anime.Image
}

type encodedPosterFetcher interface {
	GetEncodedImageByURL(url string) (string, error)
}

type RefreshPostersSubCommand struct {
	encodedPosterFetcher encodedPosterFetcher
}

func NewRefreshPostersSubCommand(
	encodedPosterFetcher encodedPosterFetcher,
) RefreshPostersSubCommand {
	return RefreshPostersSubCommand{
		encodedPosterFetcher: encodedPosterFetcher,
	}
}

func (subcmd RefreshPostersSubCommand) Run(input RefreshPostersSubCommandInput) (RefreshPostersSubCommandOutput, map[core.SeriesId]error) {
	errors := map[core.SeriesId]error{}

	for _, updatedCrAnime := range input.UpdatedCrAnime {
		localAnime, exists := input.LocalAnime[updatedCrAnime.SeriesId()]
		if !exists {
			errors[updatedCrAnime.SeriesId()] = fmt.Errorf("no local anime found")
			continue
		}

		newPosters := make([]anime.Image, anime.NumPostersPerAnime)
		for i, poster := range localAnime.Posters() {
			posterUrl := ""
			switch poster.ImageType() {
			case core.ImageTypePosterTall:
				posterUrl = updatedCrAnime.TallPoster().Source()
			case core.ImageTypePosterWide:
				posterUrl = updatedCrAnime.WidePoster().Source()
			default:
				posterUrl = updatedCrAnime.TallPoster().Source()
			}

			if poster.Url() == posterUrl {
				newPosters[i] = poster
				continue
			}

			newEncodedImage, err := subcmd.encodedPosterFetcher.GetEncodedImageByURL(posterUrl)
			if err != nil {
				errors[updatedCrAnime.SeriesId()] = err
				break
			}

			newPoster, err := anime.NewImage(anime.ImageDto{
				AnimeId:       poster.AnimeId().Int(),
				ImageType:     poster.ImageType().Int(),
				SeasonNumber:  0,
				EpisodeNumber: 0,
				Url:           posterUrl,
				Encoded:       newEncodedImage,
			})
			if err != nil {
				errors[updatedCrAnime.SeriesId()] = err
				break
			}

			newPosters[i] = newPoster
		}

		if _, errored := errors[updatedCrAnime.SeriesId()]; errored {
			continue
		}

		localAnime.UpdatePosters(newPosters)
		input.LocalAnime[updatedCrAnime.SeriesId()] = localAnime
	}

	newPosters := make(map[core.SeriesId][]anime.Image, len(input.NewCrAnime))
	for _, newCrAnime := range input.NewCrAnime {
		posters := make([]anime.Image, anime.NumPostersPerAnime)

		encodedTallPoster, err := subcmd.encodedPosterFetcher.GetEncodedImageByURL(newCrAnime.TallPoster().Source())
		if err != nil {
			errors[newCrAnime.SeriesId()] = err
			continue
		}

		encodedWidePoster, err := subcmd.encodedPosterFetcher.GetEncodedImageByURL(newCrAnime.WidePoster().Source())
		if err != nil {
			errors[newCrAnime.SeriesId()] = err
			continue
		}

		posters[0], err = anime.NewImage(anime.ImageDto{
			AnimeId:       0,
			ImageType:     core.ImageTypePosterTall.Int(),
			SeasonNumber:  0,
			EpisodeNumber: 0,
			Url:           newCrAnime.TallPoster().Source(),
			Encoded:       encodedTallPoster,
		})
		if err != nil {
			errors[newCrAnime.SeriesId()] = err
			continue
		}

		posters[1], err = anime.NewImage(anime.ImageDto{
			AnimeId:       0,
			ImageType:     core.ImageTypePosterWide.Int(),
			SeasonNumber:  0,
			EpisodeNumber: 0,
			Url:           newCrAnime.WidePoster().Source(),
			Encoded:       encodedWidePoster,
		})
		if err != nil {
			errors[newCrAnime.SeriesId()] = err
			continue
		}

		newPosters[newCrAnime.SeriesId()] = posters
	}

	return RefreshPostersSubCommandOutput{
		UpdatedLocalAnime: input.LocalAnime,
		NewPosters:        newPosters,
	}, errors
}
