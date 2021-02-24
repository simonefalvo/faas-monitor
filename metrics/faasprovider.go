package metrics

import (
	"github.com/smvfal/faas-monitor/metrics/gateway"
	"github.com/smvfal/faas-monitor/metrics/prometheus"
)

type FaasProvider struct {
}

func (*FaasProvider) Functions() []string {
	return gateway.Functions()
}

func (*FaasProvider) FunctionReplicas(functionName string) (int, error) {
	return prometheus.FunctionReplicas(functionName)
}
