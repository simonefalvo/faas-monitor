package main

import (
	"encoding/json"
	"fmt"
	"github.com/smvfal/faas-monitor/pkg/metrics"
	"github.com/smvfal/faas-monitor/pkg/nats"
	"log"
	"os"
	"strconv"
	"time"
)

type function struct {
	Name         string           `json:"name"`
	Replicas     int              `json:"replicas"`
	ResponseTime float64          `json:"response_time"`
	Throughput   float64          `json:"throughput"`
	ColdStart    float64          `json:"cold_start"`
	Cpu          map[string]int64 `json:"cpu"`
	Mem          map[string]int64 `json:"mem"`
}

var scrapePeriod int

func init() {
	env, ok := os.LookupEnv("SCRAPE_PERIOD")
	if !ok {
		log.Fatal("$SCRAPE_PERIOD not set")
	}
	var err error
	scrapePeriod, err = strconv.Atoi(env)
	if err != nil {
		log.Fatal(err.Error())
	}
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
				log.Printf("WARNING: %s", err.Error())
			}
			log.Printf("%s replicas: %d\n", f.Name, f.Replicas)

			f.ResponseTime, err = p.ResponseTime(f.Name, int64(scrapePeriod))
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
			}
			log.Printf("%s response time: %v", f.Name, f.ResponseTime)

			f.Throughput, err = p.Throughput(f.Name, int64(scrapePeriod))
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
			}
			log.Printf("%s Throughput: %v", f.Name, f.Throughput)

			f.ColdStart, err = p.ColdStart(f.Name, int64(scrapePeriod))
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
			}
			log.Printf("%s cold start time: %v", f.Name, f.ColdStart)

			f.Cpu, f.Mem, err = p.Top(f.Name)
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
			}
			log.Printf("%s CPU usage: %s", f.Name, sPrintMap(f.Cpu))
			log.Printf("%s memory usage: %s", f.Name, sPrintMap(f.Mem))

			// marshal to json
			fjson, err := json.Marshal(f)
			if err != nil {
				log.Fatal(err.Error())
			}

			nats.Publish(fjson)
		}

		time.Sleep(time.Duration(scrapePeriod) * time.Second)
	}
}

func sPrintMap(m map[string]int64) string {
	s := ""
	for key, val := range m {
		s += fmt.Sprintf("\n%s: %d", key, val)
	}
	return s
}
