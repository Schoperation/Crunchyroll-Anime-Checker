package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"schoperation/crunchyrollanimestatus/domain/core"
	"schoperation/crunchyrollanimestatus/domain/crunchyroll"
	"time"
)

// CrunchyrollClient represents a client capable of sending HTTP requests to Crunchyroll.
type CrunchyrollClient struct {
	credFilePath string
	accessToken  string
	lastLogin    time.Time
	cache        crunchyrollCache
}

type crunchyrollCreds struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	BasicAuthKey string `json:"basic_auth_key"`
}

func NewCrunchyrollClient(credFilePath string) CrunchyrollClient {
	return CrunchyrollClient{
		credFilePath: credFilePath,
		accessToken:  "",
		lastLogin:    time.Now().Add(time.Hour * -1),
		cache:        newCrunchyrollCache(),
	}
}

func (client *CrunchyrollClient) openCredFile() (crunchyrollCreds, error) {
	file, err := os.Open(client.credFilePath)
	if err != nil {
		return crunchyrollCreds{}, err
	}

	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return crunchyrollCreds{}, err
	}

	var creds crunchyrollCreds
	err = json.Unmarshal(bytes, &creds)
	if err != nil {
		return crunchyrollCreds{}, err
	}

	return creds, nil
}

func (client *CrunchyrollClient) Login() error {
	creds, err := client.openCredFile()
	if err != nil {
		return err
	}

	formValues := url.Values{}
	formValues.Set("username", creds.Username)
	formValues.Set("password", creds.Password)
	formValues.Set("grant_type", "password")
	formValues.Set("scope", "offline_access")

	request, err := http.NewRequest(http.MethodPost, "https://beta-api.crunchyroll.com/auth/v1/token", bytes.NewBufferString(formValues.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", creds.BasicAuthKey))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Crunchyroll/3.41.1 Android/1.0 okhttp/4.11.0")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode >= http.StatusBadRequest {
		if response.StatusCode == http.StatusNotAcceptable {
			return fmt.Errorf("could not get new token; being rate limited")
		} else {
			return fmt.Errorf("could not get new token; error code %d", response.StatusCode)
		}
	}

	var responseStruct struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
		Country      string `json:"country"`
		AccountId    string `json:"account_id"`
		ProfileId    string `json:"profile_id"`
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &responseStruct)
	if err != nil {
		return err
	}

	client.accessToken = responseStruct.AccessToken
	client.lastLogin = time.Now()
	return nil
}

func (client *CrunchyrollClient) GetAllAnime(locale string) ([]crunchyroll.AnimeDto, error) {
	var allAnime allAnimeResponse
	err := client.get("content/v2/discover/browse", &allAnime, map[string]string{
		"start":                    "0",
		"n":                        "2000",
		"type":                     "series",
		"sort_by":                  "alphabetical",
		"ratings":                  "true",
		"locale":                   locale,
		"preferred_audio_language": locale,
	})
	if err != nil {
		return nil, err
	}

	dtos := make([]crunchyroll.AnimeDto, len(allAnime.Data))
	for i, anime := range allAnime.Data {
		var tallPosters []crunchyroll.ImageDto
		if len(anime.Images.PosterTall) > 0 {
			tallPosters = make([]crunchyroll.ImageDto, len(anime.Images.PosterTall[0]))
			for j, image := range anime.Images.PosterTall[0] {
				tallPosters[j] = crunchyroll.ImageDto{
					Width:     image.Width,
					Height:    image.Height,
					ImageType: image.Type,
					Source:    image.Source,
				}
			}
		}

		var widePosters []crunchyroll.ImageDto
		if len(anime.Images.PosterWide) > 0 {
			widePosters = make([]crunchyroll.ImageDto, len(anime.Images.PosterWide[0]))
			for k, image := range anime.Images.PosterWide[0] {
				widePosters[k] = crunchyroll.ImageDto{
					Width:     image.Width,
					Height:    image.Height,
					ImageType: image.Type,
					Source:    image.Source,
				}
			}
		}

		dtos[i] = crunchyroll.AnimeDto{
			SeriesId:     anime.Id,
			SlugTitle:    anime.SlugTitle,
			Title:        anime.Title,
			New:          anime.New,
			LastUpdated:  anime.LastPublic,
			SeasonCount:  anime.SeriesMetaData.SeasonCount,
			EpisodeCount: anime.SeriesMetaData.EpisodeCount,
			TallPosters:  tallPosters,
			WidePosters:  widePosters,
		}
	}

	return dtos, nil
}

func (client *CrunchyrollClient) GetAllSeasonsBySeriesId(seriesId string) ([]crunchyroll.SeasonDto, error) {
	// Use Japanese as the preferred audio language, so we can get the full lists of subs and dubs.
	var seasonsResponse seasonsResponse
	err := client.get(fmt.Sprintf("content/v2/cms/series/%s/seasons", seriesId), &seasonsResponse, map[string]string{
		"locale":                   core.NewEnglishLocale().Name(),
		"preferred_audio_language": core.NewJapaneseLocale().Name(),
	})
	if err != nil {
		return nil, err
	}

	seasons := make([]crunchyroll.SeasonDto, len(seasonsResponse.Data))
	for i, season := range seasonsResponse.Data {
		dubs := make([]crunchyroll.DubDto, len(season.Versions))
		for j, version := range season.Versions {
			dubs[j] = crunchyroll.DubDto{
				AudioLocale: version.AudioLocale,
				GUID:        version.GUID,
				Original:    version.Original,
			}
		}

		seasons[i] = crunchyroll.SeasonDto{
			Id:              season.Id,
			Number:          season.SeasonNumber,
			SequenceNumber:  season.SeasonSequenceNumber,
			DisplayNumber:   season.SeasonDisplayNumber,
			Keywords:        season.Keywords,
			Identifier:      season.Identifier,
			IsSubbed:        season.IsSubbed,
			SubtitleLocales: season.SubtitleLocales,
			Dubs:            dubs,
		}
	}

	return seasons, nil
}

func (client *CrunchyrollClient) GetAllEpisodesBySeasonId(locale, seasonId string) ([]crunchyroll.EpisodeDto, error) {
	var episodesResponse episodesResponse
	if cachedResponse, ok := client.cache.GetEpisodesResponse(seasonId); ok {
		episodesResponse = cachedResponse
	} else {
		err := client.get(fmt.Sprintf("content/v2/cms/seasons/%s/episodes", seasonId), &episodesResponse, map[string]string{
			"locale": locale,
		})
		if err != nil {
			return nil, err
		}
	}

	episodes := make([]crunchyroll.EpisodeDto, len(episodesResponse.Data))
	for i, episode := range episodesResponse.Data {
		dubs := make([]crunchyroll.DubDto, len(episode.Versions))
		for j, version := range episode.Versions {
			dubs[j] = crunchyroll.DubDto{
				AudioLocale: version.AudioLocale,
				GUID:        version.GUID,
				Original:    version.Original,
			}
		}

		var thumbnails []crunchyroll.ImageDto
		if len(episode.Images.Thumbnail) > 0 {
			thumbnails = make([]crunchyroll.ImageDto, len(episode.Images.Thumbnail[0]))
			for k, image := range episode.Images.Thumbnail[0] {
				thumbnails[k] = crunchyroll.ImageDto{
					Width:     image.Width,
					Height:    image.Height,
					ImageType: image.Type,
					Source:    image.Source,
				}
			}
		}

		episodes[i] = crunchyroll.EpisodeDto{
			Number:          episode.Number,
			Season:          episode.SeasonNumber,
			Title:           episode.Title,
			SeasonId:        seasonId,
			IsSubbed:        episode.IsSubbed,
			SubtitleLocales: episode.SubtitleLocales,
			Dubs:            dubs,
			Thumbnails:      thumbnails,
		}
	}

	client.cache.SaveEpisodesResponse(seasonId, episodesResponse)
	return episodes, nil
}

func (client *CrunchyrollClient) get(path string, responseStruct any, queryParams map[string]string) error {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://beta-api.crunchyroll.com/%s", path), nil)
	if err != nil {
		return err
	}

	// Re-login if we believe the token is expired about now.
	// Crunchyroll does return 300 seconds for when it expires, but I've seen it expire before then, so I don't trust it...
	if time.Since(client.lastLogin).Seconds() > 180 {
		err = client.Login()
		if err != nil {
			return err
		}
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.accessToken))
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "Crunchyroll/3.41.1 Android/1.0 okhttp/4.11.0")

	values := url.Values{}
	for key, value := range queryParams {
		values.Set(key, value)
	}

	request.URL.RawQuery = values.Encode()

	backOffSchedule := []time.Duration{
		1 * time.Second,
		5 * time.Second,
		10 * time.Second,
	}

	var response *http.Response
	for _, backoff := range backOffSchedule {
		response, err = http.DefaultClient.Do(request)
		if err != nil {
			return err
		}

		if response.StatusCode == http.StatusOK {
			break
		} else if response.StatusCode == http.StatusGatewayTimeout || response.StatusCode == http.StatusServiceUnavailable || response.StatusCode == http.StatusRequestTimeout {
			time.Sleep(backoff)
			continue
		} else {
			return fmt.Errorf("failed to get response from crunchyroll; status code %d", response.StatusCode)
		}
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, responseStruct)
	if err != nil {
		return err
	}

	return nil
}
