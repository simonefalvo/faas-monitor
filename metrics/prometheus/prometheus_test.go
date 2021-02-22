package prometheus

import (
	"testing"
)

// test with an always running function
func TestFunctionReplicasMicroservice(t *testing.T) {
	name := "nodeinfo"
	want := 1
	got, err := FunctionReplicas(name)
	if err != nil {
		t.Errorf("Error occurred: %v\n", err)
	}
	if got < want {
		t.Errorf("FunctionReplicas(%s) = %d, that is less than %d", name, got, want)
	}
}

// test with a function scaled to zero
func TestFunctionReplicasZero(t *testing.T) {
	name := "figlet"
	want := 0
	got, err := FunctionReplicas(name)
	if got != want || err != nil {
		t.Errorf("FunctionReplicas(%s) = (%d, %v), want (%d, nil)", name, err, got, want)
	}
}

// test with a not existing function
func TestFunctionReplicasBad(t *testing.T) {
	name := "missingFunction"
	got, err := FunctionReplicas(name)
	if got != 0 || err == nil {
		t.Errorf("FunctionReplicas(%s) = (%d, %v), want (0, error)", name, got, err)
	}
}

// test with an empty name
func TestFunctionReplicasEmpty(t *testing.T) {
	name := ""
	got, err := FunctionReplicas(name)
	if got != 0 || err == nil {
		t.Errorf(`FunctionReplicas("") = (%d, %v), want (0, error)`, got, err)
	}
}
