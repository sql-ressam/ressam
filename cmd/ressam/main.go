package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/sql-ressam/ressam/server"
)

var (
	httpPort = flag.String("http", "localhost:5510", "http port")
	dsn      = flag.String("dsn", "", "data source name, required")
	driver   = flag.String("driver", "postgres", "database driver")
)

func main() {
	flag.Parse()
	if *dsn == "" {
		log.Fatalln("dsn flag required")
	}

	conn, err := sql.Open(*driver, *dsn)
	if err != nil {
		log.Fatalln("open connection:", err.Error())
	}
	if err := conn.Ping(); err != nil {
		log.Fatalln("ping:", err.Error())
	}

	s := server.Init(*httpPort, conn)
	go func() {
		if err := s.Run(); err != nil {
			log.Fatalln("can't run server:", err.Error())
		}
	}()
	s.Wait()
}
