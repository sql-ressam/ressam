package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/urfave/cli/v2"

	"github.com/sql-ressam/ressam/server"
)

func init() {
	commands = append(commands, &cli.Command{
		Name: "diagram",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "dsn",
				Required: true,
				EnvVars:  []string{"RESSAM_DSN"},
			},
			&cli.StringFlag{
				Name:  "http",
				Value: ":5555",
			},
			&cli.StringFlag{
				Name:     "driver",
				Value:    "",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "use-client",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			conn, err := sql.Open(c.String("driver"), c.String("dsn"))
			if err != nil {
				return fmt.Errorf("open connection: %w", err)
			}

			if err := conn.PingContext(c.Context); err != nil {
				return fmt.Errorf("ping: %w", err)
			}

			s := server.New(c.String("http"))
			s.InitAPI(conn)
			s.InitClient()

			go func(s server.Server) {
				if err := s.Run(); err != nil {
					log.Fatalln("can't run server:", err.Error())
				}
			}(s)

			return s.Wait()
		},
	})
}
