package main

import (
	"database/sql"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/sql-ressam/ressam/server"
)

func init() {
	dsn := &cli.StringFlag{
		Name:     "dsn",
		Required: true,
		EnvVars:  []string{"RESSAM_DSN"},
	}
	http := &cli.StringFlag{
		Name:  "http",
		Value: ":5555",
	}
	driver := &cli.StringFlag{
		Name:     "driver",
		Value:    "",
		Required: true,
	}
	useClient := &cli.BoolFlag{
		Name:  "use-client",
		Value: false,
	}
	commands = append(commands, &cli.Command{
		Name:  "diagram",
		Flags: []cli.Flag{dsn, http, driver, useClient},
		Action: func(c *cli.Context) error {
			conn, err := sql.Open(c.String(driver.Name), c.String(dsn.Name))
			if err != nil {
				return fmt.Errorf("open connection: %w", err)
			}

			if err := conn.PingContext(c.Context); err != nil {
				return fmt.Errorf("ping: %w", err)
			}

			s := server.New(c.String(http.Name))
			s.InitAPI(conn)
			s.InitClient()

			return s.Run(c.Context)
		},
	})
}
