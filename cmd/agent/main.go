package main

import (
	"github.com/nameewgeniy/go-metrics/internal/agent"
	"github.com/nameewgeniy/go-metrics/internal/agent/conf"
	"github.com/nameewgeniy/go-metrics/internal/agent/service"
	"log"
)

func main() {

	scf := conf.NewSenderConfig("localhost:8080")
	snd := service.NewMetricSender(scf)
	rm := service.NewRuntimeMetrics(snd)

	cf := conf.NewAgentConf(2, 10)
	a := agent.NewAgent(cf, rm)

	if err := a.Do(); err != nil {
		log.Fatal(err)
	}
}
