package metrics

import "github.com/smvfal/faas-monitor/pkg/types"

type Provider interface {

	// get the names of the deployed functions
	Functions() ([]string, error)

	// get the function replicas' number
	FunctionReplicas(functionName string) (int, error)

	// get function's average response time
	ResponseTime(functionName string, sinceSeconds int64) (float64, error)

	// get function's average processing time
	ProcessingTime(functionName string, sinceSeconds int64) (float64, error)

	// get function's throughput
	Throughput(functionName string, sinceSeconds int64) (float64, error)

	// get function's cold start time
	ColdStart(functionName string, SinceSeconds int64) (float64, error)

	// get function's current CPU and memory usage for each replica
	TopPods(functionName string) (map[string]int64, map[string]int64, error)

	// get nodes current CPU and memory usage
	TopNodes() ([]types.Node, error)
}
