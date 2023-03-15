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
	t.Run("Parallel", func(t *testing.T) {
		t.Run("respMsg", testResponseMessage)
		t.Run("version", testVersion)
		t.Run("statusCode", testStatusCode)
	})

	// =============================================================
	k8s.KubectlDelete(t, kubectlOptions, "sample-podinfo-app.yaml")
}

func testStatusCode(t *testing.T) {
	http_helper.HTTPDoWithRetry(
		t,
		"GET",
		fmt.Sprintf("http://%s/status/504", os.Getenv("KKA_LB_ENDPOINT")),
		nil,
		map[string]string{"Host": "sample-podinfo-app.default.example.com"},
		504,
		30,
		3*time.Second,
		nil,
	)
}

func testVersion(t *testing.T) {
	out := http_helper.HTTPDoWithRetry(
		t,
		"GET",
		fmt.Sprintf("http://%s/version", os.Getenv("KKA_LB_ENDPOINT")),
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

	actualValue, exists := actual["version"]

	assert.True(t, exists)
	assert.Equal(t, "6.3.5", actualValue)

	actualValue, exists = actual["commit"]

	assert.True(t, exists)
	assert.Equal(t, "67e2c98a60dc92283531412a9e604dd4bae005a9", actualValue)
}

func testResponseMessage(t *testing.T) {
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
}
