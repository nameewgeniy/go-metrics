package main

import (
	"github.com/nameewgeniy/go-metrics/internal/agent"
	"github.com/nameewgeniy/go-metrics/internal/agent/conf"
	"github.com/nameewgeniy/go-metrics/internal/agent/service"
	"log"
)

func main() {

	mcf := conf.NewMetricsConf()
	rm := service.NewRuntimeMetrics(mcf)

	cf := conf.NewAgentConf()
	a := agent.NewAgent(cf, rm)

	if err := a.Do(); err != nil {
		log.Fatal(err)
	}
}
