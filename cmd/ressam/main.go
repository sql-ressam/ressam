package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/urfave/cli/v2"

	// postgres driver
	_ "github.com/lib/pq"
)

var commands []*cli.Command

func main() {
	app := cli.NewApp()
	app.Usage = "show, modify, export database diagram tool"
	app.Commands = commands

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := app.RunContext(ctx, os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
