package main

import (
	"context"
	"fmt"
	"os"
	"schoperation/crunchyrollanimestatus/command"
	"schoperation/crunchyrollanimestatus/infrastructure/postgres"
	"schoperation/crunchyrollanimestatus/infrastructure/rest"
	anime_translator "schoperation/crunchyrollanimestatus/translator/anime"
	crunchyroll_translator "schoperation/crunchyrollanimestatus/translator/crunchyroll"
	"strings"

	"github.com/jackc/pgx/v5"
)

type arguments struct {
	credFilePath string
	listsPath    string
	cmd          string
}

func main() {
	args, err := parseArgs()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cmds := map[string]bool{
		"refresh-anime": true,
	}
	_, ok := cmds[args.cmd]
	if !ok {
		fmt.Println(fmt.Errorf("unknown cmd %s", args.cmd))
		return
	}

	db, err := pgx.Connect(context.Background(), fmt.Sprintf(""))
	if err != nil {
		fmt.Println(err)
		return
	}

	crunchyrollClient := rest.NewCrunchyrollClient(args.credFilePath)

	animeDao := postgres.NewAnimeDao(db)

	animeTranslator := anime_translator.NewAnimeTranslator(animeDao)
	crunchyrollAnimeTranslator := crunchyroll_translator.NewAnimeTranslator(&crunchyrollClient)

	refreshAnimeCommand := command.NewRefreshAnimeCommand(crunchyrollAnimeTranslator, animeTranslator)

	output, err := refreshAnimeCommand.Run(command.RefreshAnimeCommandInput{})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("new anime: ", output.NewAnimeCount)
	fmt.Println("updated anime: ", output.UpdatedAnimeCount)
}

func parseArgs() (arguments, error) {
	if len(os.Args) != 7 {
		return arguments{}, fmt.Errorf("need 6 arguments, have %d", len(os.Args)-1)
	}

	refinedArgs := arguments{}
	for i, arg := range os.Args {
		switch arg {
		case "-c":
			refinedArgs.credFilePath = os.Args[i+1]
		case "-l":
			refinedArgs.listsPath = os.Args[i+1]
		case "-cmd":
			refinedArgs.cmd = strings.ToLower(os.Args[i+1])
		default:
			continue
		}
	}

	return refinedArgs, nil
}
