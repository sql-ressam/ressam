package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/sql-ressam/ressam/api/db"
	"github.com/sql-ressam/ressam/pg"
)

type Server struct {
	httpServer *http.Server
	mux        *chi.Mux
}

func New(addr string) Server {
	return Server{
		httpServer: &http.Server{
			Addr:              addr,
			ReadTimeout:       time.Second * 60,
			ReadHeaderTimeout: time.Second * 5,
			WriteTimeout:      time.Second * 10,
			IdleTimeout:       time.Second * 5,
			MaxHeaderBytes:    2 << 15,
		},
		mux: chi.NewMux(),
	}
}

func (s Server) InitAPI(conn *sql.DB) {
	exp := pg.NewExporter(conn)
	dbAPI := db.NewAPI(exp)

	s.mux.Use(cors.AllowAll().Handler)

	s.mux.Post("/api/db/info", dbAPI.DBInfo)
}

func (s Server) InitClient() {
	s.mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	})
	// todo add client
}

func (s *Server) Run() error {
	s.httpServer.Handler = s.mux

	return s.httpServer.ListenAndServe()
}

func (s Server) Wait() error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	sig := <-sigCh

	fmt.Println(sig.String())

	return s.httpServer.Close()
}
