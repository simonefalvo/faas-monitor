package metrics

import (
	"github.com/smvfal/faas-monitor/pkg/metrics/apiserver"
	"github.com/smvfal/faas-monitor/pkg/metrics/gateway"
	"github.com/smvfal/faas-monitor/pkg/metrics/metricsserver"
	"github.com/smvfal/faas-monitor/pkg/metrics/prometheus"
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

func (*FaasProvider) ColdStart(functionName string, sinceSeconds int64) (float64, error) {
	return apiserver.ColdStart(functionName, sinceSeconds)
}
