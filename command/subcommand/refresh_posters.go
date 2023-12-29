package subcommand

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"
	"schoperation/crunchyrollanimestatus/domain/crunchyroll"
)

type RefreshPostersSubCommandInput struct {
	NewCrAnime     []crunchyroll.Anime
	UpdatedCrAnime map[string]crunchyroll.Anime
	SavedAnime     []anime.Anime
}

type RefreshPostersSubCommandOutput struct {
	UpdatedSavedAnime []anime.Anime
	NewPosters        map[string][]anime.Image
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
	for _, savedAnime := range input.SavedAnime {
		crAnime, exists := input.UpdatedCrAnime[savedAnime.SeriesId()]
		if !exists {
			return RefreshPostersSubCommandOutput{}, fmt.Errorf("couldn't match crunchyroll anime with saved anime: series ID %s", savedAnime.SeriesId())
		}

		for _, poster := range savedAnime.Posters() {
			posterUrl := ""
			switch poster.ImageType() {
			case core.ImageTypePosterTall:
				posterUrl = crAnime.TallPoster().Source()
			case core.ImageTypePosterWide:
				posterUrl = crAnime.WidePoster().Source()
			default:
				posterUrl = crAnime.TallPoster().Source()
			}

			if poster.Url() == posterUrl {
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
		}
	}

	newPosters := make(map[string][]anime.Image, len(input.NewCrAnime))
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
		UpdatedSavedAnime: input.SavedAnime,
		NewPosters:        newPosters,
	}, nil
}
