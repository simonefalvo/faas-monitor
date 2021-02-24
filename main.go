package main

import (
	"fmt"
	"github.com/smvfal/faas-monitor/metrics"
	"log"
	"time"
)

func main() {
	var p metrics.Provider
	p = &metrics.FaasProvider{}

	for {
		functions := p.Functions()
		for _, function := range functions {
			replicas, err := p.FunctionReplicas(function)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s: %d replicas\n", function, replicas)
		}
		time.Sleep(10 * time.Second)
	}
}
