package main

import (
	"github.com/nameewgeniy/go-metrics/internal/agent"
	"github.com/nameewgeniy/go-metrics/internal/agent/conf"
	"github.com/nameewgeniy/go-metrics/internal/agent/service"
	"log"
)

func main() {

	cf := conf.NewAgentConf()
	rm := service.NewRuntimeMetrics()

	a := agent.NewAgent(cf, rm)

	if err := a.Do(); err != nil {
		log.Fatal(err)
	}
}
