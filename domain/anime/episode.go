package anime

import "schoperation/crunchyrollanimestatus/domain/crunchyroll"

type EpisodeDto struct {
}

type Episode struct {
	number       int
	seasonNumber int
	title        string
	thumbnail    crunchyroll.Image
}
