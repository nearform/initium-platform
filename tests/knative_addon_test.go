package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
)

func TestKnativeOperatorAddon(t *testing.T) {
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
}
