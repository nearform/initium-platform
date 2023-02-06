package test

import (
	"log"
	"testing"

	_ "github.com/stretchr/testify/assert"
)

func TestHelmLokiAddon(t *testing.T) {
	addonData := HelmAddonData{
		namespaceName:   "",
		releaseName:     "",
		dependencyRepo:  "",
		addonName:       "loki",
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
		"loki",
		"loki-headless",
		"loki-memberlist",
	})

	// ----------------------------------

	destroyHelmEnvironment(t, addonData, helmOptions)
}
