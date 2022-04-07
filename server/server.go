package server

import (
	"context"
	"database/sql"
	"net"
	"net/http"
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

func New(addr string) *Server {
	return &Server{
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

func (s *Server) InitAPI(conn *sql.DB) {
	exp := pg.NewExporter(conn)
	dbAPI := db.NewAPI(exp)

	s.mux.Use(cors.AllowAll().Handler)

	s.mux.Post("/api/db/info", dbAPI.DBInfo)
}

func (s *Server) InitClient() {
	s.mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	})
	// todo add client
}

func (s *Server) Run(ctx context.Context) error {
	s.httpServer.Handler = s.mux
	s.httpServer.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}

	errCh := make(chan error, 1)

	go func() {
		err := s.httpServer.ListenAndServe()
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		timeout := time.Second * 10
		stopCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		return s.httpServer.Shutdown(stopCtx)
	case err := <-errCh:
		return err
	}
}
