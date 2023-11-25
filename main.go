package main

import (
	"fmt"
	"os"
	"schoperation/crunchyrollanimestatus/script"
	"strings"
)

type arguments struct {
	credFilePath string
	listsPath    string
	locale       string
	script       string
}

func main() {
	args, err := parseArgs()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	locale, err := script.NewLocale(args.locale)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client, err := script.NewCrunchyrollClient(args.credFilePath, args.listsPath, locale)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cmds := initCmds()
	cmd, ok := cmds[strings.ToLower(args.script)]
	if !ok {
		fmt.Println(fmt.Errorf("unknown cmd %s", args.script))
		return
	}

	err = cmd.Run(client)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func parseArgs() (arguments, error) {
	if len(os.Args) != 7 {
		return arguments{}, fmt.Errorf("need 6 arguments, have %d", len(os.Args))
	}

	refinedArgs := arguments{}
	for i, arg := range os.Args {
		switch arg {
		case "-c":
			refinedArgs.credFilePath = os.Args[i+1]
		case "-l":
			refinedArgs.listsPath = os.Args[i+1]
		case "-lo":
			refinedArgs.locale = os.Args[i+1]
		case "-s":
			refinedArgs.script = os.Args[i+1]
		default:
			continue
		}
	}

	return refinedArgs, nil
}

func initCmds() map[string]script.Command {
	return map[string]script.Command{
		script.NewRefreshAnimeCmd().Name(): script.NewRefreshAnimeCmd(),
	}
}
