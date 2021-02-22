package prometheus

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

var v1api v1.API

func init() {
	prometheusUrl := os.Getenv("PROMETHEUS_URL")
	if prometheusUrl == "" {
		log.Fatal("$PROMETHEUS_URL not set")
	}
	client, err := api.NewClient(api.Config{
		Address: prometheusUrl,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api = v1.NewAPI(client)
}

func FunctionReplicas(functionName string) int {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := fmt.Sprintf(`gateway_service_count{function_name="%v.openfaas-fn"}`, functionName)

	result, warnings, err := v1api.Query(ctx, query, time.Now())
	if err != nil {

		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	return int(extractValue(result.String()))
}

func extractValue(queryResult string) float64 {
	re := regexp.MustCompile(`=> (.*?) @`)
	rs := re.FindStringSubmatch(queryResult)
	stringVal := rs[1]
	val, err := strconv.ParseFloat(stringVal, 64)
	if err != nil {
		log.Fatalf("Unable to convert string %v to float\n", stringVal)
	}
	return val
}