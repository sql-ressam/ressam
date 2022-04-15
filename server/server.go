package server

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/sql-ressam/ressam"
	"github.com/sql-ressam/ressam/api/db"
	"github.com/sql-ressam/ressam/pg"
)

type Settings struct {
	Addr  string
	Debug bool
}

type Server struct {
	httpServer *http.Server
	mux        *chi.Mux
	settings   *Settings
}

func New(ctx context.Context, s *Settings) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:              s.Addr,
			ReadTimeout:       time.Second * 60,
			ReadHeaderTimeout: time.Second * 5,
			WriteTimeout:      time.Second * 10,
			IdleTimeout:       time.Second * 5,
			MaxHeaderBytes:    2 << 15,
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		},
		mux:      chi.NewMux(),
		settings: s,
	}
}

func (s *Server) InitAPI(ctx context.Context, driver, dsn string) error {
	s.mux.Use(cors.AllowAll().Handler)

	if s.settings.Debug {
		s.mux.Post("/api/db/info", db.FakeDBInfo)
		return nil
	}

	switch driver {
	case "postgres": // todo: add aliases
		conn, err := sql.Open("postgres", dsn)
		if err != nil {
			return fmt.Errorf("open connection: %w", err)
		}

		if err := conn.PingContext(ctx); err != nil {
			return fmt.Errorf("ping: %w", err)
		}

		exp := pg.NewExporter(conn)
		dbAPI := db.NewAPI(exp)
		s.mux.Post("/api/db/info", dbAPI.DBInfo)
	default:
		return fmt.Errorf("unsupported driver: %v", driver)
	}

	return nil
}

func (s *Server) InitClient() {
	s.mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})

	// use dynamic folder instead of embedded
	//if s.settings.Debug {
	//	s.mux.Mount("/", http.FileServer())
	//} else {
	s.mux.Mount("/", http.FileServer(http.FS(ressam.GetClientFS())))
	//}
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
