package script

type RefreshAnimeCmd struct{}

func NewRefreshAnimeCmd() Command {
	return RefreshAnimeCmd{}
}

func (cmd RefreshAnimeCmd) Name() string {
	return "refresh-anime"
}

func (cmd RefreshAnimeCmd) Run(client CrunchyrollClient) error {
	return nil
}
