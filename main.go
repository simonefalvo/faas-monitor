package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/smvfal/faas-monitor/pkg/metrics"
	"github.com/smvfal/faas-monitor/pkg/nats"
	"github.com/smvfal/faas-monitor/pkg/types"
)

var scrapePeriod int64

func init() {
	env, ok := os.LookupEnv("SCRAPE_PERIOD")
	if !ok {
		log.Fatal("$SCRAPE_PERIOD not set")
	}
	var err error
	val, err := strconv.Atoi(env)
	if err != nil {
		log.Fatal(err.Error())
	}
	scrapePeriod = int64(val)
}

func main() {

	var p metrics.Provider
	p = &metrics.FaasProvider{}

	for {

		var functions []types.Function
		var nodes []types.Node

		functionNames, err := p.Functions()
		if err != nil {
			log.Fatal(err.Error())
		}

		for _, fname := range functionNames {

			f := types.Function{Name: fname}

			f.Replicas, err = p.FunctionReplicas(f.Name)
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
			}
			log.Printf("%s replicas: %d\n", f.Name, f.Replicas)

			f.InvocationRate, err = p.FunctionInvocationRate(f.Name, scrapePeriod)
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
			}
			log.Printf("%s invocation rate: %d\n", f.Name, f.InvocationRate)

			f.ResponseTime, err = p.ResponseTime(f.Name, scrapePeriod)
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
			}
			log.Printf("%s response time: %v", f.Name, f.ResponseTime)

			f.ProcessingTime, err = p.ProcessingTime(f.Name, scrapePeriod)
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
			}
			log.Printf("%s processing time: %v", f.Name, f.ProcessingTime)

			f.Throughput, err = p.Throughput(f.Name, scrapePeriod)
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
			}
			log.Printf("%s Throughput: %v", f.Name, f.Throughput)

			f.ColdStart, err = p.ColdStart(f.Name, scrapePeriod)
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
			}
			log.Printf("%s cold start time: %v", f.Name, f.ColdStart)

			f.Cpu, f.Mem, err = p.TopPods(f.Name)
			if err != nil {
				log.Printf("WARNING: %s", err.Error())
			}
			log.Printf("%s CPU usage: %s", f.Name, sPrintMap(f.Cpu))
			log.Printf("%s memory usage: %s", f.Name, sPrintMap(f.Mem))

			functions = append(functions, f)

		}

		nodes, err = p.TopNodes()
		if err != nil {
			log.Printf("WARNING: %s", err.Error())
		}
		for _, n := range nodes {
			log.Printf("Node %s CPU usage: %v", n.Name, n.Cpu)
			log.Printf("Node %s memory usage: %v", n.Name, n.Mem)
		}

		msg := types.Message{Functions: functions, Nodes: nodes, Timestamp: time.Now().Unix()}

		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Fatal(err.Error())
		}

		nats.Publish(jsonMsg)

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
