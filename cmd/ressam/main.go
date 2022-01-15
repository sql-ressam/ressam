package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	// postgres driver
	_ "github.com/lib/pq"
)

var commands []*cli.Command

func main() {
	app := cli.NewApp()
	app.Usage = "show, modify, export database diagram tool"
	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Fatalln("app run:", err.Error())
	}
}
