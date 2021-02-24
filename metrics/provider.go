package metrics

type Provider interface {

	// get the names of the deployed functions
	Functions() ([]string, error)

	// get the function replicas' number
	FunctionReplicas(functionName string) (int, error)
}
