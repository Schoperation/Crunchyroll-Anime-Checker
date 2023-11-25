package script

type Command interface {
	Name() string
	Run(client CrunchyrollClient) error
}
