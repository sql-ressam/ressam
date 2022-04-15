package main

import (
	"github.com/urfave/cli/v2"

	"github.com/sql-ressam/ressam/server"
)

func init() {
	dsn := &cli.StringFlag{
		Name:     "dsn",
		EnvVars:  []string{"RESSAM_DSN"},
		Required: false,
	}
	http := &cli.StringFlag{
		Name:  "http",
		Value: "127.0.0.1:5555",
	}
	driver := &cli.StringFlag{
		Name:     "driver",
		Value:    "",
		Required: true,
	}
	debug := &cli.BoolFlag{
		Hidden: true,
		Name:   "debug",
		Value:  false,
	}
	commands = append(commands, &cli.Command{
		Name: "diagram",
		Action: func(c *cli.Context) error {
			s := server.New(c.Context, &server.Settings{
				Addr:  c.String(http.Name),
				Debug: c.Bool(debug.Name),
			})

			if err := s.InitAPI(c.Context, c.String(driver.Name), c.String(dsn.Name)); err != nil {
				return err
			}

			s.InitClient()

			return s.Run(c.Context)
		},
		Flags: []cli.Flag{dsn, http, driver, debug},
	})
}
