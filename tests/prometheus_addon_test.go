package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	_ "github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type PrometheusValues struct {
	KubePrometheusStackSource struct {
		Prometheus struct {
			Service struct {
				Port int `yaml:"port"`
			} `yaml:"service"`
		} `yaml:"prometheus"`
	} `yaml:"kube-prometheus-stack-source"`
}

type PrometheusAPIResponse struct {
	Status string `json:"status"`
	Data   struct {
		Result []struct {
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func TestHelmPrometheusServerAddon(t *testing.T) {
	addonData := HelmAddonData{
		namespaceName:   "monitoring",
		releaseName:     "",
		dependencyRepo:  "",
		addonName:       "kube-prometheus-stack",
		addonAlias:      "",
		chartPath:       "",
		hasCustomValues: true,
		manageNamespace: true,
	}

	// add Grafana namespace as it is expected by Prometheus stack
	kubectlOptions := k8s.NewKubectlOptions("", "", addonData.namespaceName)
	k8s.CreateNamespace(t, kubectlOptions, "grafana")

	helmOptions, err := prepareHelmEnvironment(t, &addonData)

	if err != nil {
		log.Fatal(err)
	}

	waitUntilHelmFormattedServicesAvailable(t, addonData, helmOptions, []string{
		"kube-state-metrics",
		"prometheus-node-exporter",
	})

	prometheusServiceName := "kube-prometheus-stack-test-prometheus"
	waitUntilServicesAvailable(t, *helmOptions.KubectlOptions, []string{
		"kube-prometheus-stack-test-alertmanager",
		"kube-prometheus-stack-test-operator",
		prometheusServiceName,
	})

	prometheusStatefulSetName := fmt.Sprintf("prometheus-%s", prometheusServiceName)
	statefulSetAvailable := waitUntilStatefulSetsAvailable(t, *kubectlOptions, []string{prometheusStatefulSetName})
	if !statefulSetAvailable {
		log.Fatalf("Error waiting for StatefulSet %s to become ready.", prometheusStatefulSetName)
	}

	valuesFile, err := os.ReadFile("../addons/kube-prometheus-stack/values.yaml")

	if err != nil {
		log.Fatal(err)
	}

	var yamlData PrometheusValues

	err = yaml.Unmarshal(valuesFile, &yamlData)

	targetPort := 9090
	if err == nil && yamlData.KubePrometheusStackSource.Prometheus.Service.Port != 0 {
		targetPort = yamlData.KubePrometheusStackSource.Prometheus.Service.Port
	}

	tunnel := k8s.NewTunnel(kubectlOptions, k8s.ResourceTypeService, prometheusServiceName, 0, targetPort)
	defer tunnel.Close()

	tunnel.ForwardPort(t)

	metricCheck, err := http.NewRequest("GET", fmt.Sprintf("http://%s/api/v1/query?query=prometheus_tsdb_head_series", tunnel.Endpoint()), nil)
	if err != nil {
		log.Fatalf("Error when building the request: %s", err)
	}

	resp, reqErr := http.DefaultClient.Do(metricCheck)
	if reqErr != nil {
		log.Fatalf("Error when making the request: %s", reqErr)
	}

	respDataFromHttp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error converting response into JSON: %s", err)
	}

	var jsonResp PrometheusAPIResponse
	err = json.Unmarshal([]byte(respDataFromHttp), &jsonResp)
	if err != nil {
		log.Fatalf("Error when unmarshalling response: %s", err)
	}

	if jsonResp.Status != "success" || jsonResp.Data.Result[0].Value[1].(string) == "0" {
		log.Fatalf("Error when querying the Prometheus API. status %s, result %s", jsonResp.Status, jsonResp.Data.Result)
	} else if jsonResp.Status != "error" {
		log.Printf("Prometheus has received %s metrics!", jsonResp.Data.Result[0].Value[1].(string))
	}

	// ----------------------------------

	destroyHelmEnvironment(t, addonData, helmOptions)
	k8s.DeleteNamespace(t, kubectlOptions, "grafana")
}
