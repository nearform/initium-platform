package test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"

	_ "github.com/stretchr/testify/assert"

	"github.com/gruntwork-io/terratest/modules/k8s"

	http_helper "github.com/gruntwork-io/terratest/modules/http-helper"
)

func TestSamplePodInfoApp(t *testing.T) {
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
	out := http_helper.HTTPDoWithRetry(
		t,
		"GET",
		fmt.Sprintf("http://%s", os.Getenv("KKA_LB_ENDPOINT")),
		nil,
		map[string]string{"Host": "sample-podinfo-app.default.example.com"},
		200,
		30,
		3*time.Second,
		nil,
	)

	var actual map[string]string
	err := json.Unmarshal([]byte(out), &actual)
	if err != nil {
		t.Fatal("Failed to unmarshal response body", err)
	}

	actualValue, exists := actual["message"]

	assert.True(t, exists)
	assert.Equal(t, "greetings from podinfo v6.3.5", actualValue)

	// =============================================================
	k8s.KubectlDelete(t, kubectlOptions, "sample-podinfo-app.yaml")
}
