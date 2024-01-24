package file

import (
	"fmt"
	"unicode"
)

// fileId is used in generating poster and thumbnail files, since having
// one master list of those would be too big for the Tidbyt http client to handle.
// Typically this is just the first letter of the slug title.
func fileId(slugtitle string) string {
	if !unicode.IsLetter(rune(slugtitle[0])) {
		return "zmisc"
	}

	return fmt.Sprintf("%c", slugtitle[0])
}

func getFileIds() []string {
	return []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"zmisc",
	}
}
