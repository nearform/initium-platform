package test

import (
	"log"
	"testing"

	_ "github.com/stretchr/testify/assert"
)

func TestHelmDexAddon(t *testing.T) {
	// add Istio CRDs to test Istio virtual service
	istioBaseAddonData := HelmAddonData{
		namespaceName:   "dex-test",
		releaseName:     "",
		dependencyRepo:  "",
		addonName:       "istio-base",
		addonAlias:      "",
		chartPath:       "../addons/istio/base",
		manageNamespace: true,
		overrideValues: map[string]string{
			"global.istioNamespace": "dex-test",
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
		addonName:       "dex",
		addonAlias:      "",
		chartPath:       "",
		hasCustomValues: true,
		manageNamespace: true,
	}

	helmOptions, err := prepareHelmEnvironment(t, &addonData)

	if err != nil {
		log.Fatal(err)
	}

	waitUntilServicesAvailable(t, *helmOptions.KubectlOptions, []string{"dex"})

	// ----------------------------------

	destroyHelmEnvironment(t, addonData, helmOptions)
	destroyHelmEnvironment(t, istioBaseAddonData, istioBaseHelmOptions)
}
