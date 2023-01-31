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

	istioDiscoveryAddonData := HelmAddonData{
		namespaceName:  "istio-system",
		releaseName:    "",
		dependencyRepo: "",
		addonName:      "istiod",
		addonAlias:     "",
		chartPath:      "../addons/istio/istiod",
	}

	istioDiscoveryHelmOptions, err := prepareHelmEnvironment(t, &istioDiscoveryAddonData)

	if err != nil {
		log.Fatal(err)
	}

	waitUntilServicesAvailable(t, *istioDiscoveryHelmOptions.KubectlOptions, []string{"istiod"})

	// ----------------------------------

	istioIngressAddonData := HelmAddonData{
		namespaceName:   "istio-ingress",
		releaseName:     "",
		dependencyRepo:  "",
		addonName:       "istio-ingressgateway",
		addonAlias:      "",
		chartPath:       "../addons/istio/ingressgateway",
		manageNamespace: true,
	}

	istioIngressHelmOptions, err := prepareHelmEnvironment(t, &istioIngressAddonData)

	if err != nil {
		log.Fatal(err)
	}

	waitUntilLoadBalancerAvailable(t, *istioIngressHelmOptions.KubectlOptions)

	// ----------------------------------

	destroyHelmEnvironment(t, istioIngressAddonData, istioIngressHelmOptions)
	destroyHelmEnvironment(t, istioDiscoveryAddonData, istioDiscoveryHelmOptions)
	destroyHelmEnvironment(t, istioBaseAddonData, istioBaseHelmOptions)
}
