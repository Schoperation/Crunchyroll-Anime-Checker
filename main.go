package main

import (
	"fmt"
	"os"
	"schoperation/crunchyrollanimestatus/script"
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

}

func parseArgs() (arguments, error) {
	if len(os.Args) != 7 {
		return arguments{}, fmt.Errorf("Need 6 arguments, have %d", len(os.Args))
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
