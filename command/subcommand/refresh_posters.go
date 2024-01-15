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

type getEncodedImageTranslator interface {
	GetEncodedImageByURL(url string) (string, error)
}

type RefreshPostersSubCommand struct {
	getEncodedImageTranslator getEncodedImageTranslator
}

func NewRefreshPostersSubCommand(
	getEncodedImageTranslator getEncodedImageTranslator,
) RefreshPostersSubCommand {
	return RefreshPostersSubCommand{
		getEncodedImageTranslator: getEncodedImageTranslator,
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

			newEncodedImage, err := subcmd.getEncodedImageTranslator.GetEncodedImageByURL(posterUrl)
			if err != nil {
				return RefreshPostersSubCommandOutput{}, err
			}

			err = poster.UpdatePoster(posterUrl, newEncodedImage)
			if err != nil {
				return RefreshPostersSubCommandOutput{}, err
			}

			newPosters[i] = poster
		}

		localAnime.UpdatePosters(newPosters)
		input.LocalAnime[updatedCrAnime.SeriesId()] = localAnime
	}

	newPosters := make(map[core.SeriesId][]anime.Image, len(input.NewCrAnime))
	for _, newCrAnime := range input.NewCrAnime {
		posters := make([]anime.Image, 2)

		encodedTallPoster, err := subcmd.getEncodedImageTranslator.GetEncodedImageByURL(newCrAnime.TallPoster().Source())
		if err != nil {
			return RefreshPostersSubCommandOutput{}, err
		}

		encodedWidePoster, err := subcmd.getEncodedImageTranslator.GetEncodedImageByURL(newCrAnime.WidePoster().Source())
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
