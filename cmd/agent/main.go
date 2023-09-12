package main

import (
	"github.com/nameewgeniy/go-metrics/internal/agent"
	"github.com/nameewgeniy/go-metrics/internal/agent/conf"
	"github.com/nameewgeniy/go-metrics/internal/agent/service"
	"log"
)

func main() {

	f, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}

	scf := conf.NewSenderConfig(f.pushAddress)
	snd := service.NewMetricSender(scf)
	rm := service.NewRuntimeMetrics(snd)

	cf := conf.NewAgentConf(f.pollIntervalSec, f.reportIntervalSec)
	a := agent.NewAgent(cf, rm)

	if err = a.Do(); err != nil {
		log.Fatal(err)
	}
}
