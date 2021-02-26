package prometheus

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/smvfal/faas-monitor/pkg/util"
)

var v1api v1.API

func init() {
	prometheusUrl, ok := os.LookupEnv("PROMETHEUS_URL")
	if !ok {
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
		msg := fmt.Sprintf("Function %v not found in openfaas-fn namespace.", functionName)
		return 0, errors.New(msg)
	}

	replicas, err := util.ExtractValueBetween(stringResult, `=> `, ` @`)
	return int(replicas), err
}
