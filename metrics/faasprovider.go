package metrics

import (
	"github.com/smvfal/faas-monitor/metrics/gateway"
	"github.com/smvfal/faas-monitor/metrics/metricsserver"
	"github.com/smvfal/faas-monitor/metrics/prometheus"
)

type FaasProvider struct{}

func (*FaasProvider) Functions() ([]string, error) {
	return gateway.Functions()
}

func (*FaasProvider) FunctionReplicas(functionName string) (int, error) {
	return prometheus.FunctionReplicas(functionName)
}

func (*FaasProvider) Top(functionName string) (map[string]int64, map[string]int64, error) {
	return metricsserver.Top(functionName)
}
