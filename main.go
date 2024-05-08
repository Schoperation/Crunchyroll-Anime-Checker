/*
	    Crunchyroll (R) Anime Checker
	    Copyright (C) 2024 Schoperation

	    This program is free software: you can redistribute it and/or modify
	    it under the terms of the GNU General Public License as published by
	    the Free Software Foundation, either version 3 of the License, or
	    (at your option) any later version.

	    This program is distributed in the hope that it will be useful,
	    but WITHOUT ANY WARRANTY; without even the implied warranty of
	    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	    GNU General Public License for more details.

	    You should have received a copy of the GNU General Public License
	    along with this program.  If not, see <https://www.gnu.org/licenses/>.

		Crunchyroll (R) is a licensed trademark of Crunchyroll, LLC.
*/
package main

import (
	"database/sql"
	"fmt"
	"os"
	"schoperation/crunchyroll-anime-checker/command"
	"schoperation/crunchyroll-anime-checker/command/subcommand"
	"schoperation/crunchyroll-anime-checker/factory"
	"schoperation/crunchyroll-anime-checker/infrastructure/file"
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
		"refresh-anime":        true,
		"generate-anime-files": true,
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

	animeSenseiListWriter := file.NewAnimeSenseiListWriter(args.listsPath)
	latestEpisodesWriter := file.NewLatestEpisodesWriter(args.listsPath)
	posterWriter := file.NewPosterWriter(args.listsPath)
	thumbnailWriter := file.NewThumbnailWriter(args.listsPath)

	animeDao := sqlite.NewAnimeDao(goquDb)
	latestEpisodesDao := sqlite.NewLatestEpisodesDao(goquDb)
	posterDao := sqlite.NewPosterDao(goquDb)
	thumbnailDao := sqlite.NewThumbnailDao(goquDb)

	crunchyrollClient := rest.NewCrunchyrollClient(args.credFilePath)
	imageClient := rest.NewImageClient()

	latestEpisodesTranslator := anime_translator.NewLatestEpisodesTranslator(latestEpisodesDao, latestEpisodesWriter)
	posterTranslator := anime_translator.NewPosterTranslator(posterDao, posterWriter)
	thumbnailTranslator := anime_translator.NewThumbnailTranslator(thumbnailDao, thumbnailWriter)

	imageTranslator := core_translator.NewImageTranslator(imageClient)

	crunchyrollAnimeTranslator := crunchyroll_translator.NewAnimeTranslator(&crunchyrollClient)
	crunchyrollSeasonTranslator := crunchyroll_translator.NewSeasonTranslator(&crunchyrollClient)
	crunchyrollEpisodeTranslator := crunchyroll_translator.NewEpisodeTranslator(&crunchyrollClient)

	animeFactory := factory.NewAnimeFactory(posterTranslator, latestEpisodesTranslator, thumbnailTranslator)
	animeTranslator := anime_translator.NewAnimeTranslator(animeDao, animeSenseiListWriter, animeFactory)
	animeSaver := saver.NewAnimeSaver(animeTranslator, posterTranslator, latestEpisodesTranslator, thumbnailTranslator)

	refreshBasicInfoSubCommand := subcommand.NewRefreshBasicInfoSubCommand()
	refreshPostersSubCommand := subcommand.NewRefreshPostersSubCommand(imageTranslator)
	refreshLatestEpisodesSubCommand := subcommand.NewRefreshLatestEpisodesSubCommand(crunchyrollSeasonTranslator, crunchyrollEpisodeTranslator, imageTranslator)

	refreshAnimeCommand := command.NewRefreshAnimeCommand(crunchyrollAnimeTranslator, animeTranslator, refreshBasicInfoSubCommand, refreshPostersSubCommand, refreshLatestEpisodesSubCommand, animeSaver)
	generateAnimeFilesCommand := command.NewGenerateAnimeFilesCommand(animeTranslator, animeTranslator, latestEpisodesTranslator, posterTranslator, thumbnailTranslator)

	switch args.cmd {
	case "refresh-anime":
		output, err := refreshAnimeCommand.Run(command.RefreshAnimeCommandInput{})
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("new anime: ", output.NewAnimeCount)
		fmt.Println("updated anime: ", output.UpdatedAnimeCount)
	case "generate-anime-files":
		_, err := generateAnimeFilesCommand.Run(command.GenerateAnimeFilesCommandInput{})
		if err != nil {
			fmt.Println(err)
			return
		}
	}

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
