package test

import (
	"log"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
)

func TestKnativeAndHelmIstioAddons(t *testing.T) {
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
	operatorResourcePath := "../addons/knative/templates/operator.yaml"

	operatorOptions := k8s.NewKubectlOptions("", "", "default")

	k8s.KubectlApply(t, operatorOptions, operatorResourcePath)

	waitUntilServicesAvailable(t, *operatorOptions, []string{"operator-webhook"})

	// ----------------------------------

	servingResourcePath := "../addons/knative/templates/serving.yaml"

	servingOptions := k8s.NewKubectlOptions("", "", "knative-serving")

	k8s.KubectlApply(t, servingOptions, servingResourcePath)

	waitUntilServicesAvailable(t, *servingOptions, []string{
		"activator-service",
		"autoscaler",
		"controller",
		"webhook",
	})

	// ----------------------------------

	eventingResourcePath := "../addons/knative/templates/eventing.yaml"

	eventingOptions := k8s.NewKubectlOptions("", "", "knative-eventing")

	k8s.KubectlApply(t, eventingOptions, eventingResourcePath)

	waitUntilServicesAvailable(t, *eventingOptions, []string{
		"broker-filter",
		"broker-ingress",
		"eventing-webhook",
		"imc-dispatcher",
		"inmemorychannel-webhook",
	})

	// ----------------------------------

	k8s.KubectlDelete(t, eventingOptions, eventingResourcePath)
	k8s.KubectlDelete(t, servingOptions, servingResourcePath)
	k8s.KubectlDelete(t, operatorOptions, operatorResourcePath)

	destroyHelmEnvironment(t, istioIngressAddonData, istioIngressHelmOptions)
	destroyHelmEnvironment(t, istioDiscoveryAddonData, istioDiscoveryHelmOptions)
	destroyHelmEnvironment(t, istioBaseAddonData, istioBaseHelmOptions)
}
