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
		log.Fatalf("Error creating client: %v", err)
	}

	v1api = v1.NewAPI(client)
}

func FunctionReplicas(functionName string) (int, error) {

	if len(functionName) == 0 {
		msg := "empty function name"
		return 0, errors.New(msg)
	}

	q := fmt.Sprintf(`sum by (function_name) (gateway_service_count{function_name="%v.openfaas-fn"})`,
		functionName)

	stringResult, err := query(q)
	if err != nil {
		return 0, err
	}

	if len(stringResult) == 0 {
		msg := fmt.Sprintf("function %s not found in the openfaas-fn namespace", functionName)
		return 0, errors.New(msg)
	}

	replicas, err := util.ExtractValueBetween(stringResult, `=> `, ` @`)
	if err != nil {
		return 0, err
	}

	return int(replicas), err
}

func ResponseTime(functionName string, sinceSeconds int64) (float64, error) {

	if len(functionName) == 0 {
		msg := "empty function name"
		return 0, errors.New(msg)
	}

	q := fmt.Sprintf(
		`sum by (function_name)`+
			`(rate(gateway_functions_seconds_sum{function_name="%s.openfaas-fn"}[%ds]) > 0) `+
			`/ `+
			`sum by (function_name)`+
			`(rate(gateway_functions_seconds_count{function_name="%s.openfaas-fn"}[%ds]) > 0)`,
		functionName, sinceSeconds, functionName, sinceSeconds,
	)

	stringResult, err := querySince(q, functionName, sinceSeconds)
	if err != nil {
		return 0, err
	}

	rt, err := util.ExtractValueBetween(stringResult, `=> `, ` @`)
	if err != nil {
		return 0, err
	}

	return rt, nil
}

func querySince(q, functionName string, sinceSeconds int64) (string, error) {

	result, err := query(q)
	if len(result) == 0 {
		msg := fmt.Sprintf("function %s not invoked in the last %v seconds",
			functionName, sinceSeconds)
		return "", errors.New(msg)
	}

	return result, err
}

func query(q string) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, warnings, err := v1api.Query(ctx, q, time.Now())
	if err != nil {
		return "", err
	}
	if len(warnings) > 0 {
		log.Printf("Warnings: %v", warnings)
	}

	stringResult := result.String()
	//fmt.Println(stringResult)

	return stringResult, nil
}
