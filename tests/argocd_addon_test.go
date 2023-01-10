package test

import (
	"log"
	"testing"

	_ "github.com/stretchr/testify/assert"
)

func TestHelmArgoCdServerAddon(t *testing.T) {
	addonData := HelmAddonData{
		namespaceName:   "",
		releaseName:     "",
		dependencyRepo:  "",
		addonName:       "argocd",
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
		"redis",
		"repo-server",
		"server",
	})

	// ----------------------------------

	destroyHelmEnvironment(t, addonData, helmOptions)
}
