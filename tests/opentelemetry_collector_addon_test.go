package test

import (
	"log"
	"testing"

	_ "github.com/stretchr/testify/assert"
)

func TestHelmOpentelemetryCollectorAddon(t *testing.T) {
	addonData := HelmAddonData{
		namespaceName:   "opentelemetry",
		releaseName:     "",
		dependencyRepo:  "",
		addonName:       "opentelemetry-collector",
		addonAlias:      "",
		chartPath:       "",
		hasCustomValues: true,
		manageNamespace: true,
	}

	helmOptions, err := prepareHelmEnvironment(t, &addonData)

	if err != nil {
		log.Fatal(err)
	}

	waitUntilDaemonSetAvailable(t, *helmOptions.KubectlOptions, "otlp-agent")

	// ----------------------------------

	destroyHelmEnvironment(t, addonData, helmOptions)
}
