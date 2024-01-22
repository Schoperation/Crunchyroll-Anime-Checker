package script

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"schoperation/crunchyroll-anime-checker/domain/core"
	"strconv"
	"strings"
	"time"
)

type RefreshAnimeCmd struct {
	name string
}

func NewRefreshAnimeCmd(name string) Command {
	return RefreshAnimeCmd{
		name: name,
	}
}

func (cmd RefreshAnimeCmd) Name() string {
	return cmd.name
}

func (cmd RefreshAnimeCmd) Run(client CrunchyrollClient) error {
	fmt.Printf("Retrieving all anime from Crunchyroll...\n")
	startTime := time.Now()

	var allAnime allAnimeResponse
	err := client.GetWithQueryParams("content/v2/discover/browse", &allAnime, map[string]string{
		"start":                    "0",
		"n":                        "2000",
		"type":                     "series",
		"sort_by":                  "alphabetical",
		"ratings":                  "true",
		"locale":                   client.locale.Name(),
		"preferred_audio_language": client.locale.Name(),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Took %s\n", time.Since(startTime).String())
	fmt.Printf("Refreshing Sensei CSV...\n")
	startTime = time.Now()

	err = cmd.refreshAnimeSenseiList(client, allAnime)
	if err != nil {
		return err
	}

	fmt.Printf("Took %s\n", time.Since(startTime).String())
	fmt.Printf("Refreshing Posters...\n")
	startTime = time.Now()

	err = cmd.refreshAnimePosters(client, allAnime)
	if err != nil {
		return err
	}

	/*
		fmt.Printf("Took %s\n", time.Since(startTime).String())
		fmt.Printf("Refreshing %s Atlas...\n", client.locale.Name())
		startTime = time.Now()

		err = cmd.refreshAnimeAtlas(client, allAnime)
		if err != nil {
			return err
		}
	*/

	fmt.Printf("Took %s\n", time.Since(startTime).String())

	return nil
}

// This refreshes the CSV file which is used to populate the dropdown in the app's config. And as a smaller master list without any additional info.
func (cmd RefreshAnimeCmd) refreshAnimeSenseiList(client CrunchyrollClient, allAnime allAnimeResponse) error {
	newList := [][]string{}
	for _, series := range allAnime.Data {
		if !cmd.shouldAddSeries(series) {
			continue
		}
		// series_id, slug_title, title
		newList = append(newList, []string{
			series.Id,
			series.SlugTitle,
			series.Title,
		})
	}

	newAnimeSenseiList, err := os.Create(fmt.Sprintf("%s/anime_sensei_list_new.csv", client.listsPath))
	if err != nil {
		return err
	}

	newListAsCsv := csv.NewWriter(newAnimeSenseiList)
	newListAsCsv.Comma = '|'
	newListAsCsv.Write([]string{"series_id", "slug_title", "title"})
	newListAsCsv.WriteAll(newList)
	newAnimeSenseiList.Close()

	err = os.Rename(fmt.Sprintf("%s/anime_sensei_list_new.csv", client.listsPath), fmt.Sprintf("%s/anime_sensei_list.csv", client.listsPath))
	if err != nil {
		return err
	}

	return nil
}

func (cmd RefreshAnimeCmd) refreshAnimeAtlas(client CrunchyrollClient, allAnime allAnimeResponse) error {
	animeAtlasFile, err := os.Open(fmt.Sprintf("%s/anime_atlas_%s.json", client.listsPath, client.locale.Name()))
	if err != nil {
		return err
	}

	defer animeAtlasFile.Close()

	bytes, err := io.ReadAll(animeAtlasFile)
	if err != nil {
		return err
	}

	var animeAtlas AnimeAtlas
	err = json.Unmarshal(bytes, &animeAtlas)
	if err != nil {
		return err
	}

	isDirty := false
	totalCount := 0
	errors := make(map[string]error)
	emptyAnime := []string{}
	for _, series := range allAnime.Data {
		if !cmd.shouldAddSeries(series) {
			delete(animeAtlas.Anime, series.SlugTitle)
			continue
		}

		if animeEntry, ok := animeAtlas.Anime[series.SlugTitle]; ok {
			if animeEntry.LastUpdated.After(series.LastPublic) {
				continue
			}
		}

		anime, err := cmd.fetchAnime(client, series, client.locale)
		if err != nil {
			errors[series.Id] = err
			if strings.Contains(err.Error(), "406") {
				fmt.Println("Being rate-limited, stopping...")
				break
			}
			continue
		}

		if anime.Sub.Season == 0 && anime.Sub.Episode == 0 &&
			anime.Dub.Season == 0 && anime.Dub.Episode == 0 {
			emptyAnime = append(emptyAnime, series.SlugTitle)
		}

		animeAtlas.Anime[series.SlugTitle] = anime
		isDirty = true
		totalCount++

		// Avoid rate limiting
		time.Sleep(5 * time.Second)
	}

	animeAtlas.TotalCount = totalCount
	fmt.Println("The following anime were empty: ", emptyAnime)
	fmt.Println("Errors: ", errors)

	if !isDirty {
		return nil
	}

	newAtlasFile, err := os.Create(fmt.Sprintf("%s/anime_atlas_%s_new.json", client.listsPath, client.locale.Name()))
	if err != nil {
		return err
	}

	newBytes, err := json.MarshalIndent(animeAtlas, "", "    ")
	if err != nil {
		return err
	}

	_, err = newAtlasFile.Write(newBytes)
	if err != nil {
		return err
	}

	err = os.Rename(fmt.Sprintf("%s/anime_atlas_%s_new.json", client.listsPath, client.locale.Name()), fmt.Sprintf("%s/anime_atlas_%s.json", client.listsPath, client.locale.Name()))
	if err != nil {
		return err
	}

	return nil
}

func (cmd RefreshAnimeCmd) fetchAnime(client CrunchyrollClient, series series, locale core.Locale) (Anime, error) {
	var animeSeasons seasonsResponse
	fmt.Println(series.Id, series.SlugTitle)
	err := client.GetWithQueryParams(fmt.Sprintf("content/v2/cms/series/%s/seasons", series.Id), &animeSeasons, map[string]string{
		// To ensure we get the full list of available subtitle locales, set the audio language to Japanese.
		// This even works for Korean and Chinese works... for the most part. If Japanese isn't applicable, it'll default to the original locale.
		// There can be, say, a Chinese work dubbed to Japanese (A Herbivorous Dragon...), but that one still has the same subtitle locales. And that's the only exception I've found so far...
		"preferred_audio_language": "ja-JP",
		"locale":                   client.locale.Name(),
	})
	if err != nil {
		return Anime{}, err
	}

	// TODO move somewhere else
	identifierOverride := map[string]bool{
		"a3":                          true,
		"arakawa-under-the-bridge":    true,
		"case-closed-detective-conan": true,
	}

	latestSubSeasonNum := 0
	latestSubSeasonId := ""
	latestDubSeasonNum := 0
	latestDubSeasonId := ""
	for _, season := range animeSeasons.Data {
		// We'll need to look at the identifier field to ensure we get a "real" season, and not just some TV specials or whatever.
		// The identifier should be [series_id]|S[season number]. However... of course there are exceptions.
		// - Season number can be off by one (like in One Piece)
		// - Identifier can be something else (like SP1 for Detective Conan, or S1C1 for A3)
		// So we'll check if the identifier fits a format, and includes number(s).
		if _, ok := identifierOverride[series.SlugTitle]; !ok && season.Identifier != "" {
			identifier := strings.Split(season.Identifier, "|")
			if len(identifier) != 2 {
				continue
			}

			if identifier[0] != series.Id {
				continue
			}

			if _, err := strconv.Atoi(strings.Trim(identifier[1], "S")); err != nil {
				continue
			}
		}

		for _, subLoc := range season.SubtitleLocales {
			if subLoc == locale.Name() && latestSubSeasonNum < season.SeasonNumber {
				latestSubSeasonNum = season.SeasonNumber
				latestSubSeasonId = season.Id
				break
			}
		}

		for _, dubVer := range season.Versions {
			if dubVer.AudioLocale == locale.Name() && latestDubSeasonNum < season.SeasonNumber {
				latestDubSeasonNum = season.SeasonNumber
				latestDubSeasonId = season.Id
				break
			}
		}
	}

	// Save an API call if they're the same season.
	if latestSubSeasonNum == latestDubSeasonNum {
		latestDubSeasonId = latestSubSeasonId
	}

	latestSub := Episode{Season: latestSubSeasonNum, Episode: 0, Title: ""}
	latestDub := Episode{Season: latestDubSeasonNum, Episode: 0, Title: ""}

	if latestSubSeasonNum != 0 {
		var subEpisodes seasonEpisodesResponse
		err = client.Get(fmt.Sprintf("content/v2/cms/seasons/%s/episodes", latestSubSeasonId), &subEpisodes)
		if err != nil {
			return Anime{}, err
		}

		for _, episode := range subEpisodes.Data {
			for _, subLoc := range episode.SubtitleLocales {
				if subLoc == locale.Name() && latestSub.Episode < episode.Number {
					latestSub.Episode = episode.Number
					latestSub.Title = episode.Title
					break
				}
			}

			if latestDubSeasonId == latestSubSeasonId {
				for _, dubVer := range episode.Versions {
					if dubVer.AudioLocale == locale.Name() && latestDub.Episode < episode.Number {
						latestDub.Episode = episode.Number
						latestDub.Title = episode.Title
						break
					}
				}
			}
		}
	}

	if latestDubSeasonId != latestSubSeasonId && latestDubSeasonNum != 0 {
		var dubEpisodes seasonEpisodesResponse
		err = client.Get(fmt.Sprintf("content/v2/cms/seasons/%s/episodes", latestDubSeasonId), &dubEpisodes)
		if err != nil {
			return Anime{}, err
		}

		for _, episode := range dubEpisodes.Data {
			for _, dubVer := range episode.Versions {
				if dubVer.AudioLocale == locale.Name() && latestDub.Episode < episode.Number {
					latestDub.Episode = episode.Number
					latestDub.Title = episode.Title
					break
				}
			}
		}
	}

	return Anime{
		Name:        series.Title,
		LastUpdated: time.Now().UTC(),
		Sub:         latestSub,
		Dub:         latestDub,
	}, nil
}

func (cmd RefreshAnimeCmd) refreshAnimePosters(client CrunchyrollClient, allAnime allAnimeResponse) error {
	animePostersFile, err := os.Open(fmt.Sprintf("%s/anime_posters.json", client.listsPath))
	if err != nil {
		return err
	}

	defer animePostersFile.Close()

	bytes, err := io.ReadAll(animePostersFile)
	if err != nil {
		return err
	}

	var animePosters AnimePosters
	err = json.Unmarshal(bytes, &animePosters)
	if err != nil {
		return err
	}

	totalCount := 0
	isDirty := false
	for _, series := range allAnime.Data {
		fmt.Println(series.Id, series.SlugTitle)
		if !cmd.shouldAddSeries(series) {
			delete(animePosters.Posters, series.SlugTitle)
			continue
		}

		tallPoster := image{}
		for _, poster := range series.Images.PosterTall[0] {
			if poster.Width == 60 && poster.Height == 90 {
				tallPoster = poster
				break
			}
		}

		widePoster := image{}
		for _, poster := range series.Images.PosterWide[0] {
			if poster.Width == 320 && poster.Height == 180 {
				widePoster = poster
				break
			}
		}

		if tallPoster.Source == "" && widePoster.Source == "" {
			continue
		}

		if savedPosters, ok := animePosters.Posters[series.SlugTitle]; ok {
			tallHash := sha256.Sum256([]byte(tallPoster.Source))
			if savedPosters.PosterTallHash != fmt.Sprintf("%x", tallHash) {
				encodedPoster, err := cmd.getEncodedPicture(tallPoster.Source)
				if err != nil {
					return err
				}

				savedPosters.PosterTallHash = fmt.Sprintf("%x", tallHash)
				savedPosters.PosterTallEncoded = encodedPoster
				isDirty = true
			}

			wideHash := sha256.Sum256([]byte(widePoster.Source))
			if savedPosters.PosterWideHash != fmt.Sprintf("%x", wideHash) {
				encodedPoster, err := cmd.getEncodedPicture(widePoster.Source)
				if err != nil {
					return err
				}

				savedPosters.PosterWideHash = fmt.Sprintf("%x", wideHash)
				savedPosters.PosterWideEncoded = encodedPoster
				isDirty = true
			}

			animePosters.Posters[series.SlugTitle] = savedPosters

		} else {
			tallHash := sha256.Sum256([]byte(tallPoster.Source))
			wideHash := sha256.Sum256([]byte(widePoster.Source))

			encodedTallPoster, err := cmd.getEncodedPicture(tallPoster.Source)
			if err != nil {
				return err
			}

			encodedWidePoster, err := cmd.getEncodedPicture(widePoster.Source)
			if err != nil {
				return err
			}

			savedPosters := Poster{
				PosterTallHash:    fmt.Sprintf("%x", tallHash),
				PosterTallEncoded: encodedTallPoster,
				PosterWideHash:    fmt.Sprintf("%x", wideHash),
				PosterWideEncoded: encodedWidePoster,
			}
			isDirty = true
			animePosters.Posters[series.SlugTitle] = savedPosters
		}

		totalCount++
	}

	animePosters.TotalCount = totalCount

	if !isDirty {
		return nil
	}

	newPostersFile, err := os.Create(fmt.Sprintf("%s/anime_posters_new.json", client.listsPath))
	if err != nil {
		return err
	}

	newBytes, err := json.MarshalIndent(animePosters, "", "    ")
	if err != nil {
		return err
	}

	_, err = newPostersFile.Write(newBytes)
	if err != nil {
		return err
	}

	err = os.Rename(fmt.Sprintf("%s/anime_posters_new.json", client.listsPath), fmt.Sprintf("%s/anime_posters.json", client.listsPath))
	if err != nil {
		return err
	}

	return nil
}

func (cmd RefreshAnimeCmd) getEncodedPicture(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotModified {
		return "", fmt.Errorf("got bad response when retrieving image: %d", response.StatusCode)
	}

	defer response.Body.Close()

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func (cmd RefreshAnimeCmd) refreshAnimeEpisodeThumbnails(client CrunchyrollClient, allAnime allAnimeResponse) error {
	return nil
}

// shouldAddSeries acts as a single place to store the logic to determine which anime we should even consider.
func (cmd RefreshAnimeCmd) shouldAddSeries(series series) bool {
	// Sometimes Crunchyroll marks a movie as a series. Lovely...
	// Usually they're one season with one episode.
	if series.SeriesMetaData.SeasonCount == 1 && series.SeriesMetaData.EpisodeCount == 1 {
		return false
	}

	// Or... the slug ends in -movies
	if strings.HasSuffix(series.SlugTitle, "-movies") || strings.HasSuffix(series.SlugTitle, "-movie") {
		return false
	}

	// Try not to include OVAs; since they're basically one-time
	if strings.HasSuffix(series.SlugTitle, "-ova") {
		return false
	}

	return true
}
