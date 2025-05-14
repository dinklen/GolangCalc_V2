package main

import (
	"log"
	"os"

	"github.com/dinklen/GolangCalc_V2/internal/application"
)

func main() {
	app := application.NewApplication()
	if err := app.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
