package main

import (
	"github.com/nameewgeniy/go-metrics/internal/server"
	"github.com/nameewgeniy/go-metrics/internal/server/conf"
	"github.com/nameewgeniy/go-metrics/internal/server/handlers"
	"github.com/nameewgeniy/go-metrics/internal/server/storage/memory"
	"log"
)

func main() {

	store := memory.NewMemory()
	handler := handlers.NewMuxHandlers(store)

	cnf := conf.NewServerConf(":8080")
	srv := server.NewServer(cnf, handler)

	if err := srv.Listen(); err != nil {
		log.Fatal(err)
	}
}
