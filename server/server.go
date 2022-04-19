package server

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/sql-ressam/ressam"
	"github.com/sql-ressam/ressam/api/db"
	"github.com/sql-ressam/ressam/database/pg"
)

// Settings contains server options.
type Settings struct {
	Addr  string
	Debug bool
}

// Server accepts HTTP requests.
type Server struct {
	httpServer *http.Server
	mux        *chi.Mux
	settings   *Settings
}

// New returns new Server instance.
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

const (
	dbInfoPath = "/api/db/info"
)

// InitAPI inits http API.
func (s *Server) InitAPI(ctx context.Context, driver, dsn string) error {
	s.mux.Use(cors.AllowAll().Handler)

	if s.settings.Debug {
		s.mux.Post(dbInfoPath, db.FakeDBInfo)
		return nil
	}

	switch strings.ToLower(driver) {
	case "postgres", "postgresql", "pg", "postgre":
		conn, err := sql.Open("postgres", dsn)
		if err != nil {
			return fmt.Errorf("open connection: %w", err)
		}

		if err := conn.PingContext(ctx); err != nil {
			return fmt.Errorf("can't ping pg: %w", err)
		}

		exp := pg.NewExporter(conn)
		dbAPI := db.NewAPI(exp)
		s.mux.Post(dbInfoPath, dbAPI.DBInfo)
	default:
		return fmt.Errorf("unsupported driver: %v", driver)
	}

	return nil
}

// InitClient inits embedded ui client.
func (s *Server) InitClient() {
	s.mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})

	// use os folder instead of embed
	if s.settings.Debug {
		s.mux.Mount("/", http.FileServer(http.FS(ressam.GetClientFS())))
	} else {
		s.mux.Mount("/", http.FileServer(http.FS(ressam.GetEmbeddedClientFS())))
	}
}

// Run the server until the context is canceled.
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
		// .Shutdown() here is redundant
		if err := s.httpServer.Close(); err != nil {
			return fmt.Errorf("can't close http server: %w", err)
		}
		return <-errCh
	case err := <-errCh:
		return fmt.Errorf("can't listen addr: %w", err)
	}
}
