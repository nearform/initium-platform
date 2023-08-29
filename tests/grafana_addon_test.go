package test

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	_ "github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/gruntwork-io/terratest/modules/k8s"
)

type GrafanaValues struct {
	GrafanaSource struct {
		Service struct {
			TargetPort int `yaml:"targetPort"`
		} `yaml:"service"`
	} `yaml:"grafana-source"`
}

func TestHelmGrafanaAddon(t *testing.T) {
	// add Istio CRDs to test Istio virtual service
	istioBaseAddonData := HelmAddonData{
		namespaceName:   "grafana-test",
		releaseName:     "",
		dependencyRepo:  "",
		addonName:       "istio-base",
		addonAlias:      "",
		chartPath:       "../addons/istio/base",
		manageNamespace: true,
		overrideValues: map[string]string{
			"global.istioNamespace": "grafana-test",
		},
	}

	istioBaseHelmOptions, err := prepareHelmEnvironment(t, &istioBaseAddonData)

	if err != nil {
		log.Fatal(err)
	}

	grafanaNamespaceAndAddonName := "grafana"

	addonData := HelmAddonData{
		namespaceName:   "",
		releaseName:     "",
		dependencyRepo:  "",
		addonName:       grafanaNamespaceAndAddonName,
		addonAlias:      "",
		chartPath:       "",
		hasCustomValues: true,
		manageNamespace: true,
	}

	helmOptions, err := prepareHelmEnvironment(t, &addonData)

	if err != nil {
		log.Fatal(err)
	}

	waitUntilServicesAvailable(t, *helmOptions.KubectlOptions, []string{grafanaNamespaceAndAddonName})
	waitUntilDeploymentsAvailable(t, *helmOptions.KubectlOptions, []string{grafanaNamespaceAndAddonName})

	valuesFile, err := os.ReadFile("../addons/grafana/values.yaml")

	if err != nil {
		log.Fatal(err)
	}

	var yamlData GrafanaValues

	err = yaml.Unmarshal(valuesFile, &yamlData)

	targetPort := 3000
	if err == nil && yamlData.GrafanaSource.Service.TargetPort != 0 {
		targetPort = yamlData.GrafanaSource.Service.TargetPort
	}

	grafanaKubectlOptions := k8s.NewKubectlOptions("", "", grafanaNamespaceAndAddonName)

	tunnel := k8s.NewTunnel(grafanaKubectlOptions, k8s.ResourceTypeService, grafanaNamespaceAndAddonName, 0, targetPort)
	defer tunnel.Close()

	tunnel.ForwardPort(t)

	healthCheck, err := http.NewRequest("GET", fmt.Sprintf("http://%s/healthz", tunnel.Endpoint()), nil)
	if err != nil {
		log.Fatalf("Error when building the request: %s", reqErr)
	}

	resp, reqErr := http.DefaultClient.Do(healthCheck)
	if reqErr != nil {
		log.Fatalf("Error when making the request: %s", reqErr)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Grafana is not ready for metrics: status code %d", resp.StatusCode)
	} else {
		log.Print("Grafana is ready for metrics!")
	}
	// ----------------------------------

	destroyHelmEnvironment(t, addonData, helmOptions)
	destroyHelmEnvironment(t, istioBaseAddonData, istioBaseHelmOptions)
}
