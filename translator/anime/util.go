package anime

import (
	"fmt"
	"schoperation/crunchyroll-anime-checker/domain/anime"
	"unicode"
)

// fileId is used in generating poster and thumbnail files, since having
// one master list of those would be too big for the Tidbyt http client to handle.
// Typically this is just the first letter of the slug title.
func fileId(localAnime anime.Anime) string {
	if !unicode.IsLetter(rune(localAnime.SlugTitle()[0])) {
		return "zmisc"
	}

	return fmt.Sprintf("%c", localAnime.SlugTitle()[0])
}
