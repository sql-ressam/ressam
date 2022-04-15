package main

import (
	"context"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
	"github.com/urfave/cli/v2"
)

var commands []*cli.Command

func main() {
	app := cli.NewApp()
	app.Usage = "show, modify, export database diagram tool"
	app.Commands = commands

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := app.RunContext(ctx, os.Args); err != nil {
		panic(err)
	}
}
