package subcommand

import (
	"schoperation/crunchyrollanimestatus/domain/anime"
	"schoperation/crunchyrollanimestatus/domain/crunchyroll"
)

type RefreshPostersSubCommandInput struct {
	NewAnime     []crunchyroll.Anime
	UpdatedAnime []anime.Anime
}

type RefreshPostersSubCommandOutput struct {
}

type RefreshPostersSubCommand struct {
}

func NewRefreshPostersSubCommand() RefreshPostersSubCommand {
	return RefreshPostersSubCommand{}
}

func (subcmd RefreshPostersSubCommand) Run(input RefreshPostersSubCommandInput) (RefreshPostersSubCommandOutput, error) {
	return RefreshPostersSubCommandOutput{}, nil
}
