package apiserver

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/smvfal/faas-monitor/pkg/util"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"regexp"
)

var clientset *kubernetes.Clientset

func init() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func ColdStart(function string, sinceSeconds int64) (float64, error) {

	re := regexp.MustCompile("gateway")

	scaleLine := fmt.Sprintf(`\[Scale\] function=%s 0 => 1 successful`, function)
	scaleRe := regexp.MustCompile(scaleLine)
	var scaleSum, scaleCount float64

	podLogOpts := v1.PodLogOptions{
		Container:    "gateway",
		SinceSeconds: &sinceSeconds,
	}

	pods, err := clientset.CoreV1().Pods("openfaas").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return 0, err
	}

	for _, pod := range pods.Items {
		podName := pod.Name
		if re.MatchString(podName) { // match gateway pod
			req := clientset.CoreV1().Pods("openfaas").GetLogs(podName, &podLogOpts)
			podLogs, err := req.Stream(context.TODO())
			if err != nil {
				return 0, err
			}

			scanner := bufio.NewScanner(podLogs)
			for scanner.Scan() {
				line := scanner.Text()
				if scaleRe.MatchString(line) { // match scale line
					val, err := util.ExtractValueBetween(line, `- `, `s`)
					if err != nil {
						return 0, err
					}
					scaleSum += val
					scaleCount++
				}
			}
			err = scanner.Err()
			if err != nil {
				return 0, err
			}

			err = podLogs.Close()
			if err != nil {
				return 0, err
			}
		}
	}

	if scaleCount == 0 {
		return 0, errors.New("no cold starts occurred")
	}

	return scaleSum / scaleCount, nil
}
