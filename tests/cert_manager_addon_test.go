package test

import (
	"log"
	"testing"

	_ "github.com/stretchr/testify/assert"
)

func TestHelmCertManagerAddon(t *testing.T) {
	addonData := HelmAddonData{
		namespaceName:   "",
		releaseName:     "",
		dependencyRepo:  "",
		addonName:       "cert-manager",
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
		"cert-manager-source",
		"cert-manager-source-webhook",
	})

	// ----------------------------------

	destroyHelmEnvironment(t, addonData, helmOptions)
}
