package metrics

type Provider interface {
	Functions() []string
	FunctionReplicas(functionName string) (int, error)
}
