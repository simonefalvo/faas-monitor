package prometheus

import (
	"context"
	"errors"
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
		log.Fatalf("Error creating client: %v\n", err)
	}

	v1api = v1.NewAPI(client)
}

func FunctionReplicas(functionName string) (int, error) {

	if len(functionName) == 0 {
		msg := "Empty function name\n"
		return 0, errors.New(msg)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := fmt.Sprintf(`gateway_service_count{function_name="%v.openfaas-fn"}`, functionName)

	result, warnings, err := v1api.Query(ctx, query, time.Now())
	if err != nil {
		return 0, err
	}
	if len(warnings) > 0 {
		log.Printf("Warnings: %v\n\n", warnings)
	}
	stringResult := result.String()
	if len(stringResult) == 0 {
		msg := fmt.Sprintf("Function %v not found in openfaas-fn namespace.\n", functionName)
		return 0, errors.New(msg)
	}

	replicas, err := extractValue(stringResult)
	return int(replicas), err
}

func extractValue(queryResult string) (float64, error) {
	re := regexp.MustCompile(`=> (.*?) @`)
	rs := re.FindStringSubmatch(queryResult)
	stringVal := rs[1]
	val, err := strconv.ParseFloat(stringVal, 64)
	if err != nil {
		msg := fmt.Sprintf("Unable to convert string %v to float\n", stringVal)
		return 0, errors.New(msg)
	}
	return val, nil
}
