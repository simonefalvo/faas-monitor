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
		t.Errorf("FunctionReplicas(%s) = (%d, %v), want (%d, nil)", name, got, err, want)
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

func TestResponseTimeMicroservice(t *testing.T) {
	name := "sleep"
	sinceSeconds := int64(600)
	want := 2.0
	got, err := ResponseTime(name, sinceSeconds)
	if err != nil {
		t.Errorf("Error occurred: %v\n", err)
	}
	if got < want {
		t.Errorf("ResponseTime(%s) = %v, that is less than %v", name, got, want)
	}
}

// test with a function scaled to zero
func TestResponseTimeZero(t *testing.T) {
	name := "figlet"
	sinceSeconds := int64(600)
	want := 2.0
	got, err := ResponseTime(name, sinceSeconds)
	if got < want || err != nil {
		t.Errorf("ResponseTime(%s) = (%v, %v), want (<time>, nil)", name, got, err)
	}
}

// test with a not existing function
func TestResponseTimeBad(t *testing.T) {
	name := "missingFunction"
	sinceSeconds := int64(600)
	got, err := ResponseTime(name, sinceSeconds)
	if got != 0 || err == nil {
		t.Errorf("ResponseTime(%s) = (%v, %v), want (0, error)", name, got, err)
	}
}

// test with an empty name
func TestResponseTimeEmpty(t *testing.T) {
	name := ""
	sinceSeconds := int64(600)
	got, err := ResponseTime(name, sinceSeconds)
	if got != 0 || err == nil {
		t.Errorf(`ResponseTime("") = (%v, %v), want (0, error)`, got, err)
	}
}
