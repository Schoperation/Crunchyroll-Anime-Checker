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

func (subcmd RefreshPostersSubCommand) Run(input RefreshPostersSubCommandInput) (RefreshPostersSubCommandOutput, error) {
	for _, updatedCrAnime := range input.UpdatedCrAnime {
		localAnime, exists := input.LocalAnime[updatedCrAnime.SeriesId()]
		if !exists {
			return RefreshPostersSubCommandOutput{}, fmt.Errorf("no local anime with series ID %s", updatedCrAnime.SeriesId())
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
				return RefreshPostersSubCommandOutput{}, err
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
				return RefreshPostersSubCommandOutput{}, err
			}

			newPosters[i] = newPoster
		}

		localAnime.UpdatePosters(newPosters)
		input.LocalAnime[updatedCrAnime.SeriesId()] = localAnime
	}

	newPosters := make(map[core.SeriesId][]anime.Image, len(input.NewCrAnime))
	for _, newCrAnime := range input.NewCrAnime {

		// TODO temp testing
		if newCrAnime.SeriesId().String() != "G1XHJV0KV" {
			continue
		}

		fmt.Printf("%s - %s\n", newCrAnime.SeriesId(), newCrAnime.SlugTitle())

		posters := make([]anime.Image, 2)

		encodedTallPoster, err := subcmd.encodedPosterFetcher.GetEncodedImageByURL(newCrAnime.TallPoster().Source())
		if err != nil {
			return RefreshPostersSubCommandOutput{}, err
		}

		encodedWidePoster, err := subcmd.encodedPosterFetcher.GetEncodedImageByURL(newCrAnime.WidePoster().Source())
		if err != nil {
			return RefreshPostersSubCommandOutput{}, err
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
			return RefreshPostersSubCommandOutput{}, err
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
			return RefreshPostersSubCommandOutput{}, err
		}

		newPosters[newCrAnime.SeriesId()] = posters
	}

	return RefreshPostersSubCommandOutput{
		UpdatedLocalAnime: input.LocalAnime,
		NewPosters:        newPosters,
	}, nil
}
