package script

type Command interface {
	Run(client CrunchyrollClient) error
}
