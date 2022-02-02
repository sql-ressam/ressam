package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"

	// postgres driver
	_ "github.com/lib/pq"
)

var commands []*cli.Command

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGINT)
	defer stop()

	app := cli.NewApp()
	app.Usage = "show, modify, export database diagram tool"
	app.Commands = commands

	if err := app.RunContext(ctx, os.Args); err != nil {
		_, _ = fmt.Fprint(os.Stderr, "app run:", err.Error())
		os.Exit(2)
	}
}
