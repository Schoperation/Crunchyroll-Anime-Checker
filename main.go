package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"schoperation/crunchyrollanimestatus/command"
	"schoperation/crunchyrollanimestatus/command/subcommand"
	"schoperation/crunchyrollanimestatus/factory"
	"schoperation/crunchyrollanimestatus/infrastructure/postgres"
	"schoperation/crunchyrollanimestatus/infrastructure/rest"
	anime_translator "schoperation/crunchyrollanimestatus/translator/anime"
	core_translator "schoperation/crunchyrollanimestatus/translator/core"
	crunchyroll_translator "schoperation/crunchyrollanimestatus/translator/crunchyroll"
	"strings"

	"github.com/jackc/pgx/v5"
)

type arguments struct {
	dbCredFilePath string
	credFilePath   string
	listsPath      string
	cmd            string
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

	db, err := openDb(args.dbCredFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close(context.Background())

	animeDao := postgres.NewAnimeDao(db)
	latestEpisodesDao := postgres.NewLatestEpisodesDao(db)
	posterDao := postgres.NewPosterDao(db)
	thumbnailDao := postgres.NewThumbnailDao(db)

	crunchyrollClient := rest.NewCrunchyrollClient(args.credFilePath)
	imageClient := rest.NewImageClient()

	latestEpisodesTranslator := anime_translator.NewLatestEpisodesTranslator(latestEpisodesDao)
	posterTranslator := anime_translator.NewPosterTranslator(posterDao)
	thumbnailTranslator := anime_translator.NewThumbnailTranslator(thumbnailDao)

	imageTranslator := core_translator.NewImageTranslator(imageClient)

	crunchyrollAnimeTranslator := crunchyroll_translator.NewAnimeTranslator(&crunchyrollClient)

	animeFactory := factory.NewAnimeFactory(posterTranslator, latestEpisodesTranslator, thumbnailTranslator)
	animeTranslator := anime_translator.NewAnimeTranslator(animeDao, animeFactory)

	refreshPostersSubCommand := subcommand.NewRefreshPostersSubCommand(imageTranslator)

	refreshAnimeCommand := command.NewRefreshAnimeCommand(crunchyrollAnimeTranslator, animeTranslator, refreshPostersSubCommand)

	output, err := refreshAnimeCommand.Run(command.RefreshAnimeCommandInput{})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("new anime: ", output.NewAnimeCount)
	fmt.Println("updated anime: ", output.UpdatedAnimeCount)
}

func parseArgs() (arguments, error) {
	if len(os.Args) != 9 {
		return arguments{}, fmt.Errorf("need 8 arguments, have %d", len(os.Args)-1)
	}

	refinedArgs := arguments{}
	for i, arg := range os.Args {
		switch arg {
		case "-db":
			refinedArgs.dbCredFilePath = os.Args[i+1]
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

func openDb(dbCredFilePath string) (*pgx.Conn, error) {
	file, err := os.Open(dbCredFilePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Database string `json:"database"`
	}

	err = json.Unmarshal(bytes, &creds)
	if err != nil {
		return nil, err
	}

	db, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s", creds.Username, creds.Password, creds.Host, creds.Port, creds.Database))
	if err != nil {
		return nil, err
	}

	return db, nil
}
