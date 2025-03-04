package main

import (
	"finalProject/internal/agent"
	"finalProject/internal/config"
	"sync"
)

func main() {
	cfg := config.MustLoad()
	wg := &sync.WaitGroup{}
	wg.Add(cfg.ComputingPower)
	for _ = range cfg.ComputingPower {
		go agent.StartAgent()
	}
	wg.Wait()
}
