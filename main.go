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

		functions, err := p.Functions()
		if err != nil {
			log.Fatal(err.Error())
		}

		for _, function := range functions {

			replicas, err := p.FunctionReplicas(function)
			if err != nil {
				log.Println(err.Error())
			}

			fmt.Printf("%s: %d replicas\n", function, replicas)
		}

		time.Sleep(10 * time.Second)
	}
}
