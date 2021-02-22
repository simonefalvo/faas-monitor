package main

import (
	"fmt"
	"github.com/smvfal/faas-monitor/metrics"
	"time"
)

func main() {
	var p metrics.Provider
	p = &metrics.FaasProvider{}

	for {
		fmt.Println(p.FunctionReplicas("nodeinfo"))
		time.Sleep(10 * time.Duration(time.Second))
	}
}