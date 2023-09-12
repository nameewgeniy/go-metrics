package main

import (
	"github.com/nameewgeniy/go-metrics/internal/server"
	"github.com/nameewgeniy/go-metrics/internal/server/conf"
	"github.com/nameewgeniy/go-metrics/internal/server/handlers"
	"github.com/nameewgeniy/go-metrics/internal/server/storage/memory"
	"log"
)

func main() {

	f, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}

	store := memory.NewMemory()
	handler := handlers.NewMuxHandlers(store)

	cnf := conf.NewServerConf(f.addr)
	srv := server.NewServer(cnf, handler)

	if err = srv.Listen(); err != nil {
		log.Fatal(err)
	}
}
