package prometheus

import "testing"

func TestFunctionReplicas(t *testing.T) {
	function := "nodeinfo"
	want := 1
	got := FunctionReplicas(function)
	if got != want {
		t.Errorf("FunctionReplicas(%s) == %d, want %d", function, got, want)
	}
}