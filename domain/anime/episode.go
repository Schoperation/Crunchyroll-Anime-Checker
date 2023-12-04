package anime

import (
	"schoperation/crunchyrollanimestatus/domain/crunchyroll"
)

type EpisodeDto struct {
	number       int
	seasonNumber int
	title        string
	thumbnail    crunchyroll.Image
}

type Episode struct {
	number       int
	seasonNumber int
	title        string
	thumbnail    crunchyroll.Image
}
