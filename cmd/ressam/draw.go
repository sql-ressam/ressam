package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"

	"github.com/sql-ressam/ressam/server"
)

func init() {
	dsnFlag := &cli.StringFlag{
		Name:     "dsn",
		EnvVars:  []string{"RESSAM_DSN"},
		Required: true,
	}
	driverFlag := &cli.StringFlag{
		Name:     "driver",
		EnvVars:  []string{"RESSAM_DRIVER"},
		Required: true,
		Value:    "", // todo(aleksvdim): try to parse from DSN
	}
	httpFlag := &cli.StringFlag{
		Name:  "http",
		Value: "127.0.0.1:3939",
	}
	debugFlag := &cli.BoolFlag{
		Name:   "debug",
		Hidden: true,
		Value:  false,
	}
	commands = append(commands, &cli.Command{
		Name:  "draw",
		Flags: []cli.Flag{dsnFlag, httpFlag, driverFlag, debugFlag},
		Action: func(c *cli.Context) error {
			s := server.New(c.Context, &server.Settings{
				Addr:  c.String(httpFlag.Name),
				Debug: c.Bool(debugFlag.Name),
			})

			if err := s.InitAPI(c.Context, c.String(driverFlag.Name), c.String(dsnFlag.Name)); err != nil {
				return err
			}

			s.InitClient()

			errCh := make(chan error, 1)
			go func() {
				errCh <- s.Run(c.Context)
			}()

			webAppUrl := fmt.Sprintf("http://%s/", httpFlag.Value)
			if err := server.WaitStarts(webAppUrl, errCh); err != nil {
				return err
			}

			if err := open.Run(webAppUrl); err != nil {
				return fmt.Errorf("can't open web browser: %w", err)
			}

			select {
			case err := <-errCh:
				if errors.Is(err, http.ErrServerClosed) {
					return nil
				}
				return err
			}
		},
	})
}
