package main

import (
	"database/sql"
	"fmt"
	"os"
	"schoperation/crunchyroll-anime-checker/command"
	"schoperation/crunchyroll-anime-checker/command/subcommand"
	"schoperation/crunchyroll-anime-checker/factory"
	"schoperation/crunchyroll-anime-checker/infrastructure/rest"
	"schoperation/crunchyroll-anime-checker/infrastructure/sqlite"
	"schoperation/crunchyroll-anime-checker/saver"
	anime_translator "schoperation/crunchyroll-anime-checker/translator/anime"
	core_translator "schoperation/crunchyroll-anime-checker/translator/core"
	crunchyroll_translator "schoperation/crunchyroll-anime-checker/translator/crunchyroll"
	"strings"

	"github.com/doug-martin/goqu/v9"
	goqu_sqlite3_dialect "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "modernc.org/sqlite"
)

type arguments struct {
	dbPath       string
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

	db, err := sql.Open("sqlite", args.dbPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// Pending new release from goqu to fix this...
	// In the meantime, creating new dialect
	goquDialectOpts := goqu_sqlite3_dialect.DialectOptions()
	goquDialectOpts.SupportsReturn = true
	goqu.RegisterDialect(sqlite.Dialect, goquDialectOpts)
	goquDb := goqu.New(sqlite.Dialect, db)

	animeDao := sqlite.NewAnimeDao(goquDb)
	latestEpisodesDao := sqlite.NewLatestEpisodesDao(goquDb)
	posterDao := sqlite.NewPosterDao(goquDb)
	thumbnailDao := sqlite.NewThumbnailDao(goquDb)

	crunchyrollClient := rest.NewCrunchyrollClient(args.credFilePath)
	imageClient := rest.NewImageClient()

	latestEpisodesTranslator := anime_translator.NewLatestEpisodesTranslator(latestEpisodesDao)
	posterTranslator := anime_translator.NewPosterTranslator(posterDao)
	thumbnailTranslator := anime_translator.NewThumbnailTranslator(thumbnailDao)

	imageTranslator := core_translator.NewImageTranslator(imageClient)

	crunchyrollAnimeTranslator := crunchyroll_translator.NewAnimeTranslator(&crunchyrollClient)
	crunchyrollSeasonTranslator := crunchyroll_translator.NewSeasonTranslator(&crunchyrollClient)
	crunchyrollEpisodeTranslator := crunchyroll_translator.NewEpisodeTranslator(&crunchyrollClient)

	animeFactory := factory.NewAnimeFactory(posterTranslator, latestEpisodesTranslator, thumbnailTranslator)
	animeTranslator := anime_translator.NewAnimeTranslator(animeDao, animeFactory)
	animeSaver := saver.NewAnimeSaver(animeTranslator, posterTranslator, latestEpisodesTranslator, thumbnailTranslator)

	refreshBasicInfoSubCommand := subcommand.NewRefreshBasicInfoSubCommand()
	refreshPostersSubCommand := subcommand.NewRefreshPostersSubCommand(imageTranslator)
	refreshLatestEpisodesSubCommand := subcommand.NewRefreshLatestEpisodesSubCommand(crunchyrollSeasonTranslator, crunchyrollEpisodeTranslator, imageTranslator)

	refreshAnimeCommand := command.NewRefreshAnimeCommand(crunchyrollAnimeTranslator, animeTranslator, refreshBasicInfoSubCommand, refreshPostersSubCommand, refreshLatestEpisodesSubCommand, animeSaver)

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
			refinedArgs.dbPath = os.Args[i+1]
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
