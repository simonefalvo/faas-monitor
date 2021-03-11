package metricsserver

import (
	"context"
	"errors"
	"fmt"
	"github.com/smvfal/faas-monitor/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	"log"
	"regexp"
)

var mc *metrics.Clientset

const namespace = "openfaas-fn"

func init() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	mc, err = metrics.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func TopPods(function string) (map[string]int64, map[string]int64, error) {

	cpu := make(map[string]int64)
	mem := make(map[string]int64)
	re := regexp.MustCompile(function)

	podMetrics, err := mc.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}

	for _, podMetric := range podMetrics.Items {
		podName := podMetric.Name
		if re.MatchString(podName) {
			podContainers := podMetric.Containers
			cpu[podName] = 0 // initialize cpu counter
			mem[podName] = 0 // initialize mem counter
			for _, container := range podContainers {
				cpu[podName] += container.Usage.Cpu().MilliValue() // add container cpu quantity
				mem[podName] += container.Usage.Memory().Value()   // add container memory quantity
			}
		}
	}

	if len(cpu) == 0 {
		msg := fmt.Sprintf("Function %s not found for resources utilization", function)
		return cpu, mem, errors.New(msg)
	}

	return cpu, mem, nil
}

func TopNodes() ([]types.Node, error) {

	var nodes []types.Node

	nodeMetrics, err := mc.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, nodeMetric := range nodeMetrics.Items {
		nodeName := nodeMetric.Name
		cpu := nodeMetric.Usage.Cpu().MilliValue()
		mem := nodeMetric.Usage.Memory().Value()
		nodes = append(nodes, types.Node{Name: nodeName, Cpu: cpu, Mem: mem})
	}

	return nodes, nil
}
