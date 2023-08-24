package test

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	_ "github.com/stretchr/testify/assert"

	"github.com/gruntwork-io/terratest/modules/k8s"
)

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

	addonData := HelmAddonData{
		namespaceName:   "",
		releaseName:     "",
		dependencyRepo:  "",
		addonName:       "grafana",
		addonAlias:      "",
		chartPath:       "",
		hasCustomValues: true,
		manageNamespace: true,
	}

	helmOptions, err := prepareHelmEnvironment(t, &addonData)

	if err != nil {
		log.Fatal(err)
	}

	waitUntilServicesAvailable(t, *helmOptions.KubectlOptions, []string{"grafana"})
	waitUntilDeploymentsAvailable(t, *helmOptions.KubectlOptions, []string{"grafana"})

	// valuesFile, err := os.ReadFile("../addons/grafana/values.yaml")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// yamlData := make(map[string]interface{})

	// err = yaml.Unmarshal(valuesFile, &yamlData)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// grafanaSource, ok := yamlData["grafana-source"].(map[string]interface{})
	// if !ok {
	// 	log.Fatalf("Error accessing 'grafana-source' key")
	// }

	// adminPassword, ok := grafanaSource["adminPassword"].(string)
	// if !ok {
	// 	log.Fatalf("Error accessing 'adminPassword' key")
	// }

	grafanaKubectlOptions := k8s.NewKubectlOptions("", "", "grafana")

	tunnel := k8s.NewTunnel(grafanaKubectlOptions, k8s.ResourceTypeService, "grafana", 0, 3000)
	defer tunnel.Close()

	tunnel.ForwardPort(t)

	healthCheck, err := http.NewRequest("GET", fmt.Sprintf("http://%s/healthz", tunnel.Endpoint()), nil)
	// healthCheck.Header.Add("Authorization", fmt.Sprintf("Basic %s", b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("admin:%s", adminPassword)))))

	resp, reqErr := http.DefaultClient.Do(healthCheck)
	if reqErr != nil {
		log.Fatalf("Error when making the request %s", reqErr)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Grafana is not ready for metrics: status code %d", resp.StatusCode)
	} else {
		log.Print("Grafana is ready for metrics!")
	}
	// ----------------------------------

	destroyHelmEnvironment(t, addonData, helmOptions)
	destroyHelmEnvironment(t, istioBaseAddonData, istioBaseHelmOptions)
}
