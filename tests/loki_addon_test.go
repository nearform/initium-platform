package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"

	_ "github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type LokiValues struct {
	LokiSource struct {
		Loki struct {
			Server struct {
				HttpListenPort int `yaml:"http_listen_port"`
			} `yaml:"server"`
		} `yaml:"loki"`
	} `yaml:"loki-source"`
}

type LokiAPIResponse struct {
	Status string `json:"status"`
	Data   struct {
		Result []struct {
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func TestHelmLokiAddon(t *testing.T) {
	lokiNamespaceAndAddonName := "loki"
	addonData := HelmAddonData{
		namespaceName:   "",
		releaseName:     "",
		dependencyRepo:  "",
		addonName:       lokiNamespaceAndAddonName,
		addonAlias:      "",
		chartPath:       "",
		hasCustomValues: true,
		manageNamespace: true,
	}

	helmOptions, err := prepareHelmEnvironment(t, &addonData)

	if err != nil {
		log.Fatal(err)
	}

	waitUntilServicesAvailable(t, *helmOptions.KubectlOptions, []string{
		lokiNamespaceAndAddonName,
		fmt.Sprintf("%s-headless", lokiNamespaceAndAddonName),
		fmt.Sprintf("%s-memberlist", lokiNamespaceAndAddonName),
	})

	statefulSetAvailable := waitUntilStatefulSetsAvailable(t, *helmOptions.KubectlOptions, []string{lokiNamespaceAndAddonName})
	if !statefulSetAvailable {
		log.Fatalf("Error waiting for StatefulSet %s to become ready.", lokiNamespaceAndAddonName)
	}

	// Wait for the applications to send logs to Loki
	time.Sleep(120 * time.Second)

	valuesFile, err := os.ReadFile("../addons/loki/values.yaml")

	if err != nil {
		log.Fatal(err)
	}

	var yamlData LokiValues

	err = yaml.Unmarshal(valuesFile, &yamlData)

	targetPort := 3100
	if err == nil && yamlData.LokiSource.Loki.Server.HttpListenPort != 0 {
		targetPort = yamlData.LokiSource.Loki.Server.HttpListenPort
	}

	tunnel := k8s.NewTunnel(helmOptions.KubectlOptions, k8s.ResourceTypeService, lokiNamespaceAndAddonName, 0, targetPort)
	defer tunnel.Close()

	tunnel.ForwardPort(t)

	healthCheck, err := http.NewRequest("GET", fmt.Sprintf("http://%s/loki/api/v1/query?query=%s", tunnel.Endpoint(), url.QueryEscape("sum by (msg) (rate({container_name=~\".+\"}[5m]))")), nil)
	if err != nil {
		log.Fatalf("Error when building the request: %s", err)
	}

	resp, reqErr := http.DefaultClient.Do(healthCheck)
	if reqErr != nil {
		log.Fatalf("Error when making the request: %s", reqErr)
	}

	respDataFromHttp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error converting response into JSON: %s", err)
	}

	var jsonResp LokiAPIResponse
	err = json.Unmarshal([]byte(respDataFromHttp), &jsonResp)
	if err != nil {
		log.Fatalf("Error when unmarshalling response: %s", err)
	}

	log.Print(jsonResp)
	if jsonResp.Status != "success" {
		log.Fatalf("Error when querying the Loki API. Response: %s.", jsonResp)
	} else if jsonResp.Status != "error" {
		log.Printf("Loki is ready. API response: %s!", jsonResp.Data.Result)
	}
	// ----------------------------------

	destroyHelmEnvironment(t, addonData, helmOptions)
}
