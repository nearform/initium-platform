package test

import (
	"log"
	"testing"

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

	helmOptions, err := prepareHelmEnvironment(t, &addonData)

	if err != nil {
		log.Fatal(err)
	}

	waitUntilHelmFormattedServicesAvailable(t, addonData, helmOptions, []string{
		"grafana",
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
