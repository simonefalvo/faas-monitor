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
		replicas, err := p.FunctionReplicas("nodeinfo")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(replicas)
		time.Sleep(10 * time.Second)
	}
}
