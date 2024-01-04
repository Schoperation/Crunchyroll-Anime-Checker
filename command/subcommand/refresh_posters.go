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
	for _, crAnime := range input.UpdatedCrAnime {
		localAnime, exists := input.LocalAnime[crAnime.SeriesId()]
		if !exists {
			return RefreshPostersSubCommandOutput{}, fmt.Errorf("couldn't match crunchyroll anime with saved anime: series ID %s", localAnime.SeriesId())
		}

		for _, poster := range localAnime.Posters() {
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

		input.LocalAnime[crAnime.SeriesId()] = localAnime
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
