package main

import (
	"fmt"
	"schoperation/crunchyrollanimestatus/scripts"
)

func main() {
	_, err := scripts.NewCrunchyRollClient()
	fmt.Println(err)
}
