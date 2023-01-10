package test

import (
	"log"
	"testing"
)

func TestHelmIstioAddon(t *testing.T) {
	istioBaseAddonData := HelmAddonData{
		namespaceName:   "istio-system",
		releaseName:     "",
		dependencyRepo:  "",
		addonName:       "istio-base",
		addonAlias:      "",
		chartPath:       "../addons/istio/base",
		manageNamespace: true,
	}

	istioBaseHelmOptions, err := prepareHelmEnvironment(t, &istioBaseAddonData)

	if err != nil {
		log.Fatal(err)
	}

	// ----------------------------------

	istioDaemonAddonData := HelmAddonData{
		namespaceName:  "istio-system",
		releaseName:    "",
		dependencyRepo: "",
		addonName:      "istiod",
		addonAlias:     "",
		chartPath:      "../addons/istio/istiod",
	}

	istioDaemonHelmOptions, err := prepareHelmEnvironment(t, &istioDaemonAddonData)

	if err != nil {
		log.Fatal(err)
	}

	waitUntilServicesAvailable(t, *istioDaemonHelmOptions.KubectlOptions, []string{"istiod"})

	// ----------------------------------

	istioIngressAddonData := HelmAddonData{
		namespaceName:  "istio-system",
		releaseName:    "",
		dependencyRepo: "",
		addonName:      "istio-ingressgateway",
		addonAlias:     "",
		chartPath:      "../addons/istio/ingressgateway",
	}

	istioIngressHelmOptions, err := prepareHelmEnvironment(t, &istioIngressAddonData)

	if err != nil {
		log.Fatal(err)
	}

	waitUntilLoadBalancerAvailable(t, *istioIngressHelmOptions.KubectlOptions)

	// ----------------------------------

	destroyHelmEnvironment(t, istioIngressAddonData, istioIngressHelmOptions)
	destroyHelmEnvironment(t, istioDaemonAddonData, istioDaemonHelmOptions)
	destroyHelmEnvironment(t, istioBaseAddonData, istioBaseHelmOptions)
}
