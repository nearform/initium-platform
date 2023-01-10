package test

import (
	"log"
	"testing"

	_ "github.com/stretchr/testify/assert"
)

func TestHelmDexAddon(t *testing.T) {
	addonData := HelmAddonData{
		namespaceName:   "",
		releaseName:     "",
		dependencyRepo:  "",
		addonName:       "dex",
		addonAlias:      "",
		hasCustomValues: true,
		manageNamespace: true,
	}

	helmOptions, err := prepareHelmEnvironment(t, &addonData)

	if err != nil {
		log.Fatal(err)
	}

	waitUntilHelmFormattedServicesAvailable(t, addonData, helmOptions, []string{"dex-source"})

	// ----------------------------------

	destroyHelmEnvironment(t, addonData, helmOptions)
}
