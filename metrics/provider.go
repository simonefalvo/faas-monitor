package metrics

type Provider interface {
	FunctionReplicas(functionName string) (int, error)
}
