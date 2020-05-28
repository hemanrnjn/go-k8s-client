package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	serviceHost := os.Getenv("KUBERNETES_SERVICE_HOST")
	servicePort := os.Getenv("KUBERNETES_SERVICE_PORT")
	apiVersion := "v1alpha1"            // your api version name
	apiGroup := "example.com.my.domain" // your api group name
	namespace := "default"              // k8s namespace
	resource := "backupjobs"            // your custom resource definition name
	resourceName := "backupjob-sample"  // your custom resource name

	url := fmt.Sprintf("https://%s:%s/apis/%s/%s/namespaces/%s/%s/%s", serviceHost, servicePort, apiGroup, apiVersion, namespace, resource, resourceName)

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err)
	}

	caToken, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		panic(err) // cannot find token file
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", string(caToken)))

	caCertPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
	if err != nil {
		panic(err) // Can't find cert file
	}
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %s", err.Error())
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("get:\n", string(body))
}
