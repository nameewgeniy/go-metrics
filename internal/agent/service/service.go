package service

import "fmt"

type RuntimeMetrics struct {
}

func NewRuntimeMetrics() *RuntimeMetrics {
	return &RuntimeMetrics{}
}

func (m RuntimeMetrics) Push() {
	fmt.Print("Push")
}

func (m RuntimeMetrics) Sync() {
	fmt.Print("Sync")
}
