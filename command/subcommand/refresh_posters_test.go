package subcommand

import (
	"errors"
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/core"
	"schoperation/crunchyrollanimestatus/domain/crunchyroll"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRefreshPostersSubCommand(t *testing.T) {
	seriesId := core.SeriesId("G20")
	testResources := refreshPostersSubCommandTestResources{}

	getEncodedImageTranslator := dummyGetEncodedImageTranslator{}
	subCommand := NewRefreshPostersSubCommand(getEncodedImageTranslator)

	testCases := []struct {
		name           string
		input          RefreshPostersSubCommandInput
		expectedOutput RefreshPostersSubCommandOutput
		expectedError  error
	}{
		{
			name: "with_updated_posters_returns_success",
			input: RefreshPostersSubCommandInput{
				NewCrAnime: nil,
				UpdatedCrAnime: []crunchyroll.Anime{
					testResources.getDummyCrunchyrollAnime("updated_posters"),
				},
				LocalAnime: map[core.SeriesId]anime.Anime{
					seriesId: testResources.getDummyLocalAnime("outdated_posters"),
				},
			},
			expectedOutput: RefreshPostersSubCommandOutput{
				UpdatedLocalAnime: map[core.SeriesId]anime.Anime{
					seriesId: testResources.getDummyLocalAnime("updated_posters"),
				},
				NewPosters: map[core.SeriesId][]anime.Image{},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output, err := subCommand.Run(tc.input)

			require.ErrorIs(t, errors.Unwrap(err), errors.Unwrap(tc.expectedError))
			require.EqualValues(t, tc.expectedOutput.NewPosters, output.NewPosters)
			require.EqualValues(t, tc.expectedOutput.UpdatedLocalAnime, output.UpdatedLocalAnime)
		})
	}
}

type dummyGetEncodedImageTranslator struct{}

func (translator dummyGetEncodedImageTranslator) GetEncodedImageByURL(url string) (string, error) {
	return "new picture", nil
}

type refreshPostersSubCommandTestResources struct{}

func (resources refreshPostersSubCommandTestResources) getDummyCrunchyrollAnime(id string) crunchyroll.Anime {
	crAnimeMap := map[string]crunchyroll.Anime{
		"updated_posters": crunchyroll.ReformAnime(crunchyroll.AnimeDto{
			SeriesId:  "G20",
			SlugTitle: "i-would-do-a-very-long-name-but-that-might-break-the-character-limit",
			Title:     "I Would Do a Very Long Name But That Might Break the Character Limit",
			TallPosters: []crunchyroll.ImageDto{
				{
					Width:     60,
					Height:    90,
					ImageType: core.ImageTypePosterTall.Name(),
					Source:    "http://www.example.com/newsourcetall",
				},
			},
			WidePosters: []crunchyroll.ImageDto{
				{
					Width:     320,
					Height:    180,
					ImageType: core.ImageTypePosterWide.Name(),
					Source:    "http://www.example.com/newsourcewide",
				},
			},
		}),
	}

	return crAnimeMap[id]
}

func (resources refreshPostersSubCommandTestResources) getDummyLocalAnime(id string) anime.Anime {
	localAnimeMap := map[string]anime.Anime{
		"outdated_posters": anime.ReformAnime(anime.AnimeDto{
			AnimeId:   1,
			SeriesId:  "G20",
			SlugTitle: "i-would-do-a-very-long-name-but-that-might-break-the-character-limit",
			Title:     "I Would Do a Very Long Name But That Might Break the Character Limit",
			IsDirty:   false,
		},
			[]anime.Image{
				anime.ReformImage(anime.ImageDto{
					AnimeId:   1,
					ImageType: core.ImageTypePosterTall.Int(),
					Url:       "http://www.example.com/oldsourcetall",
					Encoded:   "old picture",
					IsDirty:   false,
				}),
				anime.ReformImage(anime.ImageDto{
					AnimeId:   1,
					ImageType: core.ImageTypePosterWide.Int(),
					Url:       "http://www.example.com/oldsourcewide",
					Encoded:   "old picture",
					IsDirty:   false,
				}),
			},
			anime.EpisodeCollection{}),
		"updated_posters": anime.ReformAnime(anime.AnimeDto{
			AnimeId:   1,
			SeriesId:  "G20",
			SlugTitle: "i-would-do-a-very-long-name-but-that-might-break-the-character-limit",
			Title:     "I Would Do a Very Long Name But That Might Break the Character Limit",
			IsDirty:   true,
		},
			[]anime.Image{
				anime.ReformImage(anime.ImageDto{
					AnimeId:   1,
					ImageType: core.ImageTypePosterTall.Int(),
					Url:       "http://www.example.com/newsourcetall",
					Encoded:   "new picture",
					IsDirty:   true,
				}),
				anime.ReformImage(anime.ImageDto{
					AnimeId:   1,
					ImageType: core.ImageTypePosterWide.Int(),
					Url:       "http://www.example.com/newsourcewide",
					Encoded:   "new picture",
					IsDirty:   true,
				}),
			},
			anime.EpisodeCollection{}),
	}

	return localAnimeMap[id]
}
