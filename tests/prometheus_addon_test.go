package test

import (
	"log"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	_ "github.com/stretchr/testify/assert"
)

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
	defer func() {
		k8s.DeleteNamespace(t, kubectlOptions, "grafana")
	}()

	helmOptions, err := prepareHelmEnvironment(t, &addonData)

	if err != nil {
		log.Fatal(err)
	}

	waitUntilHelmFormattedServicesAvailable(t, addonData, helmOptions, []string{
		"kube-state-metrics",
		"prometheus-node-exporter",
	})

	waitUntilServicesAvailable(t, *helmOptions.KubectlOptions, []string{
		"kube-prometheus-stack-test-alertmanager",
		"kube-prometheus-stack-test-operator",
		"kube-prometheus-stack-test-prometheus",
	})

	// ----------------------------------

	destroyHelmEnvironment(t, addonData, helmOptions)
}
