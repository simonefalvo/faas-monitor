package metrics

type Provider interface {

	// get the names of the deployed functions
	Functions() ([]string, error)

	// get the function replicas' number
	FunctionReplicas(functionName string) (int, error)

	// get function's current resources usage for each replica
	Top(functionName string) (map[string]int64, map[string]int64, error)
}
