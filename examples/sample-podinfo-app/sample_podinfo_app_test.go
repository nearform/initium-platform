package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/stretchr/testify/assert"

	"github.com/gruntwork-io/terratest/modules/k8s"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
)

func TestExampleApp(t *testing.T) {
	// =============================================================
	kubectlOptions := k8s.NewKubectlOptions("", "", "istio-ingress")

	k8s.WaitUntilServiceAvailable(t, kubectlOptions, "istio-ingressgateway", 10, 30*time.Second)

	// =============================================================
	kubectlOptions = k8s.NewKubectlOptions("", "", "knative-serving")

	k8s.WaitUntilServiceAvailable(t, kubectlOptions, "controller", 10, 5*time.Second)

	// =============================================================
	kubectlOptions = k8s.NewKubectlOptions("", "", "default")

	k8s.KubectlApply(t, kubectlOptions, "sample-podinfo-app.yaml")

	// =============================================================
	http_helper.HTTPDoWithRetry(
		t,
		"GET",
		fmt.Sprintf("http://%s", os.Getenv("KKA_LB_ENDPOINT")),
		[]byte("Hello world from k8s-kurated-addons!"),
		map[string]string{"Host": "sample-podinfo-app.default.example.com"},
		200,
		30,
		3*time.Second,
		nil,
	)

	// =============================================================
	k8s.KubectlDelete(t, kubectlOptions, "sample-podinfo-app.yaml")
}
