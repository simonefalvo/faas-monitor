package main

import (
	"encoding/json"
	"fmt"
	"github.com/smvfal/faas-monitor/metrics"
	"github.com/smvfal/faas-monitor/nats"
	"log"
	"time"
)

type function struct {
	Name     string `json:"name"`
	Replicas int    `json:"replicas"`
}

func main() {

	var p metrics.Provider
	p = &metrics.FaasProvider{}

	for {

		functions, err := p.Functions()
		if err != nil {
			log.Fatal(err.Error())
		}

		for _, fname := range functions {

			f := function{Name: fname}
			f.Replicas, err = p.FunctionReplicas(f.Name)
			if err != nil {
				log.Println(err.Error())
			}

			fmt.Printf("%s: %d replicas\n", f.Name, f.Replicas)

			// marshal to json
			fjson, err := json.Marshal(f)
			if err != nil {
				log.Fatal(err.Error())
			}
			fmt.Println(string(fjson))

			nats.Publish(fjson)
		}

		time.Sleep(10 * time.Second)
	}
}
