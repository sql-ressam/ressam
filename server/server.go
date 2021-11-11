package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/sql-ressam/ressam/api/db"
	"github.com/sql-ressam/ressam/pg"
)

type Server struct {
	httpServer *http.Server
}

func Init(addr string, conn *sql.DB) Server {
	mux := chi.NewMux()
	exp := pg.NewExporter(conn)
	dbAPI := db.NewAPI(exp)

	mux.Post("/api/db/info", dbAPI.DBInfo)

	return Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadTimeout:       time.Second * 60,
			ReadHeaderTimeout: time.Second * 5,
			WriteTimeout:      time.Second * 10,
			IdleTimeout:       time.Second * 5,
			MaxHeaderBytes:    2 << 15,
		},
	}
}

func (s Server) Run() error {
	err := s.httpServer.ListenAndServe()
	return err
}

func (s Server) Wait() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	sig := <-sigCh

	fmt.Println(sig.String())

	_ = s.httpServer.Close()
}
