package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"

	"github.com/sql-ressam/ressam/server"
)

func drawCommand() *cli.Command {
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
	httpFlag := &cli.IntFlag{
		Name:  "port",
		Value: 33939,
	}
	debugFlag := &cli.BoolFlag{
		Name:   "debug",
		Hidden: true,
		Value:  false,
	}
	return &cli.Command{
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

			if err := server.WaitStarts(httpFlag.Value, errCh); err != nil {
				return err
			}

			if err := open.Run("http://127.0.0.1:" + strconv.Itoa(httpFlag.Value)); err != nil {
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
	}
}
