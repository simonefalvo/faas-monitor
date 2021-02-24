package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	user       string
	password   string
	gatewayUrl = os.Getenv("GATEWAY_URL")
)

type function struct {
	Name string `json:"name"`
}

func init() {

	if len(gatewayUrl) == 0 {
		log.Fatal("$GATEWAY_URL not set\n")
	}

	// get and set gateway authorization credentials
	authSecret, err := secret("basic-auth", "openfaas")
	if err != nil {
		log.Fatal(err.Error())
	}
	data := authSecret.Data
	user = string(data["basic-auth-user"])
	password = string(data["basic-auth-password"])
}

func Functions() ([]string, error) {

	// make http api request
	url := gatewayUrl + "/system/functions"
	resBody, err := apiRequest(url, "GET", nil)
	if err != nil {
		return nil, err
	}

	// unmarshal the request body
	var functions []function
	err = json.Unmarshal(resBody, &functions)
	if err != nil {
		return nil, err
	}

	var fnames []string
	for _, f := range functions {
		fnames = append(fnames, f.Name)
	}

	return fnames, nil
}

func apiRequest(url, method string, body io.Reader) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, password)
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("response body:\n%s\n", resBody)
	return resBody, nil
}

func secret(name, namespace string) (*v1.Secret, error) {

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// retrieve the secret
	secrets := clientset.CoreV1().Secrets(namespace)
	secret, err := secrets.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return secret, nil
}
